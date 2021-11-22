package database

import "testing"

func Test_CreateInvalidDatabase(t *testing.T) {
	_, err := NewClient("invalid", "")
	if err == nil {
		t.Fatal("client creation for type 'invalid' did not cause error")
	}
}

func Test_CreateInaccessibleSqlite(t *testing.T) {
	_, err := NewClient("sqlite", "/this/path/should/really/not/exist.db")
	if err == nil {
		t.Fatal("client creation with unavailable sqlite path did not cause error")
	}
}

func Test_CreateInaccessibleMYSQL(t *testing.T) {
	_, err := NewClient("mysql", "user:pass@tcp(127.0.0.1:70000)/dbname?charset=utf8mb4&parseTime=True&loc=Local")
	if err == nil {
		t.Fatal("client creation with unavailable mysql did not cause error")
	}
}
