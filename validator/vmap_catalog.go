package validator

import "github.com/admein-advertising/admein-vast-generator/vast"

var supportedVMAPVersions = []vast.Version{
	"1.0",
}

var defaultVMAPCatalog = &Catalog{Nodes: map[string]*NodeSpec{
	"VMAP": {
		Name:     "VMAP",
		Versions: supportedVMAPVersions,
		Attributes: map[string]*AttributeSpec{
			"version": {Name: "version", Versions: supportedVMAPVersions, Required: true},
		},
		Children: map[string]*ChildSpec{
			"AdBreak":    {Name: "AdBreak", Versions: supportedVMAPVersions, Multiple: true},
			"Extensions": {Name: "Extensions", Versions: supportedVMAPVersions, Optional: true},
		},
	},
	"AdBreak": {
		Name:     "AdBreak",
		Versions: supportedVMAPVersions,
		Attributes: map[string]*AttributeSpec{
			"timeOffset":  {Name: "timeOffset", Versions: supportedVMAPVersions, Required: true},
			"breakType":   {Name: "breakType", Versions: supportedVMAPVersions},
			"breakId":     {Name: "breakId", Versions: supportedVMAPVersions},
			"repeatAfter": {Name: "repeatAfter", Versions: supportedVMAPVersions},
		},
		Children: map[string]*ChildSpec{
			"AdSource":       {Name: "AdSource", Versions: supportedVMAPVersions, Multiple: true},
			"TrackingEvents": {Name: "TrackingEvents", Versions: supportedVMAPVersions, Optional: true},
			"Extensions":     {Name: "Extensions", Versions: supportedVMAPVersions, Optional: true},
		},
	},
	"AdSource": {
		Name:     "AdSource",
		Versions: supportedVMAPVersions,
		Attributes: map[string]*AttributeSpec{
			"id":               {Name: "id", Versions: supportedVMAPVersions},
			"allowMultipleAds": {Name: "allowMultipleAds", Versions: supportedVMAPVersions},
			"followRedirects":  {Name: "followRedirects", Versions: supportedVMAPVersions},
			"breakId":          {Name: "breakId", Versions: supportedVMAPVersions},
			"fallbackOnNoAd":   {Name: "fallbackOnNoAd", Versions: supportedVMAPVersions},
		},
		Children: map[string]*ChildSpec{
			"AdTagURI":       {Name: "AdTagURI", Versions: supportedVMAPVersions, Optional: true},
			"VASTAdData":     {Name: "VASTAdData", Versions: supportedVMAPVersions, Optional: true},
			"CustomAdData":   {Name: "CustomAdData", Versions: supportedVMAPVersions, Optional: true},
			"TrackingEvents": {Name: "TrackingEvents", Versions: supportedVMAPVersions, Optional: true},
			"Extensions":     {Name: "Extensions", Versions: supportedVMAPVersions, Optional: true},
		},
	},
	"AdTagURI": {
		Name:     "AdTagURI",
		Versions: supportedVMAPVersions,
		Attributes: map[string]*AttributeSpec{
			"templateType": {Name: "templateType", Versions: supportedVMAPVersions},
		},
	},
	"VASTAdData": {
		Name:                 "VASTAdData",
		Versions:             supportedVMAPVersions,
		AllowUnknownChildren: true,
	},
	"CustomAdData": {
		Name:                 "CustomAdData",
		Versions:             supportedVMAPVersions,
		AllowUnknownChildren: true,
	},
	"TrackingEvents": {
		Name:     "TrackingEvents",
		Versions: supportedVMAPVersions,
		Children: map[string]*ChildSpec{
			"Tracking": {Name: "Tracking", Versions: supportedVMAPVersions, Multiple: true},
		},
	},
	"Tracking": {
		Name:     "Tracking",
		Versions: supportedVMAPVersions,
		Attributes: map[string]*AttributeSpec{
			"event": {Name: "event", Versions: supportedVMAPVersions, Required: true},
		},
	},
	"Extensions": {
		Name:     "Extensions",
		Versions: supportedVMAPVersions,
		Children: map[string]*ChildSpec{
			"Extension": {Name: "Extension", Versions: supportedVMAPVersions, Optional: true, Multiple: true},
		},
	},
	"Extension": {
		Name:                 "Extension",
		Versions:             supportedVMAPVersions,
		AllowUnknownChildren: true,
		Attributes: map[string]*AttributeSpec{
			"type":            {Name: "type", Versions: supportedVMAPVersions},
			"suppress_bumper": {Name: "suppress_bumper", Versions: supportedVMAPVersions},
		},
	},
}}
