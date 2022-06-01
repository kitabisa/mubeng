package mubeng

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"

	"h12.io/socks"
)

// Transport to auto-switch transport between HTTP/S or SOCKS v4(A) & v5 proxies.
// Depending on the protocol scheme, returning value of http.Transport with Dialer or Proxy.
func Transport(p string) (tr *http.Transport, err error) {
	proxyURL, err := url.Parse(p)
	if err != nil {
		return nil, err
	}

	switch proxyURL.Scheme {
	case "socks4", "socks4a", "socks5":
		tr = &http.Transport{
			Dial: socks.Dial(p),
		}
	case "http", "https":
		tr = &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}
	default:
		return nil, fmt.Errorf("unsupported proxy protocol scheme: %s", proxyURL.Scheme)
	}

	tr.DisableKeepAlives = true
	tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	return tr, nil
}
