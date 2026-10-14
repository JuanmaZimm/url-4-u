package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/juan-manuel-zimmermann/url-4-u/backend/internal/model"
	"github.com/juan-manuel-zimmermann/url-4-u/backend/internal/repository"
	"github.com/juan-manuel-zimmermann/url-4-u/backend/pkg/base62"
)

var (
	ErrInvalidURL       = errors.New("invalid original url: must be http or https")
	ErrInvalidAlias     = errors.New("invalid custom alias: must be 3-30 alphanumeric characters, hyphens or underscores")
	ErrAliasReserved    = errors.New("custom alias is a reserved system path")
	ErrURLAlreadyExists = errors.New("alias or short code already in use")
)

var aliasRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{3,30}$`)

var reservedSlugs = map[string]bool{
	"api":       true,
	"health":    true,
	"r":         true,
	"admin":     true,
	"analytics": true,
	"qr":        true,
	"static":    true,
	"favicon":   true,
}

type URLService struct {
	repo    repository.URLRepository
	baseURL string
}

func NewURLService(repo repository.URLRepository, baseURL string) *URLService {
	return &URLService{
		repo:    repo,
		baseURL: strings.TrimRight(baseURL, "/"),
	}
}

// ValidateURL ensures the URL is well-formed and safe
func (s *URLService) ValidateURL(rawURL string) (string, error) {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return "", ErrInvalidURL
	}

	parsed, err := url.ParseRequestURI(rawURL)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" {
		return "", ErrInvalidURL
	}

	// Avoid self-referencing loops
	if baseParsed, err := url.Parse(s.baseURL); err == nil && baseParsed.Host != "" {
		if strings.EqualFold(parsed.Host, baseParsed.Host) {
			return "", errors.New("cannot shorten URL pointing to URL-4-U itself")
		}
	}

	return parsed.String(), nil
}

// Shorten creates a new short code or assigns the custom alias
func (s *URLService) Shorten(ctx context.Context, originalURL string, customAlias string) (*model.URL, error) {
	validURL, err := s.ValidateURL(originalURL)
	if err != nil {
		return nil, err
	}

	var shortCode string
	customAlias = strings.TrimSpace(customAlias)

	if customAlias != "" {
		if !aliasRegex.MatchString(customAlias) {
			return nil, ErrInvalidAlias
		}
		if reservedSlugs[strings.ToLower(customAlias)] {
			return nil, ErrAliasReserved
		}

		// Verify availability
		existing, _ := s.repo.FindByShortCode(ctx, customAlias)
		if existing != nil {
			return nil, ErrURLAlreadyExists
		}
		shortCode = customAlias
	} else {
		// Generate random Base62 slug with collision check
		for attempts := 0; attempts < 5; attempts++ {
			code, err := base62.GenerateRandom(6)
			if err != nil {
				return nil, fmt.Errorf("failed to generate random slug: %w", err)
			}
			if reservedSlugs[strings.ToLower(code)] {
				continue
			}
			existing, _ := s.repo.FindByShortCode(ctx, code)
			if existing == nil {
				shortCode = code
				break
			}
		}

		if shortCode == "" {
			return nil, errors.New("could not generate unique short code, please try again")
		}
	}

	urlModel := &model.URL{
		OriginalURL: validURL,
		ShortCode:   shortCode,
		CustomAlias: customAlias,
	}

	if err := s.repo.Save(ctx, urlModel); err != nil {
		return nil, err
	}

	return urlModel, nil
}

// Resolve fetches the destination URL and schedules click tracking
func (s *URLService) Resolve(ctx context.Context, shortCode string, rawIP, userAgent, referrer string) (string, error) {
	u, err := s.repo.FindByShortCode(ctx, shortCode)
	if err != nil {
		return "", err
	}

	// Asynchronously record the click event to keep redirection latency minimal
	go func() {
		bgCtx := context.Background()
		event := &model.ClickEvent{
			URLID:     u.ID,
			ShortCode: u.ShortCode,
			IPHash:    hashIP(rawIP),
			UserAgent: userAgent,
			Referrer:  cleanReferrer(referrer),
			Device:    parseDevice(userAgent),
			Browser:   parseBrowser(userAgent),
		}
		_ = s.repo.RecordClick(bgCtx, event)
	}()

	return u.OriginalURL, nil
}

func (s *URLService) GetDetails(ctx context.Context, shortCode string) (*model.URL, error) {
	return s.repo.FindByShortCode(ctx, shortCode)
}

func (s *URLService) GetRecent(ctx context.Context, limit, offset int) ([]*model.URL, error) {
	return s.repo.FindRecent(ctx, limit, offset)
}

func (s *URLService) GetAnalytics(ctx context.Context, shortCode string) (*model.URLAnalytics, error) {
	return s.repo.GetAnalytics(ctx, shortCode)
}

// Helpers for User-Agent / Analytics parsing
func hashIP(ip string) string {
	if ip == "" {
		return ""
	}
	hash := sha256.Sum256([]byte(ip))
	return hex.EncodeToString(hash[:8]) // First 16 hex chars for privacy-safe grouping
}

func cleanReferrer(ref string) string {
	if ref == "" {
		return "Direct"
	}
	u, err := url.Parse(ref)
	if err != nil || u.Host == "" {
		return "Direct"
	}
	return u.Host
}

func parseDevice(ua string) string {
	uaLower := strings.ToLower(ua)
	if strings.Contains(uaLower, "mobile") || strings.Contains(uaLower, "android") || strings.Contains(uaLower, "iphone") {
		return "Mobile"
	}
	if strings.Contains(uaLower, "tablet") || strings.Contains(uaLower, "ipad") {
		return "Tablet"
	}
	if ua == "" {
		return "Unknown"
	}
	return "Desktop"
}

func parseBrowser(ua string) string {
	uaLower := strings.ToLower(ua)
	if strings.Contains(uaLower, "edg/") {
		return "Edge"
	}
	if strings.Contains(uaLower, "chrome/") && !strings.Contains(uaLower, "edg/") {
		return "Chrome"
	}
	if strings.Contains(uaLower, "firefox/") {
		return "Firefox"
	}
	if strings.Contains(uaLower, "safari/") && !strings.Contains(uaLower, "chrome/") {
		return "Safari"
	}
	if strings.Contains(uaLower, "curl") || strings.Contains(uaLower, "wget") || strings.Contains(uaLower, "postman") {
		return "Bot/CLI"
	}
	if ua == "" {
		return "Unknown"
	}
	return "Other"
}
