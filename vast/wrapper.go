package vast

// Wrapper represents a VAST ad that references another VAST document for its content.
// Contains tracking and metadata while delegating actual ad content to another VAST tag.
//
// Reference: IAB VAST 4.x Section 2.3.1.2 - Wrapper Element
// Link: https://iabtechlab.com/wp-content/uploads/2019/06/VAST_4.2_final_june26.pdf#page=29
type Wrapper struct {
	AdDefinition
	AdVerifications          *AdVerifications      `xml:"AdVerifications,omitempty"`
	BlockedAdCategories      []BlockedAdCategories `xml:"BlockedAdCategories,omitempty"`
	Creatives                *Creatives            `xml:"Creatives,omitempty"`
	VASTAdTagURI             CData                 `xml:"VASTAdTagURI"`
	FollowAdditionalWrappers NumericBool           `xml:"followAdditionalWrappers,attr,omitempty"`
	AllowMultipleAds         NumericBool           `xml:"allowMultipleAds,attr,omitempty"`
	FallbackOnNoAd           NumericBool           `xml:"fallbackOnNoAd,attr,omitempty"`
}

// WrapperCreative represents a creative element within a wrapper ad structure.
// Provides tracking and extension points for creatives in wrapper ads.
type WrapperCreative struct {
	Creative
	CompanionAds *CompanionAds  `xml:"CompanionAds,omitempty"`
	Linear       *LinearWrapper `xml:"Linear,omitempty"`
	NonLinearAds *NonLinearAds  `xml:"NonLinearAds,omitempty"`
	ID           string         `xml:"id,attr,omitempty"`
}

// LinearWrapper represents linear creative content within a wrapper ad.
// Extends Linear with wrapper-specific tracking and click handling.
type LinearWrapper struct {
	Linear
	VideoClicks *VideoClicks `xml:"VideoClicks,omitempty"`
}
