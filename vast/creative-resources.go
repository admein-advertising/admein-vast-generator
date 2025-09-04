package vast

// CreativeResource contains the different types of resources that can be used in a creative.
// Provides HTML, IFrame, or static content for displaying ad creatives.
//
// Reference: IAB VAST 4.x Section 2.3.2.4 - Creative Resources
// Link: https://iabtechlab.com/wp-content/uploads/2019/06/VAST_4.2_final_june26.pdf#page=30
type CreativeResource struct {
	HTMLResource   []HTMLResource   `xml:"HTMLResource,omitempty"`
	IFrameResource []IFrameResource `xml:"IFrameResource,omitempty"`
	StaticResource []StaticResource `xml:"StaticResource,omitempty"`
}

// InteractiveCreativeFile represents an interactive creative file resource.
// Contains interactive content that can respond to user interactions.
//
// Reference: IAB VAST 4.x Section 2.3.2.5 - InteractiveCreativeFile Element
// Link: https://iabtechlab.com/wp-content/uploads/2019/06/VAST_4.2_final_june26.pdf#page=31
type InteractiveCreativeFile struct {
	Value            string      `xml:",cdata"`
	Type             string      `xml:"type,attr,omitempty"`
	APIFramework     string      `xml:"apiFramework,attr,omitempty"`
	VariableDuration NumericBool `xml:"variableDuration,attr,omitempty"`
}

// ExecutableResource represents an executable resource for ad verification.
// Contains executable verification code for measuring ad performance.
//
// Reference: IAB VAST 4.x Section 2.3.4.3 - ExecutableResource Element
// Link: https://iabtechlab.com/wp-content/uploads/2019/06/VAST_4.2_final_june26.pdf#page=43
type ExecutableResource struct {
	Value        string `xml:",cdata"`
	APIFramework string `xml:"apiFramework,attr,omitempty"`
	Type         string `xml:"type,attr,omitempty"`
}

// StaticResource represents a static creative resource like images or other media files.
// Contains static media content for display in ad creatives.
//
// Reference: IAB VAST 4.x Section 2.3.2.4 - StaticResource Element
// Link: https://iabtechlab.com/wp-content/uploads/2019/06/VAST_4.2_final_june26.pdf#page=30
type StaticResource struct {
	Value        string `xml:",cdata"`
	CreativeType string `xml:"creativeType,attr"`
}

// HTMLResource represents an HTML creative resource.
// Contains HTML content for displaying ad creatives.
//
// Reference: IAB VAST 4.x Section 2.3.2.4 - HTMLResource Element
// Link: https://iabtechlab.com/wp-content/uploads/2019/06/VAST_4.2_final_june26.pdf#page=72
type HTMLResource struct {
	Value string `xml:",cdata"`
}

// JavaScriptResource represents a JavaScript resource for ad verification or interactive functionality.
// Used for ad verification scripts and interactive ad frameworks.
//
// Reference: IAB VAST 4.x Section 2.3.4.2 - JavaScriptResource Element
// Link: https://iabtechlab.com/wp-content/uploads/2019/06/VAST_4.2_final_june26.pdf#page=42
type JavaScriptResource struct {
	Value           string      `xml:",cdata"`
	APIFramework    string      `xml:"apiFramework,attr,omitempty"`
	BrowserOptional NumericBool `xml:"browserOptional,attr,omitempty"`
}

// IFrameResource represents an IFrame creative resource.
// Contains IFrame content for displaying ad creatives.
//
// Reference: IAB VAST 4.x Section 2.3.2.4 - IFrameResource Element
// Link: https://iabtechlab.com/wp-content/uploads/2019/06/VAST_4.2_final_june26.pdf#page=73
type IFrameResource struct {
	Value string `xml:",cdata"`
}
