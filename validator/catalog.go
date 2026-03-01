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
	supported43Plus = []vast.Version{
		vast.Version43,
	}
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
			"sequence":      {Name: "sequence", Versions: supported30Plus},
			"conditionalAd": {Name: "conditionalAd", Versions: supported40Plus},
			"adType":        {Name: "adType", Versions: supported41Plus},
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
			"AdSystem":           {Name: "AdSystem", Versions: supported20Plus},
			"Error":              {Name: "Error", Versions: supported20Plus, Optional: true, Multiple: true},
			"Impression":         {Name: "Impression", Versions: supported20Plus, Multiple: true},
			"AdTitle":            {Name: "AdTitle", Versions: supported20Plus},
			"AdServingId":        {Name: "AdServingId", Versions: supported30Plus, Optional: true},
			"Advertiser":         {Name: "Advertiser", Versions: supported30Plus, Optional: true},
			"Category":           {Name: "Category", Versions: supported30Plus, Optional: true, Multiple: true},
			"Description":        {Name: "Description", Versions: supported20Plus, Optional: true},
			"Extensions":         {Name: "Extensions", Versions: supported20Plus, Optional: true},
			"Pricing":            {Name: "Pricing", Versions: supported30Plus, Optional: true},
			"ViewableImpression": {Name: "ViewableImpression", Versions: supported40Plus, Optional: true},
			"Survey":             {Name: "Survey", Versions: supported20Plus, Optional: true},
			"Expires":            {Name: "Expires", Versions: supported30Plus, Optional: true},
			"Creatives":          {Name: "Creatives", Versions: supported20Plus},
			"AdVerifications":    {Name: "AdVerifications", Versions: supported40Plus, Optional: true},
		},
	},
	"Wrapper": {
		Name:     "Wrapper",
		Versions: supported20Plus,
		Attributes: map[string]*AttributeSpec{
			"followAdditionalWrappers": {Name: "followAdditionalWrappers", Versions: supported30Plus},
			"allowMultipleAds":         {Name: "allowMultipleAds", Versions: supported30Plus},
			"fallbackOnNoAd":           {Name: "fallbackOnNoAd", Versions: supported30Plus},
		},
		Children: map[string]*ChildSpec{
			"AdSystem":            {Name: "AdSystem", Versions: supported20Plus},
			"Error":               {Name: "Error", Versions: supported20Plus, Optional: true, Multiple: true},
			"Impression":          {Name: "Impression", Versions: supported20Plus, Multiple: true},
			"VASTAdTagURI":        {Name: "VASTAdTagURI", Versions: supported20Plus},
			"Extensions":          {Name: "Extensions", Versions: supported20Plus, Optional: true},
			"Pricing":             {Name: "Pricing", Versions: supported30Plus, Optional: true},
			"ViewableImpression":  {Name: "ViewableImpression", Versions: supported40Plus, Optional: true},
			"Creatives":           {Name: "Creatives", Versions: supported20Plus, Optional: true},
			"BlockedAdCategories": {Name: "BlockedAdCategories", Versions: supported30Plus, Optional: true, Multiple: true},
			"AdVerifications":     {Name: "AdVerifications", Versions: supported40Plus, Optional: true},
		},
	},
	"AdSystem": {
		Name:     "AdSystem",
		Versions: supported20Plus,
		Attributes: map[string]*AttributeSpec{
			"version": {Name: "version", Versions: supported20Plus},
		},
	},
	"Error": {
		Name:     "Error",
		Versions: supported20Plus,
	},
	"Impression": {
		Name:     "Impression",
		Versions: supported20Plus,
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
		Versions: supported30Plus,
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
			"authority": {Name: "authority", Versions: supported30Plus, Required: true},
		},
	},
	"BlockedAdCategories": {
		Name:     "BlockedAdCategories",
		Versions: supported30Plus,
		Attributes: map[string]*AttributeSpec{
			"authority": {Name: "authority", Versions: supported30Plus},
		},
	},
	"Description": {
		Name:     "Description",
		Versions: supported20Plus,
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
		Name:                 "Extension",
		Versions:             supported20Plus,
		AllowUnknownChildren: true,
		Attributes: map[string]*AttributeSpec{
			"type": {Name: "type", Versions: supported20Plus},
		},
	},
	"Pricing": {
		Name:     "Pricing",
		Versions: supported30Plus,
		Attributes: map[string]*AttributeSpec{
			"model":    {Name: "model", Versions: supported30Plus, Required: true},
			"currency": {Name: "currency", Versions: supported30Plus, Required: true},
		},
	},
	"ViewableImpression": {
		Name:     "ViewableImpression",
		Versions: supported40Plus,
		Children: map[string]*ChildSpec{
			"Viewable":         {Name: "Viewable", Versions: supported40Plus, Optional: true, Multiple: true},
			"NotViewable":      {Name: "NotViewable", Versions: supported40Plus, Optional: true, Multiple: true},
			"ViewUndetermined": {Name: "ViewUndetermined", Versions: supported40Plus, Optional: true, Multiple: true},
		},
	},
	"Viewable": {
		Name:     "Viewable",
		Versions: supported40Plus,
	},
	"NotViewable": {
		Name:     "NotViewable",
		Versions: supported40Plus,
	},
	"ViewUndetermined": {
		Name:     "ViewUndetermined",
		Versions: supported40Plus,
	},
	"Expires": {
		Name:     "Expires",
		Versions: supported30Plus,
	},
	"UniversalAdId": {
		Name:     "UniversalAdId",
		Versions: supported40Plus,
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
		},
	},
	"JavaScriptResource": {
		Name:     "JavaScriptResource",
		Versions: supported40Plus,
		Attributes: map[string]*AttributeSpec{
			"apiFramework":    {Name: "apiFramework", Versions: supported30Plus},
			"browserOptional": {Name: "browserOptional", Versions: supported30Plus},
		},
	},
	"ExecutableResource": {
		Name:     "ExecutableResource",
		Versions: supported40Plus,
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
		Name:     "VASTAdTagURI",
		Versions: supported20Plus,
	},
	"Creatives": {
		Name:     "Creatives",
		Versions: supported20Plus,
		Children: map[string]*ChildSpec{
			"Creative": {Name: "Creative", Versions: supported20Plus, Multiple: true},
		},
	},
	"Creative": {
		Name:     "Creative",
		Versions: supported20Plus,
		Attributes: map[string]*AttributeSpec{
			"id":           {Name: "id", Versions: supported20Plus},
			"sequence":     {Name: "sequence", Versions: supported30Plus},
			"apiFramework": {Name: "apiFramework", Versions: supported30Plus},
			"adId":         {Name: "adId", Versions: supported30Plus},
		},
		Children: map[string]*ChildSpec{
			"Linear":             {Name: "Linear", Versions: supported20Plus, Optional: true},
			"NonLinearAds":       {Name: "NonLinearAds", Versions: supported20Plus, Optional: true},
			"CompanionAds":       {Name: "CompanionAds", Versions: supported20Plus, Optional: true},
			"CreativeExtensions": {Name: "CreativeExtensions", Versions: supported30Plus, Optional: true},
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
			"width":                {Name: "width", Versions: supported20Plus, Required: true},
			"height":               {Name: "height", Versions: supported20Plus, Required: true},
			"expandedWidth":        {Name: "expandedWidth", Versions: supported30Plus},
			"expandedHeight":       {Name: "expandedHeight", Versions: supported30Plus},
			"scalable":             {Name: "scalable", Versions: supported20Plus},
			"maintainAspectRatio":  {Name: "maintainAspectRatio", Versions: supported20Plus},
			"minSuggestedDuration": {Name: "minSuggestedDuration", Versions: supported20Plus},
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
		Name:     "NonLinearClickTracking",
		Versions: supported30Plus,
		Attributes: map[string]*AttributeSpec{
			"id": {Name: "id", Versions: supported30Plus},
		},
	},
	"NonLinearClickThrough": {
		Name:     "NonLinearClickThrough",
		Versions: supported20Plus,
	},
	"CompanionAds": {
		Name:     "CompanionAds",
		Versions: supported20Plus,
		Attributes: map[string]*AttributeSpec{
			"required": {Name: "required", Versions: supported30Plus},
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
			"width":          {Name: "width", Versions: supported20Plus, Required: true},
			"height":         {Name: "height", Versions: supported20Plus, Required: true},
			"assetWidth":     {Name: "assetWidth", Versions: supported30Plus},
			"assetHeight":    {Name: "assetHeight", Versions: supported30Plus},
			"expandedWidth":  {Name: "expandedWidth", Versions: supported30Plus},
			"expandedHeight": {Name: "expandedHeight", Versions: supported30Plus},
			"apiFramework":   {Name: "apiFramework", Versions: supported20Plus},
			"adSlotId":       {Name: "adSlotId", Versions: supported30Plus},
			"logoTile":       {Name: "logoTile", Versions: supported40Plus},
			"logoTitle":      {Name: "logoTitle", Versions: supported40Plus},
			"logoArtist":     {Name: "logoArtist", Versions: supported40Plus},
			"logoURL":        {Name: "logoURL", Versions: supported40Plus},
			"pxratio":        {Name: "pxratio", Versions: supported30Plus},
			"renderingMode":  {Name: "renderingMode", Versions: supported30Plus},
		},
		Children: map[string]*ChildSpec{
			"StaticResource":         {Name: "StaticResource", Versions: supported20Plus, Optional: true},
			"IFrameResource":         {Name: "IFrameResource", Versions: supported20Plus, Optional: true},
			"HTMLResource":           {Name: "HTMLResource", Versions: supported20Plus, Optional: true},
			"AdParameters":           {Name: "AdParameters", Versions: supported20Plus, Optional: true},
			"AltText":                {Name: "AltText", Versions: supported20Plus, Optional: true},
			"CompanionClickThrough":  {Name: "CompanionClickThrough", Versions: supported20Plus, Optional: true},
			"CompanionClickTracking": {Name: "CompanionClickTracking", Versions: supported30Plus, Optional: true, Multiple: true},
			"CreativeExtensions":     {Name: "CreativeExtensions", Versions: supported30Plus, Optional: true},
			"TrackingEvents":         {Name: "TrackingEvents", Versions: supported20Plus, Optional: true},
		},
	},
	"CompanionClickThrough": {
		Name:     "CompanionClickThrough",
		Versions: supported20Plus,
	},
	"CompanionClickTracking": {
		Name:     "CompanionClickTracking",
		Versions: supported30Plus,
		Attributes: map[string]*AttributeSpec{
			"id": {Name: "id", Versions: supported30Plus},
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
			"width":        {Name: "width", Versions: supported30Plus},
			"height":       {Name: "height", Versions: supported30Plus},
			"xPosition":    {Name: "xPosition", Versions: supported30Plus},
			"yPosition":    {Name: "yPosition", Versions: supported30Plus},
			"duration":     {Name: "duration", Versions: supported30Plus},
			"offset":       {Name: "offset", Versions: supported30Plus},
			"apiFramework": {Name: "apiFramework", Versions: supported30Plus},
			"pxratio":      {Name: "pxratio", Versions: supported30Plus},
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
			"IconClickFallbackImages": {Name: "IconClickFallbackImages", Versions: supported30Plus, Optional: true},
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
		Versions: supported30Plus,
		Children: map[string]*ChildSpec{
			"IconClickFallbackImage": {Name: "IconClickFallbackImage", Versions: supported30Plus, Multiple: true},
		},
	},
	"IconClickFallbackImage": {
		Name:     "IconClickFallbackImage",
		Versions: supported30Plus,
		Attributes: map[string]*AttributeSpec{
			"width":  {Name: "width", Versions: supported30Plus},
			"height": {Name: "height", Versions: supported30Plus},
		},
		Children: map[string]*ChildSpec{
			"AltText":        {Name: "AltText", Versions: supported30Plus, Optional: true},
			"StaticResource": {Name: "StaticResource", Versions: supported30Plus, Optional: true},
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
		Name:                 "CreativeExtension",
		Versions:             supported30Plus,
		AllowUnknownChildren: true,
		Attributes: map[string]*AttributeSpec{
			"type": {Name: "type", Versions: supported30Plus},
		},
	},
	"StaticResource": {
		Name:     "StaticResource",
		Versions: supported20Plus,
		Attributes: map[string]*AttributeSpec{
			"creativeType": {Name: "creativeType", Versions: supported30Plus},
		},
	},
	"HTMLResource": {
		Name:     "HTMLResource",
		Versions: supported20Plus,
	},
	"IFrameResource": {
		Name:     "IFrameResource",
		Versions: supported20Plus,
	},
	"Linear": {
		Name:     "Linear",
		Versions: supported20Plus,
		Attributes: map[string]*AttributeSpec{
			"skipoffset": {Name: "skipoffset", Versions: supported30Plus},
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
	"AdParameters": {
		Name:     "AdParameters",
		Versions: supported20Plus,
		Attributes: map[string]*AttributeSpec{
			"xmlEncoded": {Name: "xmlEncoded", Versions: supported30Plus},
		},
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
			"Mezzanine":               {Name: "Mezzanine", Versions: supported30Plus, Optional: true, Multiple: true},
			"InteractiveCreativeFile": {Name: "InteractiveCreativeFile", Versions: supported30Plus, Optional: true, Multiple: true},
		},
	},
	"MediaFile": {
		Name:     "MediaFile",
		Versions: supported20Plus,
		Attributes: map[string]*AttributeSpec{
			"id":                  {Name: "id", Versions: supported20Plus},
			"delivery":            {Name: "delivery", Versions: supported20Plus, Required: true},
			"type":                {Name: "type", Versions: supported20Plus, Required: true},
			"width":               {Name: "width", Versions: supported20Plus, Required: true},
			"height":              {Name: "height", Versions: supported20Plus, Required: true},
			"codec":               {Name: "codec", Versions: supported30Plus},
			"bitrate":             {Name: "bitrate", Versions: supported30Plus},
			"minBitrate":          {Name: "minBitrate", Versions: supported30Plus},
			"maxBitrate":          {Name: "maxBitrate", Versions: supported30Plus},
			"scalable":            {Name: "scalable", Versions: supported30Plus},
			"maintainAspectRatio": {Name: "maintainAspectRatio", Versions: supported30Plus},
			"fileSize":            {Name: "fileSize", Versions: supported30Plus},
			"mediaType":           {Name: "mediaType", Versions: supported30Plus},
			"apiFramework":        {Name: "apiFramework", Versions: supported41Plus},
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
		Name:     "ClosedCaptionFile",
		Versions: supported30Plus,
		Attributes: map[string]*AttributeSpec{
			"type":     {Name: "type", Versions: supported30Plus},
			"language": {Name: "language", Versions: supported30Plus},
		},
	},
	"Mezzanine": {
		Name:     "Mezzanine",
		Versions: supported30Plus,
		Attributes: map[string]*AttributeSpec{
			"delivery":  {Name: "delivery", Versions: supported30Plus, Required: true},
			"type":      {Name: "type", Versions: supported30Plus, Required: true},
			"width":     {Name: "width", Versions: supported30Plus, Required: true},
			"height":    {Name: "height", Versions: supported30Plus, Required: true},
			"codec":     {Name: "codec", Versions: supported30Plus},
			"fileSize":  {Name: "fileSize", Versions: supported30Plus},
			"mediaType": {Name: "mediaType", Versions: supported30Plus},
			"id":        {Name: "id", Versions: supported41Plus},
		},
	},
	"InteractiveCreativeFile": {
		Name:     "InteractiveCreativeFile",
		Versions: supported30Plus,
		Attributes: map[string]*AttributeSpec{
			"type":             {Name: "type", Versions: supported30Plus},
			"apiFramework":     {Name: "apiFramework", Versions: supported30Plus},
			"variableDuration": {Name: "variableDuration", Versions: supported30Plus},
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
		Name:     "ClickThrough",
		Versions: supported20Plus,
		Attributes: map[string]*AttributeSpec{
			"id": {Name: "id", Versions: supported30Plus},
		},
	},
	"ClickTracking": {
		Name:     "ClickTracking",
		Versions: supported20Plus,
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
		Name:     "Tracking",
		Versions: supported20Plus,
		Attributes: map[string]*AttributeSpec{
			"event":  {Name: "event", Versions: supported20Plus, Required: true},
			"offset": {Name: "offset", Versions: supported30Plus},
		},
	},
}}
