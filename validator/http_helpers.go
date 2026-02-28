package validator

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const probeRangeHeader = "bytes=0-0"

// probeMediaURL attempts to verify that the provided media URL responds to an
// HTTP HEAD request. When a server disallows HEAD it falls back to an HTTP GET
// with a byte range request to minimize transfer size.
func probeMediaURL(ctx context.Context, client *http.Client, rawURL string) (*http.Response, error) {
	normalized, err := normalizeProbeURL(rawURL)
	if err != nil {
		return nil, err
	}

	resp, err := doHTTPRequest(ctx, client, http.MethodHead, normalized, nil)
	if err == nil {
		if resp.StatusCode != http.StatusMethodNotAllowed {
			return resp, nil
		}
		resp.Body.Close()
	} else {
		return nil, err
	}

	// Fall back to a ranged GET request when HEAD is not supported.
	headers := map[string]string{"Range": probeRangeHeader}
	return doHTTPRequest(ctx, client, http.MethodGet, normalized, headers)
}

func normalizeProbeURL(raw string) (string, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "", errors.New("media URL is empty")
	}
	if strings.HasPrefix(trimmed, "//") {
		trimmed = "https:" + trimmed
	}
	parsed, err := url.Parse(trimmed)
	if err != nil {
		return "", fmt.Errorf("invalid media URL %q: %w", raw, err)
	}
	if parsed.Scheme == "" || parsed.Host == "" {
		return "", fmt.Errorf("invalid media URL %q: missing scheme or host", raw)
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return "", fmt.Errorf("unsupported media URL scheme %q", parsed.Scheme)
	}
	return parsed.String(), nil
}

func doHTTPRequest(ctx context.Context, client *http.Client, method, target string, headers map[string]string) (*http.Response, error) {
	if client == nil {
		client = http.DefaultClient
	}
	req, err := http.NewRequestWithContext(ctx, method, target, nil)
	if err != nil {
		return nil, err
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
