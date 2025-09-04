package vast

// Icons contains a collection of overlay icons displayed during ad playback.
// Provides container for multiple icon overlays shown on video ads.
//
// Reference: IAB VAST 4.x Section 2.3.3.4 - Icons Element
// Link: https://iabtechlab.com/wp-content/uploads/2019/06/VAST_4.2_final_june26.pdf#page=58
type Icons struct {
	Icon []Icon `xml:"Icon"`
}

// Icon represents an overlay icon displayed on top of video ad content.
// Defines interactive overlay elements with positioning, timing, and click tracking.
//
// Reference: IAB VAST 4.x Section 2.3.3.4 - Icon Element
// Link: https://iabtechlab.com/wp-content/uploads/2019/06/VAST_4.2_final_june26.pdf#page=59
type Icon struct {
	CreativeResource

	HTMLResource     []CData          `xml:"HTMLResource,omitempty"`
	IFrameResource   []CData          `xml:"IFrameResource,omitempty"`
	StaticResource   []StaticResource `xml:"StaticResource,omitempty"`
	IconClicks       *IconClicks      `xml:"IconClicks,omitempty"`
	IconViewTracking []string         `xml:"IconViewTracking,omitempty"`
	Program          string           `xml:"program,attr,omitempty"`
	Width            int              `xml:"width,attr,omitempty"`
	Height           int              `xml:"height,attr,omitempty"`
	XPosition        XPosition        `xml:"xPosition,attr,omitempty"`
	YPosition        YPosition        `xml:"yPosition,attr,omitempty"`
	Duration         Duration         `xml:"duration,attr,omitempty"`
	Offset           Duration         `xml:"offset,attr,omitempty"`
	APIFramework     string           `xml:"apiFramework,attr,omitempty"`
	PXRatio          float64          `xml:"pxratio,attr,omitempty"`
}

// IconClickFallbackImages contains a collection of fallback images for icon clicks.
// Holds multiple fallback image options for icon click scenarios.
//
// Reference: IAB VAST 4.x Section 2.3.3.4 - IconClickFallbackImages Element
// Link: https://iabtechlab.com/wp-content/uploads/2019/06/VAST_4.2_final_june26.pdf#page=61
type IconClickFallbackImages struct {
	IconClickFallbackImage []IconClickFallbackImage `xml:"IconClickFallbackImage"`
}

// IconClickFallbackImage represents a fallback image displayed when icon click tracking fails.
// Provides alternative visual content when primary icon click functionality is unavailable.
//
// Reference: IAB VAST 4.x Section 2.3.3.4 - IconClickFallbackImages Element
// Link: https://iabtechlab.com/wp-content/uploads/2019/06/VAST_4.2_final_june26.pdf#page=61
type IconClickFallbackImage struct {
	AltText        string `xml:"AltText,omitempty"`
	StaticResource *CData `xml:"StaticResource,omitempty"`
	Height         int    `xml:"height,attr,omitempty"`
	Width          int    `xml:"width,attr,omitempty"`
}

// IconClicks defines click behavior and tracking for overlay icons.
// Manages user interaction tracking and click-through behavior for icon overlays.
//
// Reference: IAB VAST 4.x Section 2.3.3.4 - IconClicks Element
// Link: https://iabtechlab.com/wp-content/uploads/2019/06/VAST_4.2_final_june26.pdf#page=60
type IconClicks struct {
	IconClickFallbackImages *IconClickFallbackImages `xml:"IconClickFallbackImages,omitempty"`
	IconClickThrough        string                   `xml:"IconClickThrough,omitempty"`
	IconClickTracking       []string                 `xml:"IconClickTracking,omitempty"`
}
