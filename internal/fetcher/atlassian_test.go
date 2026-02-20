package fetcher

import (
	"context"
	"testing"

	"github.com/Luzifer/go_helpers/fieldcollection"
)

func Test_AtlassianFetcher(t *testing.T) {
	attrs := fieldcollection.FieldCollectionFromData(map[string]interface{}{
		"product": "confluence",
		"edition": "Standard",
	})

	f := Get("atlassian")

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

	attrs.Set("edition", "ThisDoesNotExist")

	_, _, err = f.FetchVersion(context.Background(), attrs)
	if err == nil {
		t.Errorf("fetching non existing edition did not error")
	}
}
