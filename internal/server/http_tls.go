// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package server

import (
	"net/http"
	"strings"
)

// Proto header to read for the forwarded proto
const xForwardedProto = "X-Forwarded-Proto"

// forceTLSHandler forces TLS on a URL. This allows forwarded TLS connections
// with an X-Forwarded-Proto header.
func forceTLSHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		scheme := r.URL.Scheme
		if r.TLS != nil {
			// If we have an active TLS connection, consider this HTTPS
			scheme = "https"
		}
		if v := r.Header.Get(xForwardedProto); v != "" {
			scheme = strings.ToLower(v)
		}

		if scheme != "https" {
			if r.Method != "GET" {
				w.WriteHeader(400)
				return
			}

			url := *r.URL
			url.Scheme = "https"
			if url.Host == "" {
				url.Host = r.Host
			}

			http.Redirect(w, r, url.String(), http.StatusTemporaryRedirect)
			return
		}

		// Call the next handler in the chain.
		h.ServeHTTP(w, r)
	})
}
