package fetcher

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/Luzifer/go_helpers/fieldcollection"
	"github.com/antchfx/htmlquery"
	"github.com/antchfx/xpath"
	"golang.org/x/net/html"

	"github.com/Luzifer/go-latestver/internal/database"
)

/*
 * @module html
 * @module_desc Downloads website, selects text-node using XPath and optionally applies custom regular expression
 */

type (
	// HTMLFetcher implements the fetcher interface to monitor versions on websites by xpath queries
	HTMLFetcher struct{}
)

var htmlFetcherDefaultRegex = `(v?(?:[0-9]+\.?){2,})`

func init() { registerFetcher("html", func() Fetcher { return &HTMLFetcher{} }) }

// FetchVersion retrieves the latest version for the catalog entry
func (HTMLFetcher) FetchVersion(_ context.Context, attrs *fieldcollection.FieldCollection) (string, time.Time, error) {
	doc, err := htmlquery.LoadURL(attrs.MustString("url", nil))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("loading URL: %w", err)
	}

	node, err := htmlquery.Query(doc, attrs.MustString("xpath", nil))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("querying xpath: %w", err)
	}

	if node == nil {
		return "", time.Time{}, errors.New("xpath expression lead to nil-node")
	}

	if node.Type == html.ElementNode && node.FirstChild != nil && node.FirstChild.Type == html.TextNode {
		node = node.FirstChild
	}

	if node.Type != html.TextNode {
		return "", time.Time{}, fmt.Errorf("xpath expression lead to unexpected node type: %d", node.Type)
	}

	match := regexp.MustCompile(attrs.MustString("regex", &htmlFetcherDefaultRegex)).FindStringSubmatch(node.Data)
	if len(match) < 2 { //nolint:mnd // Simple count of fields, no need for constant
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
		return fmt.Errorf("compiling xpath expression: %w", err)
	}

	// @attr regex optional string "(v?(?:[0-9]+\.?){2,})" Regular expression to apply to the text from the XPath expression
	if attrs.CanString("regex") {
		if _, err := regexp.Compile(attrs.MustString("regex", nil)); err != nil {
			return fmt.Errorf("invalid regex given: %w", err)
		}
	}

	return nil
}
