package validator

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

func init() {
	registerBuiltInHTTPValidators()
}

func registerBuiltInHTTPValidators() {
	RegisterHTTPValidator("MediaFile", mediaFileHTTPValidator)
}

func mediaFileHTTPValidator(ctx context.Context, nodeCtx NodeContext, client *http.Client) (*NodeAnalysisResult, error) {
	url := nodeCtx.Text()
	if url == "" {
		return &NodeAnalysisResult{Category: CustomAnalysisCategory, Status: StatusFail, Reasons: []string{"media file URL is empty"}}, nil
	}

	resp, err := probeMediaURL(ctx, client, url)
	if err != nil {
		return &NodeAnalysisResult{Category: CustomAnalysisCategory, Status: StatusFail, Reasons: []string{fmt.Sprintf("media file request failed: %v", err)}}, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return &NodeAnalysisResult{Category: CustomAnalysisCategory, Status: StatusFail, Reasons: []string{fmt.Sprintf("media file responded with HTTP %d", resp.StatusCode)}}, nil
	}

	if expected, ok := nodeCtx.Attribute("type"); ok {
		expected = strings.ToLower(strings.TrimSpace(expected))
		if expected != "" {
			actual := strings.ToLower(strings.TrimSpace(resp.Header.Get("Content-Type")))
			if idx := strings.Index(actual, ";"); idx >= 0 {
				actual = strings.TrimSpace(actual[:idx])
			}
			if actual != "" && actual != expected {
				return &NodeAnalysisResult{Category: CustomAnalysisCategory, Status: StatusFail, Reasons: []string{fmt.Sprintf("content type mismatch: expected %s, got %s", expected, actual)}}, nil
			}
		}
	}

	return &NodeAnalysisResult{Category: CustomAnalysisCategory, Status: StatusPass}, nil
}
