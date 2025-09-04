package vast

// NonLinearAds contains a collection of non-linear advertisements and their tracking.
// Provides container for overlay ads that display during video content.
//
// Reference: IAB VAST 4.x Section 2.3.3.1 - NonLinearAds Element
// Link: https://iabtechlab.com/wp-content/uploads/2019/06/VAST_4.2_final_june26.pdf#page=62
type NonLinearAds struct {
	TrackingEvents *TrackingEvents `xml:"TrackingEvents,omitempty"`
	NonLinear      []NonLinearAd   `xml:"NonLinear,omitempty"`
}

// NonLinearAd represents a non-linear overlay ad displayed during video playback.
// Contains overlay content that appears on top of video without interrupting playback.
//
// Reference: IAB VAST 4.x Section 2.3.3.1 - NonLinear Element
// Link: https://iabtechlab.com/wp-content/uploads/2019/06/VAST_4.2_final_june26.pdf#page=62
type NonLinearAd struct {
	HTMLResource           []CData          `xml:"HTMLResource,omitempty"`
	IFrameResource         []CData          `xml:"IFrameResource,omitempty"`
	StaticResource         []StaticResource `xml:"StaticResource,omitempty"`
	AdParameters           *AdParameters    `xml:"AdParameters,omitempty"`
	NonLinearClickThrough  *CData           `xml:"NonLinearClickThrough,omitempty"`
	NonLinearClickTracking []CData          `xml:"NonLinearClickTracking,omitempty"`
	Width                  int              `xml:"width,attr"`
	Height                 int              `xml:"height,attr"`
	ExpandedWidth          int              `xml:"expandedWidth,attr,omitempty"`
	ExpandedHeight         int              `xml:"expandedHeight,attr,omitempty"`
	Scalable               NumericBool      `xml:"scalable,attr,omitempty"`
	MaintainAspectRatio    NumericBool      `xml:"maintainAspectRatio,attr,omitempty"`
	MinSuggestedDuration   Duration         `xml:"minSuggestedDuration,attr,omitempty"`
	APIFramework           string           `xml:"apiFramework,attr,omitempty"`
}
