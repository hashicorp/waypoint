package server

import (
	"bytes"
	"net/http"
	"strings"
)

// sanitizedResponseWriter sanitizes HTTP responses to prevent XSS attacks.

type sanitizedResponseWriter struct {
	buf bytes.Buffer
	http.ResponseWriter
}

func (w *sanitizedResponseWriter) Write(data []byte) (int, error) {
	var newData []byte
	for _, b := range data {
		switch b {
		case byte('&'):
			newData = append(newData, "&amp;"...)
		case []byte("'")[0]:
			newData = append(newData, "&#39;"...)
		case byte('<'):
			newData = append(newData, "&lt;"...)
		case byte('>'):
			newData = append(newData, "&gt;"...)
		case byte('"'):
			newData = append(newData, "&#34;"...)
		default:
			newData = append(newData, b)
		}
	}
	return w.buf.Write(newData)
}

func sanitizedHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "/v1") {
			h.ServeHTTP(w, r)
		}
		rw := &sanitizedResponseWriter{}
		h.ServeHTTP(rw, r)
	})
}
