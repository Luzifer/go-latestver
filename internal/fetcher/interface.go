// Package fetcher contains the implementations of fetchers to retrieve
// current versions for the catalog entries
package fetcher

import (
	"context"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/Luzifer/go-latestver/internal/database"
	"github.com/Luzifer/go_helpers/fieldcollection"
)

type (
	// Fetcher defines the interface all fetchers have to follow
	Fetcher interface {
		// FetchVersion retrieves the latest version for the catalog entry
		FetchVersion(ctx context.Context, attrs *fieldcollection.FieldCollection) (string, time.Time, error)
		// Links retrieves a collection of links for the fetcher
		Links(attrs *fieldcollection.FieldCollection) []database.CatalogLink
		// Validate validates the configuration given to the fetcher
		Validate(attrs *fieldcollection.FieldCollection) error
	}
	// Create represents a function to instantiate a Fetcher
	Create func() Fetcher
)

var (
	// ErrNoVersionFound signalizes the fetcher was not able to retrieve a version
	ErrNoVersionFound = errors.New("no version found")

	availableFetchers     = map[string]Create{}
	availableFetchersLock sync.RWMutex
)

func registerFetcher(name string, fn Create) {
	availableFetchersLock.Lock()
	defer availableFetchersLock.Unlock()

	availableFetchers[name] = fn
}

// Get retrieves an creation function for the given fetcher name
func Get(name string) Fetcher {
	availableFetchersLock.RLock()
	defer availableFetchersLock.RUnlock()

	fn, ok := availableFetchers[name]
	if !ok {
		return nil
	}

	return fn()
}
