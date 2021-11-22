package fetcher

import (
	"context"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/Luzifer/go-latestver/internal/database"
	"github.com/Luzifer/go_helpers/v2/fieldcollection"
)

type (
	Fetcher interface {
		FetchVersion(ctx context.Context, attrs *fieldcollection.FieldCollection) (string, time.Time, error)
		Links(attrs *fieldcollection.FieldCollection) []database.CatalogLink
		Validate(attrs *fieldcollection.FieldCollection) error
	}
	FetcherCreate func() Fetcher
)

var (
	ErrNoVersionFound = errors.New("no version found")

	availableFetchers     = map[string]FetcherCreate{}
	availableFetchersLock sync.RWMutex
)

func registerFetcher(name string, fn FetcherCreate) {
	availableFetchersLock.Lock()
	defer availableFetchersLock.Unlock()

	availableFetchers[name] = fn
}

func Get(name string) Fetcher {
	availableFetchersLock.RLock()
	defer availableFetchersLock.RUnlock()

	fn, ok := availableFetchers[name]
	if !ok {
		return nil
	}

	return fn()
}
