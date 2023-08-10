// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

var goodHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
})

func TestForceTLS(t *testing.T) {
	w := httptest.NewRecorder()
	url := "http://127.0.0.1/foo/bar"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	h := forceTLSHandler(goodHandler)
	h.ServeHTTP(w, req)
	if w.Code != 307 {
		t.Fatalf("bad: %d", w.Code)
	}
	if v := w.HeaderMap.Get("Location"); v != "https://127.0.0.1/foo/bar" {
		t.Fatalf("bad: %s", v)
	}
}

func TestForceTLS_valid(t *testing.T) {
	w := httptest.NewRecorder()
	url := "https://127.0.0.1/foo/bar"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	h := forceTLSHandler(goodHandler)
	h.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("bad: %d", w.Code)
	}
}

func TestForceTLS_post(t *testing.T) {
	w := httptest.NewRecorder()
	url := "http://127.0.0.1/foo/bar"
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	h := forceTLSHandler(goodHandler)
	h.ServeHTTP(w, req)
	if w.Code != 400 {
		t.Fatalf("bad: %d", w.Code)
	}
}
