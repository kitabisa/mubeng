package mubeng

import "net/http"

// Proxy define the IP address value, [http.Transport] and other additional options.
type Proxy struct {
	Address      string
	MaxRedirects int
	Transport    *http.Transport
}
