// database migration utility
package main

import (
	"fmt"

	"github.com/Luzifer/rconfig/v2"
	"github.com/sirupsen/logrus"

	"github.com/Luzifer/go-latestver/internal/database"
)

var cfg = struct {
	FromStorage    string `flag:"from-storage" description:"Storage type to migrate from" validate:"nonzero"`   //revive:disable-line:struct-tag // nonzero is valid for our validate
	FromStorageDSN string `flag:"from-storage-dsn" description:"DSN for the 'from' storage" validate:"nonzero"` //revive:disable-line:struct-tag // nonzero is valid for our validate
	ToStorage      string `flag:"to-storage" description:"Storage type to migrate to" validate:"nonzero"`       //revive:disable-line:struct-tag // nonzero is valid for our validate
	ToStorageDSN   string `flag:"to-storage-dsn" description:"DSN for the 'to' storage" validate:"nonzero"`     //revive:disable-line:struct-tag // nonzero is valid for our validate
}{}

func initApp() error {
	rconfig.AutoEnv(true)
	if err := rconfig.ParseAndValidate(&cfg); err != nil {
		return fmt.Errorf("parsing commandline options: %w", err)
	}

	return nil
}

func main() {
	var err error
	if err = initApp(); err != nil {
		logrus.WithError(err).Fatal("initializing app")
	}

	src, err := database.NewClient(cfg.FromStorage, cfg.FromStorageDSN)
	if err != nil {
		logrus.WithError(err).Fatal("opening from database")
	}

	dest, err := database.NewClient(cfg.ToStorage, cfg.ToStorageDSN)
	if err != nil {
		logrus.WithError(err).Fatal("opening to database")
	}

	if err := src.Migrate(dest); err != nil {
		logrus.WithError(err).Fatal("execute migration")
	}

	logrus.Info("your database has been migrated")
}
