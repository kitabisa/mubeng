package mubeng

import (
	"net/http"
	"net/url"

	"golang.org/x/net/proxy"
)

// Transport is switching transport between HTTP or SOCKS5 proxies
// depending on the protocol scheme, returning value of http.Transport with Dialer or Proxy.
func Transport(p string) (tr *http.Transport, err error) {
	proxyURL, err := url.Parse(p)
	if err != nil {
		return nil, err
	}

	dialer, err := proxy.SOCKS5("tcp", p, nil, nil)
	if err != nil {
		return nil, err
	}

	switch proxyURL.String() {
	case "socks5":
		tr = &http.Transport{
			Dial: dialer.Dial,
		}
	default:
		tr = &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}
	}

	tr.DisableKeepAlives = true

	return tr, nil
}
