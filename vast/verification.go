package vast

// Verification represents third-party verification resources for ad measurement and fraud detection.
// Contains verification scripts and parameters for measuring ad delivery and viewability.
//
// Reference: IAB VAST 4.x Section 2.3.4.1 - Verification Element
// Link: https://iabtechlab.com/wp-content/uploads/2019/06/VAST_4.2_final_june26.pdf#page=75
type Verification struct {
	ExecutableResource     []ExecutableResource        `xml:"ExecutableResource,omitempty"`
	JavaScriptResource     []JavaScriptResource        `xml:"JavaScriptResource,omitempty"`
	TrackingEvents         *TrackingEventsVerification `xml:"TrackingEvents,omitempty"`
	VerificationParameters string                      `xml:"VerificationParameters,omitempty"`
	Vendor                 string                      `xml:"vendor,attr,omitempty"`
}
