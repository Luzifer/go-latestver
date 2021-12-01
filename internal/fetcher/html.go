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

var htmlFetcherDefaultRegex = `(v?(?:[0-9]+\.?){2,})`

type (
	HTMLFetcher struct{}
)

func init() { registerFetcher("html", func() Fetcher { return &HTMLFetcher{} }) }

func (h HTMLFetcher) FetchVersion(ctx context.Context, attrs *fieldcollection.FieldCollection) (string, time.Time, error) {
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

func (h HTMLFetcher) Links(attrs *fieldcollection.FieldCollection) []database.CatalogLink {
	return []database.CatalogLink{
		{
			IconClass: "fas fa-globe",
			Name:      "Website",
			URL:       attrs.MustString("url", nil),
		},
	}
}

func (h HTMLFetcher) Validate(attrs *fieldcollection.FieldCollection) error {
	if v, err := attrs.String("url"); err != nil || v == "" {
		return errors.New("url is expected to be non-empty string")
	}

	if v, err := attrs.String("xpath"); err != nil || v == "" {
		return errors.New("xpath is expected to be non-empty string")
	}

	if _, err := xpath.Compile(attrs.MustString("xpath", nil)); err != nil {
		return errors.Wrap(err, "compiling xpath expression")
	}

	if attrs.CanString("regex") {
		if _, err := regexp.Compile(attrs.MustString("regex", nil)); err != nil {
			return errors.Wrap(err, "invalid regex given")
		}
	}

	return nil
}
