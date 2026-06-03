package cmd

import (
	"net/http"
	"testing"
)

func TestParseHeaders(t *testing.T) {
	tests := []struct {
		name         string
		input        []string
		wantHeaders  map[string]string
		wantWarnings int
	}{
		{
			name:        "valid single header",
			input:       []string{"X-API-Key: secret"},
			wantHeaders: map[string]string{"X-API-Key": "secret"},
		},
		{
			name:        "valid multiple headers",
			input:       []string{"X-API-Key: secret", "X-Tenant: foo"},
			wantHeaders: map[string]string{"X-API-Key": "secret", "X-Tenant": "foo"},
		},
		{
			name:        "value containing colon is preserved",
			input:       []string{"Authorization: Bearer tok:en"},
			wantHeaders: map[string]string{"Authorization": "Bearer tok:en"},
		},
		{
			name:         "malformed header is skipped with warning",
			input:        []string{"NoColonHere"},
			wantHeaders:  map[string]string{},
			wantWarnings: 1,
		},
		{
			name:         "mixed valid and malformed",
			input:        []string{"X-API-Key: secret", "BadHeader"},
			wantHeaders:  map[string]string{"X-API-Key": "secret"},
			wantWarnings: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, warns := parseHeaders(tt.input)
			if len(warns) != tt.wantWarnings {
				t.Errorf("parseHeaders() warnings = %d, want %d", len(warns), tt.wantWarnings)
			}
			for k, v := range tt.wantHeaders {
				if got[k] != v {
					t.Errorf("parseHeaders()[%q] = %q, want %q", k, got[k], v)
				}
			}
		})
	}
}

func TestDetectAuthMode(t *testing.T) {
	tests := []struct {
		name             string
		provisionURL     string
		customHeaders    []string
		basicAuthChanged bool
		want             authMode
	}{
		{"provisioning when url set", "http://localhost/provision", nil, false, authModeProvisioning},
		{"header auth when headers set", "", []string{"X-API-Key: secret"}, false, authModeHeader},
		{"basic auth when flag changed", "", nil, true, authModeBasicAuth},
		{"no auth as default", "", nil, false, authModeNone},
		{"provisioning takes priority over headers", "http://localhost/provision", []string{"X-Key: v"}, true, authModeProvisioning},
		{"header takes priority over basic auth", "", []string{"X-Key: v"}, true, authModeHeader},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := detectAuthMode(tt.provisionURL, tt.customHeaders, tt.basicAuthChanged)
			if got != tt.want {
				t.Errorf("detectAuthMode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApplyAuthHeaders(t *testing.T) {
	t.Run("header mode sets custom headers", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "http://example.com", nil)
		warns := applyAuthHeaders(req, authModeHeader, []string{"X-API-Key: secret"}, "")
		if req.Header.Get("X-API-Key") != "secret" {
			t.Errorf("expected X-API-Key: secret, got %q", req.Header.Get("X-API-Key"))
		}
		if len(warns) != 0 {
			t.Errorf("expected no warnings, got %v", warns)
		}
	})

	t.Run("basic auth mode sets Authorization header", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "http://example.com", nil)
		applyAuthHeaders(req, authModeBasicAuth, nil, "dXNlcjpwYXNz")
		if req.Header.Get("Authorization") != "Basic dXNlcjpwYXNz" {
			t.Errorf("expected Basic auth header, got %q", req.Header.Get("Authorization"))
		}
	})

	t.Run("no auth mode sets no extra headers", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "http://example.com", nil)
		applyAuthHeaders(req, authModeNone, nil, "")
		if req.Header.Get("Authorization") != "" {
			t.Errorf("expected no Authorization header, got %q", req.Header.Get("Authorization"))
		}
	})

	t.Run("provisioning mode sets no extra auth headers", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "http://example.com", nil)
		applyAuthHeaders(req, authModeProvisioning, nil, "somevalue")
		if req.Header.Get("Authorization") != "" {
			t.Errorf("expected no Authorization header in provisioning mode, got %q", req.Header.Get("Authorization"))
		}
	})

	t.Run("malformed header in header mode produces warning", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "http://example.com", nil)
		warns := applyAuthHeaders(req, authModeHeader, []string{"BadHeader"}, "")
		if len(warns) != 1 {
			t.Errorf("expected 1 warning, got %d", len(warns))
		}
	})
}
