package ui

import (
	"embed"
	"io/fs"
	"net/http"
	"os"
	"path"
)

//go:embed dist/*
var UiFS embed.FS

// HttpHandler is a custom handler meant to serve the UI's single page web application
// that need to serve index.html instead of 404s.
func HttpHandler() http.Handler {
	distSubDirectory, err := fs.Sub(UiFS, "dist")
	if err != nil {
		// It is very unlikely this would ever happen, because we're explicitly embeding this subdirectory above
		// but if for some reason dist didn't exist like we expect, we should panic
		panic(err)
	}

	rootHttpFileSystem := http.FS(distSubDirectory)
	rootHttpFileServer := http.FileServer(rootHttpFileSystem)

	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		_, err := rootHttpFileSystem.Open(path.Clean(request.URL.Path))
		// force any not found 404's back to /, which is index.html
		if os.IsNotExist(err) {
			request.URL.Path = "/"
			writer.Header().Set("Content-Type", "text/html")
		}
		rootHttpFileServer.ServeHTTP(writer, request)
	})
}
