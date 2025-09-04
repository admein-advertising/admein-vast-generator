package vast

import (
	"errors"
	"regexp"
	"strings"
)

// Linear represents the base structure for linear video advertisements.
// Contains common elements shared between inline and wrapper linear ads.
//
// Reference: IAB VAST 4.x Section 2.3.2.2 - Linear Element
// Link: https://iabtechlab.com/wp-content/uploads/2019/06/VAST_4.2_final_june26.pdf#page=51
type Linear struct {
	Icons          *Icons          `xml:"Icons,omitempty"`
	TrackingEvents *TrackingEvents `xml:"TrackingEvents,omitempty"`
	SkipOffset     SkipOffset      `xml:"skipoffset,attr,omitempty"`
}

// LinearInLine represents a complete linear ad with all media files and interactions.
// Extends Linear with media files, duration, and click tracking for inline ads.
type LinearInLine struct {
	Linear
	AdParameters *AdParameters `xml:"AdParameters,omitempty"`
	Duration     Duration      `xml:"Duration"`
	MediaFiles   MediaFiles    `xml:"MediaFiles"`
	VideoClicks  *VideoClicks  `xml:"VideoClicks,omitempty"`
}

// VideoClicks contains click-through and tracking URLs for linear video ads.
// Manages user interaction tracking and click-through behavior for video content.
//
// Reference: IAB VAST 4.x Section 2.3.2.2 - VideoClicks Element
// Link: https://iabtechlab.com/wp-content/uploads/2019/06/VAST_4.2_final_june26.pdf#page=57
type VideoClicks struct {
	ClickThrough  ClickThrough `xml:"ClickThrough,omitempty"`
	ClickTracking []CData      `xml:"ClickTracking,omitempty"`
	CustomClick   []string     `xml:"CustomClick,omitempty"`
}

// SkipOffset defines when a skip button becomes available during ad playback.
// Must match the pattern `(\d{2}:[0-5]\d:[0-5]\d(\.\d\d\d)?|1?\d?\d(\.?\d)*%)`.
//
// Reference: IAB VAST 4.x Section 2.3.2.2 - Linear Element
// Link: https://iabtechlab.com/wp-content/uploads/2019/06/VAST_4.2_final_june26.pdf#page=51
type SkipOffset string

// NewSkipOffset creates a new SkipOffset with validation.
func NewSkipOffset(value string) (SkipOffset, error) {
	offset := SkipOffset(value)
	if err := offset.Validate(); err != nil {
		return "", err
	}
	return offset, nil
}

// Validate checks if the SkipOffset matches the required pattern.
func (s SkipOffset) Validate() error {
	str := string(s)
	if str == "" {
		return nil // empty is valid
	}

	// Check for percentage format: 1?\d?\d(\.?\d)*%
	if strings.HasSuffix(str, "%") {
		percentStr := strings.TrimSuffix(str, "%")
		if matched, _ := regexp.MatchString(`^1?\d?\d(\.?\d)*$`, percentStr); matched {
			return nil
		}
	}

	// Check for time format: \d{2}:[0-5]\d:[0-5]\d(\.\d\d\d)?
	if matched, _ := regexp.MatchString(`^\d{2}:[0-5]\d:[0-5]\d(\.\d\d\d)?$`, str); matched {
		return nil
	}

	return errors.New("SkipOffset must match pattern (HH:MM:SS[.fff] or percentage)")
}
