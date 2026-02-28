package validator

import (
	"context"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/admein-advertising/admein-vast-generator/vast"
)

const (
	// IABAnalysisCategory is used for spec-compliance checks encoded in the catalog.
	IABAnalysisCategory = "iab.analysis"
	// CustomAnalysisCategory is the default bucket for caller-supplied validators.
	CustomAnalysisCategory = "custom.analysis"
)

// NodeContext provides context to custom validators.
type NodeContext struct {
	Node    *genericNode
	Version vast.Version
}

// Text returns the trimmed character data contained within the node.
func (ctx NodeContext) Text() string {
	if ctx.Node == nil {
		return ""
	}
	return strings.TrimSpace(ctx.Node.Content)
}

// Attribute fetches the value of a node attribute by name.
func (ctx NodeContext) Attribute(name string) (string, bool) {
	if ctx.Node == nil {
		return "", false
	}
	return ctx.Node.attrValue(name)
}

// NodeValidatorFunc runs custom validation logic on a node.
type NodeValidatorFunc func(ctx NodeContext) *NodeAnalysisResult

// HTTPValidatorFunc represents a validator that performs HTTP requests (e.g., HEAD checks).
type HTTPValidatorFunc func(ctx context.Context, nodeCtx NodeContext, client *http.Client) (*NodeAnalysisResult, error)

var (
	customMu         sync.RWMutex
	customValidators = map[string][]NodeValidatorFunc{}
)

// RegisterCustomValidator registers a custom validator for the given node name.
// It can be called from init() functions in other packages.
func RegisterCustomValidator(nodeName string, validator NodeValidatorFunc) {
	if validator == nil {
		return
	}
	customMu.Lock()
	defer customMu.Unlock()
	key := strings.ToLower(nodeName)
	customValidators[key] = append(customValidators[key], validator)
}

func getCustomValidators(nodeName string) []NodeValidatorFunc {
	customMu.RLock()
	defer customMu.RUnlock()
	return append([]NodeValidatorFunc(nil), customValidators[strings.ToLower(nodeName)]...)
}

// HTTPValidatorRegistry stores HTTP-based validators keyed by node name.
var HTTPValidatorRegistry = struct {
	mu    sync.RWMutex
	store map[string][]HTTPValidatorFunc
}{store: map[string][]HTTPValidatorFunc{}}

// RegisterHTTPValidator registers an HTTP-based validator for the given node name.
func RegisterHTTPValidator(nodeName string, validator HTTPValidatorFunc) {
	if validator == nil {
		return
	}
	HTTPValidatorRegistry.mu.Lock()
	defer HTTPValidatorRegistry.mu.Unlock()
	key := strings.ToLower(nodeName)
	HTTPValidatorRegistry.store[key] = append(HTTPValidatorRegistry.store[key], validator)
}

func getHTTPValidators(nodeName string) []HTTPValidatorFunc {
	HTTPValidatorRegistry.mu.RLock()
	defer HTTPValidatorRegistry.mu.RUnlock()
	return append([]HTTPValidatorFunc(nil), HTTPValidatorRegistry.store[strings.ToLower(nodeName)]...)
}

// HTTPValidationOptions configure HTTP-based custom validator behavior.
type HTTPValidationOptions struct {
	Client  *http.Client
	Timeout time.Duration
}

func (opts *HTTPValidationOptions) client() *http.Client {
	if opts == nil || opts.Client == nil {
		return http.DefaultClient
	}
	return opts.Client
}
