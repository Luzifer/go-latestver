package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/feeds"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/Luzifer/go-latestver/internal/config"
	"github.com/Luzifer/go-latestver/internal/database"
	"github.com/Luzifer/go-latestver/internal/fetcher"
)

type (
	APICatalogEntry struct {
		database.CatalogEntry
		database.CatalogMeta
	}
)

func buildFullURL(u *url.URL) string {
	return strings.Join([]string{
		strings.TrimRight(cfg.BaseURL, "/"),
		strings.TrimLeft(u.String(), "/"),
	}, "/")
}

func handleCatalogGet(w http.ResponseWriter, r *http.Request) {
	var (
		vars      = mux.Vars(r)
		name, tag = vars["name"], vars["tag"]
	)

	ce, err := configFile.CatalogEntryByTag(name, tag)
	if errors.Is(err, config.ErrCatalogEntryNotFound) {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	cm, err := storage.Catalog.GetMeta(&ce)
	if err != nil {
		log.WithError(err).Error("Unable to fetch catalog meta")
		http.Error(w, "Unable to fetch catalog meta", http.StatusInternalServerError)
		return
	}

	ce.Links = append(
		ce.Links,
		fetcher.Get(ce.Fetcher).Links(&ce.FetcherConfig)...,
	)

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(APICatalogEntry{CatalogEntry: ce, CatalogMeta: *cm}); err != nil {
		log.WithError(err).Error("Unable to encode catalog entry")
		http.Error(w, "Unable to encode catalog meta", http.StatusInternalServerError)
		return
	}
}

func handleCatalogGetVersion(w http.ResponseWriter, r *http.Request) {
	var (
		vars      = mux.Vars(r)
		name, tag = vars["name"], vars["tag"]
	)

	ce, err := configFile.CatalogEntryByTag(name, tag)
	if errors.Is(err, config.ErrCatalogEntryNotFound) {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	cm, err := storage.Catalog.GetMeta(&ce)
	if err != nil {
		log.WithError(err).Error("Unable to fetch catalog meta")
		http.Error(w, "Unable to fetch catalog meta", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, cm.CurrentVersion)
}

func handleCatalogList(w http.ResponseWriter, r *http.Request) {
	out := make([]APICatalogEntry, len(configFile.Catalog))

	for i := range configFile.Catalog {
		ce := configFile.Catalog[i]

		cm, err := storage.Catalog.GetMeta(&ce)
		if err != nil {
			log.WithError(err).Error("Unable to fetch catalog meta")
			http.Error(w, "Unable to fetch catalog meta", http.StatusInternalServerError)
			return
		}

		ce.Links = append(
			ce.Links,
			fetcher.Get(ce.Fetcher).Links(&ce.FetcherConfig)...,
		)

		out[i] = APICatalogEntry{CatalogEntry: ce, CatalogMeta: *cm}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(out); err != nil {
		log.WithError(err).Error("Unable to encode catalog entry list")
		http.Error(w, "Unable to encode catalog meta", http.StatusInternalServerError)
		return
	}
}

func handleLog(w http.ResponseWriter, r *http.Request) {
	logs, err := prepareLogForRequest(r)
	switch err {
	case nil:
		// This is fine

	case config.ErrCatalogEntryNotFound:
		http.Error(w, "Not found", http.StatusNotFound)

	default:
		http.Error(w, "Unable to fetch logs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(logs); err != nil {
		log.WithError(err).Error("Unable to encode logs")
		http.Error(w, "Unable to encode logs", http.StatusInternalServerError)
		return
	}
}

func handleLogFeed(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	logs, err := prepareLogForRequest(r)
	switch err {
	case nil:
		// This is fine

	case config.ErrCatalogEntryNotFound:
		http.Error(w, "Not found", http.StatusNotFound)

	default:
		http.Error(w, "Unable to fetch logs", http.StatusInternalServerError)
		return
	}

	var (
		feedTitle  = "Latestver Update Log"
		feedURL, _ = router.Get("catalog").URL()
	)
	if vars["name"] != "" {
		feedTitle = fmt.Sprintf("Latestver Update Log of %s:%s", vars["name"], vars["tag"])
		feedURL, _ = router.Get("catalog-entry").URL("name", vars["name"], "tag", vars["tag"])
	}

	feed := &feeds.Feed{
		Description: "Generated by go-latestver: https://github.com/Luzifer/go-latestver",
		Link:        &feeds.Link{Href: buildFullURL(feedURL)},
		Title:       feedTitle,
	}

	for _, le := range logs {
		catalogEntryURL, _ := router.Get("catalog-entry").URL("name", le.CatalogName, "tag", le.CatalogTag)
		feed.Add(&feeds.Item{
			Created:     le.Timestamp.UTC(),
			Description: fmt.Sprintf("%s:%s updated to version %s from %s", le.CatalogName, le.CatalogTag, le.VersionTo, le.VersionFrom),
			Id:          fmt.Sprintf("%s:%s-%s", le.CatalogName, le.CatalogTag, le.Timestamp.UTC().Format(time.RFC3339)),
			Link:        &feeds.Link{Href: buildFullURL(catalogEntryURL)},
			Title:       fmt.Sprintf("%s:%s %s", le.CatalogName, le.CatalogTag, le.VersionTo),
		})
	}

	w.Header().Set("Content-Type", "application/rss+xml; charset=utf-8")
	if err = feed.WriteRss(w); err != nil {
		log.WithError(err).Error("Unable to render RSS")
		http.Error(w, "Unable to render RSS", http.StatusInternalServerError)
		return
	}
}

func prepareLogForRequest(r *http.Request) ([]database.LogEntry, error) {
	var (
		vars      = mux.Vars(r)
		name, tag = vars["name"], vars["tag"]

		num, page = 25, 0

		err  error
		logs []database.LogEntry
	)

	if v, err := strconv.Atoi(r.FormValue("num")); err == nil && v > 0 && v < 100 {
		num = v
	}

	if v, err := strconv.Atoi(r.FormValue("page")); err == nil && v >= 0 {
		page = v
	}

	if name == "" && tag == "" {
		logs, err = storage.Logs.List(num, page)
	} else {
		ce, err := configFile.CatalogEntryByTag(name, tag)
		if errors.Is(err, config.ErrCatalogEntryNotFound) {
			return nil, config.ErrCatalogEntryNotFound
		}

		logs, err = storage.Logs.ListForCatalogEntry(&ce, num, page)
	}

	return logs, errors.Wrap(err, "listing log entries")
}
