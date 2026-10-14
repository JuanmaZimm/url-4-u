package model

import "time"

// URL represents a shortened link in the system
type URL struct {
	ID          int64     `json:"id"`
	OriginalURL string    `json:"original_url"`
	ShortCode   string    `json:"short_code"`
	CustomAlias string    `json:"custom_alias,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	ClickCount  int64     `json:"click_count"`
}

// CreateURLRequest represents payload received to create a short link
type CreateURLRequest struct {
	OriginalURL string `json:"original_url"`
	CustomAlias string `json:"custom_alias,omitempty"`
}

// CreateURLResponse represents the output returned to clients
type CreateURLResponse struct {
	ShortCode   string `json:"short_code"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
	QRCodeURL   string `json:"qr_code_url"`
	CreatedAt   string `json:"created_at"`
}
