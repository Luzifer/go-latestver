package main

import (
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/Luzifer/go-latestver/internal/helpers"
)

func handleSinglePage(w http.ResponseWriter, r *http.Request) {
	// get the absolute path to prevent directory traversal
	urlPath, err := filepath.Abs(r.URL.Path)
	if err != nil {
		// if we failed to get the absolute path respond with a 400 bad request
		// and stop
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, frontendPath := range []string{
		filepath.Join("frontend", path.Base(urlPath)),
		filepath.Join("frontend", "index.html"),
	} {
		f, err := frontendFS.Open(frontendPath)
		switch {
		case err == nil:
			// file is opened, serve it
			defer func() { helpers.LogIfErr(f.Close(), "closing frontend file after serve") }() //revive:disable-line:defer Fine here as it will only open one file

			stat, err := f.Stat()
			if err != nil {
				http.Error(w, errors.Wrap(err, "stating opened file").Error(), http.StatusInternalServerError)
				return
			}

			if stat.IsDir() {
				continue
			}

			http.ServeContent(w, r, stat.Name(), stat.ModTime(), f)
			return

		case os.IsNotExist(err):
			// file does not exist, try next
			continue

		default:
			// if we got an error (that wasn't that the file doesn't exist) stating the
			// file, return a 500 internal server error and stop
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// No path could be opened
	http.NotFound(w, r)
}
