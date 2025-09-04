package vast

// Namespace represents XML namespace declarations for VAST documents.
// Defines XML Schema Instance namespace for validation purposes.
//
// Reference: IAB VAST 4.x Section 2.3.1 - VAST Element
// Link: https://iabtechlab.com/wp-content/uploads/2019/06/VAST_4.2_final_june26.pdf#page=37
type Namespace string

// NamespaceSchemaLocation specifies the location of the VAST XSD schema file.
// Points to the schema definition for VAST document validation.
type NamespaceSchemaLocation string

// Default VAST namespace and schema location
const VASTNamespace Namespace = "http://www.w3.org/2001/XMLSchema-instance"
const VASTSchemaLocation NamespaceSchemaLocation = "vast.xsd"
