package fetcher

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"regexp"
	"time"

	"github.com/antchfx/jsonquery"
	"github.com/antchfx/xpath"
	"github.com/pkg/errors"

	"github.com/Luzifer/go-latestver/internal/database"
	"github.com/Luzifer/go-latestver/internal/helpers"
	"github.com/Luzifer/go_helpers/fieldcollection"
)

/*
 * @module json
 * @module_desc Fetches a JSON / JSONP file from remote source and traverses it using XPath expression
 */

var (
	jsonFetcherDefaultRegex = `(v?(?:[0-9]+\.?){2,})`
	jsonpStripRegex         = regexp.MustCompile(`(?m)^[^\(]+\((.*)\)$`)
	ptrBoolFalse            = func(v bool) *bool { return &v }(false)
)

type (
	// JSONFetcher implements the fetcher interface to retrieve a version from a JSON document
	JSONFetcher struct{}
)

func init() { registerFetcher("json", func() Fetcher { return &JSONFetcher{} }) }

// FetchVersion retrieves the latest version for the catalog entry
func (JSONFetcher) FetchVersion(ctx context.Context, attrs *fieldcollection.FieldCollection) (string, time.Time, error) {
	var (
		doc *jsonquery.Node
		err error
	)

	// @attr jsonp optional boolean "false" File contains JSONP function, strip it to get the raw JSON
	if attrs.MustBool("jsonp", ptrBoolFalse) {
		var (
			body []byte
			req  *http.Request
			resp *http.Response
		)
		req, err = http.NewRequestWithContext(ctx, http.MethodGet, attrs.MustString("url", nil), nil)
		if err != nil {
			return "", time.Time{}, errors.Wrap(err, "creating request")
		}

		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			return "", time.Time{}, errors.Wrap(err, "executing request")
		}
		defer func() { helpers.LogIfErr(resp.Body.Close(), "closing response body after read") }()

		body, err = io.ReadAll(resp.Body)
		if err != nil {
			return "", time.Time{}, errors.Wrap(err, "reading response body")
		}

		matches := jsonpStripRegex.FindSubmatch(body)
		if matches == nil {
			return "", time.Time{}, errors.New("document does not match jsonp syntax")
		}

		doc, err = jsonquery.Parse(bytes.NewReader(matches[1]))
	} else {
		doc, err = jsonquery.LoadURL(attrs.MustString("url", nil))
	}

	if err != nil {
		return "", time.Time{}, errors.New("parsing JSON document")
	}

	node, err := jsonquery.Query(doc, attrs.MustString("xpath", nil))
	if err != nil {
		return "", time.Time{}, errors.Wrap(err, "querying xpath")
	}

	if node == nil {
		return "", time.Time{}, errors.New("xpath expression lead to nil-node")
	}

	if node.Type == jsonquery.ElementNode && node.FirstChild != nil && node.FirstChild.Type == jsonquery.TextNode {
		node = node.FirstChild
	}

	if node.Type != jsonquery.TextNode {
		return "", time.Time{}, errors.Errorf("xpath expression lead to unexpected node type: %d", node.Type)
	}

	match := regexp.MustCompile(attrs.MustString("regex", &jsonFetcherDefaultRegex)).FindStringSubmatch(node.Data)
	if len(match) < 2 { //nolint:mnd // Simple count of fields, no need for constant
		return "", time.Time{}, errors.New("regular expression did not yield version")
	}

	return match[1], time.Now(), nil
}

// Links retrieves a collection of links for the fetcher
func (JSONFetcher) Links(_ *fieldcollection.FieldCollection) []database.CatalogLink { return nil }

// Validate validates the configuration given to the fetcher
func (JSONFetcher) Validate(attrs *fieldcollection.FieldCollection) error {
	// @attr url required string "" URL to fetch the HTML from
	if v, err := attrs.String("url"); err != nil || v == "" {
		return errors.New("url is expected to be non-empty string")
	}

	// @attr xpath required string "" XPath expression leading to the text-node containing the version
	if v, err := attrs.String("xpath"); err != nil || v == "" {
		return errors.New("xpath is expected to be non-empty string")
	}

	if _, err := xpath.Compile(attrs.MustString("xpath", nil)); err != nil {
		return errors.Wrap(err, "compiling xpath expression")
	}

	return nil
}
