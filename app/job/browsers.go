package job

import (
	"fmt"
	"strings"
)

// TODO brave, chrome, chromium, edge, firefox,
// TODO opera, safari, vivaldi, whale

// BrowserForCookies is type to represent browser code
type BrowserForCookies string

const (
	// Chrome --cookies-from-browser chrome
	Chrome BrowserForCookies = "chrome"
)

// ParseCookiesFromBrowser converts a string to CookiesFromBrowser enum value
// Returns error if the input string doesn't match any supported browser
func ParseCookiesFromBrowser(value string) (BrowserForCookies, error) {
	browser := BrowserForCookies(value)
	switch browser {
	case Chrome:
		return browser, nil
	default:
		return "", fmt.Errorf("unsupported browser: %s. Valid values: %s",
			value, strings.Join([]string{string(Chrome)}, ", "))
	}
}
