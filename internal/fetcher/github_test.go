package fetcher

import (
	"context"
	"testing"

	"github.com/Luzifer/go_helpers/fieldcollection"
)

func Test_GithubReleaseFetcher(t *testing.T) {
	attrs := fieldcollection.FieldCollectionFromData(map[string]interface{}{
		"repository": "Luzifer/korvike",
	})

	f := Get("github_release")

	if err := f.Validate(attrs); err != nil {
		t.Fatalf("validating attributes: %s", err)
	}

	ver, _, err := f.FetchVersion(context.Background(), attrs)
	if err != nil {
		t.Fatalf("fetching version: %s", err)
	}

	// Uses tag format: v1.0.0
	if len(ver) < 6 || ver[0] != 'v' {
		t.Errorf("version has unexpected format: %s != vX.X.X", ver)
	}

	t.Logf("found release: %s", ver)
}

func Test_GithubReleaseFetcherInvalid(t *testing.T) {
	attrs := fieldcollection.FieldCollectionFromData(map[string]interface{}{
		"repository": "Luzifer/thiswillneverexist",
	})

	f := Get("github_release")

	_, _, err := f.FetchVersion(context.Background(), attrs)
	if err == nil {
		t.Fatalf("fetching version dit not cause error")
	}
}
