package repository

import (
	"context"
	"github.com/JuanmaZimm/url-4-u/backend/internal/model"
)

// URLRepository defines the storage interface for URLs and analytics
type URLRepository interface {
	Save(ctx context.Context, url *model.URL) error
	FindByShortCode(ctx context.Context, shortCode string) (*model.URL, error)
	FindByOriginalURL(ctx context.Context, originalURL string) (*model.URL, error)
	FindRecent(ctx context.Context, limit, offset int) ([]*model.URL, error)
	RecordClick(ctx context.Context, event *model.ClickEvent) error
	GetAnalytics(ctx context.Context, shortCode string) (*model.URLAnalytics, error)
	Close() error
}
