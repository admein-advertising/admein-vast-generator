// Package vast-errors defines error variables used throughout the VAST library.
// These errors provide standardized error handling for VAST XML operations
// including reading, parsing, and generating VAST documents.
package vast

import "errors"

// ErrReadVAST indicates a failure when reading VAST XML content from an input source.
// This error typically occurs during I/O operations when fetching VAST documents.
var ErrReadVAST = errors.New("there was an issue trying to read the VAST XML")

// ErrMarshalVAST indicates a failure when converting Go structures to VAST XML.
// This error occurs during XML generation when the data cannot be properly serialized.
var ErrMarshalVAST = errors.New("there was an issue trying to marshal the VAST XML")

// ErrUnmarshalVAST indicates a failure when parsing VAST XML into Go structures.
// This error occurs when the XML content is malformed or doesn't conform to VAST schema.
var ErrUnmarshalVAST = errors.New("there was an issue trying to unmarshal the VAST XML")
