package vast

// Version represents the VAST specification version number.
// Indicates which version of the VAST specification the document conforms to.
//
// Reference: IAB VAST 4.x Section 2.3.1 - VAST Element
// Link: https://iabtechlab.com/wp-content/uploads/2019/06/VAST_4.2_final_june26.pdf
type Version string

const (
	Version20 Version = "2.0"
	Version30 Version = "3.0"
	Version40 Version = "4.0"
	Version41 Version = "4.1"
	Version42 Version = "4.2"
	// New macro only support
	Version43 Version = "4.3"
)
