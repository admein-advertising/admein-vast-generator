# VAST Generator for Go

A comprehensive Go library for generating and parsing IAB VAST (Video Ad Serving Template) XML documents. This package provides full support for VAST 4.2 specification with type-safe structures and validation.

## Features

- **Full VAST 4.2 Support** - Complete implementation of IAB VAST 4.2 specification
- **Type Safety** - Strongly typed Go structures with XML marshaling/unmarshaling
- **Pattern Validation** - Built-in validation for VAST-specific patterns and constraints
- **Memory Efficient** - Stream-based XML parsing to handle large VAST documents
- **IAB Compliant** - Follows official IAB VAST guidelines and best practices
- **Documentation** - Comprehensive documentation with IAB specification references
- **Examples** - Ready-to-use examples for common VAST use cases

## Installation

```bash
go get github.com/admein-advertising/admein-vast-generator
```

## Quick Start

### Creating a Simple VAST Document

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/admein-advertising/admein-vast-generator/vast"
)

func main() {
    // Create a new VAST document
    v := vast.New()
    
    // Add an inline ad
    v.Ad = append(v.Ad, vast.Ad{
        ID:       "ad-001",
        Sequence: 1,
        InLine: &vast.InLine{
            AdTitle: "Sample Video Ad",
            Creatives: vast.InLineCreatives{
                Creative: []vast.InLineCreative{
                    {
                        Creative: vast.Creative{Sequence: 1},
                        Linear: &vast.LinearInLine{
                            Duration: "00:00:30",
                            MediaFiles: vast.MediaFiles{
                                MediaFile: []vast.MediaFile{
                                    {
                                        Value:    "https://example.com/video.mp4",
                                        Delivery: vast.ProgressiveDelivery,
                                        Type:     "video/mp4",
                                        Width:    1920,
                                        Height:   1080,
                                    },
                                },
                            },
                        },
                    },
                },
            },
        },
    })
    
    // Generate XML
    xmlBytes, err := v.Bytes()
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println(string(xmlBytes))
}
```

### Parsing Existing VAST XML

```go
package main

import (
    "log"
    "net/http"
    
    "github.com/admein-advertising/admein-vast-generator/vast"
)

func main() {
    // Fetch VAST from URL
    resp, err := http.Get("https://example.com/vast.xml")
    if err != nil {
        log.Fatal(err)
    }
    
    // Parse VAST document
    vastDoc, err := vast.Read(resp.Body)
    if err != nil {
        log.Fatal(err)
    }
    
    // Access parsed data
    for _, ad := range vastDoc.Ad {
        if ad.InLine != nil {
            fmt.Printf("Ad Title: %s\n", ad.InLine.AdTitle)
        }
    }
}
```

## VAST Specification Support

This library implements the complete [IAB VAST 4.2 specification](https://iabtechlab.com/wp-content/uploads/2019/06/VAST_4.2_final_june26.pdf) including:

### Ad Types
- **InLine Ads** - Complete ad content within the VAST document
- **Wrapper Ads** - References to external VAST documents with additional tracking
- **Ad Pods** - Sequential ad playback with multiple ads

### Creative Types
- **Linear Ads** - Video ads that play before, during, or after content
- **Non-Linear Ads** - Overlay ads displayed on top of content
- **Companion Ads** - Banner ads displayed alongside video content

### Advanced Features
- **Ad Verification** - Third-party verification scripts for fraud detection
- **Tracking Events** - Comprehensive event tracking for analytics
- **Interactive Creatives** - Support for interactive ad experiences
- **Closed Captions** - Accessibility support with caption files
- **Universal Ad IDs** - Cross-platform ad identification

## API Reference

### Core Types

#### VASTXML
The root VAST document structure:

```go
type VASTXML struct {
    Ad      []Ad    // Collection of advertisements
    Error   []CData // Error tracking URLs
    Version Version // VAST specification version
}
```

#### Ad
Individual advertisement container:

```go
type Ad struct {
    InLine        *InLine     // Inline ad content
    Wrapper       *Wrapper    // Wrapper ad reference
    ID            string      // Unique ad identifier
    Sequence      int         // Ad sequence in pod
    ConditionalAd NumericBool // Conditional ad flag
    AdType        AdType      // Ad content type
}
```

### Validation

The library includes built-in validation for VAST-specific patterns:

```go
// Skip offset validation
skipOffset, err := vast.NewSkipOffset("00:00:05")
if err != nil {
    log.Fatal("Invalid skip offset format")
}

// Duration validation
duration, err := vast.NewDuration("00:00:06")
if err != nil {
    log.Fatal("Invalid duration format")
}

// Position validation
xPos, errX := vast.NewXPosition("left")
yPos, errY := vast.NewYPosition("top")
if errX != nil {
    log.Fatal("Invalid X position format")
}
if errY != nil {
    log.Fatal("Invalid Y position format")
}
```

## Examples

### Linear Video Ad with Tracking

```go
v := vast.New()
v.Ad = append(v.Ad, vast.Ad{
    ID: "linear-ad-001",
    InLine: &vast.InLine{
        AdTitle: "Brand Campaign Video",
        Creatives: vast.InLineCreatives{
            Creative: []vast.InLineCreative{
                {
                    Linear: &vast.LinearInLine{
                        Duration: "00:00:15",
                        Linear: vast.Linear{
                            TrackingEvents: &vast.TrackingEvents{
                                Tracking: []vast.Tracking{
                                    {
                                        Event: string(vast.StartEvent),
                                        Value: "https://analytics.example.com/start",
                                    },
                                    {
                                        Event: string(vast.CompleteEvent),
                                        Value: "https://analytics.example.com/complete",
                                    },
                                },
                            },
                        },
                        VideoClicks: &vast.VideoClicks{
                            ClickThrough: vast.ClickThrough{
                                Value: "https://advertiser.example.com/landing",
                            },
                        },
                        MediaFiles: vast.MediaFiles{
                            MediaFile: []vast.MediaFile{
                                {
                                    Value:    "https://cdn.example.com/ad.mp4",
                                    Delivery: vast.ProgressiveDelivery,
                                    Type:     "video/mp4",
                                    Width:    1280,
                                    Height:   720,
                                    Bitrate:  2000,
                                },
                            },
                        },
                    },
                },
            },
        },
    },
})
```

### Wrapper Ad with Verification

```go
v := vast.New()
v.Ad = append(v.Ad, vast.Ad{
    ID: "wrapper-ad-001",
    Wrapper: &vast.Wrapper{
        AdDefinition: vast.AdDefinition{
            AdSystem: vast.AdSystem{
                Value:   "Example Ad Server",
                Version: "1.0",
            },
            Impression: []vast.Impression{
                {Value: "https://impression.example.com/track"},
            },
        },
        VASTAdTagURI: vast.CData{
            Value: "https://downstream.example.com/vast.xml",
        },
        AdVerifications: &vast.AdVerifications{
            Verification: []vast.Verification{
                {
                    Vendor: "verification-vendor",
                    JavaScriptResource: []vast.JavaScriptResource{
                        {
                            Value: "https://verification.example.com/script.js",
                        },
                    },
                },
            },
        },
    },
})
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -am 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Standards Compliance

This library is built to comply with:

- [IAB VAST 4.2 Specification](https://iabtechlab.com/standards/vast/)
- [IAB Tech Lab Guidelines](https://iabtechlab.com/)
- Go best practices and idioms
- XML Schema validation

## License

This project is licensed under the MIT License.

## Keywords

AdMeIn VAST Testing, VAST, Video Ad Serving Template, IAB, Go, Golang, Video Advertising, Ad Server, Programmatic Advertising, XML, Ad Tech, Digital Advertising, Video Ads, Linear Ads, Non-Linear Ads, Wrapper Ads, Ad Verification, Tracking Events, Media Files, Ad Pod, Companion Ads, Interactive Ads, VPAID, Ad Serving, Real-time Bidding, RTB, Header Bidding, Programmatic Media, Video Marketing, Ad Analytics, Impression Tracking, Click Tracking, Viewability, Ad Fraud Prevention

## AdMeIn

AdMeIn is a dedicated Video Ad testing platform made for developers and AdOps teams to test their VAST implementations with advanced feature results.
- [AdMeIn VAST Tester / Validator](https://admein.in/vast-tester)
- [AdMeIn VAST examples](https://admein.in/help/article/how-to-test-vast-url)

---

Made with ❤️ for the digital advertising community
Author of the AdMeIn project [@vince-scarpa](https://github.com/vince-scarpa)