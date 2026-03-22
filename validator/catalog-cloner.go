package validator

import "github.com/admein-advertising/admein-vast-generator/vast"

func cloneCatalog(src *Catalog) *Catalog {
	if src == nil {
		return nil
	}
	cloned := &Catalog{Nodes: make(map[string]*NodeSpec, len(src.Nodes))}
	for name, node := range src.Nodes {
		cloned.Nodes[name] = cloneNodeSpec(node)
	}
	return cloned
}

func cloneNodeSpec(src *NodeSpec) *NodeSpec {
	if src == nil {
		return nil
	}
	cloned := &NodeSpec{
		Name:                   src.Name,
		Versions:               cloneVersions(src.Versions),
		AllowUnknownChildren:   src.AllowUnknownChildren,
		AllowUnknownAttributes: src.AllowUnknownAttributes,
		SupportsExtensions:     src.SupportsExtensions,
		NeedsCDATA:             src.NeedsCDATA,
		Documentation:          cloneDocumentation(src.Documentation),
	}
	if len(src.Attributes) > 0 {
		cloned.Attributes = make(map[string]*AttributeSpec, len(src.Attributes))
		for name, attr := range src.Attributes {
			cloned.Attributes[name] = cloneAttributeSpec(attr)
		}
	}
	if len(src.Children) > 0 {
		cloned.Children = make(map[string]*ChildSpec, len(src.Children))
		for name, child := range src.Children {
			cloned.Children[name] = cloneChildSpec(child)
		}
	}
	return cloned
}

func cloneAttributeSpec(src *AttributeSpec) *AttributeSpec {
	if src == nil {
		return nil
	}
	return &AttributeSpec{
		Name:          src.Name,
		Versions:      cloneVersions(src.Versions),
		Required:      src.Required,
		AllowEmpty:    src.AllowEmpty,
		Value:         cloneAttributeValueSpec(src.Value),
		Documentation: cloneDocumentation(src.Documentation),
	}
}

func cloneAttributeValueSpec(src *AttributeValueSpec) *AttributeValueSpec {
	if src == nil {
		return nil
	}
	cloned := &AttributeValueSpec{Type: src.Type, Pattern: src.Pattern, Documentation: cloneDocumentation(src.Documentation)}
	if len(src.AllowedValues) > 0 {
		cloned.AllowedValues = make([]string, len(src.AllowedValues))
		copy(cloned.AllowedValues, src.AllowedValues)
	}
	return cloned
}

func cloneChildSpec(src *ChildSpec) *ChildSpec {
	if src == nil {
		return nil
	}
	return &ChildSpec{
		Name:          src.Name,
		Versions:      cloneVersions(src.Versions),
		Optional:      src.Optional,
		Multiple:      src.Multiple,
		Documentation: cloneDocumentation(src.Documentation),
	}
}

func cloneDocumentation(src *Documentation) *Documentation {
	if src == nil {
		return nil
	}
	copy := *src
	return &copy
}

func cloneVersions(src []vast.Version) []vast.Version {
	if len(src) == 0 {
		return nil
	}
	copied := make([]vast.Version, len(src))
	copy(copied, src)
	return copied
}
