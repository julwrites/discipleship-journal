package services

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"discipleship_journal_api/database"

	"golang.org/x/net/html"
)

// HTTPClient interface for mocking
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// BibleVersionService handles bible version operations.
type BibleVersionService interface {
	GetVersions(ctx context.Context) ([]BibleVersion, error)
	SyncVersions(ctx context.Context) error
	ScrapeVersions(ctx context.Context) (map[string]string, error)
}

type bibleVersionService struct {
	db         database.DBInterface
	httpClient HTTPClient
	scrapeURL  string
}

// NewBibleVersionService creates a new BibleVersionService.
func NewBibleVersionService(db database.DBInterface) BibleVersionService {
	return &bibleVersionService{
		db:         db,
		httpClient: &http.Client{Timeout: 30 * time.Second},
		scrapeURL:  "https://classic.biblegateway.com/versions/",
	}
}

// ScrapeVersions scrapes bible versions from BibleGateway.
func (s *bibleVersionService) ScrapeVersions(ctx context.Context) (map[string]string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.scrapeURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch versions page: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	versions := make(map[string]string)
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "select" {
			// Check if it's the correct select element
			isVersionSelect := false
			for _, a := range n.Attr {
				if (a.Key == "name" && a.Val == "version") || (a.Key == "class" && strings.Contains(a.Val, "search-dropdown")) {
					isVersionSelect = true
					break
				}
			}

			if isVersionSelect {
				// Iterate over options
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					if c.Type == html.ElementNode && c.Data == "option" {
						var value, label string
						for _, a := range c.Attr {
							if a.Key == "value" {
								value = a.Val
							}
						}
						// Get label from text node
						if c.FirstChild != nil && c.FirstChild.Type == html.TextNode {
							label = c.FirstChild.Data
						}

						if value != "" && label != "" {
							versions[value] = label
						}
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return versions, nil
}

// SyncVersions syncs scraped versions to the database.
func (s *bibleVersionService) SyncVersions(ctx context.Context) error {
	versions, err := s.ScrapeVersions(ctx)
	if err != nil {
		return err
	}

	// Use a transaction for bulk update
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	// We'll prepare a statement or just loop exec. Loop exec is fine for < 1000 items.
	// We use ON CONFLICT to update existing entries or insert new ones.
	query := `
		INSERT INTO bible_versions (name, abbreviation, updated_at)
		VALUES ($1, $2, NOW())
		ON CONFLICT (abbreviation) DO UPDATE
		SET name = EXCLUDED.name, updated_at = NOW();
	`

	for abbr, name := range versions {
		_, err := tx.Exec(ctx, query, name, abbr)
		if err != nil {
			return fmt.Errorf("failed to upsert version %s: %w", abbr, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// BibleVersion represents a bible version model.
type BibleVersion struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Abbreviation string    `json:"abbreviation"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// GetVersions returns all bible versions sorted by name.
func (s *bibleVersionService) GetVersions(ctx context.Context) ([]BibleVersion, error) {
	query := `SELECT id, name, abbreviation, created_at, updated_at FROM bible_versions ORDER BY name ASC`
	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query versions: %w", err)
	}
	defer rows.Close()

	var versions []BibleVersion
	for rows.Next() {
		var v BibleVersion
		err := rows.Scan(&v.ID, &v.Name, &v.Abbreviation, &v.CreatedAt, &v.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan version: %w", err)
		}
		versions = append(versions, v)
	}

	return versions, nil
}
