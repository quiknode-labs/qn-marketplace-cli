package cmd

import (
	"fmt"
	"net/http"
	"strings"
)

type authMode int

const (
	authModeProvisioning authMode = iota
	authModeHeader
	authModeBasicAuth
	authModeNone
)

// detectAuthMode returns the active auth mode based on which flags were provided.
// Priority: provisioning > header > basic auth > none.
func detectAuthMode(provisionURL string, customHeaders []string, basicAuthChanged bool) authMode {
	switch {
	case provisionURL != "":
		return authModeProvisioning
	case len(customHeaders) > 0:
		return authModeHeader
	case basicAuthChanged:
		return authModeBasicAuth
	default:
		return authModeNone
	}
}

// parseHeaders parses a slice of "Key: Value" strings into a map.
// Malformed entries are skipped and returned as warning strings.
func parseHeaders(headers []string) (map[string]string, []string) {
	result := make(map[string]string)
	var warnings []string
	for _, h := range headers {
		parts := strings.SplitN(h, ": ", 2)
		if len(parts) != 2 {
			warnings = append(warnings, fmt.Sprintf("skipping malformed header %q: expected \"Key: Value\" format", h))
			continue
		}
		result[parts[0]] = parts[1]
	}
	return result, warnings
}

// applyAuthHeaders adds auth headers to req based on mode.
// Returns any warnings from malformed --header values.
func applyAuthHeaders(req *http.Request, mode authMode, customHeaders []string, basicAuth string) []string {
	switch mode {
	case authModeHeader:
		parsed, warnings := parseHeaders(customHeaders)
		for k, v := range parsed {
			req.Header.Set(k, v)
		}
		return warnings
	case authModeBasicAuth:
		req.Header.Set("Authorization", "Basic "+basicAuth)
	}
	return nil
}
