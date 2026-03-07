package validator

import (
	"strings"
	"sync"

	"github.com/admein-advertising/admein-vast-generator/vast"
)

type ExtensionValidationContext struct {
	NodeContext
}

type ExtensionMatchFunc func(ctx ExtensionValidationContext) bool

type ExtensionValidatorFunc func(ctx ExtensionValidationContext) *NodeAnalysisResult

type ExtensionValidatorConfig struct {
	Name     string
	Types    []string
	Match    ExtensionMatchFunc
	Validate ExtensionValidatorFunc
}

type extensionValidatorEntry struct {
	name  string
	types []string
	match ExtensionMatchFunc
	fn    ExtensionValidatorFunc
}

var (
	extensionValidatorsMu sync.RWMutex
	extensionValidators   []extensionValidatorEntry
)

func RegisterExtensionValidator(cfg ExtensionValidatorConfig) {
	if cfg.Validate == nil {
		return
	}

	entry := extensionValidatorEntry{
		name: cfg.Name,
		fn:   cfg.Validate,
	}

	if len(cfg.Types) > 0 {
		entry.types = make([]string, 0, len(cfg.Types))
		for _, t := range cfg.Types {
			trimmed := strings.TrimSpace(t)
			if trimmed == "" {
				continue
			}
			entry.types = append(entry.types, strings.ToLower(trimmed))
		}
	}

	entry.match = cfg.Match

	extensionValidatorsMu.Lock()
	extensionValidators = append(extensionValidators, entry)
	extensionValidatorsMu.Unlock()
}

func applyExtensionValidators(nodeResult *NodeResult, node *genericNode, version vast.Version) {
	validators := snapshotExtensionValidators()
	if len(validators) == 0 {
		return
	}

	ctx := ExtensionValidationContext{NodeContext: NodeContext{Node: node, Version: version}}
	for _, validator := range validators {
		if !validator.matches(ctx) {
			continue
		}
		analysis := validator.fn(ctx)
		if analysis == nil {
			continue
		}
		if analysis.Category == "" {
			analysis.Category = IABAnalysisCategory
		}
		mergeAnalysis(nodeResult, analysis)
	}
}

func snapshotExtensionValidators() []extensionValidatorEntry {
	extensionValidatorsMu.RLock()
	defer extensionValidatorsMu.RUnlock()
	if len(extensionValidators) == 0 {
		return nil
	}
	out := make([]extensionValidatorEntry, len(extensionValidators))
	copy(out, extensionValidators)
	return out
}

func (entry extensionValidatorEntry) matches(ctx ExtensionValidationContext) bool {
	if len(entry.types) > 0 {
		extType := strings.ToLower(ctx.Type())
		for _, t := range entry.types {
			if extType != "" && extType == t {
				return true
			}
		}
	}
	if entry.match != nil && entry.match(ctx) {
		return true
	}
	return len(entry.types) == 0 && entry.match == nil
}

func (ctx ExtensionValidationContext) Type() string {
	value, _ := ctx.Attribute("type")
	return strings.TrimSpace(value)
}

func (ctx ExtensionValidationContext) ChildrenNamed(name string) []*genericNode {
	if ctx.Node == nil {
		return nil
	}
	var children []*genericNode
	for _, child := range ctx.Node.Children {
		if strings.EqualFold(child.localName(), name) {
			children = append(children, child)
		}
	}
	return children
}

func (ctx ExtensionValidationContext) HasChildNamed(name string) bool {
	return len(ctx.ChildrenNamed(name)) > 0
}

func resetExtensionValidators() {
	extensionValidatorsMu.Lock()
	extensionValidators = nil
	extensionValidatorsMu.Unlock()
	registerBuiltInExtensionValidators()
}

func init() {
	registerBuiltInExtensionValidators()
}

func registerBuiltInExtensionValidators() {
	RegisterExtensionValidator(ExtensionValidatorConfig{
		Name:  "UniversalAdIdExtension",
		Types: []string{"UniversalAdId"},
		Match: func(ctx ExtensionValidationContext) bool {
			return ctx.HasChildNamed("UniversalAdId")
		},
		Validate: universalAdIDExtensionValidator,
	})

	RegisterExtensionValidator(ExtensionValidatorConfig{
		Name:  "InteractiveCreativeFileExtension",
		Types: []string{"InteractiveCreativeFile"},
		Match: func(ctx ExtensionValidationContext) bool {
			return ctx.HasChildNamed("InteractiveCreativeFile")
		},
		Validate: interactiveCreativeFileExtensionValidator,
	})

	RegisterExtensionValidator(ExtensionValidatorConfig{
		Name:  "MezzanineExtension",
		Types: []string{"Mezzanine"},
		Match: func(ctx ExtensionValidationContext) bool {
			return ctx.HasChildNamed("Mezzanine")
		},
		Validate: mezzanineExtensionValidator,
	})
}

func universalAdIDExtensionValidator(ctx ExtensionValidationContext) *NodeAnalysisResult {
	nodes := ctx.ChildrenNamed("UniversalAdId")
	var report *NodeAnalysisResult

	ensureReport := func() *NodeAnalysisResult {
		if report == nil {
			report = &NodeAnalysisResult{Category: IABAnalysisCategory, Status: StatusPass}
		}
		return report
	}

	if len(nodes) == 0 {
		analysis := ensureReport()
		markFailure(analysis, "UniversalAdId extension must include at least one UniversalAdId node")
		return analysis
	}

	extType := ctx.Type()
	if extType == "" {
		analysis := ensureReport()
		markFailure(analysis, "UniversalAdId extension should declare type=\"UniversalAdId\"")
	} else if !strings.EqualFold(extType, "UniversalAdId") {
		analysis := ensureReport()
		markFailure(analysis, "UniversalAdId extension type attribute value should be \"UniversalAdId\"")
	}

	for _, node := range nodes {
		if strings.TrimSpace(node.Content) == "" {
			analysis := ensureReport()
			markFailure(analysis, "UniversalAdId value must not be empty")
			break
		}
	}

	if report != nil && report.Status == StatusPass && len(report.Reasons) == 0 && len(report.Attributes) == 0 {
		return nil
	}
	return report
}

func interactiveCreativeFileExtensionValidator(ctx ExtensionValidationContext) *NodeAnalysisResult {
	nodes := ctx.ChildrenNamed("InteractiveCreativeFile")
	var report *NodeAnalysisResult

	ensureReport := func() *NodeAnalysisResult {
		if report == nil {
			report = &NodeAnalysisResult{Category: IABAnalysisCategory, Status: StatusPass}
		}
		return report
	}

	if len(nodes) == 0 {
		analysis := ensureReport()
		markFailure(analysis, "InteractiveCreativeFile extension must include at least one InteractiveCreativeFile node")
		return analysis
	}

	extType := ctx.Type()
	if extType == "" {
		analysis := ensureReport()
		markFailure(analysis, "InteractiveCreativeFile extension should declare type=\"InteractiveCreativeFile\"")
	} else if !strings.EqualFold(extType, "InteractiveCreativeFile") {
		analysis := ensureReport()
		markFailure(analysis, "InteractiveCreativeFile extension type attribute value should be \"InteractiveCreativeFile\"")
	}

	for _, node := range nodes {
		if strings.TrimSpace(node.Content) == "" {
			analysis := ensureReport()
			markFailure(analysis, "InteractiveCreativeFile must include executable content or a URL")
			break
		}
	}

	if report != nil && report.Status == StatusPass && len(report.Reasons) == 0 && len(report.Attributes) == 0 {
		return nil
	}
	return report
}

func mezzanineExtensionValidator(ctx ExtensionValidationContext) *NodeAnalysisResult {
	nodes := ctx.ChildrenNamed("Mezzanine")
	var report *NodeAnalysisResult

	ensureReport := func() *NodeAnalysisResult {
		if report == nil {
			report = &NodeAnalysisResult{Category: IABAnalysisCategory, Status: StatusPass}
		}
		return report
	}

	if len(nodes) == 0 {
		analysis := ensureReport()
		markFailure(analysis, "Mezzanine extension must include at least one Mezzanine node")
		return analysis
	}

	extType := ctx.Type()
	if extType == "" {
		analysis := ensureReport()
		markFailure(analysis, "Mezzanine extension should declare type=\"Mezzanine\"")
	} else if !strings.EqualFold(extType, "Mezzanine") {
		analysis := ensureReport()
		markFailure(analysis, "Mezzanine extension type attribute value should be \"Mezzanine\"")
	}

	for _, node := range nodes {
		if strings.TrimSpace(node.Content) == "" {
			analysis := ensureReport()
			markFailure(analysis, "Mezzanine value must not be empty")
			break
		}
	}

	if report != nil && report.Status == StatusPass && len(report.Reasons) == 0 && len(report.Attributes) == 0 {
		return nil
	}
	return report
}
