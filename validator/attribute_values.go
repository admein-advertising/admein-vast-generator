package validator

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

var (
	durationPattern = regexp.MustCompile(`^\d{2}:\d{2}:\d{2}(\.\d{3})?$`)
	percentPattern  = regexp.MustCompile(`^\d{1,3}(\.\d+)?%$`)
	patternCache    sync.Map
)

func validateAttributeValue(attrName, value string, spec *AttributeSpec) []string {
	if spec == nil || spec.Value == nil {
		return nil
	}
	valSpec := spec.Value
	var errs []string
	if valSpec.Type != "" {
		if err := validateAttributeType(value, valSpec.Type); err != "" {
			errs = append(errs, fmt.Sprintf("attribute %s expects %s: %s", attrName, valSpec.Type, err))
		}
	}
	if len(valSpec.AllowedValues) > 0 {
		if !containsString(valSpec.AllowedValues, value) {
			errs = append(errs, fmt.Sprintf("attribute %s must be one of %v", attrName, valSpec.AllowedValues))
		}
	}
	if valSpec.Pattern != "" {
		re, err := getCachedPattern(valSpec.Pattern)
		if err != nil {
			errs = append(errs, fmt.Sprintf("attribute %s misconfigured pattern: %v", attrName, err))
		} else if !re.MatchString(value) {
			errs = append(errs, fmt.Sprintf("attribute %s must match pattern %s", attrName, valSpec.Pattern))
		}
	}
	return errs
}

func validateAttributeType(value string, attrType AttributeType) string {
	switch attrType {
	case AttributeTypeString, AttributeTypeToken:
		return ""
	case AttributeTypeBoolean:
		lower := strings.ToLower(value)
		if lower == "true" || lower == "false" || lower == "1" || lower == "0" {
			return ""
		}
		return "value must be true, false, 1, or 0"
	case AttributeTypeInteger:
		if _, err := strconv.Atoi(value); err != nil {
			return "value must be an integer"
		}
		return ""
	case AttributeTypeNonNegativeInteger:
		if err := expectIntRange(value, 0, -1); err != "" {
			return err
		}
		return ""
	case AttributeTypePositiveInteger:
		if err := expectIntRange(value, 1, -1); err != "" {
			return err
		}
		return ""
	case AttributeTypeFloat:
		if _, err := strconv.ParseFloat(value, 64); err != nil {
			return "value must be a decimal number"
		}
		return ""
	case AttributeTypeDuration:
		if durationPattern.MatchString(value) {
			return ""
		}
		return "duration must use HH:MM:SS or HH:MM:SS.mmm"
	case AttributeTypeTimecode:
		if durationPattern.MatchString(value) {
			return ""
		}
		return "timecode must use HH:MM:SS or HH:MM:SS.mmm"
	case AttributeTypeTimeOffset:
		if durationPattern.MatchString(value) || percentPattern.MatchString(value) || isKeyword(value, []string{"start", "end"}) {
			return ""
		}
		return "timeOffset must be HH:MM:SS(.mmm), percentage, start, or end"
	case AttributeTypeURI:
		if parsed, err := url.Parse(value); err == nil && parsed.Scheme != "" && parsed.Host != "" {
			return ""
		}
		return "value must be a valid URI"
	default:
		return ""
	}
}

func containsString(values []string, needle string) bool {
	for _, candidate := range values {
		if candidate == needle {
			return true
		}
	}
	return false
}

func isKeyword(value string, accepted []string) bool {
	lower := strings.ToLower(value)
	for _, candidate := range accepted {
		if lower == strings.ToLower(candidate) {
			return true
		}
	}
	return false
}

func expectIntRange(value string, min int, max int) string {
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return "value must be an integer"
	}
	if min != -1 && parsed < min {
		return fmt.Sprintf("value must be >= %d", min)
	}
	if max != -1 && parsed > max {
		return fmt.Sprintf("value must be <= %d", max)
	}
	return ""
}

func getCachedPattern(pattern string) (*regexp.Regexp, error) {
	if cached, ok := patternCache.Load(pattern); ok {
		if compiled, ok := cached.(*regexp.Regexp); ok {
			return compiled, nil
		}
	}
	compiled, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	patternCache.Store(pattern, compiled)
	return compiled, nil
}
