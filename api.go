package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/feeds"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/Luzifer/go-latestver/internal/badge"
	"github.com/Luzifer/go-latestver/internal/config"
	"github.com/Luzifer/go-latestver/internal/database"
	"github.com/Luzifer/go-latestver/internal/fetcher"
)

type (
	apiCatalogEntry struct {
		database.CatalogEntry
		database.CatalogMeta
	}
)

func buildFullURL(u *url.URL, _ error) string {
	return strings.Join([]string{
		strings.TrimRight(cfg.BaseURL, "/"),
		strings.TrimLeft(u.String(), "/"),
	}, "/")
}

func catalogEntryToAPICatalogEntry(ce database.CatalogEntry) (apiCatalogEntry, error) {
	cm, err := storage.Catalog.GetMeta(&ce)
	if err != nil {
		return apiCatalogEntry{}, errors.Wrap(err, "fetching catalog meta")
	}

	for _, l := range fetcher.Get(ce.Fetcher).Links(ce.FetcherConfig) {
		var found bool
		for _, el := range ce.Links {
			if l.Name == el.Name {
				found = true
				break
			}
		}

		if !found {
			ce.Links = append(ce.Links, l)
		}
	}

	return apiCatalogEntry{CatalogEntry: ce, CatalogMeta: *cm}, nil
}

func handleBadge(w http.ResponseWriter, r *http.Request) {
	var (
		compare   = r.FormValue("compare")
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
		http.Error(w, "Unable to fetch catalog data", http.StatusInternalServerError)
		return
	}

	color := "green"
	if compare != "" && compare != cm.CurrentVersion {
		color = "red"
	}

	svg := badge.Create(ce.Key(), cm.CurrentVersion, color)
	w.Header().Add("Content-Type", "image/svg+xml")
	if _, err = w.Write(svg); err != nil {
		logrus.WithError(err).Error("writing SVG response")
	}
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

	ae, err := catalogEntryToAPICatalogEntry(ce)
	if err != nil {
		logrus.WithError(err).Error("Unable to fetch catalog data")
		http.Error(w, "Unable to fetch catalog data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(ae); err != nil {
		logrus.WithError(err).Error("Unable to encode catalog entry")
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
		logrus.WithError(err).Error("Unable to fetch catalog meta")
		http.Error(w, "Unable to fetch catalog meta", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, cm.CurrentVersion)
}

func handleCatalogList(w http.ResponseWriter, _ *http.Request) {
	out := make([]apiCatalogEntry, len(configFile.Catalog))

	for i := range configFile.Catalog {
		ce := configFile.Catalog[i]

		ae, err := catalogEntryToAPICatalogEntry(ce)
		if err != nil {
			logrus.WithError(err).Error("Unable to fetch catalog data")
			http.Error(w, "Unable to fetch catalog data", http.StatusInternalServerError)
			return
		}

		out[i] = ae
	}

	sort.Slice(out, func(i, j int) bool { return out[i].Key() < out[j].Key() })

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(out); err != nil {
		logrus.WithError(err).Error("Unable to encode catalog entry list")
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
		logrus.WithError(err).Error("Unable to encode logs")
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

	feed := &feeds.Feed{
		Description: "Generated by go-latestver: https://github.com/Luzifer/go-latestver",
		Link:        &feeds.Link{Href: buildFullURL(router.Get("catalog").URL())},
		Title:       "Latestver Update Log",
	}

	if vars["name"] != "" {
		feed.Title = fmt.Sprintf("Latestver Update Log of %s:%s", vars["name"], vars["tag"])
		feed.Link.Href = buildFullURL(router.Get("catalog-entry").URL("name", vars["name"], "tag", vars["tag"]))
	}

	for _, le := range logs {
		feed.Add(&feeds.Item{
			Created:     le.Timestamp.UTC(),
			Description: fmt.Sprintf("%s:%s updated to version %s from %s", le.CatalogName, le.CatalogTag, le.VersionTo, le.VersionFrom),
			Id:          fmt.Sprintf("%s:%s-%s", le.CatalogName, le.CatalogTag, le.Timestamp.UTC().Format(time.RFC3339)),
			Link:        &feeds.Link{Href: buildFullURL(router.Get("catalog-entry").URL("name", le.CatalogName, "tag", le.CatalogTag))},
			Title:       fmt.Sprintf("%s:%s %s", le.CatalogName, le.CatalogTag, le.VersionTo),
		})
	}

	w.Header().Set("Content-Type", "application/rss+xml; charset=utf-8")
	if err = feed.WriteRss(w); err != nil {
		logrus.WithError(err).Error("Unable to render RSS")
		http.Error(w, "Unable to render RSS", http.StatusInternalServerError)
		return
	}
}

func prepareLogForRequest(r *http.Request) ([]database.LogEntry, error) {
	var (
		vars      = mux.Vars(r)
		name, tag = vars["name"], vars["tag"]

		num, page = 25, 0

		ce   database.CatalogEntry
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
		ce, err = configFile.CatalogEntryByTag(name, tag)
		if errors.Is(err, config.ErrCatalogEntryNotFound) {
			return nil, config.ErrCatalogEntryNotFound
		}

		logs, err = storage.Logs.ListForCatalogEntry(&ce, num, page)
	}

	return logs, errors.Wrap(err, "listing log entries")
}
