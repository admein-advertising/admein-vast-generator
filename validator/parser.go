package validator

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strings"
)

var errEmptyXML = errors.New("validator: empty XML document")

// genericNode represents a light-weight XML node used for validation traversal.
type genericNode struct {
	Name     xml.Name
	Attrs    []xml.Attr
	Children []*genericNode
	Content  string
}

func (n *genericNode) localName() string {
	return n.Name.Local
}

func (n *genericNode) attrValue(name string) (string, bool) {
	for _, attr := range n.Attrs {
		if attr.Name.Local == name {
			return attr.Value, true
		}
	}
	return "", false
}

// buildNodeTree parses raw XML bytes into a tree of genericNode instances.
func buildNodeTree(raw []byte) (*genericNode, error) {
	decoder := xml.NewDecoder(bytes.NewReader(raw))
	var stack []*genericNode
	var root *genericNode

	for {
		token, err := decoder.Token()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, fmt.Errorf("validator: parse XML: %w", err)
		}

		switch typed := token.(type) {
		case xml.StartElement:
			node := &genericNode{Name: typed.Name, Attrs: typed.Attr}
			if len(stack) == 0 {
				root = node
			} else {
				parent := stack[len(stack)-1]
				parent.Children = append(parent.Children, node)
			}
			stack = append(stack, node)

		case xml.EndElement:
			if len(stack) == 0 {
				return nil, fmt.Errorf("validator: unexpected closing tag %q", typed.Name.Local)
			}
			stack = stack[:len(stack)-1]

		case xml.CharData:
			if len(stack) == 0 {
				continue
			}
			trimmed := strings.TrimSpace(string(typed))
			if trimmed == "" {
				continue
			}
			current := stack[len(stack)-1]
			if current.Content != "" {
				current.Content += " "
			}
			current.Content += trimmed
		}
	}

	if root == nil {
		return nil, errEmptyXML
	}

	return root, nil
}
