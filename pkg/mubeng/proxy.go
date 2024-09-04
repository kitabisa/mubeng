package mubeng

import (
	"net/http"
	"time"
)

// Proxy define the IP address value, [http.Transport] and other additional options.
type Proxy struct {
	Address      string
	MaxRedirects int
	Timeout      time.Duration
	Transport    *http.Transport
}
