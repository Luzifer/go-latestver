package fetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/Luzifer/go-latestver/internal/database"
	"github.com/Luzifer/go-latestver/internal/helpers"
	"github.com/Luzifer/go_helpers/v2/fieldcollection"
)

/*
 * @module atlassian
 * @module_desc Fetches latest version of an Atlassian product
 */

var (
	atlassianDefaultEdition = ""
	atlassianDefaultSearch  = "TAR.GZ"
)

type (
	// AtlassianFetcher implements the fetcher interface to monitor Atlassian products
	AtlassianFetcher struct{}
)

func init() { registerFetcher("atlassian", func() Fetcher { return &AtlassianFetcher{} }) }

// FetchVersion retrieves the latest version for the catalog entry
func (AtlassianFetcher) FetchVersion(ctx context.Context, attrs *fieldcollection.FieldCollection) (string, time.Time, error) {
	url := fmt.Sprintf("https://my.atlassian.com/download/feeds/current/%s.json", attrs.MustString("product", nil))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", time.Time{}, errors.Wrap(err, "creating request")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", time.Time{}, errors.Wrap(err, "executing request")
	}
	defer func() { helpers.LogIfErr(resp.Body.Close(), "closing response body after read") }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", time.Time{}, errors.Wrap(err, "reading response body")
	}

	matches := jsonpStripRegex.FindSubmatch(body)
	if matches == nil {
		return "", time.Time{}, errors.New("document does not match jsonp syntax")
	}

	var payload []struct {
		Description  string `json:"description"`
		Edition      string `json:"edition"`
		Zipurl       string `json:"zipUrl"`
		Md5          string `json:"md5"`
		Size         string `json:"size"`
		Released     string `json:"released"`
		Type         string `json:"type"`
		Platform     string `json:"platform"`
		Version      string `json:"version"`
		Releasenotes string `json:"releaseNotes"`
		Upgradenotes string `json:"upgradeNotes"`
	}

	if err = json.Unmarshal(matches[1], &payload); err != nil {
		return "", time.Time{}, errors.Wrap(err, "parsing response JSON")
	}

	sort.Slice(payload, func(j, i int) bool { // j, i -> Reverse sort, biggest date at the top
		iRelease, _ := time.Parse("02-Jan-2006", payload[i].Released)
		jRelease, _ := time.Parse("02-Jan-2006", payload[j].Released)
		return iRelease.Before(jRelease)
	})

	var (
		// @attr edition optional string "" Filter down the versions according to its edition (e.g. "Enterprise" or "Standard" for Confluence)
		edition = attrs.MustString("edition", &atlassianDefaultEdition)
		// @attr search optional string "TAR.GZ" What to search in the download description: default is to search for the standalone .tar.gz file
		search = attrs.MustString("search", &atlassianDefaultSearch)
	)

	for _, r := range payload {
		if edition != "" && !strings.Contains(r.Edition, edition) {
			continue
		}

		if search != "" && !strings.Contains(r.Description, search) {
			continue
		}

		rt, _ := time.Parse("02-Jan-2006", r.Released)
		return r.Version, rt, nil
	}

	return "", time.Time{}, ErrNoVersionFound
}

// Links retrieves a collection of links for the fetcher
func (AtlassianFetcher) Links(_ *fieldcollection.FieldCollection) []database.CatalogLink {
	return nil
}

// Validate validates the configuration given to the fetcher
func (AtlassianFetcher) Validate(attrs *fieldcollection.FieldCollection) error {
	// @attr product required string "" Lowercase name of the product to fetch (e.g. confluence, crowd, jira-software, ...)
	if v, err := attrs.String("product"); err != nil || v == "" {
		return errors.New("product is expected to be non-empty string")
	}

	return nil
}
