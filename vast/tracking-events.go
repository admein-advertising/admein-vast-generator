package vast

// Event represents the type of tracking event that triggers URL calls.
//
// Reference: IAB VAST 4.x Section 2.3.2.1 - TrackingEvents Element
// Link: https://iabtechlab.com/wp-content/uploads/2019/06/VAST_4.2_final_june26.pdf#page=68
//
// 4.1 Updates
// Link: https://iabtechlab.com/wp-content/uploads/2019/06/VAST_4.2_final_june26.pdf#page=11
//
// 4.2 Updates
// Link: https://iabtechlab.com/wp-content/uploads/2019/06/VAST_4.2_final_june26.pdf#page=12
type TrackingEvent string

const (
	MuteEvent                TrackingEvent = "mute"
	UnmuteEvent              TrackingEvent = "unmute"
	PauseEvent               TrackingEvent = "pause"
	ResumeEvent              TrackingEvent = "resume"
	RewindEvent              TrackingEvent = "rewind"
	SkipEvent                TrackingEvent = "skip"
	PlayerExpandEvent        TrackingEvent = "playerExpand"
	PlayerCollapseEvent      TrackingEvent = "playerCollapse"
	LoadedEvent              TrackingEvent = "loaded"
	StartEvent               TrackingEvent = "start"
	FirstQuartileEvent       TrackingEvent = "firstQuartile"
	MidpointEvent            TrackingEvent = "midpoint"
	ThirdQuartileEvent       TrackingEvent = "thirdQuartile"
	CompleteEvent            TrackingEvent = "complete"
	ProgressEvent            TrackingEvent = "progress"
	CloseLinearEvent         TrackingEvent = "closeLinear"
	CreativeViewEvent        TrackingEvent = "creativeView"
	AcceptInvitationEvent    TrackingEvent = "acceptInvitation"
	AdExpandEvent            TrackingEvent = "adExpand"
	AdCollapseEvent          TrackingEvent = "adCollapse"
	MinimizeEvent            TrackingEvent = "minimize"
	CloseEvent               TrackingEvent = "close"
	OverlayViewDurationEvent TrackingEvent = "overlayViewDuration"
	OtherAdInteraction       TrackingEvent = "otherAdInteraction"
	InteractiveStart         TrackingEvent = "interactiveStart"
)

// TrackingEvents contains a collection of tracking URLs for ad measurement.
// Provides URLs to be called when specific ad events occur during playback.
//
// Reference: IAB VAST 4.x Section 2.3.2.1 - TrackingEvents Element
// Link: https://iabtechlab.com/wp-content/uploads/2019/06/VAST_4.2_final_june26.pdf#page=25
type TrackingEvents struct {
	Tracking []Tracking `xml:"Tracking,omitempty"`
}

// Offset defines when a tracking event should fire during ad playback.
// Must match the pattern `(\d{2}:[0-5]\d:[0-5]\d(\.\d\d\d)?|1?\d?\d(\.?\d)*%)`.
//
// Reference: IAB VAST 4.x Section 2.3.2.1 - Tracking Element
// Link: https://iabtechlab.com/wp-content/uploads/2019/06/VAST_4.2_final_june26.pdf#page=26
type Offset string

// Tracking represents a single tracking URL associated with an ad event.
// Reference: IAB VAST 4.x Section 2.3.2.1 - Tracking Element
// Link: https://iabtechlab.com/wp-content/uploads/2019/06/VAST_4.2_final_june26.pdf#page=26
// Contains the URL and event details for ad measurement and analytics.
type Tracking struct {
	Value  string `xml:",cdata"`
	Event  string `xml:"event,attr"`
	Offset Offset `xml:"offset,attr,omitempty"`
}

// TrackingEventsVerification contains tracking events specific to ad verification.
// Reference: IAB VAST 4.x Section 2.3.4.1 - AdVerifications Element
// Link: https://iabtechlab.com/wp-content/uploads/2019/06/VAST_4.2_final_june26.pdf#page=41
// Provides verification-specific tracking for fraud detection and measurement.
type TrackingEventsVerification struct {
	Tracking []Tracking `xml:"Tracking,omitempty"`
}

// Impression represents an impression tracking URL for ad delivery confirmation.
// Reference: IAB VAST 4.x Section 2.3.1.3 - Impression Element
// Link: https://iabtechlab.com/wp-content/uploads/2019/06/VAST_4.2_final_june26.pdf#page=16
// Contains URLs called when the ad is displayed to confirm impression delivery.
type Impression struct {
	Value string `xml:",cdata"`
	ID    string `xml:"id,attr,omitempty"`
}
