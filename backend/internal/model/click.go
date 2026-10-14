package model

import "time"

// ClickEvent tracks individual visits to shortened links
type ClickEvent struct {
	ID        int64     `json:"id"`
	URLID     int64     `json:"url_id"`
	ShortCode string    `json:"short_code"`
	ClickedAt time.Time `json:"clicked_at"`
	IPHash    string    `json:"ip_hash,omitempty"`
	UserAgent string    `json:"user_agent,omitempty"`
	Referrer  string    `json:"referrer,omitempty"`
	Device    string    `json:"device,omitempty"`
	Browser   string    `json:"browser,omitempty"`
}

// URLAnalytics represents aggregated metrics for a specific URL
type URLAnalytics struct {
	ShortCode    string            `json:"short_code"`
	OriginalURL  string            `json:"original_url"`
	TotalClicks  int64             `json:"total_clicks"`
	ClicksByDate map[string]int64  `json:"clicks_by_date"`
	TopReferrers map[string]int64  `json:"top_referrers"`
	TopDevices   map[string]int64  `json:"top_devices"`
	TopBrowsers  map[string]int64  `json:"top_browsers"`
	LastClicked  *time.Time        `json:"last_clicked,omitempty"`
}
