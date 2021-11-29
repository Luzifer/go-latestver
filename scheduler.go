package main

import (
	"context"
	"crypto/md5"
	"fmt"
	"math"
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

	nct := nextCheckTime(ce, cm.LastChecked)
	logger = logger.WithFields(log.Fields{
		"last": cm.LastChecked,
		"next": nct,
	})
	logger.Trace("Next check time found")
	if nct.After(time.Now()) {
		// Not yet ready to check
		return nil
	}

	logger.Debug("Checking for updates")

	ver, vertime, err := fetcher.Get(ce.Fetcher).FetchVersion(context.Background(), &ce.FetcherConfig)
	ver = strings.TrimPrefix(ver, "v")

	switch {

	case err != nil:
		log.WithField("entry", ce.Key()).WithError(err).Error("Fetcher caused error, error is stored in entry")
		cm.Error = err.Error()

	case cm.CurrentVersion != ver:

		logger.WithFields(log.Fields{
			"from": cm.CurrentVersion,
			"to":   ver,
		}).Info("Entry had version update")

		if err = storage.Logs.Add(&database.LogEntry{
			CatalogName: ce.Name,
			CatalogTag:  ce.Tag,
			Timestamp:   vertime,
			VersionTo:   ver,
			VersionFrom: cm.CurrentVersion,
		}); err != nil {
			return errors.Wrap(err, "adding log entry")
		}

		cm.VersionTime = ptrTime(vertime)
		cm.CurrentVersion = ver
		fallthrough

	default:
		cm.Error = ""

	}

	cm.LastChecked = ptrTime(time.Now())
	return errors.Wrap(storage.Catalog.PutMeta(cm), "updating meta entry")
}

func nextCheckTime(ce *database.CatalogEntry, lastCheck *time.Time) time.Time {
	hash := md5.New()
	fmt.Fprint(hash, ce.Key())

	var jitter int64
	for i, c := range hash.Sum(nil) {
		jitter += int64(c) * int64(math.Pow(10, float64(i)))
	}

	if lastCheck == nil {
		lastCheck = ptrTime(processStart)
	}

	next := lastCheck.
		Truncate(cfg.CheckDistribution).
		Add(time.Duration(jitter) % cfg.CheckDistribution)

	if next.Before(*lastCheck) {
		next = next.Add(cfg.CheckDistribution)
	}

	return next.Truncate(time.Second)
}

func ptrTime(t time.Time) *time.Time { return &t }
