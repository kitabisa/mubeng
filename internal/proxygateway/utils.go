package proxygateway

import (
	"fmt"
	"strings"

	"net/url"
)

// isHTTPURL checks if a given URL string is HTTP or HTTPS.
func isHTTPURL(s string) bool {
	return (strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://"))
}

// GetBaseURL returns the base URL and parsed URL from a given URL string.
func GetBaseURL(s string) (string, *url.URL, error) {
	var u string

	if !isHTTPURL(s) {
		return u, nil, fmt.Errorf("must start with http:// or https://")
	}

	parsedURL, err := url.Parse(s)
	if err != nil {
		return u, nil, err
	}

	u = fmt.Sprintf("%s://%s", parsedURL.Scheme, parsedURL.Host)

	return u, parsedURL, nil
}
