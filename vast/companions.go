package vast

// Required specifies how many companion ads are required to be displayed.
// Reference: IAB VAST 4.x Section 2.3.3.2 - CompanionAds Element
// Link: https://iabtechlab.com/wp-content/uploads/2019/06/VAST_4.2_final_june26.pdf#page=35
type Required string

const (
	AllRequired  Required = "all"
	AnyRequired  Required = "any"
	NoneRequired Required = "none"
)

// RenderingMode defines when and how companion ads should be rendered.
// Reference: IAB VAST 4.x Section 2.3.3.3 - Companion Element
// Link: https://iabtechlab.com/wp-content/uploads/2019/06/VAST_4.2_final_june26.pdf#page=36
type RenderingMode string

const (
	DefaultRenderingMode    RenderingMode = "default"
	EndCardRenderingMode    RenderingMode = "end-card"
	ConcurrentRenderingMode RenderingMode = "concurrent"
)

// CompanionAds contains a collection of companion advertisements.
// Companion ads are displayed alongside or after the main video ad content.
//
// Reference: IAB VAST 4.x Section 2.3.3.2 - CompanionAds Element
// Link: https://iabtechlab.com/wp-content/uploads/2019/06/VAST_4.2_final_june26.pdf#page=63
type CompanionAds struct {
	Companion []CompanionAd `xml:"Companion,omitempty"`
	Required  Required      `xml:"required,attr,omitempty"`
}

// CompanionAd represents a single companion advertisement with its resources and tracking.
// Companion ads provide additional advertising space and can contain various media types.
//
// Reference: IAB VAST 4.x Section 2.3.3.3 - Companion Element
// Link: https://iabtechlab.com/wp-content/uploads/2019/06/VAST_4.2_final_june26.pdf#page=66
type CompanionAd struct {
	HTMLResource           []CData             `xml:"HTMLResource,omitempty"`
	IFrameResource         []CData             `xml:"IFrameResource,omitempty"`
	StaticResource         []StaticResource    `xml:"StaticResource,omitempty"`
	AdParameters           *AdParameters       `xml:"AdParameters,omitempty"`
	AltText                string              `xml:"AltText,omitempty"`
	CompanionClickThrough  *CData              `xml:"CompanionClickThrough,omitempty"`
	CompanionClickTracking []string            `xml:"CompanionClickTracking,omitempty"`
	CreativeExtensions     *CreativeExtensions `xml:"CreativeExtensions,omitempty"`
	TrackingEvents         *TrackingEvents     `xml:"TrackingEvents,omitempty"`
	ID                     string              `xml:"id,attr,omitempty"`
	Width                  int                 `xml:"width,attr"`
	Height                 int                 `xml:"height,attr"`
	AssetWidth             int                 `xml:"assetWidth,attr,omitempty"`
	AssetHeight            int                 `xml:"assetHeight,attr,omitempty"`
	ExpandedWidth          int                 `xml:"expandedWidth,attr,omitempty"`
	ExpandedHeight         int                 `xml:"expandedHeight,attr,omitempty"`
	APIFramework           string              `xml:"apiFramework,attr,omitempty"`
	AdSlotID               string              `xml:"adSlotId,attr,omitempty"`
	PXRatio                float64             `xml:"pxratio,attr,omitempty"`
	RenderingMode          RenderingMode       `xml:"renderingMode,attr,omitempty"`
}
