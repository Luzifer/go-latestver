package main

import (
	"github.com/pkg/errors"

	"github.com/Luzifer/go-latestver/internal/database"
)

func main() {
	src, err := database.NewClient("sqlite", "")
	if err != nil {
		panic(errors.Wrap(err, "opening src database"))
	}
	dest, err := database.NewClient("mysql", "")
	if err != nil {
		panic(errors.Wrap(err, "opening dest database"))
	}

	if err := src.Migrate(dest); err != nil {
		panic(errors.Wrap(err, "migrating to dest database"))
	}
}
