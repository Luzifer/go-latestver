package fetcher

import (
	"context"
	"io"
	"net/http"
	"regexp"
	"time"

	"github.com/pkg/errors"

	"github.com/Luzifer/go-latestver/internal/database"
	"github.com/Luzifer/go-latestver/internal/helpers"
	"github.com/Luzifer/go_helpers/v2/fieldcollection"
)

/*
 * @module regex
 * @module_desc Fetches URL and applies a regular expression to extract a version from it
 */

const (
	httpStatus3xx               = 300
	regexpFetcherExpectedLength = 2
)

type (
	// RegexFetcher implements the fetcher interface to monitor versions on a web request by regex
	RegexFetcher struct{}
)

func init() { registerFetcher("regex", func() Fetcher { return &RegexFetcher{} }) }

// FetchVersion retrieves the latest version for the catalog entry
func (RegexFetcher) FetchVersion(ctx context.Context, attrs *fieldcollection.FieldCollection) (string, time.Time, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, attrs.MustString("url", nil), nil)
	if err != nil {
		return "", time.Time{}, errors.Wrap(err, "creating request")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", time.Time{}, errors.Wrap(err, "executing request")
	}
	defer func() { helpers.LogIfErr(resp.Body.Close(), "closing response body after read") }()

	if resp.StatusCode >= httpStatus3xx {
		return "", time.Time{}, errors.Errorf("HTTP status %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", time.Time{}, errors.Wrap(err, "reading response body")
	}

	matches := regexp.MustCompile(attrs.MustString("regex", nil)).FindStringSubmatch(string(data))
	if matches == nil {
		return "", time.Time{}, errors.New("regex did not match the response")
	}

	if l := len(matches); l != regexpFetcherExpectedLength {
		return "", time.Time{}, errors.Errorf("unpexected number of matches: %d != 2", l)
	}

	return matches[1], time.Now(), nil
}

// Links retrieves a collection of links for the fetcher
func (RegexFetcher) Links(attrs *fieldcollection.FieldCollection) []database.CatalogLink {
	return []database.CatalogLink{
		{
			IconClass: "fas fa-globe",
			Name:      "Website",
			URL:       attrs.MustString("url", nil),
		},
	}
}

// Validate validates the configuration given to the fetcher
func (RegexFetcher) Validate(attrs *fieldcollection.FieldCollection) error {
	// @attr url required string "" URL to fetch the content from
	if v, err := attrs.String("url"); err != nil || v == "" {
		return errors.New("url is expected to be non-empty string")
	}

	// @attr regex required string "" Regular expression (RE2) to apply to the text fetched from the URL. The regex MUST have exactly one submatch containing the version.
	if v, err := attrs.String("regex"); err != nil || v == "" {
		return errors.New("regex is expected to be non-empty string")
	}

	r, err := regexp.Compile(attrs.MustString("regex", nil))
	if err != nil {
		return errors.Wrap(err, "compiling regex expression")
	}

	if n := r.NumSubexp(); n != 1 {
		return errors.Errorf("regex must have 1 submatch, has %d", n)
	}

	return nil
}
