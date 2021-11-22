package main

import (
	"context"
	"strings"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/Luzifer/go-latestver/internal/database"
	"github.com/Luzifer/go-latestver/internal/fetcher"
)

var schedulerRunActive bool

func schedulerRun() {
	if schedulerRunActive {
		return
	}
	schedulerRunActive = true
	defer func() { schedulerRunActive = false }()

	for _, ce := range configFile.Catalog {
		if err := checkForUpdates(&ce); err != nil {
			log.WithField("entry", ce.Key()).WithError(err).Error("Unable to update entry")
		}
	}
}

func checkForUpdates(ce *database.CatalogEntry) error {
	logger := log.WithField("entry", ce.Key())

	cm, err := storage.Catalog.GetMeta(ce)
	if err != nil {
		return errors.Wrap(err, "getting catalog meta")
	}

	checkTime := time.Now()
	if cm.LastChecked != nil {
		checkTime = *cm.LastChecked

		switch {
		case ce.CheckInterval > 0:
			checkTime = checkTime.Add(ce.CheckInterval)

		case configFile.CheckInterval > 0:
			checkTime = checkTime.Add(configFile.CheckInterval)

		default:
			checkTime = checkTime.Add(time.Hour)
		}
	}

	if checkTime.After(time.Now()) {
		// Not yet ready to check
		return nil
	}

	logger.Debug("Checking for updates")

	ver, vertime, err := fetcher.Get(ce.Fetcher).FetchVersion(context.Background(), &ce.FetcherConfig)
	if err != nil {
		return errors.Wrap(err, "fetching version")
	}

	ver = strings.TrimPrefix(ver, "v")

	if cm.CurrentVersion != ver {
		if err = storage.Logs.Add(&database.LogEntry{
			CatalogName: ce.Name,
			CatalogTag:  ce.Tag,
			Timestamp:   vertime,
			VersionTo:   ver,
			VersionFrom: cm.CurrentVersion,
		}); err != nil {
			return errors.Wrap(err, "adding log entry")
		}
		logger.WithFields(log.Fields{
			"from": cm.CurrentVersion,
			"to":   ver,
		}).Info("Entry had version update")
		cm.VersionTime = func(v time.Time) *time.Time { return &v }(vertime)
	}

	cm.CurrentVersion = ver
	cm.LastChecked = func(v time.Time) *time.Time { return &v }(time.Now())

	return errors.Wrap(storage.Catalog.PutMeta(cm), "updating meta entry")
}
