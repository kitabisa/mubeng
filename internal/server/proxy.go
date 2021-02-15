package server

import "time"

// Proxy as ServeMux in proxy server handler.
type Proxy struct {
	list    []string
	rotate  int
	timeout time.Duration
	verbose bool
}
