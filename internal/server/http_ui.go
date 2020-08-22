package server

import (
	"net/http"

	assetfs "github.com/elazarl/go-bindata-assetfs"
)

// httpUIHandler returns an http.Handler that registers the UI path
func httpUIHandler(handler http.Handler) http.Handler {
	mux := http.NewServeMux()

	uifs := http.FileServer(&assetfs.AssetFS{
		Asset:     Asset,
		AssetDir:  AssetDir,
		AssetInfo: AssetInfo,
		Prefix:    "ui/dist",
		Fallback:  "index.html",
	})

	mux.Handle("/", uifs)

	return mux
}
