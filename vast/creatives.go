package vast

// Creatives contains a collection of creative elements for wrapper ads.
// Holds multiple creative definitions within a wrapper ad structure.
//
// Reference: IAB VAST 4.x Section 2.3.2.1 - Creatives Element
// Link: https://iabtechlab.com/wp-content/uploads/2019/06/VAST_4.2_final_june26.pdf#page=48
type Creatives struct {
	Creative []WrapperCreative `xml:"Creative"`
}

// Creative contains common attributes shared across all creative types.
// Provides foundational properties for all creative implementations.
//
// Reference: IAB VAST 4.x Section 2.3.2.1 - Creative Element
// Link: https://iabtechlab.com/wp-content/uploads/2019/06/VAST_4.2_final_june26.pdf#page=49
type Creative struct {
	Sequence     int    `xml:"sequence,attr,omitempty"`
	APIFramework string `xml:"apiFramework,attr,omitempty"`
	ID           string `xml:"id,attr,omitempty"`
	AdID         string `xml:"adId,attr,omitempty"`
}

// CreativeExtensions contains a collection of creative extensions.
// Provides a container for multiple custom creative extensions.
//
// Reference: IAB VAST 4.x Section 2.3.2.6 - CreativeExtensions Element
// Link: https://iabtechlab.com/wp-content/uploads/2019/06/VAST_4.2_final_june26.pdf#page=49
type CreativeExtensions struct {
	CreativeExtension []CreativeExtension `xml:"CreativeExtension,omitempty"`
}

// CreativeExtension represents custom extensions to the creative element.
// Allows for custom data and functionality beyond the standard VAST specification.
//
// Reference: IAB VAST 4.x Section 2.3.2.6 - CreativeExtensions Element
// Link: https://iabtechlab.com/wp-content/uploads/2019/06/VAST_4.2_final_june26.pdf#page=50
type CreativeExtension struct {
	Items []string `xml:",any"`
	Type  string   `xml:"type,attr,omitempty"`
}
