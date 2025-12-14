package util

import (
	"regexp"
	"strings"
)

type BrowserOS struct {
	Browser string
	OS      string
}

func ParseUserAgent(userAgent string) BrowserOS {
	ua := strings.ToLower(userAgent)

	browserPatterns := []struct {
		name    string
		pattern *regexp.Regexp
	}{
		{"Edge", regexp.MustCompile(`edg(?:e|a?)?/`)},
		{"Opera", regexp.MustCompile(`op(?:era|r)[/\s]`)},
		{"Samsung", regexp.MustCompile(`samsungbrowser`)},
		{"UCBrowser", regexp.MustCompile(`ucbrowser`)},
		{"Chrome", regexp.MustCompile(`chrome|crios`)},
		{"Firefox", regexp.MustCompile(`firefox|fxios`)},
		{"IE", regexp.MustCompile(`trident|msie`)},
		{"Safari", regexp.MustCompile(`safari|version/`)},
	}

	var browser string
	for _, pattern := range browserPatterns {
		if pattern.pattern.MatchString(ua) {
			browser = pattern.name
			break
		}
	}

	if browser == "" {
		browser = "Others"
	}

	osPatterns := []struct {
		name    string
		pattern *regexp.Regexp
	}{
		{"iOS", regexp.MustCompile(`iphone|ipad|ipod`)},
		{"Android", regexp.MustCompile(`android`)},
		{"Windows", regexp.MustCompile(`windows nt`)},
		{"macOS", regexp.MustCompile(`mac os x`)},
		{"ChromeOS", regexp.MustCompile(`cros`)},
		{"Linux", regexp.MustCompile(`linux`)},
	}

	var os string
	for _, pattern := range osPatterns {
		if pattern.pattern.MatchString(ua) {
			os = pattern.name
			break
		}
	}

	if os == "" {
		os = "Others"
	}

	if strings.Contains(ua, "mobile") && os == "Windows" {
		os = "Windows Mobile"
	}

	return BrowserOS{
		Browser: browser,
		OS:      os,
	}
}
