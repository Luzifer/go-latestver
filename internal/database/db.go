package database

import (
	"fmt"
	"io"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type (
	Client struct {
		Catalog CatalogMetaStore
		Logs    LogStore

		db *gorm.DB
	}

	logwrap struct {
		l *io.PipeWriter
	}
)

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
			return nil, errors.Wrap(err, "opening mysql database")
		}
		c.db = db

	case "sqlite", "sqlite3":
		db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
			Logger: dbLogger,
		})
		if err != nil {
			return nil, errors.Wrap(err, "opening sqlite3 database")
		}
		c.db = db

	default:
		return nil, errors.Errorf("invalid db type: %s", dbtype)
	}

	if err := c.initDB(); err != nil {
		return nil, errors.Wrap(err, "initializing database")
	}

	return c, nil
}

func (c Client) initDB() error {
	for name, fn := range map[string]func() error{
		"catalogMeta": c.Catalog.ensureTable,
		"log":         c.Logs.ensureTable,
	} {
		if err := fn(); err != nil {
			return errors.Wrapf(err, "ensuring tables: %s", name)
		}
	}

	return nil
}

func (l logwrap) Printf(f string, v ...interface{}) { fmt.Fprintf(l.l, f, v...) }
