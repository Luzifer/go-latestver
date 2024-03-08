package fetcher

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Luzifer/go-latestver/internal/database"
	"github.com/Luzifer/go-latestver/internal/helpers"
	"github.com/Luzifer/go_helpers/v2/fieldcollection"
	"github.com/pkg/errors"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/repo"
	"sigs.k8s.io/yaml"
)

/*
 * @module helm
 * @module_desc Fetches the index file of a Helm Repo and yields the latest Helm-Chart version
 */

type (
	// HELMFetcher implements the fetcher interface to retrieve a version from a Helm Repo
	HELMFetcher struct{}
)

func init() { registerFetcher("helm", func() Fetcher { return &HELMFetcher{} }) }

// FetchVersion retrieves the latest version for the catalog entry
func (h HELMFetcher) FetchVersion(ctx context.Context, attrs *fieldcollection.FieldCollection) (string, time.Time, error) {
	vers, err := h.getChartVersionsFromRepo(ctx, attrs.MustString("repo", nil), attrs.MustString("chart", nil))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("fetching chart versions: %w", err)
	}

	if vers == nil {
		return "", time.Time{}, fmt.Errorf("chart not found in repo")
	}

	return vers[0].Version, vers[0].Created, nil
}

// Links retrieves a collection of links for the fetcher
func (HELMFetcher) Links(_ *fieldcollection.FieldCollection) []database.CatalogLink { return nil }

// Validate validates the configuration given to the fetcher
func (HELMFetcher) Validate(attrs *fieldcollection.FieldCollection) error {
	// @attr repo required string "" URL of the repo (i.e. "https://grafana.github.io/helm-charts")
	if v, err := attrs.String("repo"); err != nil || v == "" {
		return errors.New("repo is expected to be non-empty string")
	}

	// @attr chart required string "" Chart to fetch the version of (i.e. "grafana")
	if v, err := attrs.String("chart"); err != nil || v == "" {
		return errors.New("chart is expected to be non-empty string")
	}

	return nil
}

func (h HELMFetcher) getChartVersionsFromRepo(ctx context.Context, repoURL, chartName string) (repo.ChartVersions, error) {
	if !strings.HasSuffix(repoURL, "/index.yaml") {
		repoURL = strings.Join([]string{
			strings.TrimRight(repoURL, "/"),
			"index.yaml",
		}, "/")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, repoURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer func() { helpers.LogIfErr(resp.Body.Close(), "closing response body after read") }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected HTTP status %d", resp.StatusCode)
	}

	buf := new(bytes.Buffer)
	if _, err = io.Copy(buf, resp.Body); err != nil {
		return nil, fmt.Errorf("reading index content: %w", err)
	}

	index, err := h.loadIndex(buf.Bytes(), repoURL)
	if err != nil {
		return nil, fmt.Errorf("parsing index file: %w", err)
	}

	return index.Entries[chartName], nil
}

/*
 * Load functions taken from Helm v3 library as they are only defined
 * internally and not exposed
 *
 * https://github.com/helm/helm/blob/v3.14.2/pkg/repo/index.go#L341-L394
 */

// loadIndex loads an index file and does minimal validity checking.
//
// The source parameter is only used for logging.
// This will fail if API Version is not set (ErrNoAPIVersion) or if the unmarshal fails.
func (h HELMFetcher) loadIndex(data []byte, source string) (*repo.IndexFile, error) {
	i := &repo.IndexFile{}

	if len(data) == 0 {
		return i, repo.ErrEmptyIndexYaml
	}

	if err := h.jsonOrYamlUnmarshal(data, i); err != nil {
		return i, err
	}

	for name, cvs := range i.Entries {
		for idx := len(cvs) - 1; idx >= 0; idx-- {
			if cvs[idx] == nil {
				log.Printf("skipping loading invalid entry for chart %q from %s: empty entry", name, source)
				continue
			}
			// When metadata section missing, initialize with no data
			if cvs[idx].Metadata == nil {
				cvs[idx].Metadata = &chart.Metadata{}
			}
			if cvs[idx].APIVersion == "" {
				cvs[idx].APIVersion = chart.APIVersionV1
			}
			if err := cvs[idx].Validate(); err != nil {
				log.Printf("skipping loading invalid entry for chart %q %q from %s: %s", name, cvs[idx].Version, source, err)
				cvs = append(cvs[:idx], cvs[idx+1:]...)
			}
		}
	}
	i.SortEntries()
	if i.APIVersion == "" {
		return i, repo.ErrNoAPIVersion
	}
	return i, nil
}

// jsonOrYamlUnmarshal unmarshals the given byte slice containing JSON or YAML
// into the provided interface.
//
// It automatically detects whether the data is in JSON or YAML format by
// checking its validity as JSON. If the data is valid JSON, it will use the
// `encoding/json` package to unmarshal it. Otherwise, it will use the
// `sigs.k8s.io/yaml` package to unmarshal the YAML data.
func (HELMFetcher) jsonOrYamlUnmarshal(b []byte, i interface{}) error {
	if json.Valid(b) {
		return json.Unmarshal(b, i) //nolint:wrapcheck // Fine at this point
	}
	return yaml.UnmarshalStrict(b, i) //nolint:wrapcheck // Fine at this point
}
