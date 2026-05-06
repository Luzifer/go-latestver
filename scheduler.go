package main

import (
	"context"
	"crypto/md5" //#nosec:G501 // Used to derive a static jitter checksum, not cryptographically
	"fmt"
	"math"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/Luzifer/go-latestver/internal/database"
	"github.com/Luzifer/go-latestver/internal/fetcher"
)

const schedulerInterval = time.Minute

var schedulerRunActive bool

func schedulerRun() {
	if schedulerRunActive {
		return
	}
	schedulerRunActive = true
	defer func() { schedulerRunActive = false }()

	for i := range configFile.Catalog {
		ce := configFile.Catalog[i]
		if err := checkForUpdates(&ce); err != nil {
			log.WithField("entry", ce.Key()).WithError(err).Error("Unable to update entry")
		}
	}
}

func checkForUpdates(ce *database.CatalogEntry) error {
	logger := log.WithField("entry", ce.Key())

	cm, err := storage.Catalog.GetMeta(ce)
	if err != nil {
		return fmt.Errorf("getting catalog meta: %w", err)
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

	ver, vertime, err := fetcher.Get(ce.Fetcher).FetchVersion(context.Background(), ce.FetcherConfig)
	ver = strings.TrimPrefix(ver, "v")
	vertime = vertime.Truncate(time.Second).UTC()

	logger = logger.WithFields(log.Fields{
		"from": cm.CurrentVersion,
		"to":   ver,
	})

	var (
		compareErr   error
		shouldUpdate = true
	)
	if ce.VersionConstraint != nil {
		shouldUpdate, compareErr = ce.VersionConstraint.ShouldApply(cm.CurrentVersion, ver)
	}

	switch {
	case err != nil:
		logger.WithError(err).Error("Fetcher caused error, error is stored in entry")
		cm.Error = err.Error()

	case compareErr != nil:
		logger.WithError(compareErr).Error("Version compare caused error, error is stored in entry")
		cm.Error = compareErr.Error()

	case cm.CurrentVersion != ver && !shouldUpdate:
		logger.Info("Version-updated prevented by constraints")
		cm.Error = ""

	case cm.CurrentVersion != ver && shouldUpdate:
		logger.Info("Entry had version update")

		if err = storage.Logs.Add(&database.LogEntry{
			CatalogName: ce.Name,
			CatalogTag:  ce.Tag,
			Timestamp:   time.Now().Truncate(time.Second).UTC(),
			VersionTo:   ver,
			VersionFrom: cm.CurrentVersion,
		}); err != nil {
			return fmt.Errorf("adding log entry: %w", err)
		}

		cm.VersionTime = new(vertime)
		cm.CurrentVersion = ver
		cm.Error = ""

	case cm.CurrentVersion == ver:
		logger.Debug("Version did not change")
		cm.Error = ""

	default:
		cm.Error = ""
	}

	cm.LastChecked = new(time.Now().Truncate(time.Second).UTC())
	if err = storage.Catalog.PutMeta(cm); err != nil {
		return fmt.Errorf("updating meta entry: %w", err)
	}

	return nil
}

func nextCheckTime(ce *database.CatalogEntry, lastCheck *time.Time) time.Time {
	if lastCheck == nil {
		// Has never been checked, check ASAP
		return time.Now()
	}

	var jitter int64
	//#nosec:G401 // Used to derive a static jitter checksum, not cryptographically
	for i, c := range md5.Sum([]byte(ce.Key())) {
		jitter += int64(c) * int64(math.Pow(10, float64(i))) //revive:disable-line:add-constant // No need for constant here
	}

	next := lastCheck.
		Truncate(cfg.CheckDistribution).
		Add(time.Duration(jitter) % cfg.CheckDistribution)

	if next.Before(lastCheck.Add(schedulerInterval)) {
		next = next.Add(cfg.CheckDistribution)
	}

	return next.Truncate(time.Second)
}
