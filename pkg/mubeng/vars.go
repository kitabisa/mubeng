package mubeng

import "net/http"

var (
	client *http.Client

	// HopHeaders are meaningful only for a single transport-level connection, and are not stored by caches or forwarded by proxies.
	HopHeaders = []string{
		"Connection",
		"Keep-Alive",
		"Proxy-Authenticate",
		"Proxy-Authorization",
		"Proxy-Connection",
		"Te", // canonicalized version of "TE"
		"Trailers",
		"Transfer-Encoding",
		"Upgrade",
	}
)
