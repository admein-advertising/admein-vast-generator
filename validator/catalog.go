package validator

import "github.com/admein-advertising/admein-vast-generator/vast"

// AttributeSpec describes a valid attribute for a node.
type AttributeSpec struct {
	Name       string
	Versions   []vast.Version
	Required   bool
	AllowEmpty bool
}

// ChildSpec describes a valid child node relationship.
type ChildSpec struct {
	Name     string
	Versions []vast.Version
	Optional bool
	Multiple bool
}

// NodeSpec defines the validation metadata for a node.
type NodeSpec struct {
	Name                 string
	Versions             []vast.Version
	Attributes           map[string]*AttributeSpec
	Children             map[string]*ChildSpec
	AllowUnknownChildren bool
}

// Catalog stores node specifications keyed by node name.
type Catalog struct {
	Nodes map[string]*NodeSpec
}

func (c *Catalog) node(name string) (*NodeSpec, bool) {
	if c == nil {
		return nil, false
	}
	spec, ok := c.Nodes[name]
	return spec, ok
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

func (spec *NodeSpec) child(name string) (*ChildSpec, bool) {
	if spec == nil {
		return nil, false
	}
	child, ok := spec.Children[name]
	return child, ok
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
	allSupportedVersions = []vast.Version{
		vast.Version30,
		vast.Version40,
		vast.Version41,
		vast.Version42,
		vast.Version43,
	}
	version40Plus = []vast.Version{
		vast.Version40,
		vast.Version41,
		vast.Version42,
		vast.Version43,
	}
)

// defaultCatalog contains a subset of the IAB VAST specification, focused on the
// most common nodes used by this project. Additional nodes can be appended over
// time without changing the validator API.
var defaultCatalog = &Catalog{Nodes: map[string]*NodeSpec{
	"VAST": {
		Name:     "VAST",
		Versions: allSupportedVersions,
		Attributes: map[string]*AttributeSpec{
			"version": {Name: "version", Versions: allSupportedVersions, Required: true},
		},
		Children: map[string]*ChildSpec{
			"Ad":    {Name: "Ad", Versions: allSupportedVersions, Multiple: true},
			"Error": {Name: "Error", Versions: allSupportedVersions, Optional: true, Multiple: true},
		},
	},
	"Ad": {
		Name:     "Ad",
		Versions: allSupportedVersions,
		Attributes: map[string]*AttributeSpec{
			"id":            {Name: "id", Versions: allSupportedVersions},
			"sequence":      {Name: "sequence", Versions: allSupportedVersions},
			"conditionalAd": {Name: "conditionalAd", Versions: allSupportedVersions},
			"adType":        {Name: "adType", Versions: allSupportedVersions},
		},
		Children: map[string]*ChildSpec{
			"InLine":  {Name: "InLine", Versions: allSupportedVersions, Optional: true},
			"Wrapper": {Name: "Wrapper", Versions: allSupportedVersions, Optional: true},
		},
	},
	"InLine": {
		Name:     "InLine",
		Versions: allSupportedVersions,
		Children: map[string]*ChildSpec{
			"AdSystem":           {Name: "AdSystem", Versions: allSupportedVersions},
			"Error":              {Name: "Error", Versions: allSupportedVersions, Optional: true, Multiple: true},
			"Impression":         {Name: "Impression", Versions: allSupportedVersions, Multiple: true},
			"AdTitle":            {Name: "AdTitle", Versions: allSupportedVersions},
			"AdServingId":        {Name: "AdServingId", Versions: allSupportedVersions, Optional: true},
			"Advertiser":         {Name: "Advertiser", Versions: allSupportedVersions, Optional: true},
			"Category":           {Name: "Category", Versions: allSupportedVersions, Optional: true, Multiple: true},
			"Description":        {Name: "Description", Versions: allSupportedVersions, Optional: true},
			"Extensions":         {Name: "Extensions", Versions: allSupportedVersions, Optional: true},
			"Pricing":            {Name: "Pricing", Versions: allSupportedVersions, Optional: true},
			"ViewableImpression": {Name: "ViewableImpression", Versions: version40Plus, Optional: true},
			"Survey":             {Name: "Survey", Versions: allSupportedVersions, Optional: true},
			"Expires":            {Name: "Expires", Versions: allSupportedVersions, Optional: true},
			"Creatives":          {Name: "Creatives", Versions: allSupportedVersions},
			"AdVerifications":    {Name: "AdVerifications", Versions: version40Plus, Optional: true},
		},
	},
	"Wrapper": {
		Name:     "Wrapper",
		Versions: allSupportedVersions,
		Attributes: map[string]*AttributeSpec{
			"followAdditionalWrappers": {Name: "followAdditionalWrappers", Versions: allSupportedVersions},
			"allowMultipleAds":         {Name: "allowMultipleAds", Versions: allSupportedVersions},
			"fallbackOnNoAd":           {Name: "fallbackOnNoAd", Versions: allSupportedVersions},
		},
		Children: map[string]*ChildSpec{
			"AdSystem":            {Name: "AdSystem", Versions: allSupportedVersions},
			"Error":               {Name: "Error", Versions: allSupportedVersions, Optional: true, Multiple: true},
			"Impression":          {Name: "Impression", Versions: allSupportedVersions, Multiple: true},
			"VASTAdTagURI":        {Name: "VASTAdTagURI", Versions: allSupportedVersions},
			"Extensions":          {Name: "Extensions", Versions: allSupportedVersions, Optional: true},
			"Pricing":             {Name: "Pricing", Versions: allSupportedVersions, Optional: true},
			"ViewableImpression":  {Name: "ViewableImpression", Versions: version40Plus, Optional: true},
			"Creatives":           {Name: "Creatives", Versions: allSupportedVersions, Optional: true},
			"BlockedAdCategories": {Name: "BlockedAdCategories", Versions: allSupportedVersions, Optional: true, Multiple: true},
			"AdVerifications":     {Name: "AdVerifications", Versions: version40Plus, Optional: true},
		},
	},
	"AdSystem": {
		Name:     "AdSystem",
		Versions: allSupportedVersions,
		Attributes: map[string]*AttributeSpec{
			"version": {Name: "version", Versions: allSupportedVersions},
		},
	},
	"Error": {
		Name:     "Error",
		Versions: allSupportedVersions,
	},
	"Impression": {
		Name:     "Impression",
		Versions: allSupportedVersions,
		Attributes: map[string]*AttributeSpec{
			"id": {Name: "id", Versions: allSupportedVersions},
		},
	},
	"AdTitle": {
		Name:     "AdTitle",
		Versions: allSupportedVersions,
	},
	"AdServingId": {
		Name:     "AdServingId",
		Versions: allSupportedVersions,
	},
	"Advertiser": {
		Name:     "Advertiser",
		Versions: allSupportedVersions,
	},
	"Category": {
		Name:     "Category",
		Versions: allSupportedVersions,
		Attributes: map[string]*AttributeSpec{
			"authority": {Name: "authority", Versions: allSupportedVersions, Required: true},
		},
	},
	"BlockedAdCategories": {
		Name:     "BlockedAdCategories",
		Versions: allSupportedVersions,
		Attributes: map[string]*AttributeSpec{
			"authority": {Name: "authority", Versions: allSupportedVersions},
		},
	},
	"Description": {
		Name:     "Description",
		Versions: allSupportedVersions,
	},
	"Survey": {
		Name:     "Survey",
		Versions: allSupportedVersions,
		Attributes: map[string]*AttributeSpec{
			"type": {Name: "type", Versions: allSupportedVersions},
		},
	},
	"Extensions": {
		Name:     "Extensions",
		Versions: allSupportedVersions,
		Children: map[string]*ChildSpec{
			"Extension": {Name: "Extension", Versions: allSupportedVersions, Optional: true, Multiple: true},
		},
	},
	"Extension": {
		Name:                 "Extension",
		Versions:             allSupportedVersions,
		AllowUnknownChildren: true,
		Attributes: map[string]*AttributeSpec{
			"type": {Name: "type", Versions: allSupportedVersions},
		},
	},
	"Pricing": {
		Name:     "Pricing",
		Versions: allSupportedVersions,
		Attributes: map[string]*AttributeSpec{
			"model":    {Name: "model", Versions: allSupportedVersions, Required: true},
			"currency": {Name: "currency", Versions: allSupportedVersions, Required: true},
		},
	},
	"ViewableImpression": {
		Name:     "ViewableImpression",
		Versions: version40Plus,
		Children: map[string]*ChildSpec{
			"Viewable":         {Name: "Viewable", Versions: version40Plus, Optional: true, Multiple: true},
			"NotViewable":      {Name: "NotViewable", Versions: version40Plus, Optional: true, Multiple: true},
			"ViewUndetermined": {Name: "ViewUndetermined", Versions: version40Plus, Optional: true, Multiple: true},
		},
	},
	"Viewable": {
		Name:     "Viewable",
		Versions: version40Plus,
	},
	"NotViewable": {
		Name:     "NotViewable",
		Versions: version40Plus,
	},
	"ViewUndetermined": {
		Name:     "ViewUndetermined",
		Versions: version40Plus,
	},
	"Expires": {
		Name:     "Expires",
		Versions: allSupportedVersions,
	},
	"UniversalAdId": {
		Name:     "UniversalAdId",
		Versions: version40Plus,
		Attributes: map[string]*AttributeSpec{
			"idRegistry": {Name: "idRegistry", Versions: version40Plus, Required: true},
		},
	},
	"AdVerifications": {
		Name:     "AdVerifications",
		Versions: version40Plus,
		Children: map[string]*ChildSpec{
			"Verification": {Name: "Verification", Versions: version40Plus, Multiple: true},
		},
	},
	"Verification": {
		Name:     "Verification",
		Versions: version40Plus,
		Attributes: map[string]*AttributeSpec{
			"vendor": {Name: "vendor", Versions: allSupportedVersions},
		},
		Children: map[string]*ChildSpec{
			"JavaScriptResource":     {Name: "JavaScriptResource", Versions: version40Plus, Optional: true, Multiple: true},
			"ExecutableResource":     {Name: "ExecutableResource", Versions: version40Plus, Optional: true, Multiple: true},
			"TrackingEvents":         {Name: "TrackingEvents", Versions: version40Plus, Optional: true},
			"VerificationParameters": {Name: "VerificationParameters", Versions: version40Plus, Optional: true},
		},
	},
	"JavaScriptResource": {
		Name:     "JavaScriptResource",
		Versions: version40Plus,
		Attributes: map[string]*AttributeSpec{
			"apiFramework":    {Name: "apiFramework", Versions: allSupportedVersions},
			"browserOptional": {Name: "browserOptional", Versions: allSupportedVersions},
		},
	},
	"ExecutableResource": {
		Name:     "ExecutableResource",
		Versions: version40Plus,
		Attributes: map[string]*AttributeSpec{
			"apiFramework": {Name: "apiFramework", Versions: allSupportedVersions},
			"type":         {Name: "type", Versions: allSupportedVersions},
		},
	},
	"VerificationParameters": {
		Name:     "VerificationParameters",
		Versions: version40Plus,
	},
	"VASTAdTagURI": {
		Name:     "VASTAdTagURI",
		Versions: allSupportedVersions,
	},
	"Creatives": {
		Name:     "Creatives",
		Versions: allSupportedVersions,
		Children: map[string]*ChildSpec{
			"Creative": {Name: "Creative", Versions: allSupportedVersions, Multiple: true},
		},
	},
	"Creative": {
		Name:     "Creative",
		Versions: allSupportedVersions,
		Attributes: map[string]*AttributeSpec{
			"id":           {Name: "id", Versions: allSupportedVersions},
			"sequence":     {Name: "sequence", Versions: allSupportedVersions},
			"apiFramework": {Name: "apiFramework", Versions: allSupportedVersions},
			"adId":         {Name: "adId", Versions: allSupportedVersions},
		},
		Children: map[string]*ChildSpec{
			"Linear":             {Name: "Linear", Versions: allSupportedVersions, Optional: true},
			"NonLinearAds":       {Name: "NonLinearAds", Versions: allSupportedVersions, Optional: true},
			"CompanionAds":       {Name: "CompanionAds", Versions: allSupportedVersions, Optional: true},
			"CreativeExtensions": {Name: "CreativeExtensions", Versions: allSupportedVersions, Optional: true},
			"UniversalAdId":      {Name: "UniversalAdId", Versions: version40Plus, Optional: true, Multiple: true},
		},
	},
	"NonLinearAds": {
		Name:     "NonLinearAds",
		Versions: allSupportedVersions,
		Children: map[string]*ChildSpec{
			"NonLinear":      {Name: "NonLinear", Versions: allSupportedVersions, Multiple: true},
			"TrackingEvents": {Name: "TrackingEvents", Versions: allSupportedVersions, Optional: true},
		},
	},
	"NonLinear": {
		Name:     "NonLinear",
		Versions: allSupportedVersions,
		Attributes: map[string]*AttributeSpec{
			"id":                   {Name: "id", Versions: allSupportedVersions},
			"width":                {Name: "width", Versions: allSupportedVersions, Required: true},
			"height":               {Name: "height", Versions: allSupportedVersions, Required: true},
			"expandedWidth":        {Name: "expandedWidth", Versions: allSupportedVersions},
			"expandedHeight":       {Name: "expandedHeight", Versions: allSupportedVersions},
			"scalable":             {Name: "scalable", Versions: allSupportedVersions},
			"maintainAspectRatio":  {Name: "maintainAspectRatio", Versions: allSupportedVersions},
			"minSuggestedDuration": {Name: "minSuggestedDuration", Versions: allSupportedVersions},
			"apiFramework":         {Name: "apiFramework", Versions: allSupportedVersions},
		},
		Children: map[string]*ChildSpec{
			"StaticResource":         {Name: "StaticResource", Versions: allSupportedVersions, Optional: true},
			"IFrameResource":         {Name: "IFrameResource", Versions: allSupportedVersions, Optional: true},
			"HTMLResource":           {Name: "HTMLResource", Versions: allSupportedVersions, Optional: true},
			"AdParameters":           {Name: "AdParameters", Versions: allSupportedVersions, Optional: true},
			"NonLinearClickTracking": {Name: "NonLinearClickTracking", Versions: allSupportedVersions, Optional: true, Multiple: true},
			"NonLinearClickThrough":  {Name: "NonLinearClickThrough", Versions: allSupportedVersions, Optional: true},
		},
	},
	"NonLinearClickTracking": {
		Name:     "NonLinearClickTracking",
		Versions: allSupportedVersions,
	},
	"NonLinearClickThrough": {
		Name:     "NonLinearClickThrough",
		Versions: allSupportedVersions,
	},
	"CompanionAds": {
		Name:     "CompanionAds",
		Versions: allSupportedVersions,
		Attributes: map[string]*AttributeSpec{
			"required": {Name: "required", Versions: allSupportedVersions},
		},
		Children: map[string]*ChildSpec{
			"Companion": {Name: "Companion", Versions: allSupportedVersions, Multiple: true},
		},
	},
	"Companion": {
		Name:     "Companion",
		Versions: allSupportedVersions,
		Attributes: map[string]*AttributeSpec{
			"id":             {Name: "id", Versions: allSupportedVersions},
			"width":          {Name: "width", Versions: allSupportedVersions, Required: true},
			"height":         {Name: "height", Versions: allSupportedVersions, Required: true},
			"assetWidth":     {Name: "assetWidth", Versions: allSupportedVersions},
			"assetHeight":    {Name: "assetHeight", Versions: allSupportedVersions},
			"expandedWidth":  {Name: "expandedWidth", Versions: allSupportedVersions},
			"expandedHeight": {Name: "expandedHeight", Versions: allSupportedVersions},
			"apiFramework":   {Name: "apiFramework", Versions: allSupportedVersions},
			"adSlotId":       {Name: "adSlotId", Versions: allSupportedVersions},
			"pxratio":        {Name: "pxratio", Versions: allSupportedVersions},
			"renderingMode":  {Name: "renderingMode", Versions: allSupportedVersions},
		},
		Children: map[string]*ChildSpec{
			"StaticResource":         {Name: "StaticResource", Versions: allSupportedVersions, Optional: true},
			"IFrameResource":         {Name: "IFrameResource", Versions: allSupportedVersions, Optional: true},
			"HTMLResource":           {Name: "HTMLResource", Versions: allSupportedVersions, Optional: true},
			"AdParameters":           {Name: "AdParameters", Versions: allSupportedVersions, Optional: true},
			"AltText":                {Name: "AltText", Versions: allSupportedVersions, Optional: true},
			"CompanionClickThrough":  {Name: "CompanionClickThrough", Versions: allSupportedVersions, Optional: true},
			"CompanionClickTracking": {Name: "CompanionClickTracking", Versions: allSupportedVersions, Optional: true, Multiple: true},
			"CreativeExtensions":     {Name: "CreativeExtensions", Versions: allSupportedVersions, Optional: true},
			"TrackingEvents":         {Name: "TrackingEvents", Versions: allSupportedVersions, Optional: true},
		},
	},
	"CompanionClickThrough": {
		Name:     "CompanionClickThrough",
		Versions: allSupportedVersions,
	},
	"CompanionClickTracking": {
		Name:     "CompanionClickTracking",
		Versions: allSupportedVersions,
	},
	"AltText": {
		Name:     "AltText",
		Versions: allSupportedVersions,
	},
	"Icons": {
		Name:     "Icons",
		Versions: allSupportedVersions,
		Children: map[string]*ChildSpec{
			"Icon": {Name: "Icon", Versions: allSupportedVersions, Multiple: true},
		},
	},
	"Icon": {
		Name:     "Icon",
		Versions: allSupportedVersions,
		Attributes: map[string]*AttributeSpec{
			"program":      {Name: "program", Versions: allSupportedVersions},
			"width":        {Name: "width", Versions: allSupportedVersions},
			"height":       {Name: "height", Versions: allSupportedVersions},
			"xPosition":    {Name: "xPosition", Versions: allSupportedVersions},
			"yPosition":    {Name: "yPosition", Versions: allSupportedVersions},
			"duration":     {Name: "duration", Versions: allSupportedVersions},
			"offset":       {Name: "offset", Versions: allSupportedVersions},
			"apiFramework": {Name: "apiFramework", Versions: allSupportedVersions},
			"pxratio":      {Name: "pxratio", Versions: allSupportedVersions},
		},
		Children: map[string]*ChildSpec{
			"StaticResource":   {Name: "StaticResource", Versions: allSupportedVersions, Optional: true},
			"IFrameResource":   {Name: "IFrameResource", Versions: allSupportedVersions, Optional: true},
			"HTMLResource":     {Name: "HTMLResource", Versions: allSupportedVersions, Optional: true},
			"IconClicks":       {Name: "IconClicks", Versions: allSupportedVersions, Optional: true},
			"IconViewTracking": {Name: "IconViewTracking", Versions: allSupportedVersions, Optional: true, Multiple: true},
		},
	},
	"IconClicks": {
		Name:     "IconClicks",
		Versions: allSupportedVersions,
		Children: map[string]*ChildSpec{
			"IconClickFallbackImages": {Name: "IconClickFallbackImages", Versions: allSupportedVersions, Optional: true},
			"IconClickThrough":        {Name: "IconClickThrough", Versions: allSupportedVersions, Optional: true},
			"IconClickTracking":       {Name: "IconClickTracking", Versions: allSupportedVersions, Optional: true, Multiple: true},
		},
	},
	"IconClickThrough": {
		Name:     "IconClickThrough",
		Versions: allSupportedVersions,
	},
	"IconClickTracking": {
		Name:     "IconClickTracking",
		Versions: allSupportedVersions,
	},
	"IconClickFallbackImages": {
		Name:     "IconClickFallbackImages",
		Versions: allSupportedVersions,
		Children: map[string]*ChildSpec{
			"IconClickFallbackImage": {Name: "IconClickFallbackImage", Versions: allSupportedVersions, Multiple: true},
		},
	},
	"IconClickFallbackImage": {
		Name:     "IconClickFallbackImage",
		Versions: allSupportedVersions,
		Attributes: map[string]*AttributeSpec{
			"width":  {Name: "width", Versions: allSupportedVersions},
			"height": {Name: "height", Versions: allSupportedVersions},
		},
		Children: map[string]*ChildSpec{
			"AltText":        {Name: "AltText", Versions: allSupportedVersions, Optional: true},
			"StaticResource": {Name: "StaticResource", Versions: allSupportedVersions, Optional: true},
		},
	},
	"IconViewTracking": {
		Name:     "IconViewTracking",
		Versions: allSupportedVersions,
	},
	"CreativeExtensions": {
		Name:     "CreativeExtensions",
		Versions: allSupportedVersions,
		Children: map[string]*ChildSpec{
			"CreativeExtension": {Name: "CreativeExtension", Versions: allSupportedVersions, Multiple: true},
		},
	},
	"CreativeExtension": {
		Name:                 "CreativeExtension",
		Versions:             allSupportedVersions,
		AllowUnknownChildren: true,
		Attributes: map[string]*AttributeSpec{
			"type": {Name: "type", Versions: allSupportedVersions},
		},
	},
	"StaticResource": {
		Name:     "StaticResource",
		Versions: allSupportedVersions,
		Attributes: map[string]*AttributeSpec{
			"creativeType": {Name: "creativeType", Versions: allSupportedVersions},
		},
	},
	"HTMLResource": {
		Name:     "HTMLResource",
		Versions: allSupportedVersions,
	},
	"IFrameResource": {
		Name:     "IFrameResource",
		Versions: allSupportedVersions,
	},
	"Linear": {
		Name:     "Linear",
		Versions: allSupportedVersions,
		Attributes: map[string]*AttributeSpec{
			"skipoffset": {Name: "skipoffset", Versions: allSupportedVersions},
		},
		Children: map[string]*ChildSpec{
			"Icons":          {Name: "Icons", Versions: allSupportedVersions, Optional: true},
			"AdParameters":   {Name: "AdParameters", Versions: allSupportedVersions, Optional: true},
			"Duration":       {Name: "Duration", Versions: allSupportedVersions},
			"MediaFiles":     {Name: "MediaFiles", Versions: allSupportedVersions},
			"VideoClicks":    {Name: "VideoClicks", Versions: allSupportedVersions, Optional: true},
			"TrackingEvents": {Name: "TrackingEvents", Versions: allSupportedVersions, Optional: true},
		},
	},
	"AdParameters": {
		Name:     "AdParameters",
		Versions: allSupportedVersions,
		Attributes: map[string]*AttributeSpec{
			"xmlEncoded": {Name: "xmlEncoded", Versions: allSupportedVersions},
		},
	},
	"Duration": {
		Name:     "Duration",
		Versions: allSupportedVersions,
	},
	"MediaFiles": {
		Name:     "MediaFiles",
		Versions: allSupportedVersions,
		Children: map[string]*ChildSpec{
			"MediaFile":               {Name: "MediaFile", Versions: allSupportedVersions, Multiple: true},
			"ClosedCaptionFiles":      {Name: "ClosedCaptionFiles", Versions: allSupportedVersions, Optional: true},
			"Mezzanine":               {Name: "Mezzanine", Versions: allSupportedVersions, Optional: true, Multiple: true},
			"InteractiveCreativeFile": {Name: "InteractiveCreativeFile", Versions: allSupportedVersions, Optional: true, Multiple: true},
		},
	},
	"MediaFile": {
		Name:     "MediaFile",
		Versions: allSupportedVersions,
		Attributes: map[string]*AttributeSpec{
			"id":                  {Name: "id", Versions: allSupportedVersions},
			"delivery":            {Name: "delivery", Versions: allSupportedVersions, Required: true},
			"type":                {Name: "type", Versions: allSupportedVersions, Required: true},
			"width":               {Name: "width", Versions: allSupportedVersions, Required: true},
			"height":              {Name: "height", Versions: allSupportedVersions, Required: true},
			"codec":               {Name: "codec", Versions: allSupportedVersions},
			"bitrate":             {Name: "bitrate", Versions: allSupportedVersions},
			"minBitrate":          {Name: "minBitrate", Versions: allSupportedVersions},
			"maxBitrate":          {Name: "maxBitrate", Versions: allSupportedVersions},
			"scalable":            {Name: "scalable", Versions: allSupportedVersions},
			"maintainAspectRatio": {Name: "maintainAspectRatio", Versions: allSupportedVersions},
			"fileSize":            {Name: "fileSize", Versions: allSupportedVersions},
			"mediaType":           {Name: "mediaType", Versions: allSupportedVersions},
			"apiFramework":        {Name: "apiFramework", Versions: allSupportedVersions},
		},
	},
	"ClosedCaptionFiles": {
		Name:     "ClosedCaptionFiles",
		Versions: allSupportedVersions,
		Children: map[string]*ChildSpec{
			"ClosedCaptionFile": {Name: "ClosedCaptionFile", Versions: allSupportedVersions, Multiple: true},
		},
	},
	"ClosedCaptionFile": {
		Name:     "ClosedCaptionFile",
		Versions: allSupportedVersions,
		Attributes: map[string]*AttributeSpec{
			"type":     {Name: "type", Versions: allSupportedVersions},
			"language": {Name: "language", Versions: allSupportedVersions},
		},
	},
	"Mezzanine": {
		Name:     "Mezzanine",
		Versions: allSupportedVersions,
		Attributes: map[string]*AttributeSpec{
			"delivery":  {Name: "delivery", Versions: allSupportedVersions, Required: true},
			"type":      {Name: "type", Versions: allSupportedVersions, Required: true},
			"width":     {Name: "width", Versions: allSupportedVersions, Required: true},
			"height":    {Name: "height", Versions: allSupportedVersions, Required: true},
			"codec":     {Name: "codec", Versions: allSupportedVersions},
			"fileSize":  {Name: "fileSize", Versions: allSupportedVersions},
			"mediaType": {Name: "mediaType", Versions: allSupportedVersions},
		},
	},
	"InteractiveCreativeFile": {
		Name:     "InteractiveCreativeFile",
		Versions: allSupportedVersions,
		Attributes: map[string]*AttributeSpec{
			"type":             {Name: "type", Versions: allSupportedVersions},
			"apiFramework":     {Name: "apiFramework", Versions: allSupportedVersions},
			"variableDuration": {Name: "variableDuration", Versions: allSupportedVersions},
		},
	},
	"VideoClicks": {
		Name:     "VideoClicks",
		Versions: allSupportedVersions,
		Children: map[string]*ChildSpec{
			"ClickThrough":  {Name: "ClickThrough", Versions: allSupportedVersions, Optional: true},
			"ClickTracking": {Name: "ClickTracking", Versions: allSupportedVersions, Optional: true, Multiple: true},
			"CustomClick":   {Name: "CustomClick", Versions: allSupportedVersions, Optional: true, Multiple: true},
		},
	},
	"ClickThrough": {
		Name:     "ClickThrough",
		Versions: allSupportedVersions,
		Attributes: map[string]*AttributeSpec{
			"id": {Name: "id", Versions: allSupportedVersions},
		},
	},
	"ClickTracking": {
		Name:     "ClickTracking",
		Versions: allSupportedVersions,
	},
	"CustomClick": {
		Name:     "CustomClick",
		Versions: allSupportedVersions,
	},
	"TrackingEvents": {
		Name:     "TrackingEvents",
		Versions: allSupportedVersions,
		Children: map[string]*ChildSpec{
			"Tracking": {Name: "Tracking", Versions: allSupportedVersions, Multiple: true},
		},
	},
	"Tracking": {
		Name:     "Tracking",
		Versions: allSupportedVersions,
		Attributes: map[string]*AttributeSpec{
			"event":  {Name: "event", Versions: allSupportedVersions, Required: true},
			"offset": {Name: "offset", Versions: allSupportedVersions},
		},
	},
}}
