package validator

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/admein-advertising/admein-vast-generator/vast"
)

var (
	// ErrInvalidRoot is returned when the XML root element is not <VAST>.
	ErrInvalidRoot = errors.New("Root element must be VAST")
	// ErrMissingVersion indicates the version attribute is missing on <VAST>.
	ErrMissingVersion = errors.New("Missing VAST version attribute")
	// ErrUnsupportedVersion indicates the provided version is not in catalog.
	ErrUnsupportedVersion = errors.New("Unsupported VAST version")
)

// Option configures the validation behavior.
type Option func(*config)

type config struct {
	catalog     *Catalog
	runCustom   bool
	runHTTP     bool
	httpOptions HTTPValidationOptions
}

func defaultConfig() *config {
	return &config{
		catalog:     defaultCatalog,
		runCustom:   true,
		runHTTP:     true,
		httpOptions: HTTPValidationOptions{Timeout: 2 * time.Second},
	}
}

// WithCatalog allows callers to substitute the catalog used for IAB analysis.
func WithCatalog(catalog *Catalog) Option {
	return func(cfg *config) {
		if catalog != nil {
			cfg.catalog = catalog
		}
	}
}

// DisableCustomValidators disables execution of registered custom validators as
// well as HTTP validators.
func DisableCustomValidators() Option {
	return func(cfg *config) {
		cfg.runCustom = false
		cfg.runHTTP = false
	}
}

// DisableHTTPValidators disables execution of registered HTTP validators while
// keeping non-networked custom validators enabled.
func DisableHTTPValidators() Option {
	return func(cfg *config) {
		cfg.runHTTP = false
	}
}

// WithHTTPValidationOptions configures the HTTP client/timeout used by HTTP validators.
func WithHTTPValidationOptions(opts HTTPValidationOptions) Option {
	return func(cfg *config) {
		cfg.httpOptions = opts
	}
}

// Validate parses and validates a VAST XML document.
func Validate(raw []byte, opts ...Option) (*ValidationResult, error) {
	if len(raw) == 0 {
		return nil, errEmptyXML
	}

	cfg := defaultConfig()
	for _, opt := range opts {
		opt(cfg)
	}

	root, err := buildNodeTree(raw)
	if err != nil {
		return nil, err
	}

	if root.localName() != "VAST" {
		return nil, ErrInvalidRoot
	}

	versionValue, ok := root.attrValue("version")
	if !ok || strings.TrimSpace(versionValue) == "" {
		return nil, ErrMissingVersion
	}
	version := vast.Version(strings.TrimSpace(versionValue))

	rootSpec, hasRootSpec := cfg.catalog.node("VAST")
	if !hasRootSpec {
		return nil, fmt.Errorf("validator: catalog missing VAST spec")
	}
	rootVersionSupported := rootSpec.supports(version)

	rootPointer := buildSourcePointer("", root.localName(), 1)
	rootResult := validateNodeRecursive(root, version, cfg, rootSpec, nil, false, rootPointer)
	if !rootVersionSupported {
		iab := rootResult.addAnalysis(IABAnalysisCategory)
		markFailure(iab, fmt.Sprintf("%s: %s", ErrUnsupportedVersion.Error(), version))
	}

	return &ValidationResult{Version: version, Root: rootResult, Summaries: summarizeCategories(rootResult)}, nil
}

func validateNodeRecursive(node *genericNode, version vast.Version, cfg *config, spec *NodeSpec, parentSpec *NodeSpec, parentAllowsUnknown bool, sourcePointer string) *NodeResult {
	result := &NodeResult{
		Node:           node.localName(),
		SourcePointer:  sourcePointer,
		VersionSupport: nil,
	}

	var nodeCaseMismatch string
	if spec == nil {
		if matchedSpec, canonicalName, ok := cfg.catalog.nodeCaseInsensitive(result.Node); ok {
			spec = matchedSpec
			nodeCaseMismatch = canonicalName
		}
	}
	extensionBackport := spec != nil && spec.SupportsExtensions && isExtensionContainerSpec(parentSpec)

	if spec != nil {
		result.VersionSupport = spec.Versions
		result.IntroducedAt = introducedAtFromVersions(spec.Versions)
	}

	iabAnalysis := result.addAnalysis(IABAnalysisCategory)
	if spec == nil {
		if !parentAllowsUnknown {
			markFailure(iabAnalysis, fmt.Sprintf("node %s is not a recognized VAST node. Check you have the correct spelling and casing.", result.Node))
		}
	} else {
		if nodeCaseMismatch != "" && nodeCaseMismatch != result.Node {
			markFailure(iabAnalysis, fmt.Sprintf("node %s casing is invalid; use %s", result.Node, nodeCaseMismatch))
		}
		if !spec.supports(version) && !extensionBackport {
			markFailure(iabAnalysis, fmt.Sprintf("node %s is not supported in VAST %s", result.Node, version))
		}
		if parentSpec != nil && !parentAllowsUnknown {
			childSpec, ok := parentSpec.child(result.Node)
			childCaseMismatch := ""
			if !ok {
				if matchedChild, canonicalName, ok := parentSpec.childCaseInsensitive(result.Node); ok {
					childSpec = matchedChild
					childCaseMismatch = canonicalName
					ok = true
				}
			}
			if !ok {
				markFailure(iabAnalysis, fmt.Sprintf("node %s is not a valid child of %s", result.Node, parentSpec.Name))
			} else {
				if childCaseMismatch != "" && childCaseMismatch != result.Node {
					markFailure(iabAnalysis, fmt.Sprintf("child node %s casing is invalid for parent %s; use %s", result.Node, parentSpec.Name, childCaseMismatch))
				}
				if !childSpec.supports(version) {
					markFailure(iabAnalysis, fmt.Sprintf("node %s is not allowed for parent %s in VAST %s", result.Node, parentSpec.Name, version))
				}
			}
		}
	}

	if !parentAllowsUnknown || extensionBackport {
		validateAttributes(node, version, spec, iabAnalysis, extensionBackport)
	}

	if isExtensionContainerSpec(spec) {
		applyExtensionValidators(result, node, version)
	}

	if cfg.runCustom {
		applyCustomValidators(result, node, version)
	}
	if cfg.runHTTP {
		applyHTTPValidators(result, node, version, cfg)
	}

	childAllowsUnknown := parentAllowsUnknown
	if spec != nil && spec.AllowUnknownChildren {
		childAllowsUnknown = true
	}

	childOccurrences := map[string]int{}
	for _, child := range node.Children {
		childName := child.localName()
		childOccurrences[childName]++
		childSpec, _ := cfg.catalog.node(childName)
		childPointer := buildSourcePointer(sourcePointer, childName, childOccurrences[childName])
		childResult := validateNodeRecursive(child, version, cfg, childSpec, spec, childAllowsUnknown, childPointer)
		result.Children = append(result.Children, childResult)
	}

	return result
}

// applyExtensionValidators executes registered extension validators that match the given node and merges their results into the provided node result.
func buildSourcePointer(parentPointer, nodeName string, occurrence int) string {
	if nodeName == "" {
		return parentPointer
	}
	if occurrence < 1 {
		occurrence = 1
	}
	if parentPointer == "" {
		return fmt.Sprintf("/%s[%d]", nodeName, occurrence)
	}
	return fmt.Sprintf("%s/%s[%d]", parentPointer, nodeName, occurrence)
}

func validateAttributes(node *genericNode, version vast.Version, spec *NodeSpec, analysis *NodeAnalysisResult, allowBackport bool) {
	seen := map[string]bool{}

	for _, attr := range node.Attrs {
		attrName := attr.Name.Local
		if attr.Name.Space != "" || strings.EqualFold(attrName, "xmlns") {
			// Skip namespace declarations or namespace-scoped attributes; they are not part of VAST validation.
			continue
		}
		attributeResult := AttributeResult{Name: attrName, Status: StatusPass}
		resolvedName := attrName

		if spec == nil {
			seen[resolvedName] = true
			attributeResult.Status = StatusFail
			msg := "node is not recognized; attribute cannot be validated"
			attributeResult.addReason(msg)
			analysis.addAttribute(attributeResult)
			markFailure(analysis, msg)
			continue
		}

		attrSpec, ok := spec.attribute(attrName)
		caseMismatchName := ""
		if !ok {
			if matchedAttr, canonicalName, matchOk := spec.attributeCaseInsensitive(attrName); matchOk {
				attrSpec = matchedAttr
				caseMismatchName = canonicalName
				resolvedName = canonicalName
				ok = true
			}
		}
		seen[resolvedName] = true

		if !ok {
			attributeResult.Status = StatusFail
			msg := fmt.Sprintf("attribute %s is not allowed on %s for version %s", attrName, spec.Name, version)
			attributeResult.addReason(msg)
			analysis.addAttribute(attributeResult)
			markFailure(analysis, msg)
			continue
		}
		attributeResult.VersionSupport = attrSpec.Versions
		attributeResult.IntroducedAt = introducedAtFromVersions(attrSpec.Versions)

		if caseMismatchName != "" && caseMismatchName != attrName {
			attributeResult.Status = StatusFail
			msg := fmt.Sprintf("attribute %s casing is invalid; use %s", attrName, caseMismatchName)
			attributeResult.addReason(msg)
			markFailure(analysis, msg)
		}

		if !attrSpec.supports(version) && !allowBackport {
			attributeResult.Status = StatusFail
			msg := fmt.Sprintf("attribute %s is not supported in VAST %s", attrName, version)
			attributeResult.addReason(msg)
			markFailure(analysis, msg)
		}

		value := strings.TrimSpace(attr.Value)
		if value == "" && !attrSpec.AllowEmpty {
			attributeResult.Status = StatusFail
			msg := fmt.Sprintf("attribute %s cannot be empty", attrName)
			attributeResult.addReason(msg)
			markFailure(analysis, msg)
		} else {
			attributeResult.Value = value
		}

		analysis.addAttribute(attributeResult)
	}

	if spec == nil {
		return
	}

	for _, attrSpec := range spec.Attributes {
		if !attrSpec.Required {
			continue
		}
		if seen[attrSpec.Name] {
			continue
		}
		msg := fmt.Sprintf("missing required attribute %s", attrSpec.Name)
		analysis.addAttribute(AttributeResult{
			Name:           attrSpec.Name,
			IntroducedAt:   introducedAtFromVersions(attrSpec.Versions),
			VersionSupport: attrSpec.Versions,
			Status:         StatusFail,
			Reasons:        []string{msg},
		})
		markFailure(analysis, msg)
	}
}

func isExtensionContainerSpec(spec *NodeSpec) bool {
	if spec == nil {
		return false
	}
	switch spec.Name {
	case "Extension", "CreativeExtension":
		return true
	default:
		return false
	}
}

func introducedAtFromVersions(versions []vast.Version) *float64 {
	if len(versions) == 0 {
		return nil
	}
	var min float64
	found := false
	for _, version := range versions {
		value, ok := vastVersionToFloat(version)
		if !ok {
			continue
		}
		if !found || value < min {
			min = value
			found = true
		}
	}
	if !found {
		return nil
	}
	return &min
}

func vastVersionToFloat(version vast.Version) (float64, bool) {
	trimmed := strings.TrimSpace(string(version))
	if trimmed == "" {
		return 0, false
	}
	value, err := strconv.ParseFloat(trimmed, 64)
	if err != nil {
		return 0, false
	}
	return value, true
}

func applyCustomValidators(nodeResult *NodeResult, node *genericNode, version vast.Version) {
	for _, validator := range getCustomValidators(nodeResult.Node) {
		analysis := validator(NodeContext{Node: node, Version: version})
		if analysis == nil {
			continue
		}
		if analysis.Category == "" {
			analysis.Category = CustomAnalysisCategory
		}
		mergeAnalysis(nodeResult, analysis)
	}
}

func applyHTTPValidators(nodeResult *NodeResult, node *genericNode, version vast.Version, cfg *config) {
	validators := getHTTPValidators(nodeResult.Node)
	if len(validators) == 0 {
		return
	}
	ctx := context.Background()
	if cfg.httpOptions.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, cfg.httpOptions.Timeout)
		defer cancel()
	}
	client := cfg.httpOptions.client()
	for _, validator := range validators {
		analysis, err := validator(ctx, NodeContext{Node: node, Version: version}, client)
		if err != nil {
			analysis = &NodeAnalysisResult{Category: CustomAnalysisCategory}
			markFailure(analysis, err.Error())
		}
		if analysis == nil {
			continue
		}
		if analysis.Category == "" {
			analysis.Category = CustomAnalysisCategory
		}
		mergeAnalysis(nodeResult, analysis)
	}
}

func mergeAnalysis(nodeResult *NodeResult, analysis *NodeAnalysisResult) {
	if nodeResult.Analyses == nil {
		nodeResult.Analyses = make(map[string]*NodeAnalysisResult)
	}
	existing := nodeResult.Analyses[analysis.Category]
	if existing == nil {
		nodeResult.Analyses[analysis.Category] = analysis
		return
	}
	existing.Attributes = append(existing.Attributes, analysis.Attributes...)
	if analysis.Status == StatusFail {
		markFailure(existing, analysis.Reasons...)
	}
}

func markFailure(analysis *NodeAnalysisResult, reasons ...string) {
	if analysis == nil {
		return
	}
	if analysis.Status != StatusFail {
		analysis.Status = StatusFail
	}
	for _, reason := range reasons {
		if reason == "" {
			continue
		}
		analysis.Reasons = append(analysis.Reasons, reason)
	}
}
