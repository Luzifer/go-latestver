//go:build dev

package main

import "net/http"

var frontendFS http.FileSystem = http.Dir("./")
