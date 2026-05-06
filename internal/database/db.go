// Package database implements a wrapper around the real database
// with some helper functions to store catalog / log entries
package database

import (
	"fmt"
	"io"
	"time"

	"github.com/glebarez/sqlite"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type (
	// Client represents a database client
	Client struct {
		Catalog CatalogMetaStore
		Logs    LogStore

		db *gorm.DB
	}

	logwrap struct {
		l *io.PipeWriter
	}

	migrator interface {
		Migrate(dest *Client) error
	}
)

// NewClient creates a new Client and connects to the database using
// some default configurations. The database is automatically
// initialized with required tables.
func NewClient(dbtype, dsn string) (*Client, error) {
	c := &Client{}
	c.Catalog = CatalogMetaStore{c}
	c.Logs = LogStore{c}

	dbLogger := logger.New(
		&logwrap{log.StandardLogger().WriterLevel(log.TraceLevel)},
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Silent,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	switch dbtype {
	case "mysql":
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: dbLogger,
		})
		if err != nil {
			return nil, fmt.Errorf("opening mysql database: %w", err)
		}
		c.db = db

	case "crdb", "postgres", "postgresql":
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: dbLogger,
		})
		if err != nil {
			return nil, fmt.Errorf("opening postgres database: %w", err)
		}
		c.db = db

	case "sqlite", "sqlite3":
		db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
			Logger: dbLogger,
		})
		if err != nil {
			return nil, fmt.Errorf("opening sqlite3 database: %w", err)
		}
		c.db = db

	default:
		return nil, fmt.Errorf("invalid db type: %s", dbtype)
	}

	if err := c.initDB(); err != nil {
		return nil, fmt.Errorf("initializing database: %w", err)
	}

	return c, nil
}

// Migrate executes database migrations for all required types
func (c Client) Migrate(dest *Client) error {
	for _, m := range []migrator{
		c.Catalog,
		c.Logs,
	} {
		if err := m.Migrate(dest); err != nil {
			return fmt.Errorf("executing migrating: %w", err)
		}
	}

	return nil
}

func (c Client) initDB() error {
	for name, fn := range map[string]func() error{
		"catalogMeta": c.Catalog.ensureTable,
		"log":         c.Logs.ensureTable,
	} {
		if err := fn(); err != nil {
			return fmt.Errorf("ensuring table %q: %w", name, err)
		}
	}

	return nil
}

func (l logwrap) Printf(f string, v ...any) {
	fmt.Fprintf(l.l, f, v...) //nolint:errcheck // only logging
}
