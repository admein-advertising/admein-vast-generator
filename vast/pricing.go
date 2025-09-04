package vast

// Currency represents a three-letter ISO currency code for pricing information.
// RegEx Pattern: [a-zA-Z]{3}.
//
// Reference: IAB VAST 4.x Section 2.3.1.4 - Pricing Element
// Link: https://iabtechlab.com/wp-content/uploads/2019/06/VAST_4.2_final_june26.pdf#page=45
type Currency string

// Model represents the pricing model used for the advertisement.
// Reference: IAB VAST 4.x Section 2.3.1.4 - Pricing Element
// Link: https://iabtechlab.com/wp-content/uploads/2019/06/VAST_4.2_final_june26.pdf#page=45
type Model string

// Pricing contains cost information for the advertisement impression.
// Provides pricing model and cost details used by real-time bidding (RTB) systems.
//
// Reference: IAB VAST 4.x Section 2.3.1.4 - Pricing Element
// Link: https://iabtechlab.com/wp-content/uploads/2019/06/VAST_4.2_final_june26.pdf#page=45
type Pricing struct {
	Value    float64  `xml:",cdata"`
	Model    Model    `xml:"model,attr"`
	Currency Currency `xml:"currency,attr"`
}

// Identifies the pricing model as one of: CPM, CPC, CPE, or CPV
// Reference: IAB VAST 4.x Section 2.3.1.4 - Pricing Element
// Link: https://iabtechlab.com/wp-content/uploads/2019/06/VAST_4.2_final_june26.pdf#page=45
const CPCModel Model = "CPC"
const CPMModel Model = "CPM"
const CPEModel Model = "CPE"
const CPVModel Model = "CPV"
