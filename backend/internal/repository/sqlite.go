package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/JuanmaZimm/url-4-u/backend/internal/model"
	_ "modernc.org/sqlite"
)

var ErrNotFound = errors.New("url not found")
var ErrAlreadyExists = errors.New("short code or alias already exists")

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(dbPath string) (*SQLiteRepository, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite database: %w", err)
	}

	repo := &SQLiteRepository{db: db}
	if err := repo.initSchema(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return repo, nil
}

func (r *SQLiteRepository) initSchema() error {
	queries := []string{
		`PRAGMA foreign_keys = ON;`,
		`PRAGMA journal_mode = WAL;`,
		`CREATE TABLE IF NOT EXISTS urls (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			original_url TEXT NOT NULL,
			short_code TEXT NOT NULL UNIQUE,
			custom_alias TEXT,
			click_count INTEGER NOT NULL DEFAULT 0,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS click_events (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			url_id INTEGER NOT NULL,
			short_code TEXT NOT NULL,
			clicked_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			ip_hash TEXT,
			user_agent TEXT,
			referrer TEXT,
			device TEXT,
			browser TEXT,
			FOREIGN KEY (url_id) REFERENCES urls(id) ON DELETE CASCADE
		);`,
		`CREATE INDEX IF NOT EXISTS idx_urls_short_code ON urls(short_code);`,
		`CREATE INDEX IF NOT EXISTS idx_click_events_short_code ON click_events(short_code);`,
		`CREATE INDEX IF NOT EXISTS idx_click_events_clicked_at ON click_events(clicked_at);`,
	}

	for _, query := range queries {
		if _, err := r.db.Exec(query); err != nil {
			return err
		}
	}

	return nil
}

func (r *SQLiteRepository) Save(ctx context.Context, u *model.URL) error {
	query := `INSERT INTO urls (original_url, short_code, custom_alias, click_count, created_at, updated_at)
			  VALUES (?, ?, ?, ?, ?, ?)`
	now := time.Now().UTC()
	u.CreatedAt = now
	u.UpdatedAt = now

	res, err := r.db.ExecContext(ctx, query, u.OriginalURL, u.ShortCode, u.CustomAlias, u.ClickCount, now, now)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return fmt.Errorf("%w: %v", ErrAlreadyExists, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	u.ID = id
	return nil
}

func (r *SQLiteRepository) FindByShortCode(ctx context.Context, shortCode string) (*model.URL, error) {
	query := `SELECT id, original_url, short_code, custom_alias, click_count, created_at, updated_at
			  FROM urls WHERE short_code = ? LIMIT 1`
	row := r.db.QueryRowContext(ctx, query, shortCode)

	var u model.URL
	var alias sql.NullString
	err := row.Scan(&u.ID, &u.OriginalURL, &u.ShortCode, &alias, &u.ClickCount, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if alias.Valid {
		u.CustomAlias = alias.String
	}
	return &u, nil
}

func (r *SQLiteRepository) FindByOriginalURL(ctx context.Context, originalURL string) (*model.URL, error) {
	query := `SELECT id, original_url, short_code, custom_alias, click_count, created_at, updated_at
			  FROM urls WHERE original_url = ? AND (custom_alias IS NULL OR custom_alias = '') ORDER BY id DESC LIMIT 1`
	row := r.db.QueryRowContext(ctx, query, originalURL)

	var u model.URL
	var alias sql.NullString
	err := row.Scan(&u.ID, &u.OriginalURL, &u.ShortCode, &alias, &u.ClickCount, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if alias.Valid {
		u.CustomAlias = alias.String
	}
	return &u, nil
}

func (r *SQLiteRepository) FindRecent(ctx context.Context, limit, offset int) ([]*model.URL, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	query := `SELECT id, original_url, short_code, custom_alias, click_count, created_at, updated_at
			  FROM urls ORDER BY id DESC LIMIT ? OFFSET ?`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var urls []*model.URL
	for rows.Next() {
		var u model.URL
		var alias sql.NullString
		if err := rows.Scan(&u.ID, &u.OriginalURL, &u.ShortCode, &alias, &u.ClickCount, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		if alias.Valid {
			u.CustomAlias = alias.String
		}
		urls = append(urls, &u)
	}

	return urls, nil
}

func (r *SQLiteRepository) RecordClick(ctx context.Context, event *model.ClickEvent) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	now := time.Now().UTC()
	event.ClickedAt = now

	// Insert click event
	insertQuery := `INSERT INTO click_events (url_id, short_code, clicked_at, ip_hash, user_agent, referrer, device, browser)
					VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	_, err = tx.ExecContext(ctx, insertQuery, event.URLID, event.ShortCode, event.ClickedAt.Format("2006-01-02 15:04:05"), event.IPHash, event.UserAgent, event.Referrer, event.Device, event.Browser)
	if err != nil {
		return err
	}

	// Increment counter
	updateQuery := `UPDATE urls SET click_count = click_count + 1, updated_at = ? WHERE id = ?`
	_, err = tx.ExecContext(ctx, updateQuery, now, event.URLID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *SQLiteRepository) GetAnalytics(ctx context.Context, shortCode string) (*model.URLAnalytics, error) {
	u, err := r.FindByShortCode(ctx, shortCode)
	if err != nil {
		return nil, err
	}

	analytics := &model.URLAnalytics{
		ShortCode:    u.ShortCode,
		OriginalURL:  u.OriginalURL,
		TotalClicks:  u.ClickCount,
		ClicksByDate: make(map[string]int64),
		TopReferrers: make(map[string]int64),
		TopDevices:   make(map[string]int64),
		TopBrowsers:  make(map[string]int64),
	}

	// 1. Clicks by date (last 30 days)
	dateQuery := `SELECT substr(clicked_at, 1, 10) AS click_date, COUNT(*)
				  FROM click_events WHERE short_code = ?
				  GROUP BY click_date ORDER BY click_date DESC LIMIT 30`
	rows, err := r.db.QueryContext(ctx, dateQuery, shortCode)
	if err == nil {
		for rows.Next() {
			var date string
			var count int64
			if err := rows.Scan(&date, &count); err == nil {
				analytics.ClicksByDate[date] = count
			}
		}
		rows.Close()
	}

	// 2. Top Referrers
	refQuery := `SELECT COALESCE(NULLIF(referrer, ''), 'Direct') as ref, COUNT(*)
				 FROM click_events WHERE short_code = ?
				 GROUP BY ref ORDER BY COUNT(*) DESC LIMIT 5`
	rows, err = r.db.QueryContext(ctx, refQuery, shortCode)
	if err == nil {
		for rows.Next() {
			var ref string
			var count int64
			if err := rows.Scan(&ref, &count); err == nil {
				analytics.TopReferrers[ref] = count
			}
		}
		rows.Close()
	}

	// 3. Devices
	devQuery := `SELECT COALESCE(NULLIF(device, ''), 'Unknown') as dev, COUNT(*)
				 FROM click_events WHERE short_code = ?
				 GROUP BY dev ORDER BY COUNT(*) DESC LIMIT 5`
	rows, err = r.db.QueryContext(ctx, devQuery, shortCode)
	if err == nil {
		for rows.Next() {
			var dev string
			var count int64
			if err := rows.Scan(&dev, &count); err == nil {
				analytics.TopDevices[dev] = count
			}
		}
		rows.Close()
	}

	// 4. Browsers
	browserQuery := `SELECT COALESCE(NULLIF(browser, ''), 'Unknown') as br, COUNT(*)
					 FROM click_events WHERE short_code = ?
					 GROUP BY br ORDER BY COUNT(*) DESC LIMIT 5`
	rows, err = r.db.QueryContext(ctx, browserQuery, shortCode)
	if err == nil {
		for rows.Next() {
			var br string
			var count int64
			if err := rows.Scan(&br, &count); err == nil {
				analytics.TopBrowsers[br] = count
			}
		}
		rows.Close()
	}

	// 5. Last clicked timestamp
	lastClickedQuery := `SELECT clicked_at FROM click_events WHERE short_code = ? ORDER BY clicked_at DESC LIMIT 1`
	var lastTime time.Time
	if err := r.db.QueryRowContext(ctx, lastClickedQuery, shortCode).Scan(&lastTime); err == nil {
		analytics.LastClicked = &lastTime
	}

	return analytics, nil
}

func (r *SQLiteRepository) Close() error {
	return r.db.Close()
}
