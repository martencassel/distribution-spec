package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	assert "github.com/stretchr/testify/assert"
)

// [a-z0-9]+((\.|_|__|-+)[a-z0-9]+)*(\/[a-z0-9]+((\.|_|__|-+)[a-z0-9]+)*)*

func TestMain(t *testing.T) {
	var handleManifest http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v2/alpine/manifests/latest", r.URL.Path)
		w.WriteHeader(http.StatusOK)
	}
	req, err := http.NewRequest("GET", "/v2/alpine/manifests/latest", nil)
	assert.NoError(t, err)
	rr := httptest.NewRecorder()
	handleManifest.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

}
