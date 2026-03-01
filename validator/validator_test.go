package validator

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

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
}
