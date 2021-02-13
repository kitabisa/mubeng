package mubeng

import (
	"net"
	"net/http"
	"strings"

	"github.com/henvic/httpretty"
)

// New define HTTP client & request of http.Request itself.
// Dump HTTP request & responses if verbose mode is enabled using httpretty package as http.Transport
//
// also removes Hop-by-hop headers when it is sent to backend (see http://www.w3.org/Protocols/rfc2616/rfc2616-sec13.html),
// then add X-Forwarded-For header value with the IP address value of rotator proxy IP.
func (proxy *Proxy) New(req *http.Request) (*http.Client, *http.Request) {
	log := &httpretty.Logger{
		ResponseHeader: proxy.Verbose,
		RequestHeader:  proxy.Verbose,
		Colors:         proxy.Color,
	}

	client = &http.Client{Transport: proxy.Transport}
	if proxy.Verbose {
		client.Transport = log.RoundTripper(proxy.Transport)
	}

	// http: Request.RequestURI can't be set in client requests.
	// http://golang.org/src/pkg/net/http/client.go
	req.RequestURI = ""

	for _, h := range HopHeaders {
		req.Header.Del(h)
	}

	if host, _, err := net.SplitHostPort(proxy.Address); err == nil {
		if prior, ok := req.Header["X-Forwarded-For"]; ok {
			host = strings.Join(prior, ", ") + ", " + host
		}
		req.Header.Set("X-Forwarded-For", host)
	}

	return client, req
}
