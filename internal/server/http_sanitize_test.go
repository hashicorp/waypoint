package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

//var goodHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//	w.WriteHeader(200)
//})

func TestSanitize_api(t *testing.T) {
	w := httptest.NewRecorder()
	url := "http://127.0.0.1/foo/bar"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	h := sanitizedHandler(goodHandler)
	h.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("bad: %d", w.Code)
	}
}

func TestSanitize_ui(t *testing.T) {
	w := httptest.NewRecorder()
	url := "http://127.0.0.1/foo/bar"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	h := sanitizedHandler(goodHandler)
	h.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("bad: %d", w.Code)
	}
}
