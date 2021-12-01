package database

import (
	"testing"
	"time"
)

const sqlliteMemoryDSN = "file::memory:?cache=shared"

func Test_CatalogMetaStorage(t *testing.T) {
	ttp := func(v time.Time) *time.Time { return &v }

	dbc, err := NewClient("sqlite3", sqlliteMemoryDSN)
	if err != nil {
		t.Fatalf("unable to create database client: %s", err)
	}

	var (
		ce = CatalogEntry{Name: "testapp", Tag: "latest"}
		cm = CatalogMeta{CatalogName: ce.Name, CatalogTag: ce.Tag, CurrentVersion: "1.0.0", LastChecked: ttp(time.Now())}
	)

	// Empty fetch
	fetchedCM, err := dbc.Catalog.GetMeta(&ce)
	if err != nil {
		t.Fatalf("unable to retrieve catalog meta: %s", err)
	}

	for name, check := range map[string]bool{
		"name match":    fetchedCM.CatalogName == ce.Name,
		"tag match":     fetchedCM.CatalogTag == ce.Tag,
		"version empty": fetchedCM.CurrentVersion == "",
		"date nil":      fetchedCM.LastChecked == nil,
	} {
		if !check {
			t.Errorf("check failed: %s", name)
		}
	}

	// Initial set
	if err = dbc.Catalog.PutMeta(&cm); err != nil {
		t.Fatalf("unable to store catalog meta: %s", err)
	}

	fetchedCM, err = dbc.Catalog.GetMeta(&ce)
	if err != nil {
		t.Fatalf("unable to retrieve catalog meta: %s", err)
	}

	for name, check := range map[string]bool{
		"name match":    fetchedCM.CatalogName == ce.Name,
		"tag match":     fetchedCM.CatalogTag == ce.Tag,
		"version match": fetchedCM.CurrentVersion == cm.CurrentVersion,
		"date match":    fetchedCM.LastChecked.Equal(*cm.LastChecked),
	} {
		if !check {
			t.Errorf("check failed: %s", name)
		}
	}

	// Update
	cm.LastChecked = ttp(time.Now().Add(time.Hour)) // Compensate test running quite fast
	cm.CurrentVersion = "1.1.0"

	if err = dbc.Catalog.PutMeta(&cm); err != nil {
		t.Fatalf("unable to update catalog meta: %s", err)
	}

	fetchedCM, err = dbc.Catalog.GetMeta(&ce)
	if err != nil {
		t.Fatalf("unable to retrieve catalog meta: %s", err)
	}

	for name, check := range map[string]bool{
		"name match":    fetchedCM.CatalogName == ce.Name,
		"tag match":     fetchedCM.CatalogTag == ce.Tag,
		"version match": fetchedCM.CurrentVersion == cm.CurrentVersion,
		"date match":    fetchedCM.LastChecked.Equal(*cm.LastChecked),
	} {
		if !check {
			t.Errorf("check failed: %s", name)
		}
	}
}

func Test_LogStorage(t *testing.T) {
	dbc, err := NewClient("sqlite3", sqlliteMemoryDSN)
	if err != nil {
		t.Fatalf("unable to create database client: %s", err)
	}

	var (
		ce = CatalogEntry{Name: "testapp", Tag: "latest"}
		rt = time.Now()
	)

	for _, le := range []LogEntry{
		{CatalogName: ce.Name, CatalogTag: ce.Tag, Timestamp: rt.Add(-3 * time.Hour), VersionFrom: "1.0.0", VersionTo: "1.1.0"},
		{CatalogName: ce.Name, CatalogTag: ce.Tag, Timestamp: rt.Add(-1 * time.Hour), VersionFrom: "1.2.0", VersionTo: "1.3.0"},
		{CatalogName: ce.Name, CatalogTag: ce.Tag, Timestamp: rt.Add(-2 * time.Hour), VersionFrom: "1.1.0", VersionTo: "1.2.0"},
		{CatalogName: "anotherapp", CatalogTag: ce.Tag, Timestamp: rt.Add(-2 * time.Hour), VersionFrom: "5.2.0", VersionTo: "5.2.1"},
		{CatalogName: "anotherapp", CatalogTag: ce.Tag, Timestamp: rt.Add(-1 * time.Hour), VersionFrom: "5.2.1", VersionTo: "6.0.0"},
	} {
		//#nosec G601 // Acceptable for test usage
		if err = dbc.Logs.Add(&le); err != nil {
			t.Fatalf("unable to add log entry: %s", err)
		}
	}

	logs, err := dbc.Logs.ListForCatalogEntry(&ce, 100, 0)
	if err != nil {
		t.Fatalf("unable to fetch log entries for entry: %s", err)
	}

	if c := len(logs); c != 3 {
		t.Errorf("got unexpected number of logs for entry: %d != 3", c)
	}

	if !logs[2].Timestamp.Before(logs[1].Timestamp) || !logs[1].Timestamp.Before(logs[0].Timestamp) {
		t.Error("log entries are not sorted descending")
	}

	logs, err = dbc.Logs.List(100, 0)
	if err != nil {
		t.Fatalf("unable to fetch log entries: %s", err)
	}

	if c := len(logs); c != 5 {
		t.Errorf("got unexpected number of logs: %d != 5", c)
	}
}
