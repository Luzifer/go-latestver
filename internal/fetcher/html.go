package fetcher

import (
	"context"
	"regexp"
	"time"

	"github.com/antchfx/htmlquery"
	"github.com/antchfx/xpath"
	"github.com/pkg/errors"
	"golang.org/x/net/html"

	"github.com/Luzifer/go-latestver/internal/database"
	"github.com/Luzifer/go_helpers/v2/fieldcollection"
)

/*
 * @module html
 * @module_desc Downloads website, selects text-node using XPath and optionally applies custom regular expression
 */

var htmlFetcherDefaultRegex = `(v?(?:[0-9]+\.?){2,})`

type (
	// HTMLFetcher implements the fetcher interface to monitor versions on websites by xpath queries
	HTMLFetcher struct{}
)

func init() { registerFetcher("html", func() Fetcher { return &HTMLFetcher{} }) }

// FetchVersion retrieves the latest version for the catalog entry
func (HTMLFetcher) FetchVersion(_ context.Context, attrs *fieldcollection.FieldCollection) (string, time.Time, error) {
	doc, err := htmlquery.LoadURL(attrs.MustString("url", nil))
	if err != nil {
		return "", time.Time{}, errors.Wrap(err, "loading URL")
	}

	node, err := htmlquery.Query(doc, attrs.MustString("xpath", nil))
	if err != nil {
		return "", time.Time{}, errors.Wrap(err, "querying xpath")
	}

	if node == nil {
		return "", time.Time{}, errors.New("xpath expression lead to nil-node")
	}

	if node.Type == html.ElementNode && node.FirstChild != nil && node.FirstChild.Type == html.TextNode {
		node = node.FirstChild
	}

	if node.Type != html.TextNode {
		return "", time.Time{}, errors.Errorf("xpath expression lead to unexpected node type: %d", node.Type)
	}

	match := regexp.MustCompile(attrs.MustString("regex", &htmlFetcherDefaultRegex)).FindStringSubmatch(node.Data)
	if len(match) < 2 { //nolint:gomnd // Simple count of fields, no need for constant
		return "", time.Time{}, errors.New("regular expression did not yield version")
	}

	return match[1], time.Now(), nil
}

// Links retrieves a collection of links for the fetcher
func (HTMLFetcher) Links(attrs *fieldcollection.FieldCollection) []database.CatalogLink {
	return []database.CatalogLink{
		{
			IconClass: "fas fa-globe",
			Name:      "Website",
			URL:       attrs.MustString("url", nil),
		},
	}
}

// Validate validates the configuration given to the fetcher
func (HTMLFetcher) Validate(attrs *fieldcollection.FieldCollection) error {
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

	// @attr regex optional string "(v?(?:[0-9]+\.?){2,})" Regular expression to apply to the text from the XPath expression
	if attrs.CanString("regex") {
		if _, err := regexp.Compile(attrs.MustString("regex", nil)); err != nil {
			return errors.Wrap(err, "invalid regex given")
		}
	}

	return nil
}
