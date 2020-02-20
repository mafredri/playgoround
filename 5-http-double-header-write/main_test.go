package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(201)
		w.WriteHeader(202)
	}))

	req, err := http.NewRequest("GET", srv.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(resp.StatusCode)
}
