package validator

import "github.com/admein-advertising/admein-vast-generator/vast"

// ResultStatus represents the outcome of a validation rule.
type ResultStatus string

const (
	StatusPass ResultStatus = "pass"
	StatusFail ResultStatus = "fail"
	StatusInfo ResultStatus = "info"
)

// AttributeResult captures the outcome of validating a single attribute.
type AttributeResult struct {
	Name           string         `json:"name"`
	VersionSupport []vast.Version `json:"versionSupport,omitempty"`
	Status         ResultStatus   `json:"status"`
	Reason         string         `json:"reason,omitempty"`
}

// NodeAnalysisResult encapsulates all results for a specific analysis category
// (e.g., "iab.analysis" or "custom.analysis") at the node level.
type NodeAnalysisResult struct {
	Category   string            `json:"category"`
	Status     ResultStatus      `json:"status"`
	Reason     string            `json:"reason,omitempty"`
	Attributes []AttributeResult `json:"attributes,omitempty"`
}

// addAttribute appends an attribute result to the analysis bucket.
func (nar *NodeAnalysisResult) addAttribute(result AttributeResult) {
	nar.Attributes = append(nar.Attributes, result)
}

// NodeResult represents the validation result of a node, including one or more
// analysis categories and nested child results.
type NodeResult struct {
	Node           string                         `json:"node"`
	VersionSupport []vast.Version                 `json:"versionSupport,omitempty"`
	Analyses       map[string]*NodeAnalysisResult `json:"analyses,omitempty"`
	Children       []*NodeResult                  `json:"children,omitempty"`
}

// addAnalysis ensures there is an analysis bucket for the given category and returns it.
func (nr *NodeResult) addAnalysis(category string) *NodeAnalysisResult {
	if nr.Analyses == nil {
		nr.Analyses = make(map[string]*NodeAnalysisResult)
	}
	analysis, ok := nr.Analyses[category]
	if !ok {
		analysis = &NodeAnalysisResult{Category: category, Status: StatusPass}
		nr.Analyses[category] = analysis
	}
	return analysis
}

// ValidationResult is the root object returned by the validator.
type ValidationResult struct {
	Version   vast.Version                `json:"version"`
	Root      *NodeResult                 `json:"root"`
	Summaries map[string]*CategorySummary `json:"summaries,omitempty"`
}

// CategorySummary aggregates node results per analysis category for quick UI consumption.
type CategorySummary struct {
	Category     string       `json:"category"`
	TotalNodes   int          `json:"totalNodes"`
	FailingNodes int          `json:"failingNodes"`
	Status       ResultStatus `json:"status"`
	Reasons      []string     `json:"reasons,omitempty"`
}

func summarizeCategories(root *NodeResult) map[string]*CategorySummary {
	if root == nil {
		return nil
	}
	summaries := map[string]*CategorySummary{}
	var walk func(node *NodeResult)
	walk = func(node *NodeResult) {
		if node == nil {
			return
		}
		for category, analysis := range node.Analyses {
			summary := summaries[category]
			if summary == nil {
				summary = &CategorySummary{Category: category, Status: StatusPass}
				summaries[category] = summary
			}
			summary.TotalNodes++
			if analysis.Status == StatusFail {
				summary.FailingNodes++
				summary.Status = StatusFail
				if analysis.Reason != "" && len(summary.Reasons) < 5 {
					summary.Reasons = append(summary.Reasons, analysis.Reason)
				}
			}
		}
		for _, child := range node.Children {
			walk(child)
		}
	}
	walk(root)
	if len(summaries) == 0 {
		return nil
	}
	return summaries
}
