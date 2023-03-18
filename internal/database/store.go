package database

import (
	"strings"
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/Luzifer/go-latestver/internal/version"
	"github.com/Luzifer/go_helpers/v2/fieldcollection"
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

	return out, errors.Wrap(err, "querying metadata")
}

// Migrate applies the updated database schema for the CatalogMetaStore
func (c CatalogMetaStore) Migrate(dest *Client) error {
	var metas []*CatalogMeta
	if err := c.c.db.Find(&metas).Error; err != nil {
		return errors.Wrap(err, "listing meta entries")
	}

	for _, m := range metas {
		if err := dest.Catalog.PutMeta(m); err != nil {
			return errors.Wrap(err, "storing meta to dest database")
		}
	}

	return nil
}

// PutMeta stores the updated CatalogMeta
func (c CatalogMetaStore) PutMeta(cm *CatalogMeta) error {
	return errors.Wrap(
		c.c.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(cm).Error,
		"writing catalog meta",
	)
}

func (c CatalogMetaStore) ensureTable() error {
	return errors.Wrap(c.c.db.AutoMigrate(&CatalogMeta{}), "applying migration")
}

// Add creates a new LogEntry inside the LogStore
func (l LogStore) Add(le *LogEntry) error {
	return errors.Wrap(
		l.c.db.Create(le).Error,
		"writing log entry",
	)
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
		return errors.Wrap(err, "listing log entries")
	}

	for _, l := range logs {
		if err := dest.Logs.Add(l); err != nil {
			return errors.Wrap(err, "storing log to dest database")
		}
	}

	return nil
}

func (LogStore) listWithFilter(filter *gorm.DB, num, page int) ([]LogEntry, error) {
	var out []LogEntry
	return out, errors.Wrap(
		filter.
			Order("timestamp desc").
			Limit(num).Offset(num*page).
			Find(&out).
			Error,
		"fetching log entries",
	)
}

func (l LogStore) ensureTable() error {
	return errors.Wrap(l.c.db.AutoMigrate(&LogEntry{}), "applying migration")
}
