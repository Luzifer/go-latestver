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
		BaseURL        string        `flag:"base-url" default:"https://example.com/" description:"Base-URL the application is reachable at"`
		Config         string        `flag:"config,c" default:"config.yaml" description:"Configuration file with catalog entries"`
		Listen         string        `flag:"listen" default:":3000" description:"Port/IP to listen on"`
		LogLevel       string        `flag:"log-level" default:"info" description:"Log level (debug, info, warn, error, fatal)"`
		MaxJitter      time.Duration `flag:"max-jitter" default:"30m" description:"Maximum jitter to add to the check interval for load balancing"`
		Storage        string        `flag:"storage" default:"sqlite" description:"Storage adapter to use (mysql, sqlite)"`
		StorageDSN     string        `flag:"storage-dsn" default:"file::memory:?cache=shared" description:"DSN to connect to the database"`
		VersionAndExit bool          `flag:"version" default:"false" description:"Prints current version and exits"`
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
	scheduler.AddFunc("@every 1m", schedulerRun)
	scheduler.Start()

	router = mux.NewRouter()
	router.HandleFunc("/v1/catalog", handleCatalogList).Methods(http.MethodGet)
	router.HandleFunc("/v1/catalog/{name}/{tag}", handleCatalogGet).Methods(http.MethodGet)
	router.HandleFunc("/v1/catalog/{name}/{tag}/log", handleLog).Methods(http.MethodGet)
	router.HandleFunc("/v1/catalog/{name}/{tag}/version", handleCatalogGetVersion).Methods(http.MethodGet)
	router.HandleFunc("/v1/log", handleLog).Methods(http.MethodGet)

	router.HandleFunc("/", nil).Methods(http.MethodGet).Name("catalog")
	router.HandleFunc("/{name}/{tag}", nil).Methods(http.MethodGet).Name("catalog-entry")
	router.HandleFunc("/{name}/{tag}/log.rss", handleLogFeed).Methods(http.MethodGet).Name("catalog-entry-rss")
	router.HandleFunc("/log", nil).Methods(http.MethodGet)
	router.HandleFunc("/log.rss", handleLogFeed).Methods(http.MethodGet).Name("log-rss")

	var handler http.Handler = router
	handler = httpHelper.GzipHandler(handler)
	handler = httpHelper.NewHTTPLogHandler(handler)

	if err := http.ListenAndServe(cfg.Listen, handler); err != nil {
		log.WithError(err).Fatal("HTTP server exited unclean")
	}
}
