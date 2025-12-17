package fetcher

import (
	"context"
	"testing"

	"github.com/Luzifer/go_helpers/v2/fieldcollection"
)

func Test_GitTagFetcher(t *testing.T) {
	attrs := fieldcollection.FieldCollectionFromData(map[string]interface{}{
		"remote": "https://github.com/Luzifer/go-latestver.git",
	})

	f := Get("git_tag")

	if err := f.Validate(attrs); err != nil {
		t.Fatalf("validating attributes: %s", err)
	}

	ver, _, err := f.FetchVersion(context.Background(), attrs)
	if err != nil {
		t.Fatalf("fetching version: %s", err)
	}

	// Uses tag format: 1.0.0
	if len(ver) < 5 {
		t.Errorf("version has unexpected format: %s != X.X.X", ver)
	}

	t.Logf("found tag: %s", ver)
}

func Test_GitTagFetcherInvalid(t *testing.T) {
	attrs := fieldcollection.FieldCollectionFromData(map[string]interface{}{
		"remote": "https://example.com/example.git",
	})

	f := Get("git_tag")

	_, _, err := f.FetchVersion(context.Background(), attrs)
	if err == nil {
		t.Fatalf("fetching version dit not cause error")
	}
}
