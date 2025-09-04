package vast

// ClickThrough represents the URL to redirect users when they click on the ad.
// Contains the landing page URL where users are taken after clicking the ad.
//
// Reference: IAB VAST 4.x Section 2.3.2.2 - ClickThrough Element
// Link: https://iabtechlab.com/wp-content/uploads/2019/06/VAST_4.2_final_june26.pdf#page=57
type ClickThrough struct {
	Value string `xml:",cdata"`
	ID    string `xml:"id,attr,omitempty"`
}
