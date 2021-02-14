package server

import "time"

// Proxy as ServeMux in proxy server handler.
type Proxy struct {
	list    []string
	timeout time.Duration
	verbose bool
}
