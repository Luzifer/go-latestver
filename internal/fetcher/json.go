package fetcher

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"

	"github.com/antchfx/jsonquery"
	"github.com/antchfx/xpath"
	"github.com/pkg/errors"

	"github.com/Luzifer/go-latestver/internal/database"
	"github.com/Luzifer/go_helpers/v2/fieldcollection"
)

var (
	jsonFetcherDefaultRegex = `(v?(?:[0-9]+\.?){2,})`
	jsonpStripRegex         = regexp.MustCompile(`(?m)^[^\(]+\((.*)\)$`)
	ptrBoolFalse            = func(v bool) *bool { return &v }(false)
)

type (
	JSONFetcher struct{}
)

func init() { registerFetcher("json", func() Fetcher { return &JSONFetcher{} }) }

func (JSONFetcher) FetchVersion(ctx context.Context, attrs *fieldcollection.FieldCollection) (string, time.Time, error) {
	var (
		doc *jsonquery.Node
		err error
	)

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
		defer resp.Body.Close()

		body, err = ioutil.ReadAll(resp.Body)
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
	if len(match) < 2 { //nolint:gomnd // Simple count of fields, no need for constant
		return "", time.Time{}, errors.New("regular expression did not yield version")
	}

	return match[1], time.Now(), nil
}

func (JSONFetcher) Links(attrs *fieldcollection.FieldCollection) []database.CatalogLink { return nil }

func (JSONFetcher) Validate(attrs *fieldcollection.FieldCollection) error {
	if v, err := attrs.String("url"); err != nil || v == "" {
		return errors.New("url is expected to be non-empty string")
	}

	if v, err := attrs.String("xpath"); err != nil || v == "" {
		return errors.New("xpath is expected to be non-empty string")
	}

	if _, err := xpath.Compile(attrs.MustString("xpath", nil)); err != nil {
		return errors.Wrap(err, "compiling xpath expression")
	}

	return nil
}
