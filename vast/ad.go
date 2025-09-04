package vast

// AdType represents the type of ad content as defined in IAB VAST 4.x specification.
// Reference: IAB VAST 4.x Section 2.3.1.1 - Ad Element
// Link: https://iabtechlab.com/wp-content/uploads/2019/06/VAST_4.2_final_june26.pdf
type AdType string

// Ad represents a single advertisement unit containing either inline content or a wrapper.
// The Ad element contains the actual ad content or a reference to another VAST document.
//
// Reference: IAB VAST 4.x Section 2.3.1.1 - Ad Element
// Link: https://iabtechlab.com/wp-content/uploads/2019/06/VAST_4.2_final_june26.pdf#page=39
type Ad struct {
	InLine        *InLine     `xml:"InLine,omitempty"`
	Wrapper       *Wrapper    `xml:"Wrapper,omitempty"`
	ID            string      `xml:"id,attr,omitempty"`
	Sequence      int         `xml:"sequence,attr,omitempty"`
	ConditionalAd NumericBool `xml:"conditionalAd,attr,omitempty"`
	AdType        AdType      `xml:"adType,attr,omitempty"`
}

// AdType constants as defined in IAB VAST 4.x specification.
// Reference: IAB VAST 4.x Section 2.3.1.1 - adType attribute
// Link: https://iabtechlab.com/wp-content/uploads/2019/06/VAST_4.2_final_june26.pdf
const (
	VideoAdType  AdType = "video"  // Video advertisement content
	AudioAdType  AdType = "audio"  // Audio-only advertisement content
	HybridAdType AdType = "hybrid" // Mixed audio and video content
)

// AdDefinition contains common elements shared between InLine and Wrapper ads.
// These elements provide essential ad metadata and tracking information.
type AdDefinition struct {
	AdSystem           AdSystem            `xml:"AdSystem"`
	Error              []CData             `xml:"Error,omitempty"`
	Impression         []Impression        `xml:"Impression"`
	Extensions         *Extensions         `xml:"Extensions,omitempty"`
	Pricing            *Pricing            `xml:"Pricing,omitempty"`
	ViewableImpression *ViewableImpression `xml:"ViewableImpression,omitempty"`
}

// AdSystem identifies the ad server that returned the ad and optionally its version.
// Reference: IAB VAST 4.x Section 2.3.1.3 - AdSystem Element
// Link: https://iabtechlab.com/wp-content/uploads/2019/06/VAST_4.2_final_june26.pdf#page=41
type AdSystem struct {
	Value   string `xml:",chardata"`
	Version string `xml:"version,attr,omitempty"`
}

// ViewableImpression provides URLs for tracking viewable impressions according to different viewability states.
// Supports MRC viewability standards and provides granular viewability tracking.
//
// Reference: IAB VAST 4.x Section 2.3.1.4 - ViewableImpression Element
// Link: https://iabtechlab.com/wp-content/uploads/2019/06/VAST_4.2_final_june26.pdf#page=46
type ViewableImpression struct {
	Viewable         []CData `xml:"Viewable,omitempty"`
	NotViewable      []CData `xml:"NotViewable,omitempty"`
	ViewUndetermined []CData `xml:"ViewUndetermined,omitempty"`
}
