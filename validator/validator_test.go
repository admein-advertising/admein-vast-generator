package validator

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func followCatalogPath(t *testing.T, cat *Catalog, start string, path ...string) *NodeSpec {
	t.Helper()
	node, ok := cat.Nodes[start]
	if !ok || node == nil {
		t.Fatalf("expected node %q in catalog", start)
	}
	for _, childName := range path {
		child, ok := node.Children[childName]
		if !ok || child == nil {
			t.Fatalf("expected child %q on node %q", childName, node.Name)
		}
		targetKey := child.Name
		if child.NodeOverride != "" {
			targetKey = child.NodeOverride
		}
		node, ok = cat.Nodes[targetKey]
		if !ok || node == nil {
			t.Fatalf("expected target node %q for child %q", targetKey, childName)
		}
	}
	return node
}

func TestDefaultVASTCatalog_ReturnsDefensiveCopy(t *testing.T) {
	cat := DefaultVASTCatalog()
	if cat == nil {
		t.Fatalf("expected catalog copy")
	}
	if cat == defaultCatalog {
		t.Fatalf("expected catalog copy to be distinct from default")
	}
	vastNode, ok := cat.Nodes["VAST"]
	if !ok {
		t.Fatalf("expected VAST node in catalog copy")
	}
	vastAttr, ok := vastNode.Attributes["version"]
	if !ok {
		t.Fatalf("expected version attr in VAST node")
	}
	vastAttr.Required = false
	delete(cat.Nodes, "Ad")
	fresh := DefaultVASTCatalog()
	freshVAST, ok := fresh.Nodes["VAST"]
	if !ok {
		t.Fatalf("expected VAST node in fresh catalog")
	}
	freshVersion, ok := freshVAST.Attributes["version"]
	if !ok {
		t.Fatalf("expected version attribute in fresh catalog")
	}
	if !freshVersion.Required {
		t.Fatalf("mutations to copy must not affect default catalog")
	}
	if _, ok := fresh.Nodes["Ad"]; !ok {
		t.Fatalf("default catalog should still expose Ad node")
	}
}

func TestDefaultVASTCatalog_DocumentationPopulated(t *testing.T) {
	cat := DefaultVASTCatalog()
	node, ok := cat.Nodes["Ad"]
	if !ok {
		t.Fatalf("expected Ad node in catalog")
	}
	if node.Documentation == nil || node.Documentation.Content == "" {
		t.Fatalf("expected Ad node documentation with content, got %+v", node.Documentation)
	}
	if node.Documentation.Source != vast42SchemaURL {
		t.Fatalf("expected Ad node documentation source %s, got %+v", vast42SchemaURL, node.Documentation)
	}
	attr, ok := node.Attributes["sequence"]
	if !ok {
		t.Fatalf("expected sequence attribute on Ad node")
	}
	if attr.Documentation == nil || !strings.Contains(strings.ToLower(attr.Documentation.Content), "sequence") {
		t.Fatalf("expected documentation for Ad/@sequence, got %+v", attr.Documentation)
	}
	if attr.Value == nil || attr.Value.Documentation == nil || attr.Value.Documentation.Content == "" {
		t.Fatalf("expected attribute value spec documentation for Ad/@sequence")
	}
	child, ok := node.Children["InLine"]
	if !ok {
		t.Fatalf("expected InLine child on Ad node")
	}
	if child.Documentation == nil || child.Documentation.Content == "" {
		t.Fatalf("expected documentation for Ad->InLine relationship, got %+v", child.Documentation)
	}
}

func TestDefaultVASTCatalog_ChildDocumentationInheritsNodeDocumentation(t *testing.T) {
	cat := DefaultVASTCatalog()
	linear, ok := cat.Nodes["Linear"]
	if !ok {
		t.Fatalf("expected Linear node in catalog")
	}
	videoClicks, ok := linear.Children["VideoClicks"]
	if !ok {
		t.Fatalf("expected VideoClicks child on Linear")
	}
	if videoClicks.Documentation == nil || strings.TrimSpace(videoClicks.Documentation.Content) == "" {
		t.Fatalf("expected documentation for Linear->VideoClicks, got %+v", videoClicks.Documentation)
	}
	if strings.Contains(videoClicks.Documentation.Content, "For Companions creativeView is the only supported event") {
		t.Fatalf("expected context-safe documentation for Linear->VideoClicks, got misattributed text: %+v", videoClicks.Documentation)
	}
	if videoClicks.Documentation.Source != vast42SchemaURL {
		t.Fatalf("expected VideoClicks child documentation source %s, got %+v", vast42SchemaURL, videoClicks.Documentation)
	}
}

func TestDefaultVASTCatalog_TrackingEventsDocsAreNotMisattributed(t *testing.T) {
	cat := DefaultVASTCatalog()
	node, ok := cat.Nodes["TrackingEvents"]
	if !ok || node == nil || node.Documentation == nil {
		t.Fatalf("expected TrackingEvents node documentation")
	}
	if strings.Contains(node.Documentation.Content, "For Companions creativeView is the only supported event") {
		t.Fatalf("TrackingEvents node should not be companion-specific: %+v", node.Documentation)
	}

	linear, ok := cat.Nodes["Linear"]
	if !ok {
		t.Fatalf("expected Linear node in catalog")
	}
	child, ok := linear.Children["TrackingEvents"]
	if !ok {
		t.Fatalf("expected TrackingEvents child on Linear")
	}
	if child.Documentation == nil || strings.TrimSpace(child.Documentation.Content) == "" {
		t.Fatalf("expected documentation for Linear->TrackingEvents")
	}
	if strings.Contains(child.Documentation.Content, "For Companions creativeView is the only supported event") {
		t.Fatalf("Linear->TrackingEvents should not inherit Companion-specific documentation: %+v", child.Documentation)
	}
}

func TestVAST42ContextElementDocs_TracksLinearPath(t *testing.T) {
	const path = "VAST>Ad>InLine>Creatives>Creative>Linear>TrackingEvents"
	if _, ok := vast42ContextElementDocs[path]; !ok {
		t.Fatalf("expected generated context docs to include %q", path)
	}
}

func TestDefaultVASTCatalog_CompanionTrackingEventsUsesContextualNode(t *testing.T) {
	cat := DefaultVASTCatalog()
	companionNode := followCatalogPath(t, cat, "VAST", "Ad", "InLine", "Creatives", "Creative", "CompanionAds", "Companion")
	child, ok := companionNode.Children["TrackingEvents"]
	if !ok {
		t.Fatalf("expected TrackingEvents child on contextual Companion node")
	}
	if child.Documentation == nil || !strings.Contains(child.Documentation.Content, "For Companions creativeView is the only supported event") {
		t.Fatalf("expected Companion TrackingEvents child documentation to remain companion-specific, got %+v", child.Documentation)
	}
}

func TestDefaultVASTCatalog_PromotesUnambiguousContextNodeDocs(t *testing.T) {
	cat := DefaultVASTCatalog()
	node, ok := cat.Nodes["AdVerifications"]
	if !ok || node == nil {
		t.Fatalf("expected AdVerifications node in catalog")
	}
	if node.Documentation == nil || strings.TrimSpace(node.Documentation.Content) == "" {
		t.Fatalf("expected documentation for AdVerifications node")
	}
	if strings.Contains(node.Documentation.Content, "Defined in VAST 4.2 XSD element <AdVerifications>") {
		t.Fatalf("expected contextual node documentation to replace generic fallback, got %+v", node.Documentation)
	}
	if !strings.Contains(node.Documentation.Content, "AdVerification element") {
		t.Fatalf("expected specific AdVerifications documentation, got %+v", node.Documentation)
	}
}

func TestDefaultVASTCatalog_VerificationTrackingUsesContextualAttributes(t *testing.T) {
	cat := DefaultVASTCatalog()
	verificationNode := followCatalogPath(t, cat, "VAST", "Ad", "InLine", "AdVerifications", "Verification")
	trackingEventsChild, ok := verificationNode.Children["TrackingEvents"]
	if !ok {
		t.Fatalf("expected TrackingEvents child on contextual Verification node")
	}
	if trackingEventsChild.Documentation == nil || !strings.Contains(strings.ToLower(trackingEventsChild.Documentation.Content), "trackingevents") {
		t.Fatalf("expected Verification->TrackingEvents child documentation, got %+v", trackingEventsChild.Documentation)
	}

	trackingNode, ok := cat.Nodes["Tracking"]
	if !ok || trackingNode == nil {
		t.Fatalf("expected base Tracking node in catalog")
	}
	if trackingNode.Attributes["offset"] == nil {
		t.Fatalf("expected base Tracking node to retain offset attribute")
	}
}

func TestDefaultVASTCatalog_WrapperLinearOverrides(t *testing.T) {
	cat := DefaultVASTCatalog()
	wrapper, ok := cat.Nodes["Wrapper"]
	if !ok {
		t.Fatalf("expected Wrapper node")
	}
	creativesChild, ok := wrapper.Children["Creatives"]
	if !ok {
		t.Fatalf("expected Creatives child on Wrapper")
	}
	if creativesChild.NodeOverride != "WrapperCreatives" {
		t.Fatalf("expected Wrapper Creatives override, got %q", creativesChild.NodeOverride)
	}
	wrapperCreatives, ok := cat.Nodes["WrapperCreatives"]
	if !ok {
		t.Fatalf("expected WrapperCreatives spec")
	}
	creativeChild, ok := wrapperCreatives.Children["Creative"]
	if !ok {
		t.Fatalf("expected Creative child on WrapperCreatives")
	}
	if creativeChild.NodeOverride != "WrapperCreative" {
		t.Fatalf("expected WrapperCreative override, got %q", creativeChild.NodeOverride)
	}
	wrapperCreative, ok := cat.Nodes["WrapperCreative"]
	if !ok {
		t.Fatalf("expected WrapperCreative spec")
	}
	linearChild, ok := wrapperCreative.Children["Linear"]
	if !ok {
		t.Fatalf("expected Linear child on WrapperCreative")
	}
	if linearChild.NodeOverride != "WrapperLinear" {
		t.Fatalf("expected WrapperLinear override, got %q", linearChild.NodeOverride)
	}
	wrapperLinear, ok := cat.Nodes["WrapperLinear"]
	if !ok {
		t.Fatalf("expected WrapperLinear spec")
	}
	if child, ok := wrapperLinear.Children["Duration"]; !ok || !child.Optional {
		t.Fatalf("expected Wrapper Linear duration to be optional")
	}
	if child, ok := wrapperLinear.Children["MediaFiles"]; !ok || !child.Optional {
		t.Fatalf("expected Wrapper Linear media files to be optional")
	}
	linear, ok := cat.Nodes["Linear"]
	if !ok {
		t.Fatalf("expected base Linear spec")
	}
	if child := linear.Children["Duration"]; child == nil || child.Optional {
		t.Fatalf("expected base Linear duration to remain required")
	}
	if child := linear.Children["MediaFiles"]; child == nil || child.Optional {
		t.Fatalf("expected base Linear media files to remain required")
	}
}

func TestDefaultVMAPCatalog_ReturnsDefensiveCopy(t *testing.T) {
	cat := DefaultVMAPCatalog()
	if cat == nil {
		t.Fatalf("expected VMAP catalog copy")
	}
	if cat == defaultVMAPCatalog {
		t.Fatalf("expected VMAP catalog copy to be distinct from default")
	}
	vmapNode, ok := cat.Nodes["VMAP"]
	if !ok {
		t.Fatalf("expected VMAP node in catalog copy")
	}
	vmapAttr, ok := vmapNode.Attributes["version"]
	if !ok {
		t.Fatalf("expected VMAP version attribute")
	}
	vmapAttr.Required = false
	fresh := DefaultVMAPCatalog()
	freshVMAP, ok := fresh.Nodes["VMAP"]
	if !ok {
		t.Fatalf("expected VMAP node in fresh catalog")
	}
	freshVersion, ok := freshVMAP.Attributes["version"]
	if !ok {
		t.Fatalf("expected version attribute in fresh VMAP catalog")
	}
	if !freshVersion.Required {
		t.Fatalf("mutations to VMAP copy must not affect default catalog")
	}
}

func TestValidate_SuccessfulInline(t *testing.T) {
	resetCustom(t)
	xml := `<?xml version="1.0" encoding="UTF-8"?>
<VAST version="4.2">
  <Ad id="ad-1">
    <InLine>
      <Creatives>
        <Creative id="c1" sequence="1">
          <Linear skipoffset="00:00:05">
            <MediaFiles>
              <MediaFile delivery="progressive" type="video/mp4" width="640" height="360">http://cdn.example.com/asset.mp4</MediaFile>
            </MediaFiles>
          </Linear>
        </Creative>
      </Creatives>
    </InLine>
  </Ad>
</VAST>`

	result, err := Validate([]byte(xml))
	if err != nil {
		t.Fatalf("validate returned error: %v", err)
	}

	assertStatus(t, result.Root, "VAST", StatusPass)
	adNode := findNode(result.Root, "Ad")
	if adNode == nil {
		t.Fatalf("expected Ad node in result tree")
	}
	analysis := adNode.Analyses[IABAnalysisCategory]
	if analysis == nil || analysis.Status != StatusPass {
		t.Fatalf("expected Ad node to pass IAB analysis, got %+v", analysis)
	}
}

func TestValidate_AssignsSourcePointers(t *testing.T) {
	resetCustom(t)
	xml := `<?xml version="1.0" encoding="UTF-8"?>
<VAST version="4.2">
  <Ad id="ad-1">
    <InLine>
      <Creatives>
        <Creative>
          <Linear>
            <TrackingEvents>
              <Tracking event="start"><![CDATA[https://example.com/start]]></Tracking>
              <Tracking event="firstQuartile"><![CDATA[https://example.com/q1]]></Tracking>
              <Tracking event="midpoint"><![CDATA[https://example.com/mid]]></Tracking>
            </TrackingEvents>
          </Linear>
        </Creative>
      </Creatives>
    </InLine>
  </Ad>
</VAST>`

	result, err := Validate([]byte(xml), DisableHTTPValidators())
	if err != nil {
		t.Fatalf("validate returned error: %v", err)
	}
	if result.Root == nil {
		t.Fatalf("expected root result")
	}
	if result.Root.SourcePointer != "/VAST[1]" {
		t.Fatalf("expected root pointer /VAST[1], got %s", result.Root.SourcePointer)
	}
	trackingEvents := findNode(result.Root, "TrackingEvents")
	if trackingEvents == nil {
		t.Fatalf("expected TrackingEvents node in result")
	}
	expectedTrackingEventsPointer := "/VAST[1]/Ad[1]/InLine[1]/Creatives[1]/Creative[1]/Linear[1]/TrackingEvents[1]"
	if trackingEvents.SourcePointer != expectedTrackingEventsPointer {
		t.Fatalf("expected pointer %s, got %s", expectedTrackingEventsPointer, trackingEvents.SourcePointer)
	}
	if len(trackingEvents.Children) != 3 {
		t.Fatalf("expected 3 Tracking children, got %d", len(trackingEvents.Children))
	}
	for index, child := range trackingEvents.Children {
		expectedChildPointer := fmt.Sprintf("%s/Tracking[%d]", expectedTrackingEventsPointer, index+1)
		if child.SourcePointer != expectedChildPointer {
			t.Fatalf("expected child pointer %s, got %s", expectedChildPointer, child.SourcePointer)
		}
	}
}

func TestValidate_TrackingEventEnumValidation(t *testing.T) {
	resetCustom(t)
	xml := `<?xml version="1.0" encoding="UTF-8"?>
<VAST version="4.2">
	<Ad id="1">
		<InLine>
			<AdSystem>Example</AdSystem>
			<AdServingId>srv-1</AdServingId>
			<Impression><![CDATA[https://example.com/imp]]></Impression>
			<AdTitle>Sample</AdTitle>
			<Creatives>
				<Creative id="c1">
					<Linear>
						<Duration>00:00:05</Duration>
						<TrackingEvents>
							<Tracking event="invalidEvent"><![CDATA[https://example.com/track]]></Tracking>
						</TrackingEvents>
						<MediaFiles>
							<MediaFile delivery="progressive" type="video/mp4" width="640" height="360">https://example.com/video.mp4</MediaFile>
						</MediaFiles>
					</Linear>
				</Creative>
			</Creatives>
		</InLine>
	</Ad>
</VAST>`

	result, err := Validate([]byte(xml), DisableHTTPValidators())
	if err != nil {
		t.Fatalf("validate returned error: %v", err)
	}

	tracking := findNode(result.Root, "Tracking")
	if tracking == nil {
		t.Fatalf("expected Tracking node in result")
	}
	iab := tracking.Analyses[IABAnalysisCategory]
	if iab == nil || iab.Status != StatusFail {
		t.Fatalf("expected Tracking analysis failure, got %+v", iab)
	}
	joined := strings.Join(iab.Reasons, ";")
	if !strings.Contains(joined, "event must be one of") {
		t.Fatalf("expected failure mentioning allowed values, got %s", joined)
	}
}

func TestValidate_TrackingEventInteractiveStartAllowed(t *testing.T) {
	resetCustom(t)
	xml := `<?xml version="1.0" encoding="UTF-8"?>
<VAST version="4.2">
	<Ad id="1">
		<InLine>
			<AdSystem>Example</AdSystem>
			<AdServingId>srv-1</AdServingId>
			<Impression><![CDATA[https://example.com/imp]]></Impression>
			<AdTitle>Sample</AdTitle>
			<Creatives>
				<Creative id="c1">
					<Linear>
						<Duration>00:00:05</Duration>
						<TrackingEvents>
							<Tracking event="interactiveStart"><![CDATA[https://example.com/interactive]]></Tracking>
						</TrackingEvents>
						<MediaFiles>
							<MediaFile delivery="progressive" type="video/mp4" width="640" height="360">https://example.com/video.mp4</MediaFile>
						</MediaFiles>
					</Linear>
				</Creative>
			</Creatives>
		</InLine>
	</Ad>
</VAST>`

	result, err := Validate([]byte(xml), DisableHTTPValidators())
	if err != nil {
		t.Fatalf("validate returned error: %v", err)
	}

	assertStatus(t, result.Root, "Tracking", StatusPass)
}

func TestValidate_UnknownNode(t *testing.T) {
	resetCustom(t)
	xml := `<VAST version="4.2"><UnknownNode /></VAST>`

	result, err := Validate([]byte(xml))
	if err != nil {
		t.Fatalf("validate returned error: %v", err)
	}

	child := findNode(result.Root, "UnknownNode")
	if child == nil {
		t.Fatalf("expected UnknownNode result to exist")
	}
	analysis := child.Analyses[IABAnalysisCategory]
	if analysis == nil || analysis.Status != StatusFail {
		t.Fatalf("expected failure for unknown node, got %+v", analysis)
	}
	if len(analysis.Reasons) == 0 {
		t.Fatalf("expected reason for unknown node failure")
	}
}

func TestValidate_AttributeBooleanValidation(t *testing.T) {
	resetCustom(t)
	xml := `<?xml version="1.0" encoding="UTF-8"?>
<VAST version="4.2">
	<Ad id="1">
		<Wrapper allowMultipleAds="maybe" followAdditionalWrappers="true" fallbackOnNoAd="false">
			<AdSystem>Example</AdSystem>
			<VASTAdTagURI><![CDATA[https://example.com/vast.xml]]></VASTAdTagURI>
			<Impression><![CDATA[https://example.com/imp]]></Impression>
		</Wrapper>
	</Ad>
</VAST>`

	result, err := Validate([]byte(xml), DisableHTTPValidators())
	if err != nil {
		t.Fatalf("validate returned error: %v", err)
	}
	wrapper := findNode(result.Root, "Wrapper")
	if wrapper == nil {
		t.Fatalf("expected Wrapper node in result")
	}
	iab := wrapper.Analyses[IABAnalysisCategory]
	if iab == nil || iab.Status != StatusFail {
		t.Fatalf("expected Wrapper validation failure, got %+v", iab)
	}
	joined := strings.Join(iab.Reasons, ";")
	if !strings.Contains(joined, "allowMultipleAds") {
		t.Fatalf("expected failure mentioning allowMultipleAds, got %s", joined)
	}
}

func TestValidate_AttributeEnumValidation(t *testing.T) {
	resetCustom(t)
	xml := `<?xml version="1.0" encoding="UTF-8"?>
<VAST version="4.2">
	<Ad id="1">
		<InLine>
			<AdSystem>Example</AdSystem>
			<AdServingId>srv-1</AdServingId>
			<Impression><![CDATA[https://example.com/imp]]></Impression>
			<Pricing model="CPL" currency="US">
				<![CDATA[1.00]]>
			</Pricing>
			<Creatives>
				<Creative id="c1">
					<Linear>
						<Duration>00:00:05</Duration>
						<MediaFiles>
							<MediaFile delivery="progressive" type="video/mp4" width="640" height="360">https://example.com/video.mp4</MediaFile>
						</MediaFiles>
					</Linear>
				</Creative>
			</Creatives>
		</InLine>
	</Ad>
</VAST>`

	result, err := Validate([]byte(xml), DisableHTTPValidators())
	if err != nil {
		t.Fatalf("validate returned error: %v", err)
	}
	pricing := findNode(result.Root, "Pricing")
	if pricing == nil {
		t.Fatalf("expected Pricing node in result")
	}
	iab := pricing.Analyses[IABAnalysisCategory]
	if iab == nil || iab.Status != StatusFail {
		t.Fatalf("expected Pricing validation failure, got %+v", iab)
	}
	joined := strings.Join(iab.Reasons, ";")
	if !strings.Contains(joined, "model must be one of") {
		t.Fatalf("expected enum failure for model, got %s", joined)
	}
	if !strings.Contains(joined, "currency must match pattern") {
		t.Fatalf("expected pattern failure for currency, got %s", joined)
	}
}

func TestValidate_WrapperBlockedAdCategoriesAuthority(t *testing.T) {
	resetCustom(t)
	xml := `<?xml version="1.0" encoding="UTF-8"?>
<VAST version="4.2">
	<Ad id="1">
		<Wrapper allowMultipleAds="false" followAdditionalWrappers="false" fallbackOnNoAd="true">
			<AdSystem>Example</AdSystem>
			<VASTAdTagURI><![CDATA[https://example.com/vast.xml]]></VASTAdTagURI>
			<Impression><![CDATA[https://example.com/imp]]></Impression>
			<BlockedAdCategories authority="example.com">cat</BlockedAdCategories>
		</Wrapper>
	</Ad>
</VAST>`

	result, err := Validate([]byte(xml), DisableHTTPValidators())
	if err != nil {
		t.Fatalf("validate returned error: %v", err)
	}
	blocked := findNode(result.Root, "BlockedAdCategories")
	if blocked == nil {
		t.Fatalf("expected BlockedAdCategories node in result")
	}
	iab := blocked.Analyses[IABAnalysisCategory]
	if iab == nil || iab.Status != StatusFail {
		t.Fatalf("expected BlockedAdCategories analysis failure, got %+v", iab)
	}
	joined := strings.Join(iab.Reasons, ";")
	if !strings.Contains(joined, "authority") {
		t.Fatalf("expected failure mentioning authority, got %s", joined)
	}
}

func TestValidate_CompanionAdsRequiredEnum(t *testing.T) {
	resetCustom(t)
	xml := `<?xml version="1.0" encoding="UTF-8"?>
<VAST version="4.2">
	<Ad id="1">
		<InLine>
			<AdSystem>Example</AdSystem>
			<AdServingId>srv-1</AdServingId>
			<Impression><![CDATA[https://example.com/imp]]></Impression>
			<AdTitle>Sample</AdTitle>
			<Creatives>
				<Creative id="c1">
					<Linear>
						<Duration>00:00:05</Duration>
						<MediaFiles>
							<MediaFile delivery="progressive" type="video/mp4" width="640" height="360">https://example.com/video.mp4</MediaFile>
						</MediaFiles>
					</Linear>
					<CompanionAds required="sometimes">
						<Companion width="300" height="250">
							<StaticResource creativeType="image/png">https://example.com/companion.png</StaticResource>
						</Companion>
					</CompanionAds>
				</Creative>
			</Creatives>
		</InLine>
	</Ad>
</VAST>`

	result, err := Validate([]byte(xml), DisableHTTPValidators())
	if err != nil {
		t.Fatalf("validate returned error: %v", err)
	}
	compAds := findNode(result.Root, "CompanionAds")
	if compAds == nil {
		t.Fatalf("expected CompanionAds node in result")
	}
	iab := compAds.Analyses[IABAnalysisCategory]
	if iab == nil || iab.Status != StatusFail {
		t.Fatalf("expected CompanionAds analysis failure, got %+v", iab)
	}
	joined := strings.Join(iab.Reasons, ";")
	if !strings.Contains(joined, "must be one of") {
		t.Fatalf("expected failure mentioning allowed values, got %s", joined)
	}
}

func TestValidate_CustomValidator(t *testing.T) {
	resetCustom(t)
	RegisterCustomValidator("MediaFile", func(ctx NodeContext) *NodeAnalysisResult {
		return &NodeAnalysisResult{
			Category: CustomAnalysisCategory,
			Status:   StatusFail,
			Reasons:  []string{"media file URL failed custom check"},
		}
	})

	xml := `<VAST version="4.2"><Ad><InLine><Creatives><Creative><Linear><MediaFiles><MediaFile delivery="progressive" type="video/mp4" width="1" height="1">http://invalid</MediaFile></MediaFiles></Linear></Creative></Creatives></InLine></Ad></VAST>`

	result, err := Validate([]byte(xml))
	if err != nil {
		t.Fatalf("validate returned error: %v", err)
	}

	mediaFile := findNode(result.Root, "MediaFile")
	if mediaFile == nil {
		t.Fatalf("expected MediaFile node in result")
	}
	analysis := mediaFile.Analyses[CustomAnalysisCategory]
	if analysis == nil {
		t.Fatalf("expected custom analysis result for MediaFile, got %+v", mediaFile.Analyses)
	}
	if analysis.Status != StatusFail {
		t.Fatalf("expected custom validator to fail, got status %s", analysis.Status)
	}
}

func TestValidate_HTTPValidator(t *testing.T) {
	resetCustom(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "video/mp4")
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()
	RegisterHTTPValidator("MediaFile", func(ctx context.Context, nodeCtx NodeContext, client *http.Client) (*NodeAnalysisResult, error) {
		return &NodeAnalysisResult{
			Category: CustomAnalysisCategory,
			Status:   StatusFail,
			Reasons:  []string{"HTTP check failed"},
		}, nil
	})

	xml := fmt.Sprintf(`<VAST version="4.2"><Ad><InLine><Creatives><Creative><Linear><MediaFiles><MediaFile delivery="progressive" type="video/mp4" width="1" height="1">%s/video.mp4</MediaFile></MediaFiles></Linear></Creative></Creatives></InLine></Ad></VAST>`, ts.URL)

	result, err := Validate([]byte(xml))
	if err != nil {
		t.Fatalf("validate returned error: %v", err)
	}

	mediaFile := findNode(result.Root, "MediaFile")
	if mediaFile == nil {
		t.Fatalf("expected MediaFile node in result")
	}
	analysis := mediaFile.Analyses[CustomAnalysisCategory]
	if analysis == nil {
		t.Fatalf("expected custom analysis result for MediaFile, got %+v", mediaFile.Analyses)
	}
	if analysis.Status != StatusFail {
		t.Fatalf("expected HTTP validator to fail, got status %s", analysis.Status)
	}
}

func TestValidate_BuiltInMediaFileHTTPValidatorSuccess(t *testing.T) {
	resetCustom(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "video/mp4")
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()
	xml := fmt.Sprintf(`<VAST version="4.2"><Ad><InLine><Creatives><Creative><Linear><MediaFiles><MediaFile delivery="progressive" type="video/mp4" width="1" height="1">%s/video.mp4</MediaFile></MediaFiles></Linear></Creative></Creatives></InLine></Ad></VAST>`, ts.URL)

	result, err := Validate([]byte(xml))
	if err != nil {
		t.Fatalf("validate returned error: %v", err)
	}
	mediaFile := findNode(result.Root, "MediaFile")
	if mediaFile == nil {
		t.Fatalf("expected MediaFile node in result")
	}
	analysis := mediaFile.Analyses[CustomAnalysisCategory]
	if analysis == nil || analysis.Status != StatusPass {
		t.Fatalf("expected built-in HTTP validator to pass, got %+v", analysis)
	}
}

func TestValidate_BuiltInMediaFileHTTPValidatorFailure(t *testing.T) {
	resetCustom(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()
	xml := fmt.Sprintf(`<VAST version="4.2"><Ad><InLine><Creatives><Creative><Linear><MediaFiles><MediaFile delivery="progressive" type="video/mp4" width="1" height="1">%s/video.mp4</MediaFile></MediaFiles></Linear></Creative></Creatives></InLine></Ad></VAST>`, ts.URL)

	result, err := Validate([]byte(xml))
	if err != nil {
		t.Fatalf("validate returned error: %v", err)
	}
	mediaFile := findNode(result.Root, "MediaFile")
	if mediaFile == nil {
		t.Fatalf("expected MediaFile node in result")
	}
	analysis := mediaFile.Analyses[CustomAnalysisCategory]
	if analysis == nil || analysis.Status != StatusFail {
		t.Fatalf("expected built-in HTTP validator to fail, got %+v", analysis)
	}
	if len(analysis.Reasons) == 0 || !strings.Contains(strings.Join(analysis.Reasons, ";"), "HTTP") {
		t.Fatalf("expected failure reason mentioning HTTP status, got %+v", analysis.Reasons)
	}
}

func TestValidate_CategorySummaries(t *testing.T) {
	resetCustom(t)
	RegisterCustomValidator("Linear", func(ctx NodeContext) *NodeAnalysisResult {
		return &NodeAnalysisResult{
			Category: CustomAnalysisCategory,
			Status:   StatusFail,
			Reasons:  []string{"linear custom failure"},
		}
	})
	xml := `<?xml version="1.0" encoding="UTF-8"?>
<VAST version="4.2">
	<Ad id="1">
		<InLine>
			<Creatives>
				<Creative id="c1">
					<Linear>
						<Duration>00:00:05</Duration>
						<MediaFiles>
							<MediaFile delivery="progressive" type="video/mp4" width="1" height="1">http://example.com/video.mp4</MediaFile>
						</MediaFiles>
					</Linear>
				</Creative>
			</Creatives>
		</InLine>
	</Ad>
</VAST>`

	result, err := Validate([]byte(xml), DisableHTTPValidators())
	if err != nil {
		t.Fatalf("validate returned error: %v", err)
	}
	if len(result.Summaries) == 0 {
		t.Fatalf("expected summaries to be populated")
	}
	iab := result.Summaries[IABAnalysisCategory]
	if iab == nil || iab.TotalNodes == 0 {
		t.Fatalf("expected IAB summary, got %+v", iab)
	}
	customSummary := result.Summaries[CustomAnalysisCategory]
	if customSummary == nil {
		t.Fatalf("expected custom summary")
	}
	if customSummary.Status != StatusFail || customSummary.FailingNodes == 0 {
		t.Fatalf("expected failing custom summary, got %+v", customSummary)
	}
	if customSummary.WarningNodes != 0 || customSummary.RecommendationNodes != 0 {
		t.Fatalf("expected no warnings or recommendations in custom summary, got %+v", customSummary)
	}
	if len(customSummary.Reasons) == 0 || customSummary.Reasons[0] != "linear custom failure" {
		t.Fatalf("expected custom failure reason recorded")
	}
}

func TestValidate_ExtensionAllowsCustomNodes(t *testing.T) {
	resetCustom(t)
	xml := `<?xml version="1.0" encoding="UTF-8"?>
<VAST version="4.2">
	<Ad id="1">
		<InLine>
			<AdSystem>Example</AdSystem>
			<Impression><![CDATA[https://example.com/imp]]></Impression>
			<AdTitle>Sample Ad</AdTitle>
			<Creatives>
				<Creative id="c1">
					<Linear>
						<Duration>00:00:05</Duration>
						<MediaFiles>
							<MediaFile delivery="progressive" type="video/mp4" width="640" height="360">https://example.com/video.mp4</MediaFile>
						</MediaFiles>
					</Linear>
				</Creative>
			</Creatives>
			<Extensions>
				<Extension type="pm">
					<Meta><![CDATA[ name=pm-forcepixel;ver=1.0 ]]></Meta>
					<Pixel loc="0">
						<Code type="1"><![CDATA[ https://example.com/pixel ]]></Code>
					</Pixel>
				</Extension>
			</Extensions>
		</InLine>
	</Ad>
</VAST>`

	result, err := Validate([]byte(xml), DisableHTTPValidators())
	if err != nil {
		t.Fatalf("validate returned error: %v", err)
	}

	for _, nodeName := range []string{"Meta", "Pixel", "Code"} {
		node := findNode(result.Root, nodeName)
		if node == nil {
			t.Fatalf("expected %s node in result", nodeName)
		}
		analysis := node.Analyses[IABAnalysisCategory]
		if analysis == nil || analysis.Status != StatusPass {
			t.Fatalf("expected %s node to pass IAB analysis, got %+v", nodeName, analysis)
		}
	}

	iabSummary := result.Summaries[IABAnalysisCategory]
	if iabSummary == nil {
		t.Fatalf("expected IAB summary")
	}
	if iabSummary.Status != StatusPass {
		t.Fatalf("expected IAB summary to pass, got %+v", iabSummary)
	}
}

func TestValidate_ExtensionAllowsUnknownAttributes(t *testing.T) {
	resetCustom(t)
	xml := `<?xml version="1.0" encoding="UTF-8"?>
<VAST version="4.0">
	<Ad id="1">
		<InLine>
			<AdSystem>Example</AdSystem>
			<Impression><![CDATA[https://example.com/imp]]></Impression>
			<AdTitle>Sample</AdTitle>
			<Creatives>
				<Creative id="c1">
					<Linear>
						<Duration>00:00:05</Duration>
						<MediaFiles>
							<MediaFile delivery="progressive" type="video/mp4" width="640" height="360">https://example.com/video.mp4</MediaFile>
						</MediaFiles>
					</Linear>
				</Creative>
			</Creatives>
			<Extensions>
				<Extension type="waterfall" fallback_index="0">
					<CustomNode>value</CustomNode>
				</Extension>
			</Extensions>
		</InLine>
	</Ad>
</VAST>`

	result, err := Validate([]byte(xml), DisableHTTPValidators())
	if err != nil {
		t.Fatalf("validate returned error: %v", err)
	}

	ext := findNode(result.Root, "Extension")
	if ext == nil {
		t.Fatalf("expected Extension node in result")
	}
	analysis := ext.Analyses[IABAnalysisCategory]
	if analysis == nil || analysis.Status != StatusPass {
		t.Fatalf("expected Extension node to pass, got %+v", analysis)
	}

	var fallbackAttr AttributeResult
	found := false
	for _, attr := range analysis.Attributes {
		if attr.Name == "fallback_index" {
			fallbackAttr = attr
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected fallback_index attribute result to be recorded")
	}
	if fallbackAttr.Status != StatusInfo {
		t.Fatalf("expected fallback_index to be treated as informational, got %s", fallbackAttr.Status)
	}
}

func TestValidate_ExtensionUniversalAdIdBackport(t *testing.T) {
	resetCustom(t)
	xml := `<?xml version="1.0" encoding="UTF-8"?>
<VAST version="2.0">
	<Ad id="1">
		<InLine>
			<AdSystem>Example</AdSystem>
			<Impression><![CDATA[https://example.com/imp]]></Impression>
			<AdTitle>Sample</AdTitle>
			<Creatives>
				<Creative id="c1">
					<Linear>
						<Duration>00:00:05</Duration>
						<MediaFiles>
							<MediaFile delivery="progressive" type="video/mp4" width="1" height="1">https://example.com/video.mp4</MediaFile>
						</MediaFiles>
					</Linear>
				</Creative>
			</Creatives>
			<Extensions>
				<Extension type="UniversalAdId">
					<UniversalAdId idRegistry="ad-id.org" idValue="campaign-123">CNPA0484000H</UniversalAdId>
				</Extension>
			</Extensions>
		</InLine>
	</Ad>
</VAST>`

	result, err := Validate([]byte(xml), DisableHTTPValidators())
	if err != nil {
		t.Fatalf("validate returned error: %v", err)
	}

	uaid := findNode(result.Root, "UniversalAdId")
	if uaid == nil {
		t.Fatalf("expected UniversalAdId node in result")
	}
	analysis := uaid.Analyses[IABAnalysisCategory]
	if analysis == nil || analysis.Status != StatusPass {
		t.Fatalf("expected UniversalAdId backport to pass, got %+v", analysis)
	}
}

func TestValidate_ExtensionAdVerificationsBackport(t *testing.T) {
	resetCustom(t)
	xml := `<?xml version="1.0" encoding="UTF-8"?>
<VAST version="2.0">
	<Ad id="1">
		<InLine>
			<AdSystem>Example</AdSystem>
			<Impression><![CDATA[https://example.com/imp]]></Impression>
			<Extensions>
				<Extension type="AdVerifications">
					<AVID>
						<AdVerifications>
							<Verification vendor="iabtechlab">
								<JavaScriptResource><![CDATA[https://example.com/verification.js]]></JavaScriptResource>
							</Verification>
						</AdVerifications>
					</AVID>
				</Extension>
			</Extensions>
		</InLine>
	</Ad>
</VAST>`

	result, err := Validate([]byte(xml), DisableHTTPValidators())
	if err != nil {
		t.Fatalf("validate returned error: %v", err)
	}

	for _, nodeName := range []string{"AdVerifications", "Verification", "JavaScriptResource"} {
		node := findNode(result.Root, nodeName)
		if node == nil {
			t.Fatalf("expected %s node in result", nodeName)
		}
		analysis := node.Analyses[IABAnalysisCategory]
		if analysis == nil || analysis.Status != StatusPass {
			t.Fatalf("expected %s backport to pass, got %+v", nodeName, analysis)
		}
	}
}

func TestValidate_ExtensionAdVerificationsTypeMismatchFails(t *testing.T) {
	resetCustom(t)
	xml := `<?xml version="1.0" encoding="UTF-8"?>
<VAST version="2.0">
	<Ad id="1">
		<InLine>
			<AdSystem>Example</AdSystem>
			<Impression><![CDATA[https://example.com/imp]]></Impression>
			<Extensions>
				<Extension type="AVID">
					<AdVerifications>
						<Verification vendor="iabtechlab">
							<JavaScriptResource><![CDATA[https://example.com/verification.js]]></JavaScriptResource>
						</Verification>
					</AdVerifications>
				</Extension>
			</Extensions>
		</InLine>
	</Ad>
</VAST>`

	result, err := Validate([]byte(xml), DisableHTTPValidators())
	if err != nil {
		t.Fatalf("validate returned error: %v", err)
	}

	adVerifications := findNode(result.Root, "AdVerifications")
	if adVerifications == nil {
		t.Fatalf("expected AdVerifications node in result")
	}
	analysis := adVerifications.Analyses[IABAnalysisCategory]
	if analysis == nil || analysis.Status != StatusFail {
		t.Fatalf("expected AdVerifications to fail with mismatched extension type, got %+v", analysis)
	}
	joined := strings.Join(analysis.Reasons, ";")
	if !strings.Contains(joined, "Extension type") {
		t.Fatalf("expected failure reason mentioning Extension type, got %s", joined)
	}
}

func TestValidate_ExtensionUniversalAdIdMissingChild(t *testing.T) {
	resetCustom(t)
	xml := `<?xml version="1.0" encoding="UTF-8"?>
<VAST version="2.0">
	<Ad id="1">
		<InLine>
			<AdSystem>Example</AdSystem>
			<Impression><![CDATA[https://example.com/imp]]></Impression>
			<Extensions>
				<Extension type="UniversalAdId"></Extension>
			</Extensions>
		</InLine>
	</Ad>
</VAST>`

	result, err := Validate([]byte(xml), DisableHTTPValidators())
	if err != nil {
		t.Fatalf("validate returned error: %v", err)
	}

	ext := findNode(result.Root, "Extension")
	if ext == nil {
		t.Fatalf("expected Extension node in result")
	}
	analysis := ext.Analyses[IABAnalysisCategory]
	if analysis == nil || analysis.Status != StatusFail {
		t.Fatalf("expected Extension validator to fail, got %+v", analysis)
	}
	joined := strings.Join(analysis.Reasons, ";")
	if !strings.Contains(joined, "UniversalAdId") {
		t.Fatalf("expected UniversalAdId failure reason, got %s", joined)
	}
}

func TestValidate_ExtensionTypeWarning(t *testing.T) {
	resetCustom(t)
	xml := `<?xml version="1.0" encoding="UTF-8"?>
<VAST version="4.2">
	<Ad id="1">
		<InLine>
			<AdSystem>Example</AdSystem>
			<Impression><![CDATA[https://example.com/imp]]></Impression>
			<Creatives>
				<Creative id="c1">
					<Linear>
						<Duration>00:00:05</Duration>
						<MediaFiles>
							<MediaFile delivery="progressive" type="video/mp4" width="640" height="360">https://example.com/video.mp4</MediaFile>
						</MediaFiles>
					</Linear>
				</Creative>
			</Creatives>
			<Extensions>
				<Extension type="GDCM">
					<UniversalAdId idRegistry="ad-id.org" idValue="campaign-123">CNPA0484000H</UniversalAdId>
				</Extension>
			</Extensions>
		</InLine>
	</Ad>
</VAST>`

	result, err := Validate([]byte(xml), DisableHTTPValidators())
	if err != nil {
		t.Fatalf("validate returned error: %v", err)
	}

	ext := findNode(result.Root, "Extension")
	if ext == nil {
		t.Fatalf("expected Extension node in result")
	}
	analysis := ext.Analyses[IABAnalysisCategory]
	if analysis == nil {
		t.Fatalf("expected IAB analysis for Extension")
	}
	if analysis.Status != StatusWarning {
		t.Fatalf("expected warning status for Extension type mismatch, got %s", analysis.Status)
	}
	if len(analysis.Reasons) == 0 || !strings.Contains(strings.Join(analysis.Reasons, ";"), "type attribute value should") {
		t.Fatalf("expected warning reason for Extension type mismatch, got %+v", analysis.Reasons)
	}

	iabSummary := result.Summaries[IABAnalysisCategory]
	if iabSummary == nil {
		t.Fatalf("expected IAB summary")
	}
	if iabSummary.Status != StatusWarning {
		t.Fatalf("expected IAB summary warning status, got %+v", iabSummary)
	}
	if iabSummary.WarningNodes == 0 {
		t.Fatalf("expected warning node count in summary, got %+v", iabSummary)
	}
	if iabSummary.FailingNodes != 0 {
		t.Fatalf("expected zero failing nodes when only warnings present, got %+v", iabSummary)
	}
}

func TestValidate_VMAPAcceptedWithWarning(t *testing.T) {
	resetCustom(t)
	xml := `<?xml version="1.0" encoding="UTF-8"?>
<vmap:VMAP xmlns:vmap="http://www.iab.net/videosuite/vmap" version="1.0">
	<vmap:AdBreak timeOffset="start" breakType="linear" breakId="preroll">
		<vmap:AdSource id="pre-1" allowMultipleAds="false" followRedirects="true">
			<vmap:AdTagURI templateType="vast3"><![CDATA[https://example.com/vast.xml]]></vmap:AdTagURI>
		</vmap:AdSource>
	</vmap:AdBreak>
</vmap:VMAP>`

	result, err := Validate([]byte(xml), DisableHTTPValidators())
	if err != nil {
		t.Fatalf("validate returned error: %v", err)
	}
	if result == nil || result.Root == nil {
		t.Fatalf("expected VMAP validation result")
	}
	if result.Root.Node != "VMAP" {
		t.Fatalf("expected root node VMAP, got %s", result.Root.Node)
	}
	iab := result.Root.Analyses[IABAnalysisCategory]
	if iab == nil {
		t.Fatalf("expected IAB analysis on VMAP root")
	}
	if iab.Status != StatusInfo {
		t.Fatalf("expected VMAP root informational status, got %s", iab.Status)
	}
	joined := strings.Join(iab.Reasons, ";")
	if !strings.Contains(joined, "VMAP validation is informational") {
		t.Fatalf("expected VMAP informational reason, got %s", joined)
	}
	iabSummary := result.Summaries[IABAnalysisCategory]
	if iabSummary == nil || iabSummary.Status != StatusInfo {
		t.Fatalf("expected IAB summary informational status for VMAP, got %+v", iabSummary)
	}
}

func TestValidate_VMAPUnknownAttributeFails(t *testing.T) {
	resetCustom(t)
	xml := `<?xml version="1.0" encoding="UTF-8"?>
<VMAP version="1.0" xmlns="http://www.iab.net/videosuite/vmap">
	<AdBreak timeOffset="start" breakType="linear" breakId="preroll">
		<AdSource id="source-1" unknownAttr="nope">
			<AdTagURI><![CDATA[https://example.com/vast.xml]]></AdTagURI>
		</AdSource>
	</AdBreak>
</VMAP>`

	result, err := Validate([]byte(xml), DisableHTTPValidators())
	if err != nil {
		t.Fatalf("validate returned error: %v", err)
	}
	adSource := findNode(result.Root, "AdSource")
	if adSource == nil {
		t.Fatalf("expected AdSource node")
	}
	analysis := adSource.Analyses[IABAnalysisCategory]
	if analysis == nil || analysis.Status != StatusFail {
		t.Fatalf("expected AdSource failure for unknown attribute, got %+v", analysis)
	}
	joined := strings.Join(analysis.Reasons, ";")
	if !strings.Contains(joined, "unknownAttr") {
		t.Fatalf("expected reason mentioning unknownAttr, got %s", joined)
	}
	iabSummary := result.Summaries[IABAnalysisCategory]
	if iabSummary == nil || iabSummary.Status != StatusFail {
		t.Fatalf("expected failing summary for invalid VMAP, got %+v", iabSummary)
	}
}

func TestValidate_ExtensionInteractiveCreativeFileBackport(t *testing.T) {
	resetCustom(t)
	xml := `<?xml version="1.0" encoding="UTF-8"?>
<VAST version="2.0">
	<Ad id="1">
		<InLine>
			<AdSystem>Example</AdSystem>
			<Impression><![CDATA[https://example.com/imp]]></Impression>
			<Extensions>
				<Extension type="InteractiveCreativeFile">
					<InteractiveCreativeFile type="text/html" apiFramework="SIMID" variableDuration="true">
						<![CDATA[https://adserver.com/creative.html]]>
					</InteractiveCreativeFile>
				</Extension>
			</Extensions>
		</InLine>
	</Ad>
</VAST>`

	result, err := Validate([]byte(xml), DisableHTTPValidators())
	if err != nil {
		t.Fatalf("validate returned error: %v", err)
	}

	icf := findNode(result.Root, "InteractiveCreativeFile")
	if icf == nil {
		t.Fatalf("expected InteractiveCreativeFile node in result")
	}
	analysis := icf.Analyses[IABAnalysisCategory]
	if analysis == nil || analysis.Status != StatusPass {
		t.Fatalf("expected InteractiveCreativeFile backport to pass, got %+v", analysis)
	}
}

func TestValidate_ExtensionMezzanineBackport(t *testing.T) {
	resetCustom(t)
	xml := `<?xml version="1.0" encoding="UTF-8"?>
<VAST version="2.0">
	<Ad id="1">
		<InLine>
			<AdSystem>Example</AdSystem>
			<Impression><![CDATA[https://example.com/imp]]></Impression>
			<Creatives>
				<Creative id="c1">
					<Linear>
						<Duration>00:00:05</Duration>
						<MediaFiles>
							<MediaFile delivery="progressive" type="video/mp4" width="1" height="1">https://example.com/video.mp4</MediaFile>
						</MediaFiles>
					</Linear>
				</Creative>
			</Creatives>
			<Extensions>
				<Extension type="Mezzanine">
					<Mezzanine delivery="streaming" type="video/mp4" width="7680" height="4320" codec="video/3gpp" fileSize="300000000" mediaType="2D">
						<![CDATA[https://creative-company.com/mezzanine.mp4]]>
					</Mezzanine>
				</Extension>
			</Extensions>
		</InLine>
	</Ad>
</VAST>`

	result, err := Validate([]byte(xml), DisableHTTPValidators())
	if err != nil {
		t.Fatalf("validate returned error: %v", err)
	}

	mezz := findNode(result.Root, "Mezzanine")
	if mezz == nil {
		t.Fatalf("expected Mezzanine node in result")
	}
	analysis := mezz.Analyses[IABAnalysisCategory]
	if analysis == nil || analysis.Status != StatusPass {
		t.Fatalf("expected Mezzanine backport to pass, got %+v", analysis)
	}
}

func TestValidate_ExtensionMezzanineMissingContent(t *testing.T) {
	resetCustom(t)
	xml := `<?xml version="1.0" encoding="UTF-8"?>
<VAST version="2.0">
	<Ad id="1">
		<InLine>
			<AdSystem>Example</AdSystem>
			<Impression><![CDATA[https://example.com/imp]]></Impression>
			<Extensions>
				<Extension type="Mezzanine">
					<Mezzanine delivery="streaming" type="video/mp4" width="7680" height="4320"></Mezzanine>
				</Extension>
			</Extensions>
		</InLine>
	</Ad>
</VAST>`

	result, err := Validate([]byte(xml), DisableHTTPValidators())
	if err != nil {
		t.Fatalf("validate returned error: %v", err)
	}

	ext := findNode(result.Root, "Extension")
	if ext == nil {
		t.Fatalf("expected Extension node in result")
	}
	analysis := ext.Analyses[IABAnalysisCategory]
	if analysis == nil || analysis.Status != StatusFail {
		t.Fatalf("expected Mezzanine extension validator to fail, got %+v", analysis)
	}
	joined := strings.Join(analysis.Reasons, ";")
	if !strings.Contains(joined, "Mezzanine") {
		t.Fatalf("expected Mezzanine failure reason, got %s", joined)
	}
}

func TestValidate_UnsupportedVersionReturnsResult(t *testing.T) {
	resetCustom(t)
	xml := `<?xml version="1.0" encoding="UTF-8"?>
	<VAST version="5.0">
	<Ad>
		<InLine>
			<AdSystem>Example</AdSystem>
			<Impression><![CDATA[https://example.com/imp]]></Impression>
			<Creatives>
				<Creative>
					<Linear>
						<Duration>00:00:05</Duration>
						<MediaFiles>
								<MediaFile delivery="progressive" type="video/mp4" width="1" height="1">http://example.com/video.mp4</MediaFile>
						</MediaFiles>
					</Linear>
				</Creative>
			</Creatives>
		</InLine>
	</Ad>
</VAST>`

	result, err := Validate([]byte(xml), DisableHTTPValidators())
	if err != nil {
		t.Fatalf("validate returned error: %v", err)
	}
	if result == nil || result.Root == nil {
		t.Fatalf("expected validation result for unsupported version")
	}
	iab := result.Root.Analyses[IABAnalysisCategory]
	if iab == nil {
		t.Fatalf("expected IAB analysis on root")
	}
	if iab.Status != StatusFail {
		t.Fatalf("expected root analysis to fail for unsupported version")
	}
	if len(iab.Reasons) == 0 || !strings.Contains(strings.Join(iab.Reasons, ";"), "Unsupported VAST version") {
		t.Fatalf("expected unsupported version reason, got %+v", iab.Reasons)
	}
}

func TestValidate_AdVerificationsUnsupportedInV3(t *testing.T) {
	resetCustom(t)
	xml := `<?xml version="1.0" encoding="UTF-8"?>
<VAST version="3.0">
	<Ad id="1">
		<InLine>
			<AdSystem>Example</AdSystem>
			<Impression><![CDATA[https://example.com/imp]]></Impression>
			<Creatives>
				<Creative>
					<Linear>
						<Duration>00:00:05</Duration>
						<MediaFiles>
							<MediaFile delivery="progressive" type="video/mp4" width="1" height="1">http://example.com/video.mp4</MediaFile>
						</MediaFiles>
					</Linear>
				</Creative>
			</Creatives>
			<AdVerifications>
				<Verification vendor="iabtechlab">
					<JavaScriptResource><![CDATA[https://example.com/verification.js]]></JavaScriptResource>
				</Verification>
			</AdVerifications>
		</InLine>
	</Ad>
</VAST>`

	result, err := Validate([]byte(xml), DisableHTTPValidators())
	if err != nil {
		t.Fatalf("validate returned error: %v", err)
	}

	adVerifications := findNode(result.Root, "AdVerifications")
	if adVerifications == nil {
		t.Fatalf("expected AdVerifications node in result")
	}
	analysis := adVerifications.Analyses[IABAnalysisCategory]
	if analysis == nil {
		t.Fatalf("expected IAB analysis for AdVerifications")
	}
	if analysis.Status != StatusFail {
		t.Fatalf("expected AdVerifications to fail for VAST 3.0, got status %s", analysis.Status)
	}
	if len(analysis.Reasons) == 0 || !strings.Contains(strings.Join(analysis.Reasons, ";"), "not supported") {
		t.Fatalf("expected failure reason mentioning support, got %+v", analysis.Reasons)
	}
}

func TestValidate_WrapperScenario(t *testing.T) {
	resetCustom(t)
	xml := `<?xml version="1.0" encoding="UTF-8"?>
<VAST version="4.2">
	<Ad id="1">
		<Wrapper followAdditionalWrappers="0" allowMultipleAds="0" fallbackOnNoAd="0">
			<AdSystem version="4.2">Example Wrapper System</AdSystem>
			<Impression><![CDATA[https://example.com/impression]]></Impression>
			<VASTAdTagURI><![CDATA[https://example.com/tag]]></VASTAdTagURI>
			<Creatives>
				<Creative id="1">
					<Linear>
						<TrackingEvents>
							<Tracking event="start"><![CDATA[https://example.com/track]]></Tracking>
						</TrackingEvents>
					</Linear>
				</Creative>
			</Creatives>
			<AdVerifications>
				<Verification vendor="iabtechlab">
					<JavaScriptResource><![CDATA[https://example.com/verification.js]]></JavaScriptResource>
				</Verification>
			</AdVerifications>
		</Wrapper>
	</Ad>
</VAST>`

	result, err := Validate([]byte(xml), DisableHTTPValidators())
	if err != nil {
		t.Fatalf("validate returned error: %v", err)
	}

	assertStatus(t, result.Root, "Wrapper", StatusPass)
	assertStatus(t, result.Root, "TrackingEvents", StatusPass)
}

func TestValidate_NonLinearAndCompanion(t *testing.T) {
	resetCustom(t)
	xml := `<?xml version="1.0" encoding="UTF-8"?>
<VAST version="4.2">
	<Ad id="1">
		<InLine>
			<AdSystem>example</AdSystem>
			<Impression><![CDATA[https://example.com/imp]]></Impression>
			<AdTitle>Example</AdTitle>
			<Creatives>
				<Creative id="nl">
					<NonLinearAds>
						<NonLinear width="300" height="250">
							<StaticResource creativeType="image/png"><![CDATA[https://example.com/banner.png]]></StaticResource>
							<NonLinearClickThrough><![CDATA[https://example.com/click]]></NonLinearClickThrough>
						</NonLinear>
					</NonLinearAds>
				</Creative>
				<Creative id="comp">
					<CompanionAds>
						<Companion width="300" height="250">
							<StaticResource creativeType="image/png"><![CDATA[https://example.com/companion.png]]></StaticResource>
							<CompanionClickThrough><![CDATA[https://example.com/comp-click]]></CompanionClickThrough>
						</Companion>
					</CompanionAds>
				</Creative>
			</Creatives>
		</InLine>
	</Ad>
</VAST>`

	result, err := Validate([]byte(xml))
	if err != nil {
		t.Fatalf("validate returned error: %v", err)
	}

	assertStatus(t, result.Root, "NonLinear", StatusPass)
	assertStatus(t, result.Root, "Companion", StatusPass)
}

func findNode(root *NodeResult, name string) *NodeResult {
	if root == nil {
		return nil
	}
	if root.Node == name {
		return root
	}
	for _, child := range root.Children {
		if found := findNode(child, name); found != nil {
			return found
		}
	}
	return nil
}

func assertStatus(t *testing.T, root *NodeResult, name string, expected ResultStatus) {
	node := findNode(root, name)
	if node == nil {
		t.Fatalf("expected to find node %s", name)
	}
	analysis := node.Analyses[IABAnalysisCategory]
	if analysis == nil {
		t.Fatalf("node %s missing IAB analysis", name)
	}
	if analysis.Status != expected {
		t.Fatalf("expected status %s for node %s, got %s", expected, name, analysis.Status)
	}
}

func resetCustom(t *testing.T) {
	t.Helper()
	customMu.Lock()
	customValidators = map[string][]NodeValidatorFunc{}
	customMu.Unlock()
	HTTPValidatorRegistry.mu.Lock()
	HTTPValidatorRegistry.store = map[string][]HTTPValidatorFunc{}
	HTTPValidatorRegistry.mu.Unlock()
	registerBuiltInHTTPValidators()
	resetExtensionValidators()
}
