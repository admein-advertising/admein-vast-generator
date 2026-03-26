package validator

import (
	"fmt"
	"strings"

	"github.com/admein-advertising/admein-vast-generator/vast"
)

// AttributeType enumerates the primitive XML Schema datatypes we validate for catalog attributes.
type AttributeType string

const (
	AttributeTypeString             AttributeType = "string"
	AttributeTypeToken              AttributeType = "token"
	AttributeTypeBoolean            AttributeType = "boolean"
	AttributeTypeInteger            AttributeType = "integer"
	AttributeTypeNonNegativeInteger AttributeType = "nonNegativeInteger"
	AttributeTypePositiveInteger    AttributeType = "positiveInteger"
	AttributeTypeFloat              AttributeType = "float"
	AttributeTypeDuration           AttributeType = "duration"
	AttributeTypeTimecode           AttributeType = "timecode"
	AttributeTypeTimeOffset         AttributeType = "timeOffset"
	AttributeTypeURI                AttributeType = "anyURI"
)

// Documentation captures human-readable description text along with its source reference.
type Documentation struct {
	Content string `json:"content,omitempty"`
	Source  string `json:"source,omitempty"`
}

// AttributeValueSpec captures datatype and restriction metadata for an attribute value.
type AttributeValueSpec struct {
	Type          AttributeType
	AllowedValues []string
	Pattern       string
	Documentation *Documentation
}

// AttributeSpec describes a valid attribute for a node.
type AttributeSpec struct {
	Name          string
	Versions      []vast.Version
	Required      bool
	AllowEmpty    bool
	Value         *AttributeValueSpec
	Documentation *Documentation
}

// ChildSpec describes a valid child node relationship.
type ChildSpec struct {
	Name          string
	Versions      []vast.Version
	Optional      bool
	Multiple      bool
	NodeOverride  string // Optional catalog node key to use instead of the child's XML name.
	Documentation *Documentation
}

// NodeSpec defines the validation metadata for a node.
type NodeSpec struct {
	Name                   string
	Versions               []vast.Version
	Attributes             map[string]*AttributeSpec
	Children               map[string]*ChildSpec
	AllowUnknownChildren   bool
	AllowUnknownAttributes bool
	SupportsExtensions     bool
	NeedsCDATA             bool // Node text content must be wrapped in CDATA when generating VAST.
	Documentation          *Documentation
}

// Catalog stores node specifications keyed by node name.
type Catalog struct {
	Nodes map[string]*NodeSpec
}

// DefaultVASTCatalog returns a defensive copy of the built-in VAST catalog so
// callers can inspect or serialize it without mutating validator defaults.
func DefaultVASTCatalog() *Catalog {
	return cloneCatalog(defaultCatalog)
}

func (c *Catalog) node(name string) (*NodeSpec, bool) {
	if c == nil {
		return nil, false
	}
	spec, ok := c.Nodes[name]
	return spec, ok
}

func (c *Catalog) nodeCaseInsensitive(name string) (*NodeSpec, string, bool) {
	if c == nil {
		return nil, "", false
	}
	for key, spec := range c.Nodes {
		if strings.EqualFold(key, name) {
			return spec, key, true
		}
	}
	return nil, "", false
}

func (spec *NodeSpec) supports(version vast.Version) bool {
	for _, v := range spec.Versions {
		if v == version {
			return true
		}
	}
	return false
}

func (spec *NodeSpec) attribute(name string) (*AttributeSpec, bool) {
	if spec == nil {
		return nil, false
	}
	attr, ok := spec.Attributes[name]
	return attr, ok
}

func (spec *NodeSpec) attributeCaseInsensitive(name string) (*AttributeSpec, string, bool) {
	if spec == nil {
		return nil, "", false
	}
	for key, attr := range spec.Attributes {
		if strings.EqualFold(key, name) {
			return attr, key, true
		}
	}
	return nil, "", false
}

func (spec *NodeSpec) child(name string) (*ChildSpec, bool) {
	if spec == nil {
		return nil, false
	}
	child, ok := spec.Children[name]
	return child, ok
}

func (spec *NodeSpec) childCaseInsensitive(name string) (*ChildSpec, string, bool) {
	if spec == nil {
		return nil, "", false
	}
	for key, child := range spec.Children {
		if strings.EqualFold(key, name) {
			return child, key, true
		}
	}
	return nil, "", false
}

func (spec *ChildSpec) supports(version vast.Version) bool {
	for _, v := range spec.Versions {
		if v == version {
			return true
		}
	}
	return false
}

func (spec *AttributeSpec) supports(version vast.Version) bool {
	for _, v := range spec.Versions {
		if v == version {
			return true
		}
	}
	return false
}

var (
	supported20Plus = []vast.Version{
		vast.Version20,
		vast.Version30,
		vast.Version40,
		vast.Version41,
		vast.Version42,
		vast.Version43,
	}
	supported30Plus = []vast.Version{
		vast.Version30,
		vast.Version40,
		vast.Version41,
		vast.Version42,
		vast.Version43,
	}
	supported40Plus = []vast.Version{
		vast.Version40,
		vast.Version41,
		vast.Version42,
		vast.Version43,
	}
	supported41Plus = []vast.Version{
		vast.Version41,
		vast.Version42,
		vast.Version43,
	}
	supported42Plus = []vast.Version{
		vast.Version42,
		vast.Version43,
	}

	vast42SchemaURL = "https://raw.githubusercontent.com/InteractiveAdvertisingBureau/vast/refs/heads/master/vast_4.2.xsd"
)

// defaultCatalog contains a subset of the IAB VAST specification, focused on the
// most common nodes used by this project. Additional nodes can be appended over
// time without changing the validator API.
var defaultCatalog = &Catalog{Nodes: map[string]*NodeSpec{
	"VAST": {
		Name:     "VAST",
		Versions: supported20Plus,
		Attributes: map[string]*AttributeSpec{
			"version": {Name: "version", Versions: supported20Plus, Required: true},
		},
		Children: map[string]*ChildSpec{
			"Ad":    {Name: "Ad", Versions: supported20Plus, Multiple: true},
			"Error": {Name: "Error", Versions: supported20Plus, Optional: true, Multiple: true},
		},
	},
	"Ad": {
		Name:     "Ad",
		Versions: supported20Plus,
		Attributes: map[string]*AttributeSpec{
			"id":            {Name: "id", Versions: supported20Plus},
			"sequence":      {Name: "sequence", Versions: supported30Plus, Value: &AttributeValueSpec{Type: AttributeTypePositiveInteger}},
			"conditionalAd": {Name: "conditionalAd", Versions: supported40Plus, Value: &AttributeValueSpec{Type: AttributeTypeBoolean}},
			"adType":        {Name: "adType", Versions: supported41Plus, Value: &AttributeValueSpec{Type: AttributeTypeToken, AllowedValues: []string{"video", "audio", "application", "hybrid"}}},
		},
		Children: map[string]*ChildSpec{
			"InLine":  {Name: "InLine", Versions: supported20Plus, Optional: true},
			"Wrapper": {Name: "Wrapper", Versions: supported20Plus, Optional: true},
		},
	},
	"InLine": {
		Name:     "InLine",
		Versions: supported20Plus,
		Children: map[string]*ChildSpec{
			"AdSystem":    {Name: "AdSystem", Versions: supported20Plus},
			"Error":       {Name: "Error", Versions: supported20Plus, Optional: true, Multiple: true},
			"Impression":  {Name: "Impression", Versions: supported20Plus, Multiple: true},
			"AdTitle":     {Name: "AdTitle", Versions: supported20Plus},
			"Description": {Name: "Description", Versions: supported20Plus, Optional: true},
			"Survey":      {Name: "Survey", Versions: supported20Plus, Optional: true},
			"Creatives":   {Name: "Creatives", Versions: supported20Plus},
			"Extensions":  {Name: "Extensions", Versions: supported20Plus, Optional: true, Multiple: true},

			"Advertiser":         {Name: "Advertiser", Versions: supported30Plus, Optional: true},
			"Pricing":            {Name: "Pricing", Versions: supported30Plus, Optional: true},
			"AdServingId":        {Name: "AdServingId", Versions: supported40Plus},
			"Category":           {Name: "Category", Versions: supported30Plus, Optional: true, Multiple: true},
			"ViewableImpression": {Name: "ViewableImpression", Versions: supported40Plus, Optional: true},
			"Expires":            {Name: "Expires", Versions: supported30Plus, Optional: true},
			"AdVerifications":    {Name: "AdVerifications", Versions: supported40Plus, Optional: true},
		},
	},
	"Wrapper": {
		Name:     "Wrapper",
		Versions: supported20Plus,
		Attributes: map[string]*AttributeSpec{
			"followAdditionalWrappers": {Name: "followAdditionalWrappers", Versions: supported40Plus, Value: &AttributeValueSpec{Type: AttributeTypeBoolean}},
			"allowMultipleAds":         {Name: "allowMultipleAds", Versions: supported40Plus, Value: &AttributeValueSpec{Type: AttributeTypeBoolean}},
			"fallbackOnNoAd":           {Name: "fallbackOnNoAd", Versions: supported40Plus, Value: &AttributeValueSpec{Type: AttributeTypeBoolean}},
		},
		Children: map[string]*ChildSpec{
			"AdSystem":     {Name: "AdSystem", Versions: supported20Plus},
			"VASTAdTagURI": {Name: "VASTAdTagURI", Versions: supported20Plus},
			"Extensions":   {Name: "Extensions", Versions: supported20Plus, Optional: true, Multiple: true},
			"Error":        {Name: "Error", Versions: supported20Plus, Optional: true, Multiple: true},
			"Impression":   {Name: "Impression", Versions: supported20Plus, Multiple: true},
			"Creatives":    {Name: "Creatives", Versions: supported20Plus, Optional: true, NodeOverride: "WrapperCreatives"},

			"Pricing":             {Name: "Pricing", Versions: supported40Plus, Optional: true},
			"ViewableImpression":  {Name: "ViewableImpression", Versions: supported40Plus, Optional: true},
			"AdVerifications":     {Name: "AdVerifications", Versions: supported40Plus, Optional: true},
			"BlockedAdCategories": {Name: "BlockedAdCategories", Versions: supported30Plus, Optional: true, Multiple: true},
		},
	},
	"AdSystem": {
		Name:     "AdSystem",
		Versions: supported20Plus,
		Attributes: map[string]*AttributeSpec{
			"version": {Name: "version", Versions: supported30Plus},
		},
	},
	"Error": {
		Name:       "Error",
		Versions:   supported20Plus,
		NeedsCDATA: true,
	},
	"Impression": {
		Name:       "Impression",
		Versions:   supported20Plus,
		NeedsCDATA: true,
		Attributes: map[string]*AttributeSpec{
			"id": {Name: "id", Versions: supported30Plus},
		},
	},
	"AdTitle": {
		Name:     "AdTitle",
		Versions: supported20Plus,
	},
	"AdServingId": {
		Name:     "AdServingId",
		Versions: supported40Plus,
	},
	"Advertiser": {
		Name:     "Advertiser",
		Versions: supported30Plus,
		Attributes: map[string]*AttributeSpec{
			"id": {Name: "id", Versions: supported41Plus},
		},
	},
	"Category": {
		Name:     "Category",
		Versions: supported30Plus,
		Attributes: map[string]*AttributeSpec{
			"authority": {Name: "authority", Versions: supported30Plus, Required: true, Value: &AttributeValueSpec{Type: AttributeTypeURI}},
		},
	},
	"BlockedAdCategories": {
		Name:     "BlockedAdCategories",
		Versions: supported30Plus,
		Attributes: map[string]*AttributeSpec{
			"authority": {Name: "authority", Versions: supported30Plus, Value: &AttributeValueSpec{Type: AttributeTypeURI}},
		},
	},
	"Description": {
		Name:       "Description",
		Versions:   supported20Plus,
		NeedsCDATA: true,
	},
	"Survey": {
		Name:     "Survey",
		Versions: supported20Plus,
		Attributes: map[string]*AttributeSpec{
			"type": {Name: "type", Versions: supported20Plus},
		},
	},
	"Extensions": {
		Name:     "Extensions",
		Versions: supported20Plus,
		Children: map[string]*ChildSpec{
			"Extension": {Name: "Extension", Versions: supported20Plus, Optional: true, Multiple: true},
		},
	},
	"Extension": {
		Name:                   "Extension",
		Versions:               supported20Plus,
		AllowUnknownChildren:   true,
		AllowUnknownAttributes: true,
		Attributes: map[string]*AttributeSpec{
			"type": {Name: "type", Versions: supported20Plus},
		},
	},
	"Pricing": {
		Name:       "Pricing",
		Versions:   supported30Plus,
		NeedsCDATA: true,
		Attributes: map[string]*AttributeSpec{
			"model":    {Name: "model", Versions: supported30Plus, Required: true, Value: &AttributeValueSpec{Type: AttributeTypeToken, AllowedValues: []string{"CPM", "CPC", "CPV", "CPA"}}},
			"currency": {Name: "currency", Versions: supported30Plus, Required: true, Value: &AttributeValueSpec{Type: AttributeTypeToken, Pattern: "^[A-Z]{3}$"}},
		},
	},
	"ViewableImpression": {
		Name:     "ViewableImpression",
		Versions: supported40Plus,
		Attributes: map[string]*AttributeSpec{
			"id": {Name: "id", Versions: supported42Plus},
		},
		Children: map[string]*ChildSpec{
			"Viewable":         {Name: "Viewable", Versions: supported40Plus, Optional: true, Multiple: true},
			"NotViewable":      {Name: "NotViewable", Versions: supported40Plus, Optional: true, Multiple: true},
			"ViewUndetermined": {Name: "ViewUndetermined", Versions: supported40Plus, Optional: true, Multiple: true},
		},
	},
	"Viewable": {
		Name:       "Viewable",
		Versions:   supported40Plus,
		NeedsCDATA: true,
	},
	"NotViewable": {
		Name:       "NotViewable",
		Versions:   supported40Plus,
		NeedsCDATA: true,
	},
	"ViewUndetermined": {
		Name:       "ViewUndetermined",
		Versions:   supported40Plus,
		NeedsCDATA: true,
	},
	"Expires": {
		Name:     "Expires",
		Versions: supported30Plus,
	},
	"UniversalAdId": {
		Name:               "UniversalAdId",
		Versions:           supported40Plus,
		SupportsExtensions: true,
		Attributes: map[string]*AttributeSpec{
			"idRegistry": {Name: "idRegistry", Versions: supported40Plus, Required: true},
			"idValue":    {Name: "idValue", Versions: supported40Plus, Required: true},
		},
	},
	"AdVerifications": {
		Name:     "AdVerifications",
		Versions: supported40Plus,
		Children: map[string]*ChildSpec{
			"Verification": {Name: "Verification", Versions: supported40Plus, Multiple: true},
		},
	},
	"Verification": {
		Name:     "Verification",
		Versions: supported40Plus,
		Attributes: map[string]*AttributeSpec{
			"vendor": {Name: "vendor", Versions: supported30Plus},
		},
		Children: map[string]*ChildSpec{
			"JavaScriptResource":     {Name: "JavaScriptResource", Versions: supported40Plus, Optional: true, Multiple: true},
			"ExecutableResource":     {Name: "ExecutableResource", Versions: supported40Plus, Optional: true, Multiple: true},
			"TrackingEvents":         {Name: "TrackingEvents", Versions: supported40Plus, Optional: true},
			"VerificationParameters": {Name: "VerificationParameters", Versions: supported40Plus, Optional: true},
			"BlockedAdCategories":    {Name: "BlockedAdCategories", Versions: supported41Plus, Optional: true, Multiple: true},
		},
	},
	"JavaScriptResource": {
		Name:       "JavaScriptResource",
		Versions:   supported40Plus,
		NeedsCDATA: true,
		Attributes: map[string]*AttributeSpec{
			"apiFramework":    {Name: "apiFramework", Versions: supported30Plus},
			"browserOptional": {Name: "browserOptional", Versions: supported30Plus},
		},
	},
	"ExecutableResource": {
		Name:       "ExecutableResource",
		Versions:   supported40Plus,
		NeedsCDATA: true,
		Attributes: map[string]*AttributeSpec{
			"apiFramework": {Name: "apiFramework", Versions: supported30Plus},
			"type":         {Name: "type", Versions: supported30Plus},
			"language":     {Name: "language", Versions: supported41Plus},
		},
	},
	"VerificationParameters": {
		Name:     "VerificationParameters",
		Versions: supported40Plus,
	},
	"VASTAdTagURI": {
		Name:       "VASTAdTagURI",
		Versions:   supported20Plus,
		NeedsCDATA: true,
	},
	"Creatives": {
		Name:     "Creatives",
		Versions: supported20Plus,
		Children: map[string]*ChildSpec{
			"Creative": {Name: "Creative", Versions: supported20Plus, Multiple: true},
		},
	},
	"WrapperCreatives": {
		Name:     "Creatives",
		Versions: supported20Plus,
		Children: map[string]*ChildSpec{
			"Creative": {Name: "Creative", Versions: supported20Plus, Multiple: true, NodeOverride: "WrapperCreative"},
		},
	},
	"Creative": {
		Name:     "Creative",
		Versions: supported20Plus,
		Attributes: map[string]*AttributeSpec{
			"id":           {Name: "id", Versions: supported20Plus},
			"sequence":     {Name: "sequence", Versions: supported20Plus, Value: &AttributeValueSpec{Type: AttributeTypePositiveInteger}},
			"apiFramework": {Name: "apiFramework", Versions: supported20Plus},
			"adId":         {Name: "adId", Versions: supported20Plus},
		},
		Children: map[string]*ChildSpec{
			"Linear":             {Name: "Linear", Versions: supported20Plus, Optional: true},
			"NonLinearAds":       {Name: "NonLinearAds", Versions: supported20Plus, Optional: true},
			"CompanionAds":       {Name: "CompanionAds", Versions: supported20Plus, Optional: true},
			"CreativeExtensions": {Name: "CreativeExtensions", Versions: supported30Plus, Optional: true, Multiple: true},
			"UniversalAdId":      {Name: "UniversalAdId", Versions: supported40Plus, Optional: true, Multiple: true},
		},
	},
	"WrapperCreative": {
		Name:     "Creative",
		Versions: supported20Plus,
		Attributes: map[string]*AttributeSpec{
			"id":           {Name: "id", Versions: supported20Plus},
			"sequence":     {Name: "sequence", Versions: supported20Plus, Value: &AttributeValueSpec{Type: AttributeTypePositiveInteger}},
			"apiFramework": {Name: "apiFramework", Versions: supported20Plus},
			"adId":         {Name: "adId", Versions: supported20Plus},
		},
		Children: map[string]*ChildSpec{
			"Linear":             {Name: "Linear", Versions: supported20Plus, Optional: true, NodeOverride: "WrapperLinear"},
			"NonLinearAds":       {Name: "NonLinearAds", Versions: supported20Plus, Optional: true},
			"CompanionAds":       {Name: "CompanionAds", Versions: supported20Plus, Optional: true},
			"CreativeExtensions": {Name: "CreativeExtensions", Versions: supported30Plus, Optional: true, Multiple: true},
			"UniversalAdId":      {Name: "UniversalAdId", Versions: supported40Plus, Optional: true, Multiple: true},
		},
	},
	"NonLinearAds": {
		Name:     "NonLinearAds",
		Versions: supported20Plus,
		Children: map[string]*ChildSpec{
			"NonLinear":      {Name: "NonLinear", Versions: supported20Plus, Multiple: true},
			"TrackingEvents": {Name: "TrackingEvents", Versions: supported20Plus, Optional: true},
		},
	},
	"NonLinear": {
		Name:     "NonLinear",
		Versions: supported20Plus,
		Attributes: map[string]*AttributeSpec{
			"id":                   {Name: "id", Versions: supported20Plus},
			"width":                {Name: "width", Versions: supported20Plus, Required: true, Value: &AttributeValueSpec{Type: AttributeTypePositiveInteger}},
			"height":               {Name: "height", Versions: supported20Plus, Required: true, Value: &AttributeValueSpec{Type: AttributeTypePositiveInteger}},
			"expandedWidth":        {Name: "expandedWidth", Versions: supported30Plus, Value: &AttributeValueSpec{Type: AttributeTypePositiveInteger}},
			"expandedHeight":       {Name: "expandedHeight", Versions: supported30Plus, Value: &AttributeValueSpec{Type: AttributeTypePositiveInteger}},
			"scalable":             {Name: "scalable", Versions: supported20Plus, Value: &AttributeValueSpec{Type: AttributeTypeBoolean}},
			"maintainAspectRatio":  {Name: "maintainAspectRatio", Versions: supported20Plus, Value: &AttributeValueSpec{Type: AttributeTypeBoolean}},
			"minSuggestedDuration": {Name: "minSuggestedDuration", Versions: supported20Plus, Value: &AttributeValueSpec{Type: AttributeTypeDuration}},
			"apiFramework":         {Name: "apiFramework", Versions: supported20Plus},
		},
		Children: map[string]*ChildSpec{
			"StaticResource":         {Name: "StaticResource", Versions: supported20Plus, Optional: true},
			"IFrameResource":         {Name: "IFrameResource", Versions: supported20Plus, Optional: true},
			"HTMLResource":           {Name: "HTMLResource", Versions: supported20Plus, Optional: true},
			"AdParameters":           {Name: "AdParameters", Versions: supported20Plus, Optional: true},
			"NonLinearClickTracking": {Name: "NonLinearClickTracking", Versions: supported30Plus, Optional: true, Multiple: true},
			"NonLinearClickThrough":  {Name: "NonLinearClickThrough", Versions: supported20Plus, Optional: true},
		},
	},
	"NonLinearClickTracking": {
		Name:       "NonLinearClickTracking",
		Versions:   supported30Plus,
		NeedsCDATA: true,
		Attributes: map[string]*AttributeSpec{
			"id": {Name: "id", Versions: supported30Plus},
		},
	},
	"NonLinearClickThrough": {
		Name:       "NonLinearClickThrough",
		Versions:   supported20Plus,
		NeedsCDATA: true,
	},
	"CompanionAds": {
		Name:     "CompanionAds",
		Versions: supported20Plus,
		Attributes: map[string]*AttributeSpec{
			"required": {Name: "required", Versions: supported30Plus, Value: &AttributeValueSpec{Type: AttributeTypeToken, AllowedValues: []string{"all", "any", "none"}}},
		},
		Children: map[string]*ChildSpec{
			"Companion": {Name: "Companion", Versions: supported20Plus, Multiple: true},
		},
	},
	"Companion": {
		Name:     "Companion",
		Versions: supported20Plus,
		Attributes: map[string]*AttributeSpec{
			"id":             {Name: "id", Versions: supported20Plus},
			"width":          {Name: "width", Versions: supported20Plus, Required: true, Value: &AttributeValueSpec{Type: AttributeTypePositiveInteger}},
			"height":         {Name: "height", Versions: supported20Plus, Required: true, Value: &AttributeValueSpec{Type: AttributeTypePositiveInteger}},
			"apiFramework":   {Name: "apiFramework", Versions: supported20Plus},
			"assetWidth":     {Name: "assetWidth", Versions: supported30Plus, Value: &AttributeValueSpec{Type: AttributeTypePositiveInteger}},
			"assetHeight":    {Name: "assetHeight", Versions: supported30Plus, Value: &AttributeValueSpec{Type: AttributeTypePositiveInteger}},
			"expandedWidth":  {Name: "expandedWidth", Versions: supported30Plus, Value: &AttributeValueSpec{Type: AttributeTypePositiveInteger}},
			"expandedHeight": {Name: "expandedHeight", Versions: supported30Plus, Value: &AttributeValueSpec{Type: AttributeTypePositiveInteger}},
			"adSlotId":       {Name: "adSlotId", Versions: supported30Plus},
			"pxratio":        {Name: "pxratio", Versions: supported30Plus, Value: &AttributeValueSpec{Type: AttributeTypeFloat}},
			"renderingMode":  {Name: "renderingMode", Versions: supported30Plus, Value: &AttributeValueSpec{Type: AttributeTypeToken, AllowedValues: []string{"default", "end-card", "concurrent"}}},
		},
		Children: map[string]*ChildSpec{
			"StaticResource":         {Name: "StaticResource", Versions: supported20Plus, Optional: true},
			"IFrameResource":         {Name: "IFrameResource", Versions: supported20Plus, Optional: true},
			"HTMLResource":           {Name: "HTMLResource", Versions: supported20Plus, Optional: true},
			"AdParameters":           {Name: "AdParameters", Versions: supported20Plus, Optional: true},
			"AltText":                {Name: "AltText", Versions: supported20Plus, Optional: true},
			"CompanionClickThrough":  {Name: "CompanionClickThrough", Versions: supported20Plus, Optional: true},
			"CompanionClickTracking": {Name: "CompanionClickTracking", Versions: supported30Plus, Optional: true, Multiple: true},
			"CreativeExtensions":     {Name: "CreativeExtensions", Versions: supported30Plus, Optional: true, Multiple: true},
			"TrackingEvents":         {Name: "TrackingEvents", Versions: supported20Plus, Optional: true},
		},
	},
	"CompanionClickThrough": {
		Name:       "CompanionClickThrough",
		Versions:   supported20Plus,
		NeedsCDATA: true,
	},
	"CompanionClickTracking": {
		Name:     "CompanionClickTracking",
		Versions: supported30Plus,
		Attributes: map[string]*AttributeSpec{
			"id": {Name: "id", Versions: supported30Plus, Required: true},
		},
	},
	"AltText": {
		Name:     "AltText",
		Versions: supported20Plus,
	},
	"Icons": {
		Name:     "Icons",
		Versions: supported30Plus,
		Children: map[string]*ChildSpec{
			"Icon": {Name: "Icon", Versions: supported30Plus, Multiple: true},
		},
	},
	"Icon": {
		Name:     "Icon",
		Versions: supported30Plus,
		Attributes: map[string]*AttributeSpec{
			"program":      {Name: "program", Versions: supported30Plus},
			"width":        {Name: "width", Versions: supported30Plus, Value: &AttributeValueSpec{Type: AttributeTypePositiveInteger}},
			"height":       {Name: "height", Versions: supported30Plus, Value: &AttributeValueSpec{Type: AttributeTypePositiveInteger}},
			"xPosition":    {Name: "xPosition", Versions: supported30Plus, Value: &AttributeValueSpec{Type: AttributeTypeToken, Pattern: "([0-9]*|left|right)"}},
			"yPosition":    {Name: "yPosition", Versions: supported30Plus, Value: &AttributeValueSpec{Type: AttributeTypeToken, Pattern: "([0-9]*|top|bottom)"}},
			"duration":     {Name: "duration", Versions: supported30Plus, Value: &AttributeValueSpec{Type: AttributeTypeDuration}},
			"offset":       {Name: "offset", Versions: supported30Plus, Value: &AttributeValueSpec{Type: AttributeTypeDuration}},
			"apiFramework": {Name: "apiFramework", Versions: supported30Plus},
			"pxratio":      {Name: "pxratio", Versions: supported30Plus, Value: &AttributeValueSpec{Type: AttributeTypeFloat}},
		},
		Children: map[string]*ChildSpec{
			"StaticResource":   {Name: "StaticResource", Versions: supported30Plus, Optional: true},
			"IFrameResource":   {Name: "IFrameResource", Versions: supported30Plus, Optional: true},
			"HTMLResource":     {Name: "HTMLResource", Versions: supported30Plus, Optional: true},
			"IconClicks":       {Name: "IconClicks", Versions: supported30Plus, Optional: true},
			"IconViewTracking": {Name: "IconViewTracking", Versions: supported30Plus, Optional: true, Multiple: true},
		},
	},
	"IconClicks": {
		Name:     "IconClicks",
		Versions: supported30Plus,
		Children: map[string]*ChildSpec{
			"IconClickFallbackImages": {Name: "IconClickFallbackImages", Versions: supported42Plus, Optional: true},
			"IconClickThrough":        {Name: "IconClickThrough", Versions: supported30Plus, Optional: true},
			"IconClickTracking":       {Name: "IconClickTracking", Versions: supported30Plus, Optional: true, Multiple: true},
		},
	},
	"IconClickThrough": {
		Name:     "IconClickThrough",
		Versions: supported30Plus,
	},
	"IconClickTracking": {
		Name:     "IconClickTracking",
		Versions: supported30Plus,
		Attributes: map[string]*AttributeSpec{
			"id": {Name: "id", Versions: supported30Plus},
		},
	},
	"IconClickFallbackImages": {
		Name:     "IconClickFallbackImages",
		Versions: supported42Plus,
		Children: map[string]*ChildSpec{
			"IconClickFallbackImage": {Name: "IconClickFallbackImage", Versions: supported42Plus, Multiple: true},
		},
	},
	"IconClickFallbackImage": {
		Name:     "IconClickFallbackImage",
		Versions: supported42Plus,
		Attributes: map[string]*AttributeSpec{
			"width":  {Name: "width", Versions: supported42Plus, Value: &AttributeValueSpec{Type: AttributeTypePositiveInteger}},
			"height": {Name: "height", Versions: supported42Plus, Value: &AttributeValueSpec{Type: AttributeTypePositiveInteger}},
		},
		Children: map[string]*ChildSpec{
			"AltText":        {Name: "AltText", Versions: supported42Plus, Optional: true},
			"StaticResource": {Name: "StaticResource", Versions: supported42Plus, Optional: true},
		},
	},
	"IconViewTracking": {
		Name:     "IconViewTracking",
		Versions: supported30Plus,
	},
	"CreativeExtensions": {
		Name:     "CreativeExtensions",
		Versions: supported30Plus,
		Children: map[string]*ChildSpec{
			"CreativeExtension": {Name: "CreativeExtension", Versions: supported30Plus, Multiple: true},
		},
	},
	"CreativeExtension": {
		Name:                   "CreativeExtension",
		Versions:               supported30Plus,
		AllowUnknownChildren:   true,
		AllowUnknownAttributes: true,
		Attributes: map[string]*AttributeSpec{
			"type": {Name: "type", Versions: supported30Plus},
		},
	},
	"StaticResource": {
		Name:       "StaticResource",
		Versions:   supported20Plus,
		NeedsCDATA: true,
		Attributes: map[string]*AttributeSpec{
			"creativeType": {Name: "creativeType", Versions: supported20Plus},
		},
	},
	"HTMLResource": {
		Name:       "HTMLResource",
		Versions:   supported20Plus,
		NeedsCDATA: true,
	},
	"IFrameResource": {
		Name:       "IFrameResource",
		Versions:   supported20Plus,
		NeedsCDATA: true,
	},
	"Linear": {
		Name:     "Linear",
		Versions: supported20Plus,
		Attributes: map[string]*AttributeSpec{
			"skipoffset": {Name: "skipoffset", Versions: supported30Plus, Value: &AttributeValueSpec{Type: AttributeTypeTimeOffset}},
		},
		Children: map[string]*ChildSpec{
			"Icons":          {Name: "Icons", Versions: supported30Plus, Optional: true},
			"AdParameters":   {Name: "AdParameters", Versions: supported20Plus, Optional: true},
			"Duration":       {Name: "Duration", Versions: supported20Plus},
			"MediaFiles":     {Name: "MediaFiles", Versions: supported20Plus},
			"VideoClicks":    {Name: "VideoClicks", Versions: supported20Plus, Optional: true},
			"TrackingEvents": {Name: "TrackingEvents", Versions: supported20Plus, Optional: true},
		},
	},
	"WrapperLinear": {
		Name:     "Linear",
		Versions: supported20Plus,
		Attributes: map[string]*AttributeSpec{
			"skipoffset": {Name: "skipoffset", Versions: supported30Plus, Value: &AttributeValueSpec{Type: AttributeTypeTimeOffset}},
		},
		Children: map[string]*ChildSpec{
			"Icons":          {Name: "Icons", Versions: supported30Plus, Optional: true},
			"AdParameters":   {Name: "AdParameters", Versions: supported20Plus, Optional: true},
			"Duration":       {Name: "Duration", Versions: supported20Plus, Optional: true},
			"MediaFiles":     {Name: "MediaFiles", Versions: supported20Plus, Optional: true},
			"VideoClicks":    {Name: "VideoClicks", Versions: supported20Plus, Optional: true},
			"TrackingEvents": {Name: "TrackingEvents", Versions: supported20Plus, Optional: true},
		},
	},
	"AdParameters": {
		Name:     "AdParameters",
		Versions: supported20Plus,
		Attributes: map[string]*AttributeSpec{
			"xmlEncoded": {Name: "xmlEncoded", Versions: supported30Plus, Value: &AttributeValueSpec{Type: AttributeTypeBoolean}},
		},
		NeedsCDATA: true,
	},
	"Duration": {
		Name:     "Duration",
		Versions: supported20Plus,
	},
	"MediaFiles": {
		Name:     "MediaFiles",
		Versions: supported20Plus,
		Children: map[string]*ChildSpec{
			"MediaFile":               {Name: "MediaFile", Versions: supported20Plus, Multiple: true},
			"ClosedCaptionFiles":      {Name: "ClosedCaptionFiles", Versions: supported30Plus, Optional: true},
			"Mezzanine":               {Name: "Mezzanine", Versions: supported40Plus, Optional: true, Multiple: true},
			"InteractiveCreativeFile": {Name: "InteractiveCreativeFile", Versions: supported30Plus, Optional: true, Multiple: true},
		},
	},
	"MediaFile": {
		Name:       "MediaFile",
		Versions:   supported20Plus,
		NeedsCDATA: true,
		Attributes: map[string]*AttributeSpec{
			"id":                  {Name: "id", Versions: supported20Plus},
			"delivery":            {Name: "delivery", Versions: supported20Plus, Required: true, Value: &AttributeValueSpec{Type: AttributeTypeToken, AllowedValues: []string{"progressive", "streaming"}}},
			"type":                {Name: "type", Versions: supported20Plus, Required: true},
			"width":               {Name: "width", Versions: supported20Plus, Required: true, Value: &AttributeValueSpec{Type: AttributeTypeNonNegativeInteger}},
			"height":              {Name: "height", Versions: supported20Plus, Required: true, Value: &AttributeValueSpec{Type: AttributeTypeNonNegativeInteger}},
			"codec":               {Name: "codec", Versions: supported30Plus},
			"bitrate":             {Name: "bitrate", Versions: supported30Plus, Value: &AttributeValueSpec{Type: AttributeTypeNonNegativeInteger}},
			"minBitrate":          {Name: "minBitrate", Versions: supported30Plus, Value: &AttributeValueSpec{Type: AttributeTypeNonNegativeInteger}},
			"maxBitrate":          {Name: "maxBitrate", Versions: supported30Plus, Value: &AttributeValueSpec{Type: AttributeTypeNonNegativeInteger}},
			"scalable":            {Name: "scalable", Versions: supported20Plus, Value: &AttributeValueSpec{Type: AttributeTypeBoolean}},
			"maintainAspectRatio": {Name: "maintainAspectRatio", Versions: supported20Plus, Value: &AttributeValueSpec{Type: AttributeTypeBoolean}},
			"fileSize":            {Name: "fileSize", Versions: supported30Plus, Value: &AttributeValueSpec{Type: AttributeTypeNonNegativeInteger}},
			"mediaType":           {Name: "mediaType", Versions: supported30Plus, Value: &AttributeValueSpec{Type: AttributeTypeToken, AllowedValues: []string{"2D", "3D"}}},
			"apiFramework":        {Name: "apiFramework", Versions: supported20Plus},
		},
	},
	"ClosedCaptionFiles": {
		Name:     "ClosedCaptionFiles",
		Versions: supported30Plus,
		Children: map[string]*ChildSpec{
			"ClosedCaptionFile": {Name: "ClosedCaptionFile", Versions: supported30Plus, Multiple: true},
		},
	},
	"ClosedCaptionFile": {
		Name:       "ClosedCaptionFile",
		Versions:   supported30Plus,
		NeedsCDATA: true,
		Attributes: map[string]*AttributeSpec{
			"type":     {Name: "type", Versions: supported30Plus},
			"language": {Name: "language", Versions: supported30Plus},
		},
	},
	"Mezzanine": {
		Name:               "Mezzanine",
		Versions:           supported40Plus,
		SupportsExtensions: true,
		NeedsCDATA:         true,
		Attributes: map[string]*AttributeSpec{
			"id":                  {Name: "id", Versions: supported40Plus},
			"delivery":            {Name: "delivery", Versions: supported40Plus, Required: true, Value: &AttributeValueSpec{Type: AttributeTypeToken, AllowedValues: []string{"progressive", "streaming"}}},
			"type":                {Name: "type", Versions: supported40Plus, Required: true},
			"bitrate":             {Name: "bitrate", Versions: supported40Plus},
			"minBitrate":          {Name: "minBitrate", Versions: supported40Plus},
			"maxBitrate":          {Name: "maxBitrate", Versions: supported40Plus},
			"width":               {Name: "width", Versions: supported40Plus, Required: true, Value: &AttributeValueSpec{Type: AttributeTypeNonNegativeInteger}},
			"height":              {Name: "height", Versions: supported40Plus, Required: true, Value: &AttributeValueSpec{Type: AttributeTypeNonNegativeInteger}},
			"scalable":            {Name: "scalable", Versions: supported40Plus},
			"maintainAspectRatio": {Name: "maintainAspectRatio", Versions: supported40Plus},
			"codec":               {Name: "codec", Versions: supported40Plus},
			"fileSize":            {Name: "fileSize", Versions: supported41Plus, Value: &AttributeValueSpec{Type: AttributeTypeNonNegativeInteger}},
			"mediaType":           {Name: "mediaType", Versions: supported40Plus},
		},
	},
	"InteractiveCreativeFile": {
		Name:               "InteractiveCreativeFile",
		Versions:           supported30Plus,
		SupportsExtensions: true,
		NeedsCDATA:         true,
		Attributes: map[string]*AttributeSpec{
			"type":             {Name: "type", Versions: supported30Plus},
			"apiFramework":     {Name: "apiFramework", Versions: supported30Plus},
			"variableDuration": {Name: "variableDuration", Versions: supported30Plus, Value: &AttributeValueSpec{Type: AttributeTypeBoolean}},
		},
	},
	"VideoClicks": {
		Name:     "VideoClicks",
		Versions: supported20Plus,
		Children: map[string]*ChildSpec{
			"ClickThrough":  {Name: "ClickThrough", Versions: supported20Plus, Optional: true},
			"ClickTracking": {Name: "ClickTracking", Versions: supported20Plus, Optional: true, Multiple: true},
			"CustomClick":   {Name: "CustomClick", Versions: supported30Plus, Optional: true, Multiple: true},
		},
	},
	"ClickThrough": {
		Name:       "ClickThrough",
		Versions:   supported20Plus,
		NeedsCDATA: true,
		Attributes: map[string]*AttributeSpec{
			"id": {Name: "id", Versions: supported30Plus},
		},
	},
	"ClickTracking": {
		Name:       "ClickTracking",
		Versions:   supported20Plus,
		NeedsCDATA: true,
		Attributes: map[string]*AttributeSpec{
			"id": {Name: "id", Versions: supported30Plus},
		},
	},
	"CustomClick": {
		Name:     "CustomClick",
		Versions: supported30Plus,
		Attributes: map[string]*AttributeSpec{
			"id": {Name: "id", Versions: supported30Plus},
		},
	},
	"TrackingEvents": {
		Name:     "TrackingEvents",
		Versions: supported20Plus,
		Children: map[string]*ChildSpec{
			"Tracking": {Name: "Tracking", Versions: supported20Plus, Multiple: true},
		},
	},
	"Tracking": {
		Name:       "Tracking",
		Versions:   supported20Plus,
		NeedsCDATA: true,
		Attributes: map[string]*AttributeSpec{
			"event": {
				Name:     "event",
				Versions: supported20Plus,
				Required: true,
				Value: &AttributeValueSpec{
					Type: AttributeTypeToken,
					AllowedValues: []string{
						string(vast.MuteEvent),
						string(vast.UnmuteEvent),
						string(vast.PauseEvent),
						string(vast.ResumeEvent),
						string(vast.RewindEvent),
						string(vast.SkipEvent),
						string(vast.PlayerExpandEvent),
						string(vast.PlayerCollapseEvent),
						string(vast.LoadedEvent),
						string(vast.StartEvent),
						string(vast.FirstQuartileEvent),
						string(vast.MidpointEvent),
						string(vast.ThirdQuartileEvent),
						string(vast.CompleteEvent),
						string(vast.ProgressEvent),
						string(vast.CloseLinearEvent),
						string(vast.CreativeViewEvent),
						string(vast.AcceptInvitationEvent),
						string(vast.AdExpandEvent),
						string(vast.AdCollapseEvent),
						string(vast.MinimizeEvent),
						string(vast.CloseEvent),
						string(vast.OverlayViewDurationEvent),
						string(vast.OtherAdInteraction),
						string(vast.InteractiveStart),
					},
				},
			},
			"offset": {Name: "offset", Versions: supported30Plus, Value: &AttributeValueSpec{Type: AttributeTypeTimeOffset}},
		},
	},
}}

func init() {
	annotateCatalogDocs(defaultCatalog)
}

func docIsEmpty(doc *Documentation) bool {
	return doc == nil || strings.TrimSpace(doc.Content) == ""
}

func schemaDocumentation(text string) *Documentation {
	trimmed := strings.TrimSpace(text)
	if trimmed == "" {
		return nil
	}
	return &Documentation{Content: trimmed, Source: vast42SchemaURL}
}

func annotateCatalogDocs(cat *Catalog) {
	if cat == nil {
		return
	}
	for _, node := range cat.Nodes {
		annotateNodeDoc(node)
	}
}

func annotateNodeDoc(node *NodeSpec) {
	if node == nil {
		return
	}
	if docIsEmpty(node.Documentation) {
		if doc := schemaDocumentation(vast42ElementDocs[node.Name]); doc != nil {
			node.Documentation = doc
		}
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
	for _, child := range node.Children {
		annotateChildDoc(node.Name, child)
	}
}

func annotateAttributeDoc(nodeName string, attr *AttributeSpec) {
	if attr == nil {
		return
	}
	if docIsEmpty(attr.Documentation) {
		if scoped, ok := vast42AttributeDocs[nodeName]; ok {
			if doc := schemaDocumentation(scoped[attr.Name]); doc != nil {
				attr.Documentation = doc
			}
		}
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

func annotateChildDoc(parent string, child *ChildSpec) {
	if child == nil {
		return
	}
	if docIsEmpty(child.Documentation) {
		if doc := schemaDocumentation(vast42ElementDocs[child.Name]); doc != nil {
			child.Documentation = doc
		}
	}
	if docIsEmpty(child.Documentation) {
		child.Documentation = &Documentation{
			Content: fmt.Sprintf("Child <%s> permitted within <%s> per VAST 4.2 XSD.", child.Name, parent),
			Source:  vast42SchemaURL,
		}
	}
}
