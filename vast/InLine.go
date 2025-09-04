package vast

// AdVerifications contains verification resources for measuring ad delivery and viewability.
// Provides third-party verification scripts for ad measurement and fraud detection.
//
// Reference: IAB VAST 4.x Section 2.3.4.1 - AdVerifications Element
// Link: https://iabtechlab.com/wp-content/uploads/2019/06/VAST_4.2_final_june26.pdf#page=74
type AdVerifications struct {
	Verification []Verification `xml:"Verification,omitempty"`
}

// Survey represents an optional survey URL for gathering viewer feedback.
// Allows collection of viewer feedback through external survey mechanisms.
//
// Reference: IAB VAST 4.x Section 2.3.1.8 - Survey Element
// Link: https://iabtechlab.com/wp-content/uploads/2019/06/VAST_4.2_final_june26.pdf#page=45
type Survey struct {
	Value string `xml:",chardata"`
	Type  string `xml:"type,attr,omitempty"`
}

// UniversalAdID provides a unique identifier for the advertisement across systems.
// Enables cross-platform ad tracking and deduplication using industry standard registries.
//
// Reference: IAB VAST 4.x Section 2.3.2.7 - UniversalAdId Element
// Link: https://iabtechlab.com/wp-content/uploads/2019/06/VAST_4.2_final_june26.pdf#page=49
type UniversalAdID struct {
	Value      string `xml:",chardata"`
	IDRegistry string `xml:"idRegistry,attr"`
}

// InLine represents a complete ad definition with all creative content included.
// Contains the full ad content without requiring additional VAST requests.
//
// Reference: IAB VAST 4.x Section 2.3.1.2 - InLine Element
// Link: https://iabtechlab.com/wp-content/uploads/2019/06/VAST_4.2_final_june26.pdf#page=41
type InLine struct {
	AdDefinition
	AdServingID     string           `xml:"AdServingId"`
	AdTitle         string           `xml:"AdTitle"`
	AdVerifications *AdVerifications `xml:"AdVerifications,omitempty"`
	Advertiser      string           `xml:"Advertiser,omitempty"`
	Category        []Category       `xml:"Category,omitempty"`
	Creatives       InLineCreatives  `xml:"Creatives"`
	Description     *CData           `xml:"Description,omitempty"`
	Expires         int              `xml:"Expires,omitempty"`
	Survey          *Survey          `xml:"Survey,omitempty"` // Deprecated since VAST 4.1
}

// InLineCreatives contains a collection of creative elements for inline ads.
// Holds multiple creative definitions within an inline ad structure.
type InLineCreatives struct {
	Creative []InLineCreative `xml:"Creative"`
}

// InLineCreative represents a single creative element within an inline ad.
// Defines individual creative content including linear, non-linear, and companion ads.
type InLineCreative struct {
	Creative
	CompanionAds       *CompanionAds       `xml:"CompanionAds,omitempty"`
	CreativeExtensions *CreativeExtensions `xml:"CreativeExtensions,omitempty"`
	Linear             *LinearInLine       `xml:"Linear,omitempty"`
	NonLinearAds       *NonLinearAds       `xml:"NonLinearAds,omitempty"`
	UniversalAdID      []UniversalAdID     `xml:"UniversalAdId"`
	ID                 string              `xml:"id,attr,omitempty"`
}
