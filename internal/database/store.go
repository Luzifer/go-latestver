package database

import (
	_ "embed"
	"strings"
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/Luzifer/go_helpers/v2/fieldcollection"
)

type (
	CatalogEntry struct {
		Name string `json:"name" yaml:"name"`
		Tag  string `json:"tag" yaml:"tag"`

		Fetcher       string                          `json:"-" yaml:"fetcher"`
		FetcherConfig fieldcollection.FieldCollection `json:"-" yaml:"fetcher_config"`

		CheckInterval time.Duration `json:"-" yaml:"check_interval"`

		Links []CatalogLink `json:"links" yaml:"links"`
	}

	CatalogLink struct {
		IconClass string `json:"icon_class" yaml:"icon_class"`
		Name      string `json:"name" yaml:"name"`
		URL       string `json:"url" yaml:"url"`
	}

	CatalogMeta struct {
		CatalogName    string     `gorm:"primaryKey" json:"-"`
		CatalogTag     string     `gorm:"primaryKey" json:"-"`
		CurrentVersion string     `json:"current_version,omitempty"`
		Error          string     `json:"error,omitempty"`
		LastChecked    *time.Time `json:"last_checked,omitempty"`
		VersionTime    *time.Time `json:"version_time,omitempty"`
	}

	LogEntry struct {
		CatalogName string    `gorm:"index:catalog_key" json:"catalog_name"`
		CatalogTag  string    `gorm:"index:catalog_key" json:"catalog_tag"`
		Timestamp   time.Time `gorm:"index:,sort:desc" json:"timestamp"`
		VersionTo   string    `json:"version_to"`
		VersionFrom string    `json:"version_from"`
	}

	CatalogMetaStore struct {
		c *Client
	}

	LogStore struct {
		c *Client
	}
)

func (c CatalogEntry) Key() string { return strings.Join([]string{c.Name, c.Tag}, ":") }

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

func (c CatalogMetaStore) PutMeta(cm *CatalogMeta) error {
	return errors.Wrap(
		c.c.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(cm).Error,
		"writing catalog meta",
	)
}

func (c CatalogMetaStore) ensureTable() error {
	return errors.Wrap(c.c.db.AutoMigrate(&CatalogMeta{}), "applying migration")
}

func (l LogStore) Add(le *LogEntry) error {
	return errors.Wrap(
		l.c.db.Create(le).Error,
		"writing log entry",
	)
}

func (l LogStore) List(num, page int) ([]LogEntry, error) {
	return l.listWithFilter(l.c.db, num, page)
}

func (l LogStore) ListForCatalogEntry(ce *CatalogEntry, num, page int) ([]LogEntry, error) {
	return l.listWithFilter(l.c.db.Where(&LogEntry{CatalogName: ce.Name, CatalogTag: ce.Tag}), num, page)
}

func (l LogStore) listWithFilter(filter *gorm.DB, num, page int) ([]LogEntry, error) {
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
