// Package config holds the definition of the configuration file and
// some methods to load and validate it
package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/Luzifer/go-latestver/internal/database"
	"github.com/Luzifer/go-latestver/internal/fetcher"
	"github.com/Luzifer/go-latestver/internal/helpers"
)

type (
	// File represents the configuration file content
	File struct {
		Catalog       []database.CatalogEntry `yaml:"catalog"`
		CheckInterval time.Duration           `yaml:"check_interval"`
	}
)

// ErrCatalogEntryNotFound signalizes a catalog entry with the given
// name was not found
var ErrCatalogEntryNotFound = errors.New("catalog entry not found")

// New creates a new empty File object with defaults
func New() *File {
	return &File{
		CheckInterval: time.Hour,
	}
}

// CatalogEntryByTag retrieves a catalog entry by its name or returns
// and ErrCatalogEntryNotFound if it does not exist
func (f File) CatalogEntryByTag(name, tag string) (database.CatalogEntry, error) {
	for i := range f.Catalog {
		ce := f.Catalog[i]
		if ce.Name == name && ce.Tag == tag {
			return ce, nil
		}
	}

	return database.CatalogEntry{}, ErrCatalogEntryNotFound
}

// Load loads the contents of the configuration file in the filesystem
// into the File object
func (f *File) Load(filepath string) error {
	fh, err := os.Open(filepath) //#nosec:G304 // As this is the config file and needs a location this is fine
	if err != nil {
		return fmt.Errorf("opening config file: %w", err)
	}
	defer func() { helpers.LogIfErr(fh.Close(), "closing config after load") }()

	dec := yaml.NewDecoder(fh)
	dec.KnownFields(true)

	if err = dec.Decode(f); err != nil {
		return fmt.Errorf("decoding config: %w", err)
	}

	return nil
}

// ValidateCatalog checks whether invalid fetchers are used or the
// configuration of the fetcher is not suitable for the given fetcher
func (f File) ValidateCatalog() error {
	for i, ce := range f.Catalog {
		fi := fetcher.Get(ce.Fetcher)
		if fi == nil {
			return fmt.Errorf("catalog entry %d has unknown fetcher", i)
		}

		if err := fi.Validate(ce.FetcherConfig); err != nil {
			return fmt.Errorf("catalog entry %d has invalid fetcher config: %w", i, err)
		}
	}

	return nil
}
