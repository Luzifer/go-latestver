package database

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Luzifer/go_helpers/fieldcollection"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/Luzifer/go-latestver/internal/version"
)

type (
	// CatalogEntry represents the entry in the config file
	CatalogEntry struct {
		Name string `json:"name" yaml:"name"`
		Tag  string `json:"tag" yaml:"tag"`

		Fetcher       string                           `json:"-" yaml:"fetcher"`
		FetcherConfig *fieldcollection.FieldCollection `json:"-" yaml:"fetcher_config"`

		VersionConstraint *version.Constraint `json:"-" yaml:"version_constraint"`

		Links []CatalogLink `json:"links" yaml:"links"`
	}

	// CatalogLink represents a link assigned to a CatalogEntry
	CatalogLink struct {
		IconClass string `json:"icon_class" yaml:"icon_class"`
		Name      string `json:"name" yaml:"name"`
		URL       string `json:"url" yaml:"url"`
	}

	// CatalogMeta contains meta-information about the catalog entry
	CatalogMeta struct {
		CatalogName    string     `gorm:"primaryKey" json:"-"`
		CatalogTag     string     `gorm:"primaryKey" json:"-"`
		CurrentVersion string     `json:"current_version,omitempty"`
		Error          string     `json:"error,omitempty"`
		LastChecked    *time.Time `json:"last_checked,omitempty"`
		VersionTime    *time.Time `json:"version_time,omitempty"`
	}

	// LogEntry represents a single version change for a given catalog entry
	LogEntry struct {
		CatalogName string    `gorm:"index:catalog_key" json:"catalog_name"`
		CatalogTag  string    `gorm:"index:catalog_key" json:"catalog_tag"`
		Timestamp   time.Time `gorm:"index:,sort:desc" json:"timestamp"`
		VersionTo   string    `json:"version_to"`
		VersionFrom string    `json:"version_from"`
	}

	// CatalogMetaStore is an accessor for the meta store and wraps a Client
	CatalogMetaStore struct {
		c *Client
	}

	// LogStore is an accessor for the log store and wraps a Client
	LogStore struct {
		c *Client
	}
)

// Key returns the name / tag combination as a single key
func (c CatalogEntry) Key() string { return strings.Join([]string{c.Name, c.Tag}, ":") }

// GetMeta fetches the current database stored CatalogMeta for the CatalogEntry
func (c CatalogMetaStore) GetMeta(ce *CatalogEntry) (*CatalogMeta, error) {
	out := &CatalogMeta{
		CatalogName: ce.Name,
		CatalogTag:  ce.Tag,
	}

	err := c.c.db.
		Where(out).
		First(out).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// If there is no meta yet we just return empty meta
		err = nil
	}

	if err != nil {
		return nil, fmt.Errorf("querying metadata: %w", err)
	}

	return out, nil
}

// Migrate applies the updated database schema for the CatalogMetaStore
func (c CatalogMetaStore) Migrate(dest *Client) error {
	var metas []*CatalogMeta
	if err := c.c.db.Find(&metas).Error; err != nil {
		return fmt.Errorf("listing meta entries: %w", err)
	}

	for _, m := range metas {
		if err := dest.Catalog.PutMeta(m); err != nil {
			return fmt.Errorf("storing meta to dest database: %w", err)
		}
	}

	return nil
}

// PutMeta stores the updated CatalogMeta
func (c CatalogMetaStore) PutMeta(cm *CatalogMeta) error {
	if err := c.c.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(cm).Error; err != nil {
		return fmt.Errorf("writing catalog meta: %w", err)
	}

	return nil
}

func (c CatalogMetaStore) ensureTable() error {
	if err := c.c.db.AutoMigrate(&CatalogMeta{}); err != nil {
		return fmt.Errorf("applying migration: %w", err)
	}

	return nil
}

// Add creates a new LogEntry inside the LogStore
func (l LogStore) Add(le *LogEntry) error {
	if err := l.c.db.Create(le).Error; err != nil {
		return fmt.Errorf("writing log entry: %w", err)
	}

	return nil
}

// List retrieves unfiltered log entries by page
func (l LogStore) List(num, page int) ([]LogEntry, error) {
	return l.listWithFilter(l.c.db, num, page)
}

// ListForCatalogEntry retrieves filered log entries by page
func (l LogStore) ListForCatalogEntry(ce *CatalogEntry, num, page int) ([]LogEntry, error) {
	return l.listWithFilter(l.c.db.Where(&LogEntry{CatalogName: ce.Name, CatalogTag: ce.Tag}), num, page)
}

// Migrate applies the updated database schema for the LogStore
func (l LogStore) Migrate(dest *Client) error {
	var logs []*LogEntry
	if err := l.c.db.Find(&logs).Error; err != nil {
		return fmt.Errorf("listing log entries: %w", err)
	}

	for _, l := range logs {
		if err := dest.Logs.Add(l); err != nil {
			return fmt.Errorf("storing log to dest database: %w", err)
		}
	}

	return nil
}

func (l LogStore) ensureTable() error {
	if err := l.c.db.AutoMigrate(&LogEntry{}); err != nil {
		return fmt.Errorf("applying migration: %w", err)
	}

	return nil
}

func (LogStore) listWithFilter(filter *gorm.DB, num, page int) (out []LogEntry, err error) {
	if err = filter.
		Order("timestamp desc").
		Limit(num).Offset(num * page).
		Find(&out).
		Error; err != nil {
		return nil, fmt.Errorf("fetching log entries: %w", err)
	}

	return out, nil
}
