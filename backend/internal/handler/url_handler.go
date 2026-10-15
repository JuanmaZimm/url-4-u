package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/JuanmaZimm/url-4-u/backend/internal/model"
	"github.com/JuanmaZimm/url-4-u/backend/internal/repository"
	"github.com/JuanmaZimm/url-4-u/backend/internal/service"
)

type URLHandler struct {
	urlService *service.URLService
	qrService  *service.QRService
	baseURL    string
}

func NewURLHandler(urlService *service.URLService, qrService *service.QRService, baseURL string) *URLHandler {
	return &URLHandler{
		urlService: urlService,
		qrService:  qrService,
		baseURL:    strings.TrimRight(baseURL, "/"),
	}
}

// JSON helpers
func respondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}

// Create handles POST /api/v1/urls
func (h *URLHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateURLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid JSON request body")
		return
	}

	u, err := h.urlService.Shorten(r.Context(), req.OriginalURL, req.CustomAlias)
	if err != nil {
		if errors.Is(err, service.ErrInvalidURL) || errors.Is(err, service.ErrInvalidAlias) || errors.Is(err, service.ErrAliasReserved) {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, service.ErrURLAlreadyExists) || errors.Is(err, repository.ErrAlreadyExists) {
			respondError(w, http.StatusConflict, "Short code or alias is already in use")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to create short URL")
		return
	}

	resp := model.CreateURLResponse{
		ShortCode:   u.ShortCode,
		ShortURL:    fmt.Sprintf("%s/r/%s", h.baseURL, u.ShortCode),
		OriginalURL: u.OriginalURL,
		QRCodeURL:   fmt.Sprintf("%s/api/v1/urls/%s/qr", h.baseURL, u.ShortCode),
		CreatedAt:   u.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	respondJSON(w, http.StatusCreated, resp)
}

// Redirect handles GET /r/{code}
func (h *URLHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	if code == "" {
		respondError(w, http.StatusBadRequest, "Missing short code")
		return
	}

	// Extract client IP (respecting proxy headers)
	clientIP := r.Header.Get("X-Forwarded-For")
	if clientIP == "" {
		clientIP = r.Header.Get("X-Real-IP")
	}
	if clientIP == "" {
		clientIP = r.RemoteAddr
	}
	if strings.Contains(clientIP, ",") {
		clientIP = strings.TrimSpace(strings.Split(clientIP, ",")[0])
	}

	userAgent := r.UserAgent()
	referrer := r.Referer()

	destURL, err := h.urlService.Resolve(r.Context(), code, clientIP, userAgent, referrer)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			http.Error(w, "URL-4-U: Short link not found or expired", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, destURL, http.StatusFound)
}

// Get handles GET /api/v1/urls/{code}
func (h *URLHandler) Get(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	u, err := h.urlService.GetDetails(r.Context(), code)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			respondError(w, http.StatusNotFound, "Short URL not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "Database error")
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{
		"short_code":   u.ShortCode,
		"short_url":    fmt.Sprintf("%s/r/%s", h.baseURL, u.ShortCode),
		"original_url": u.OriginalURL,
		"custom_alias": u.CustomAlias,
		"click_count":  u.ClickCount,
		"qr_code_url":  fmt.Sprintf("%s/api/v1/urls/%s/qr", h.baseURL, u.ShortCode),
		"created_at":   u.CreatedAt,
		"updated_at":   u.UpdatedAt,
	})
}

// List handles GET /api/v1/urls
func (h *URLHandler) List(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 10
	offset := 0

	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		limit = l
	}
	if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
		offset = o
	}

	urls, err := h.urlService.GetRecent(r.Context(), limit, offset)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to retrieve URLs")
		return
	}

	type ListItem struct {
		ShortCode   string `json:"short_code"`
		ShortURL    string `json:"short_url"`
		OriginalURL string `json:"original_url"`
		CustomAlias string `json:"custom_alias,omitempty"`
		ClickCount  int64  `json:"click_count"`
		QRCodeURL   string `json:"qr_code_url"`
		CreatedAt   string `json:"created_at"`
	}

	items := make([]ListItem, 0, len(urls))
	for _, u := range urls {
		items = append(items, ListItem{
			ShortCode:   u.ShortCode,
			ShortURL:    fmt.Sprintf("%s/r/%s", h.baseURL, u.ShortCode),
			OriginalURL: u.OriginalURL,
			CustomAlias: u.CustomAlias,
			ClickCount:  u.ClickCount,
			QRCodeURL:   fmt.Sprintf("%s/api/v1/urls/%s/qr", h.baseURL, u.ShortCode),
			CreatedAt:   u.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	respondJSON(w, http.StatusOK, map[string]any{
		"data":   items,
		"limit":  limit,
		"offset": offset,
	})
}

// GetQR handles GET /api/v1/urls/{code}/qr
func (h *URLHandler) GetQR(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	u, err := h.urlService.GetDetails(r.Context(), code)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			respondError(w, http.StatusNotFound, "Short URL not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to fetch URL")
		return
	}

	// We generate the QR code pointing to the redirect URL
	redirectURL := fmt.Sprintf("%s/r/%s", h.baseURL, u.ShortCode)

	size := 256
	if sizeParam := r.URL.Query().Get("size"); sizeParam != "" {
		if s, err := strconv.Atoi(sizeParam); err == nil && s >= 64 && s <= 1024 {
			size = s
		}
	}

	pngData, err := h.qrService.GeneratePNG(redirectURL, size)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to render QR code")
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(pngData)
}

// GetAnalytics handles GET /api/v1/urls/{code}/analytics
func (h *URLHandler) GetAnalytics(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	analytics, err := h.urlService.GetAnalytics(r.Context(), code)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			respondError(w, http.StatusNotFound, "Short URL not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to retrieve analytics")
		return
	}

	respondJSON(w, http.StatusOK, analytics)
}
