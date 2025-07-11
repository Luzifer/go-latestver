package main

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/Luzifer/go-latestver/internal/database"
	"github.com/Luzifer/rconfig/v2"
)

var cfg = struct {
	FromStorage    string `flag:"from-storage" description:"Storage type to migrate from" validate:"nonzero"`   //revive:disable-line:struct-tag
	FromStorageDSN string `flag:"from-storage-dsn" description:"DSN for the 'from' storage" validate:"nonzero"` //revive:disable-line:struct-tag
	ToStorage      string `flag:"to-storage" description:"Storage type to migrate to" validate:"nonzero"`       //revive:disable-line:struct-tag
	ToStorageDSN   string `flag:"to-storage-dsn" description:"DSN for the 'to' storage" validate:"nonzero"`     //revive:disable-line:struct-tag
}{}

func initApp() error {
	rconfig.AutoEnv(true)
	if err := rconfig.ParseAndValidate(&cfg); err != nil {
		return errors.Wrap(err, "parsing commandline options")
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
