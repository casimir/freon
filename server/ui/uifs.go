package ui

import (
	"io/fs"
	"log"
	"net/http"
)

var FS http.FileSystem

func init() {
	if IsEmbedded {
		subFS, err := fs.Sub(uiFS, "statics")
		if err != nil {
			log.Fatal(err)
		}

		FS = http.FS(subFS)
	}
}
