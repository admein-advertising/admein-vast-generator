package vast

// Delivery specifies the method of media content delivery to the player.
//
// Reference: IAB VAST 4.x Section 2.3.2.3 - MediaFile Element
// Link: https://iabtechlab.com/wp-content/uploads/2019/06/VAST_4.2_final_june26.pdf#page=52
type Delivery string

const (
	StreamingDelivery   Delivery = "streaming"
	ProgressiveDelivery Delivery = "progressive"
)

// MediaFiles contains media file resources and related content for linear ads.
// Provides video files, closed captions, and interactive content for ad playback.
//
// Reference: IAB VAST 4.x Section 2.3.2.3 - MediaFiles Element
// Link: https://iabtechlab.com/wp-content/uploads/2019/06/VAST_4.2_final_june26.pdf#page=52
type MediaFiles struct {
	ClosedCaptionFiles      *ClosedCaptionFiles       `xml:"ClosedCaptionFiles,omitempty"`
	MediaFile               []MediaFile               `xml:"MediaFile"`
	Mezzanine               []Mezzanine               `xml:"Mezzanine,omitempty"`
	InteractiveCreativeFile []InteractiveCreativeFile `xml:"InteractiveCreativeFile,omitempty"`
}

// MediaFile represents a single media file resource with delivery and encoding details.
// Contains video content URL and technical specifications for ad playback.
//
// Reference: IAB VAST 4.x Section 2.3.2.3 - MediaFile Element
// Link: https://iabtechlab.com/wp-content/uploads/2019/06/VAST_4.2_final_june26.pdf#page=53
type MediaFile struct {
	Value               string      `xml:",cdata"`
	ID                  string      `xml:"id,attr,omitempty"`
	Delivery            Delivery    `xml:"delivery,attr"`
	Type                string      `xml:"type,attr"`
	Width               int         `xml:"width,attr"`
	Height              int         `xml:"height,attr"`
	Codec               string      `xml:"codec,attr,omitempty"`
	Bitrate             int         `xml:"bitrate,attr,omitempty"`
	MinBitrate          int         `xml:"minBitrate,attr,omitempty"`
	MaxBitrate          int         `xml:"maxBitrate,attr,omitempty"`
	Scalable            NumericBool `xml:"scalable,attr,omitempty"`
	MaintainAspectRatio NumericBool `xml:"maintainAspectRatio,attr,omitempty"`
	FileSize            int         `xml:"fileSize,attr,omitempty"`
	MediaType           string      `xml:"mediaType,attr,omitempty"`
	APIFramework        string      `xml:"apiFramework,attr,omitempty"`
}

// Mezzanine represents a high-quality source file for transcoding purposes.
// Provides high-quality source content for server-side transcoding and optimization.
// Also known as ad stitching in SSAI.
//
// Reference: IAB VAST 4.x Section 2.3.2.3 - Mezzanine Element
// Link: https://iabtechlab.com/wp-content/uploads/2019/06/VAST_4.2_final_june26.pdf#page=55
type Mezzanine struct {
	Value     string   `xml:",cdata"`
	ID        string   `xml:"id,attr,omitempty"`
	Delivery  Delivery `xml:"delivery,attr"`
	Type      string   `xml:"type,attr"`
	Width     int      `xml:"width,attr"`
	Height    int      `xml:"height,attr"`
	Codec     string   `xml:"codec,attr,omitempty"`
	FileSize  int      `xml:"fileSize,attr,omitempty"`
	MediaType string   `xml:"mediaType,attr,omitempty"`
}
