package service

import (
	"context"
	"testing"

	"github.com/JuanmaZimm/url-4-u/backend/internal/model"
)

type mockRepo struct {
	saved map[string]*model.URL
}

func (m *mockRepo) Save(ctx context.Context, u *model.URL) error {
	m.saved[u.ShortCode] = u
	return nil
}

func (m *mockRepo) FindByShortCode(ctx context.Context, code string) (*model.URL, error) {
	if u, ok := m.saved[code]; ok {
		return u, nil
	}
	return nil, nil
}

func (m *mockRepo) FindByOriginalURL(ctx context.Context, orig string) (*model.URL, error) {
	return nil, nil
}

func (m *mockRepo) FindRecent(ctx context.Context, limit, offset int) ([]*model.URL, error) {
	return nil, nil
}

func (m *mockRepo) RecordClick(ctx context.Context, event *model.ClickEvent) error {
	return nil
}

func (m *mockRepo) GetAnalytics(ctx context.Context, code string) (*model.URLAnalytics, error) {
	return nil, nil
}

func (m *mockRepo) Close() error {
	return nil
}

func TestURLService_ValidateURL(t *testing.T) {
	svc := NewURLService(&mockRepo{}, "http://localhost:8080")

	tests := []struct {
		url     string
		wantErr bool
	}{
		{"https://google.com", false},
		{"http://example.org/test?q=1", false},
		{"ftp://invalid.com", true},
		{"", true},
		{"not-a-url", true},
		{"http://localhost:8080/r/abc", true}, // loop prevention
	}

	for _, tt := range tests {
		_, err := svc.ValidateURL(tt.url)
		if (err != nil) != tt.wantErr {
			t.Errorf("ValidateURL(%q) err = %v, wantErr = %v", tt.url, err, tt.wantErr)
		}
	}
}

func TestURLService_ShortenWithAlias(t *testing.T) {
	repo := &mockRepo{saved: make(map[string]*model.URL)}
	svc := NewURLService(repo, "http://localhost:8080")

	ctx := context.Background()
	u, err := svc.Shorten(ctx, "https://github.com", "my-portfolio")
	if err != nil {
		t.Fatalf("Unexpected error shortening: %v", err)
	}

	if u.ShortCode != "my-portfolio" {
		t.Errorf("Expected short code 'my-portfolio', got %s", u.ShortCode)
	}

	// Try shortening with reserved alias
	_, err = svc.Shorten(ctx, "https://github.com", "health")
	if err != ErrAliasReserved {
		t.Errorf("Expected ErrAliasReserved, got %v", err)
	}
}
