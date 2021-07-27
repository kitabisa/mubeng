package server

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/elazarl/goproxy"
	"ktbs.dev/mubeng/pkg/mubeng"
)

// onRequest handles client request
func (p *Proxy) onRequest(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	resChan := make(chan *http.Response)
	errChan := make(chan error, 1)

	// Rotate proxy IP for every AFTER request
	if (rotate == "") || (ok >= p.Options.Rotate) {
		rotate = p.Options.ProxyManager.NextProxy()
		if ok >= p.Options.Rotate {
			ok = 1
		}
	} else {
		ok++
	}

	go func() {
		if (req.URL.Scheme != "http") && (req.URL.Scheme != "https") {
			errChan <- fmt.Errorf("Unsupported protocol scheme: %s", req.URL.Scheme)
			return
		}

		log.Debugf("%s %s %s", req.RemoteAddr, req.Method, req.URL)

		tr, err := mubeng.Transport(rotate)
		if err != nil {
			errChan <- err
			return
		}

		proxy := &mubeng.Proxy{
			Address:   rotate,
			Transport: tr,
		}

		client, req = proxy.New(req)
		client.Timeout = p.Options.Timeout
		if p.Options.Verbose {
			client.Transport = dump.RoundTripper(tr)
		}

		resp, err := client.Do(req)
		if err != nil {
			errChan <- err
			return
		}
		defer resp.Body.Close()

		// Copying response body
		buf, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			errChan <- err
			return
		}

		resp.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
		resChan <- resp
	}()

	select {
	case err := <-errChan:
		log.Errorf("%s %s", req.RemoteAddr, err)
		return req, goproxy.NewResponse(req, mime, http.StatusInternalServerError, "Proxy Server Error")
	case resp := <-resChan:
		log.Debug(req.RemoteAddr, " ", resp.Status)
		return req, resp
	}
}

// onConnect handles CONNECT method
func (p *Proxy) onConnect(host string, ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
	return goproxy.MitmConnect, host
}

// onResponse handles backend responses, and removing hop-by-hop headers
func (p *Proxy) onResponse(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
	for _, h := range mubeng.HopHeaders {
		resp.Header.Del(h)
	}

	return resp
}
