package config

import (
	"os"
	"time"

	"github.com/Luzifer/go-latestver/internal/database"
	"github.com/Luzifer/go-latestver/internal/fetcher"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

var ErrCatalogEntryNotFound = errors.New("catalog entry not found")

type (
	File struct {
		Catalog       []database.CatalogEntry `yaml:"catalog"`
		CheckInterval time.Duration           `yaml:"check_interval"`
	}
)

func New() *File {
	return &File{
		CheckInterval: time.Hour,
	}
}

func (f File) CatalogEntryByTag(name, tag string) (database.CatalogEntry, error) {
	for i := range f.Catalog {
		ce := f.Catalog[i]
		if ce.Name == name && ce.Tag == tag {
			return ce, nil
		}
	}

	return database.CatalogEntry{}, ErrCatalogEntryNotFound
}

func (f *File) Load(filepath string) error {
	fh, err := os.Open(filepath)
	if err != nil {
		return errors.Wrap(err, "opening config file")
	}
	defer fh.Close()

	dec := yaml.NewDecoder(fh)
	dec.KnownFields(true)

	return errors.Wrap(dec.Decode(f), "decoding config")
}

func (f File) ValidateCatalog() error {
	for i, ce := range f.Catalog {
		f := fetcher.Get(ce.Fetcher)
		if f == nil {
			return errors.Errorf("catalog entry %d has unknown fetcher", i)
		}

		if err := f.Validate(&ce.FetcherConfig); err != nil {
			return errors.Wrapf(err, "catalog entry %d has invalid fetcher config", i)
		}
	}

	return nil
}
