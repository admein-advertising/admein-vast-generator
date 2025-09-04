package vast

// ClosedCaptionFiles contains a collection of closed caption file resources.
// Provides accessibility support through closed captioning for video ads.
//
// Reference: IAB VAST 4.x Section 2.3.2.3 - ClosedCaptionFiles Element
// Link: https://iabtechlab.com/wp-content/uploads/2019/06/VAST_4.2_final_june26.pdf#page=56
type ClosedCaptionFiles struct {
	ClosedCaptionFile []ClosedCaptionFile `xml:"ClosedCaptionFile,omitempty"`
}

// ClosedCaptionFile represents a single closed caption file resource.
// Contains the URI and metadata for closed caption content to support accessibility.
//
// Reference: IAB VAST 4.x Section 2.3.2.3 - ClosedCaptionFile Element
// Link: https://iabtechlab.com/wp-content/uploads/2019/06/VAST_4.2_final_june26.pdf#page=56
type ClosedCaptionFile struct {
	Value    string `xml:",cdata"`
	Type     string `xml:"type,attr,omitempty"`
	Language string `xml:"language,attr,omitempty"`
}
