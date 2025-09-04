package vast

// Extensions contains a collection of custom extension elements.
// Provides a container for multiple vendor-specific extensions.
//
// Reference: IAB VAST 4.x Section 2.3.1.7 - Extensions Element
// Link: https://iabtechlab.com/wp-content/uploads/2019/06/VAST_4.2_final_june26.pdf#page=78
type Extensions struct {
	Extension []Extension `xml:"Extension,omitempty"`
}

// Extension represents a custom extension element for vendor-specific data.
// Allows advertisers and vendors to include custom data beyond the VAST specification.
//
// Reference: IAB VAST 4.x Section 2.3.1.7 - Extensions Element
// Link: https://iabtechlab.com/wp-content/uploads/2019/06/VAST_4.2_final_june26.pdf#page=78
type Extension struct {
	Value string `xml:",innerxml"`
	Type  string `xml:"type,attr,omitempty"`
}
