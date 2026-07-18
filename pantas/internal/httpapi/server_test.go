package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServeSPADoesNotRedirectIndex(t *testing.T) {
	app := &App{index: []byte("<!doctype html><title>PANTAS</title>")}

	for _, path := range []string{"/", "/dashboard"} {
		t.Run(path, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, path, nil)
			response := httptest.NewRecorder()

			app.serveSPA(response, request)

			if response.Code != http.StatusOK {
				t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
			}
			if location := response.Header().Get("Location"); location != "" {
				t.Fatalf("unexpected redirect to %q", location)
			}
			if response.Body.Len() == 0 {
				t.Fatal("response body is empty")
			}
		})
	}
}

func TestServeSPAHeadHasNoBodyOrRedirect(t *testing.T) {
	app := &App{index: []byte("<!doctype html><title>PANTAS</title>")}
	request := httptest.NewRequest(http.MethodHead, "/", nil)
	response := httptest.NewRecorder()

	app.serveSPA(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
	}
	if location := response.Header().Get("Location"); location != "" {
		t.Fatalf("unexpected redirect to %q", location)
	}
	if response.Body.Len() != 0 {
		t.Fatalf("body length = %d, want 0", response.Body.Len())
	}
}
