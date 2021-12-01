package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"

	"github.com/Luzifer/go-latestver/internal/config"
	"github.com/Luzifer/go-latestver/internal/database"
	httpHelper "github.com/Luzifer/go_helpers/v2/http"
	"github.com/Luzifer/rconfig/v2"
)

var (
	cfg = struct {
		BadgeGenInstance  string        `flag:"badge-gen-instance" default:"https://badges.fyi/" description:"Where to find the badge-gen instance to use badges from"`
		BaseURL           string        `flag:"base-url" default:"https://example.com/" description:"Base-URL the application is reachable at"`
		Config            string        `flag:"config,c" default:"config.yaml" description:"Configuration file with catalog entries"`
		Listen            string        `flag:"listen" default:":3000" description:"Port/IP to listen on"`
		LogLevel          string        `flag:"log-level" default:"info" description:"Log level (debug, info, warn, error, fatal)"`
		CheckDistribution time.Duration `flag:"check-distribution" default:"1h" description:"Checks are executed at static times every [value]"`
		Storage           string        `flag:"storage" default:"sqlite" description:"Storage adapter to use (mysql, sqlite)"`
		StorageDSN        string        `flag:"storage-dsn" default:"file::memory:?cache=shared" description:"DSN to connect to the database"`
		VersionAndExit    bool          `flag:"version" default:"false" description:"Prints current version and exits"`
	}{}

	configFile = config.New()
	router     *mux.Router
	storage    *database.Client

	version = "dev"
)

func initApp() {
	rconfig.AutoEnv(true)
	if err := rconfig.ParseAndValidate(&cfg); err != nil {
		log.Fatalf("Unable to parse commandline options: %s", err)
	}

	if cfg.VersionAndExit {
		fmt.Printf("go-latestver %s\n", version)
		os.Exit(0)
	}

	if l, err := log.ParseLevel(cfg.LogLevel); err != nil {
		log.WithError(err).Fatal("Unable to parse log level")
	} else {
		log.SetLevel(l)
	}
}

func main() {
	initApp()

	var err error

	if err = configFile.Load(cfg.Config); err != nil {
		log.WithError(err).Fatal("Unable to load configuration")
	}

	if err = configFile.ValidateCatalog(); err != nil {
		log.WithError(err).Fatal("Configuration is not valid")
	}

	storage, err = database.NewClient(cfg.Storage, cfg.StorageDSN)
	if err != nil {
		log.WithError(err).Fatal("Unable to connect to database")
	}

	scheduler := cron.New()
	scheduler.AddFunc(fmt.Sprintf("@every %s", schedulerInterval), schedulerRun)
	scheduler.Start()

	router = mux.NewRouter()
	router.HandleFunc("/v1/catalog", handleCatalogList).Methods(http.MethodGet)
	router.HandleFunc("/v1/catalog/{name}/{tag}", handleCatalogGet).Methods(http.MethodGet)
	router.HandleFunc("/v1/catalog/{name}/{tag}/log", handleLog).Methods(http.MethodGet)
	router.HandleFunc("/v1/catalog/{name}/{tag}/version", handleCatalogGetVersion).Methods(http.MethodGet)
	router.HandleFunc("/v1/log", handleLog).Methods(http.MethodGet)

	router.HandleFunc("/{name}/{tag}.svg", handleBadgeRedirect).Methods(http.MethodGet).Name("catalog-entry-badge")
	router.HandleFunc("/{name}/{tag}/log.rss", handleLogFeed).Methods(http.MethodGet).Name("catalog-entry-rss")
	router.HandleFunc("/log.rss", handleLogFeed).Methods(http.MethodGet).Name("log-rss")

	router.HandleFunc("/", handleSinglePage).Methods(http.MethodGet).Name("catalog")
	router.HandleFunc("/{name}/{tag}", handleSinglePage).Methods(http.MethodGet).Name("catalog-entry")
	router.PathPrefix("/").HandlerFunc(handleSinglePage)

	var handler http.Handler = router
	handler = httpHelper.GzipHandler(handler)
	handler = httpHelper.NewHTTPLogHandler(handler)

	if err := http.ListenAndServe(cfg.Listen, handler); err != nil {
		log.WithError(err).Fatal("HTTP server exited unclean")
	}
}
