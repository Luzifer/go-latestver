package fetcher

import (
	"context"
	"testing"

	"github.com/Luzifer/go_helpers/v2/fieldcollection"
)

func Test_JSONFetcher(t *testing.T) {
	attrs := fieldcollection.FieldCollectionFromData(map[string]interface{}{
		"jsonp": false,
		"url":   "https://my.atlassian.com/download/feeds/current/crowd.json",
		"xpath": "*[1]/version",
	})

	f := Get("json")

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

	t.Logf("found version: %s", ver)
}
