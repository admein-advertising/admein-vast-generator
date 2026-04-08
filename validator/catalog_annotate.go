package validator

import (
	"fmt"
	"strings"
)

// docIsEmpty checks if a Documentation value is nil or has only whitespace.
func docIsEmpty(doc *Documentation) bool {
	return doc == nil || strings.TrimSpace(doc.Content) == ""
}

// schemaDocumentation creates a Documentation from text, returning nil if text is empty.
func schemaDocumentation(text string) *Documentation {
	trimmed := strings.TrimSpace(text)
	if trimmed == "" {
		return nil
	}
	return &Documentation{Content: trimmed, Source: vast42SchemaURL}
}

// appendContextPath appends a name to a context path, returning a new slice.
func appendContextPath(path []string, name string) []string {
	next := make([]string, len(path)+1)
	copy(next, path)
	next[len(path)] = name
	return next
}

// joinContextPath joins a context path into a ">" delimited string.
func joinContextPath(path []string) string {
	return strings.Join(path, ">")
}

// annotateCatalogDocs is the main entry point for annotation logic, ensuring all nodes and their
// relationships have appropriate documentation.
func annotateCatalogDocs(cat *Catalog) {
	if cat == nil {
		return
	}
	for _, node := range cat.Nodes {
		annotateNodeDoc(node)
	}
	promoteContextualNodeDocs(cat)
	promoteContextualAttributeDocs(cat)
	root, ok := cat.node("VAST")
	if ok && root != nil {
		annotateChildDocs(cat, "VAST", root, []string{"VAST"}, map[string]bool{})
		return
	}
	for key, node := range cat.Nodes {
		annotateChildDocs(cat, key, node, []string{node.Name}, map[string]bool{})
	}
}

// promoteContextualNodeDocs promotes context-specific documentation for nodes that currently
// only have generic XSD fallback text, if an unambiguous context-specific doc exists.
func promoteContextualNodeDocs(cat *Catalog) {
	if cat == nil {
		return
	}
	candidates := collectContextualNodeDocCandidates(cat)
	for nodeKey, node := range cat.Nodes {
		if node == nil || !isGenericElementFallbackDoc(node.Documentation) {
			continue
		}
		candidate, ok := selectUnambiguousDoc(candidates[nodeKey])
		if !ok {
			continue
		}
		if candidates[nodeKey][candidate] < 2 {
			continue
		}
		if doc := schemaDocumentation(candidate); doc != nil {
			node.Documentation = doc
		}
	}
}

// collectContextualNodeDocCandidates collects all context-specific docs for each node by traversing
// the catalog starting from root or all top-level nodes.
func collectContextualNodeDocCandidates(cat *Catalog) map[string]map[string]int {
	out := map[string]map[string]int{}
	if cat == nil {
		return out
	}
	root, ok := cat.node("VAST")
	if ok && root != nil {
		collectContextualNodeDocCandidatesFromPath(cat, "VAST", root, []string{"VAST"}, map[string]bool{}, out)
		return out
	}
	for key, node := range cat.Nodes {
		if node == nil {
			continue
		}
		collectContextualNodeDocCandidatesFromPath(cat, key, node, []string{node.Name}, map[string]bool{}, out)
	}
	return out
}

// collectContextualNodeDocCandidatesFromPath traverses the catalog hierarchy from a given node,
// collecting context-specific docs from vast42ContextElementDocs for each path.
func collectContextualNodeDocCandidatesFromPath(cat *Catalog, nodeKey string, node *NodeSpec, path []string, visiting map[string]bool, out map[string]map[string]int) {
	if cat == nil || node == nil {
		return
	}
	stateKey := nodeKey + "|" + joinContextPath(path)
	if visiting[stateKey] {
		return
	}
	visiting[stateKey] = true
	defer delete(visiting, stateKey)
	if docText := strings.TrimSpace(vast42ContextElementDocs[joinContextPath(path)]); docText != "" {
		if out[nodeKey] == nil {
			out[nodeKey] = map[string]int{}
		}
		out[nodeKey][docText]++
	}

	for _, child := range node.Children {
		if child == nil {
			continue
		}
		childPath := appendContextPath(path, child.Name)
		contextPath := joinContextPath(childPath)
		docText := strings.TrimSpace(vast42ContextElementDocs[contextPath])

		targetKey := child.Name
		if child.NodeOverride != "" {
			targetKey = child.NodeOverride
		}
		if docText != "" {
			if out[targetKey] == nil {
				out[targetKey] = map[string]int{}
			}
			out[targetKey][docText]++
		}
		if targetNode, ok := cat.node(targetKey); ok && targetNode != nil {
			collectContextualNodeDocCandidatesFromPath(cat, targetKey, targetNode, childPath, visiting, out)
		}
	}
}

// selectUnambiguousDoc returns a doc string if there is exactly one candidate, empty/false otherwise.
func selectUnambiguousDoc(candidates map[string]int) (string, bool) {
	if len(candidates) != 1 {
		return "", false
	}
	for doc := range candidates {
		trimmed := strings.TrimSpace(doc)
		if trimmed != "" {
			return trimmed, true
		}
	}
	return "", false
}

// isGenericElementFallbackDoc checks if a Documentation is a generic XSD element fallback message.
func isGenericElementFallbackDoc(doc *Documentation) bool {
	if doc == nil {
		return false
	}
	return strings.HasPrefix(strings.TrimSpace(doc.Content), "Defined in VAST 4.2 XSD element <")
}

// annotateNodeDoc ensures each node and its attributes have documentation.
func annotateNodeDoc(node *NodeSpec) {
	if node == nil {
		return
	}
	if docIsEmpty(node.Documentation) {
		node.Documentation = &Documentation{
			Content: fmt.Sprintf("Defined in VAST 4.2 XSD element <%s>.", node.Name),
			Source:  vast42SchemaURL,
		}
	}
	for _, attr := range node.Attributes {
		annotateAttributeDoc(node.Name, attr)
	}
}

// annotateChildDocs recursively ensures child relationships have appropriate documentation.
func annotateChildDocs(cat *Catalog, nodeKey string, node *NodeSpec, path []string, visiting map[string]bool) {
	if node == nil || cat == nil {
		return
	}
	stateKey := nodeKey + "|" + joinContextPath(path)
	if visiting[stateKey] {
		return
	}
	visiting[stateKey] = true
	defer delete(visiting, stateKey)
	for _, child := range node.Children {
		childPath := appendContextPath(path, child.Name)
		annotateChildDoc(cat, node.Name, child, joinContextPath(childPath))
		targetKey := child.Name
		if child.NodeOverride != "" {
			targetKey = child.NodeOverride
		}
		targetNode, ok := cat.node(targetKey)
		if ok && targetNode != nil {
			annotateChildDocs(cat, targetKey, targetNode, childPath, visiting)
		}
	}
}

// annotateAttributeDoc ensures an attribute and its value spec have documentation.
func annotateAttributeDoc(nodeName string, attr *AttributeSpec) {
	if attr == nil {
		return
	}
	if docIsEmpty(attr.Documentation) {
		attr.Documentation = &Documentation{
			Content: fmt.Sprintf("Attribute @%s on <%s> per VAST 4.2 XSD.", attr.Name, nodeName),
			Source:  vast42SchemaURL,
		}
	}
	if attr.Value != nil && docIsEmpty(attr.Value.Documentation) {
		attr.Value.Documentation = &Documentation{
			Content: fmt.Sprintf("Constraints for @%s on <%s> defined in VAST 4.2 XSD.", attr.Name, nodeName),
			Source:  vast42SchemaURL,
		}
	}
}

// annotateChildDoc ensures a child relationship has appropriate documentation from context-specific sources,
// node inheritance, or a generic fallback.
func annotateChildDoc(cat *Catalog, parent string, child *ChildSpec, contextPath string) {
	if child == nil {
		return
	}
	if docIsEmpty(child.Documentation) {
		if doc := schemaDocumentation(vast42ContextElementDocs[contextPath]); doc != nil && !isGenericElementFallbackDoc(doc) {
			child.Documentation = doc
		}
	}
	if docIsEmpty(child.Documentation) && cat != nil {
		lookupName := child.Name
		if child.NodeOverride != "" {
			lookupName = child.NodeOverride
		}
		if target, ok := cat.node(lookupName); ok && target != nil && !docIsEmpty(target.Documentation) {
			canInheritFromTarget := child.NodeOverride != ""
			if canInheritFromTarget && !strings.HasPrefix(strings.TrimSpace(target.Documentation.Content), "Defined in VAST 4.2 XSD element <") {
				child.Documentation = cloneDocumentation(target.Documentation)
			}
		}
	}
	if docIsEmpty(child.Documentation) {
		child.Documentation = &Documentation{
			Content: fmt.Sprintf("Child <%s> permitted within <%s> per VAST 4.2 XSD.", child.Name, parent),
			Source:  vast42SchemaURL,
		}
	}
}

// promoteContextualAttributeDocs promotes context-specific documentation for attributes that currently
// only have generic XSD fallback text, if an unambiguous context-specific doc exists.
func promoteContextualAttributeDocs(cat *Catalog) {
	if cat == nil {
		return
	}
	candidates := collectContextualAttributeDocCandidates(cat)
	for nodeKey, attrs := range candidates {
		node, ok := cat.node(nodeKey)
		if !ok || node == nil {
			continue
		}
		for attrName, docs := range attrs {
			attr, ok := node.attribute(attrName)
			if !ok || attr == nil || !isGenericAttributeFallbackDoc(attr.Documentation) {
				continue
			}
			candidate, ok := selectUnambiguousDoc(docs)
			if !ok {
				continue
			}
			if docs[candidate] < 2 {
				continue
			}
			if doc := schemaDocumentation(candidate); doc != nil {
				attr.Documentation = doc
			}
		}
	}
}

// collectContextualAttributeDocCandidates collects all context-specific attribute docs for each node by
// traversing the catalog starting from root or all top-level nodes.
func collectContextualAttributeDocCandidates(cat *Catalog) map[string]map[string]map[string]int {
	out := map[string]map[string]map[string]int{}
	if cat == nil {
		return out
	}
	root, ok := cat.node("VAST")
	if ok && root != nil {
		collectContextualAttributeDocCandidatesFromPath(cat, "VAST", root, []string{"VAST"}, map[string]bool{}, out)
		return out
	}
	for key, node := range cat.Nodes {
		if node == nil {
			continue
		}
		collectContextualAttributeDocCandidatesFromPath(cat, key, node, []string{node.Name}, map[string]bool{}, out)
	}
	return out
}

// collectContextualAttributeDocCandidatesFromPath traverses the catalog hierarchy from a given node,
// collecting context-specific attribute docs from vast42AttributeDocs for each path.
func collectContextualAttributeDocCandidatesFromPath(cat *Catalog, nodeKey string, node *NodeSpec, path []string, visiting map[string]bool, out map[string]map[string]map[string]int) {
	if cat == nil || node == nil {
		return
	}
	stateKey := nodeKey + "|" + joinContextPath(path)
	if visiting[stateKey] {
		return
	}
	visiting[stateKey] = true
	defer delete(visiting, stateKey)

	if scoped, ok := vast42AttributeDocs[joinContextPath(path)]; ok {
		for attrName, attrDoc := range scoped {
			trimmed := strings.TrimSpace(attrDoc)
			if trimmed == "" {
				continue
			}
			if out[nodeKey] == nil {
				out[nodeKey] = map[string]map[string]int{}
			}
			if out[nodeKey][attrName] == nil {
				out[nodeKey][attrName] = map[string]int{}
			}
			out[nodeKey][attrName][trimmed]++
		}
	}

	for _, child := range node.Children {
		if child == nil {
			continue
		}
		targetKey := child.Name
		if child.NodeOverride != "" {
			targetKey = child.NodeOverride
		}
		targetNode, ok := cat.node(targetKey)
		if !ok || targetNode == nil {
			continue
		}
		collectContextualAttributeDocCandidatesFromPath(cat, targetKey, targetNode, appendContextPath(path, child.Name), visiting, out)
	}
}

// isGenericAttributeFallbackDoc checks if a Documentation is a generic attribute fallback message.
func isGenericAttributeFallbackDoc(doc *Documentation) bool {
	if doc == nil {
		return false
	}
	return strings.HasPrefix(strings.TrimSpace(doc.Content), "Attribute @")
}
