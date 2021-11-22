package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Luzifer/go-latestver/internal/config"
	"github.com/Luzifer/go-latestver/internal/database"
	"github.com/Luzifer/go-latestver/internal/fetcher"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

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
	ce.CatalogMeta = *cm

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(ce); err != nil {
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
	out := make([]database.CatalogEntry, len(configFile.Catalog))

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
		ce.CatalogMeta = *cm

		out[i] = ce
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(out); err != nil {
		log.WithError(err).Error("Unable to encode catalog entry list")
		http.Error(w, "Unable to encode catalog meta", http.StatusInternalServerError)
		return
	}
}

func handleLog(w http.ResponseWriter, r *http.Request) {
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
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		logs, err = storage.Logs.ListForCatalogEntry(&ce, num, page)
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(logs); err != nil {
		log.WithError(err).Error("Unable to encode logs")
		http.Error(w, "Unable to encode logs", http.StatusInternalServerError)
		return
	}
}

func handleLogFeed(w http.ResponseWriter, r *http.Request) {} // FIXME
