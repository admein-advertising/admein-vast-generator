package vast

import "errors"

type NumericBool bool

func (n NumericBool) MarshalText() ([]byte, error) {
	if n {
		return []byte("1"), nil
	}
	return []byte("0"), nil
}

// CData represents XML CDATA sections which preserve text content exactly as written.
// CDATA is used in VAST to ensure URLs, tracking pixels, and other content with special
// characters (like &, <, >) are not parsed as XML markup. This prevents XML parsing
// errors and ensures the content is treated as literal text data.
//
// Reference: IAB VAST 4.x - Used throughout specification for URL and text content
// Link: https://iabtechlab.com/wp-content/uploads/2019/06/VAST_4.2_final_june26.pdf
type CData struct {
	Value string `xml:",cdata"`
}

// AdParameters contains ad-specific parameters passed to the ad creative.
// Link: https://iabtechlab.com/wp-content/uploads/2019/06/VAST_4.2_final_june26.pdf#page=52
type AdParameters struct {
	Value      string      `xml:",chardata"`
	XMLEncoded NumericBool `xml:"xmlEncoded,attr,omitempty"`
}

// Create a new duration with validation.
type Duration string

func NewDuration(value string) (Duration, error) {
	dur := Duration(value)
	if err := dur.ValidateDuration(); err != nil {
		return "", err
	}
	return dur, nil
}

// Duration must be in the format hh:mm:ss, where hh is 00-23, mm is 00-59, and ss is 00-59.
// The total duration must be at least 5 seconds (00:00:05) and cannot be 00:00:00.
func (d Duration) ValidateDuration() error {
	str := string(d)
	if len(str) != 8 {
		return errors.New("Duration must be in the format hh:mm:ss")
	}
	if str[2] != ':' || str[5] != ':' {
		return errors.New("Duration must be in the format hh:mm:ss")
	}
	for i, char := range str {
		if i == 2 || i == 5 {
			continue
		}
		if char < '0' || char > '9' {
			return errors.New("Duration must be in the format hh:mm:ss")
		}
	}

	// Validate ranges: hh (00-23), mm (00-59), ss (00-59)
	var numericValue int
	hours := numericValue / 10000
	minutes := (numericValue / 100) % 100
	seconds := numericValue % 100

	if hours < 0 || hours > 23 {
		return errors.New("hours must be between 00 and 23")
	}
	if minutes < 0 || minutes > 59 {
		return errors.New("minutes must be between 00 and 59")
	}
	if seconds < 0 || seconds > 59 {
		return errors.New("seconds must be between 00 and 59")
	}

	// Get all values and convert them to a single integer
	for i, char := range str {
		if i == 2 || i == 5 {
			continue
		}
		numericValue = numericValue*10 + int(char-'0')
	}
	// Check if the total duration is zero
	if numericValue == 0 {
		return errors.New("Duration cannot be 00:00:00")
	}
	// Check if the total duration is less than 5 seconds (000005)
	if numericValue < 5 {
		return errors.New("Duration must be at least 5 seconds")
	}
	return nil
}

// XPosition constraints ([0-9]*|left|right).
type XPosition string

// XPosition constants for named positions
const (
	XPositionLeft  XPosition = "left"
	XPositionRight XPosition = "right"
)

// NewXPosition creates a new XPosition with validation.
func NewXPosition(value string) (XPosition, error) {
	pos := XPosition(value)
	if err := pos.Validate(); err != nil {
		return "", err
	}
	return pos, nil
}

// Validate checks if the XPosition matches the required pattern ([0-9]*|left|right).
func (x XPosition) Validate() error {
	str := string(x)
	if str == "left" || str == "right" {
		return nil
	}

	// Check if it's a numeric value
	for _, char := range str {
		if char < '0' || char > '9' {
			return errors.New("XPosition must be a number, 'left', or 'right'")
		}
	}
	return nil
}

// YPosition constraints ([0-9]*|top|bottom).
type YPosition string

// YPosition constants for named positions
const (
	YPositionTop    YPosition = "top"
	YPositionBottom YPosition = "bottom"
)

// NewYPosition creates a new YPosition with validation.
func NewYPosition(value string) (YPosition, error) {
	pos := YPosition(value)
	if err := pos.Validate(); err != nil {
		return "", err
	}
	return pos, nil
}

// Validate checks if the YPosition matches the required pattern ([0-9]*|top|bottom).
func (y YPosition) Validate() error {
	str := string(y)
	if str == "top" || str == "bottom" {
		return nil
	}

	// Check if it's a numeric value
	for _, char := range str {
		if char < '0' || char > '9' {
			return errors.New("YPosition must be a number, 'top', or 'bottom'")
		}
	}
	return nil
}
