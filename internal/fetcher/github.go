package fetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/pkg/errors"

	"github.com/Luzifer/go-latestver/internal/database"
	"github.com/Luzifer/go_helpers/v2/fieldcollection"
)

const githubHTTPTimeout = 2 * time.Second

type (
	GithubReleaseFetcher struct{}

	githubRelease struct {
		TagName     string    `json:"tag_name"`
		PublishedAt time.Time `json:"published_at"`
		Prerelease  bool      `json:"prerelease"`
	}
)

func init() { registerFetcher("github_release", func() Fetcher { return &GithubReleaseFetcher{} }) }

func (g GithubReleaseFetcher) FetchVersion(ctx context.Context, attrs *fieldcollection.FieldCollection) (string, time.Time, error) {
	ctx, cancel := context.WithTimeout(ctx, githubHTTPTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("https://api.github.com/repos/%s/releases", attrs.MustString("repository", nil)),
		nil,
	)
	if err != nil {
		return "", time.Time{}, errors.Wrap(err, "creating http request")
	}
	req.Header.Set("User-Agent", "Luzifer/go-latestver GithubReleaseFetcher")

	if os.Getenv("GITHUB_CLIENT_ID") != "" && os.Getenv("GITHUB_CLIENT_SECRET") != "" {
		req.SetBasicAuth(os.Getenv("GITHUB_CLIENT_ID"), os.Getenv("GITHUB_CLIENT_SECRET"))
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", time.Time{}, errors.Wrap(err, "executing request")
	}

	if resp.StatusCode != http.StatusOK {
		return "", time.Time{}, errors.Errorf("unexpected HTTP status %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	var payload []githubRelease
	if err = json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return "", time.Time{}, errors.Wrap(err, "decoding response")
	}

	var release *githubRelease
	for i := range payload {
		if payload[i].Prerelease {
			continue
		}

		if release == nil || release.PublishedAt.Before(payload[i].PublishedAt) {
			release = &payload[i]
		}
	}

	if release == nil {
		return "", time.Time{}, ErrNoVersionFound
	}

	return release.TagName, release.PublishedAt, nil
}

func (g GithubReleaseFetcher) Links(attrs *fieldcollection.FieldCollection) []database.CatalogLink {
	return []database.CatalogLink{
		{
			IconClass: "fab fa-github",
			Name:      "Repository",
			URL:       fmt.Sprintf("https://github.com/%s", attrs.MustString("repository", nil)),
		},
	}
}

func (g GithubReleaseFetcher) Validate(attrs *fieldcollection.FieldCollection) error {
	if v, err := attrs.String("repository"); err != nil || v == "" {
		return errors.New("repository is expected to be non-empty string")
	}

	return nil
}
