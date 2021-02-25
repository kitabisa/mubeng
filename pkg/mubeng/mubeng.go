package mubeng

import (
	"net"
	"net/http"
	"net/url"
	"strings"
)

// New define HTTP client & request of http.Request itself.
//
// also removes Hop-by-hop headers when it is sent to backend (see http://www.w3.org/Protocols/rfc2616/rfc2616-sec13.html),
// then add X-Forwarded-For header value with the IP address value of rotator proxy IP.
func (proxy *Proxy) New(req *http.Request) (*http.Client, *http.Request) {
	client = &http.Client{Transport: proxy.Transport}

	// http: Request.RequestURI can't be set in client requests.
	// http://golang.org/src/pkg/net/http/client.go
	req.RequestURI = ""

	for _, h := range HopHeaders {
		req.Header.Del(h)
	}

	proxyURL, _ := url.Parse(proxy.Address)

	if host, _, err := net.SplitHostPort(proxyURL.Host); err == nil {
		if prior, ok := req.Header["X-Forwarded-For"]; ok {
			host = strings.Join(prior, ", ") + ", " + host
		}
		req.Header.Set("X-Forwarded-For", host)
	}

	req.Header.Set("X-Forwarded-Proto", req.URL.Scheme)

	return client, req
}
