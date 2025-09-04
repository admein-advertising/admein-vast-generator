package vast

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io"
)

// VAST represents the root element of a VAST document containing ads and metadata.
// Reference: IAB VAST 4.x Section 2.3.1 - VAST Element
// Link: https://iabtechlab.com/wp-content/uploads/2019/06/VAST_4.2_final_june26.pdf#page=15
// The top-level container for all VAST advertisement content and configuration.
type VAST struct {
	Ad      []Ad    `xml:"Ad,omitempty"`
	Error   []CData `xml:"Error,omitempty"`
	Version Version `xml:"version,attr"`

	// VAST Namespace. This makes up the first node of the VAST XML document.
	XMLNS                        Namespace `xml:"xmlns,attr,omitempty"`
	XMLNSXsi                     string    `xml:"xmlns:xsi,attr,omitempty"`
	XsiNoNamespaceSchemaLocation string    `xml:"xsi:noNamespaceSchemaLocation,attr,omitempty"`
}

// New creates a new instance of VAST with default values.
// To use this function, simply call it to get a new VAST instance.
func New() *VAST {
	return &VAST{
		Version:                      Version30,
		XMLNSXsi:                     string(VASTNamespace),
		XsiNoNamespaceSchemaLocation: string(VASTSchemaLocation),
	}
}

// Bytes returns the VAST XML as a byte slice.
// This function is useful for getting the raw XML representation of the VAST object.
// Using xml.MarshalIndent for pretty-printing the XML.
func (v *VAST) Bytes() ([]byte, error) {
	xmlContent, err := xml.MarshalIndent(v, "", "  ")
	if err != nil {
		return nil, errors.Join(ErrMarshalVAST, err)
	}
	// The byte representation of the VAST XML.
	var buf bytes.Buffer
	buf.WriteString(xml.Header)
	buf.Write(xmlContent)
	return buf.Bytes(), nil
}

// Read creates a new instance of VAST and reads the content from an io.ReadCloser.
// This function is useful for decoding VAST XML directly from a stream.
// It avoids loading the entire XML document into memory at once.
func Read(reader io.ReadCloser) (*VAST, error) {
	defer reader.Close()

	vast := &VAST{}
	decoder := xml.NewDecoder(reader)
	if err := decoder.Decode(vast); err != nil {
		return nil, errors.Join(ErrUnmarshalVAST, err)
	}

	return vast, nil
}
