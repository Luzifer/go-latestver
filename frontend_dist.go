//go:build !dev

package main

import (
	"embed"
	"net/http"
)

var (
	//go:embed frontend/**
	embeddedFrontend embed.FS

	frontendFS http.FileSystem = http.FS(embeddedFrontend)
)
