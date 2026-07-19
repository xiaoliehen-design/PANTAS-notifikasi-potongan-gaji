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

func TestServeSPAStaticAssetsRequireRevalidation(t *testing.T) {
	app := &App{static: http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
		response.WriteHeader(http.StatusOK)
		_, _ = response.Write([]byte("asset"))
	})}
	request := httptest.NewRequest(http.MethodGet, "/app.js", nil)
	response := httptest.NewRecorder()

	app.serveSPA(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
	}
	if cacheControl := response.Header().Get("Cache-Control"); cacheControl != "no-cache" {
		t.Fatalf("Cache-Control = %q, want no-cache", cacheControl)
	}
}

func TestDeductionRuleInputValidation(t *testing.T) {
	for _, source := range []string{"late", "early_leave", "leave", "status", "shift"} {
		if !validRuleSource(source) {
			t.Errorf("valid source %q rejected", source)
		}
	}
	for _, source := range []string{"", "overtime", "late "} {
		if validRuleSource(source) {
			t.Errorf("invalid source %q accepted", source)
		}
	}
	for _, code := range []string{"TL4", "Cuti Khusus Dipotong", "PSW-5"} {
		if !validRuleCode(code) {
			t.Errorf("valid code %q rejected", code)
		}
	}
	for _, code := range []string{"", " TL4", "TL4\n"} {
		if validRuleCode(code) {
			t.Errorf("invalid code %q accepted", code)
		}
	}
}
