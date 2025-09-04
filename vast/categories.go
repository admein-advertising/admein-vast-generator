package vast

// BlockedAdCategories represents categories of ads that should be blocked from serving.
// Used by publishers to prevent certain types of ads from being displayed.
//
// Reference: IAB VAST 4.x Section 3.19.2 - BlockedAdCategories Element
// Link: https://iabtechlab.com/wp-content/uploads/2019/06/VAST_4.2_final_june26.pdf#page=80
type BlockedAdCategories struct {
	Value     string `xml:",chardata"`
	Authority string `xml:"authority,attr,omitempty"`
}

// Category represents the category classification of an advertisement.
// Allows for content categorization using industry standard taxonomies.
//
// Reference: IAB VAST 4.x Section 3.4.5 - Category Element
// Link: https://iabtechlab.com/wp-content/uploads/2019/06/VAST_4.2_final_june26.pdf#page=43
type Category struct {
	Value     string `xml:",chardata"`
	Authority string `xml:"authority,attr"`
}
