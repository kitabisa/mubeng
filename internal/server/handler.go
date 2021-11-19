package server

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/elazarl/goproxy"
	"ktbs.dev/mubeng/pkg/mubeng"
)

// onRequest handles client request
func (p *Proxy) onRequest(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	if p.Options.Sync {
		mutex.Lock()
		defer mutex.Unlock()
	}

	// Rotate proxy IP for every AFTER request
	if (rotate == "") || (ok >= p.Options.Rotate) {
		if p.Options.Method == "sequent" {
			rotate = p.Options.ProxyManager.NextProxy()
		}

		if p.Options.Method == "random" {
			rotate = p.Options.ProxyManager.RandomProxy()
		}

		if ok >= p.Options.Rotate {
			ok = 1
		}
	} else {
		ok++
	}

	resChan := make(chan *http.Response)
	errChan := make(chan error, 1)

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
		return req, goproxy.NewResponse(req, mime, http.StatusBadGateway, "Proxy server error")
	case resp := <-resChan:
		log.Debug(req.RemoteAddr, " ", resp.Status)
		return req, resp
	}
}

// onConnect handles CONNECT method
func (p *Proxy) onConnect(host string, ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
	if p.Options.Auth != "" {
		auth := ctx.Req.Header.Get("Proxy-Authorization")
		if auth != "" {
			creds := strings.SplitN(auth, " ", 2)
			if len(creds) != 2 {
				return goproxy.RejectConnect, host
			}

			auth, err := base64.StdEncoding.DecodeString(creds[1])
			if err != nil {
				log.Warnf("%s: Error decoding proxy authorization", ctx.Req.RemoteAddr)
				return goproxy.RejectConnect, host
			}

			if string(auth) != p.Options.Auth {
				log.Errorf("%s: Invalid proxy authorization", ctx.Req.RemoteAddr)
				return goproxy.RejectConnect, host
			}
		} else {
			log.Warnf("%s: Unathorized proxy request to %s", ctx.Req.RemoteAddr, host)
			return goproxy.RejectConnect, host
		}
	}

	return goproxy.MitmConnect, host
}

// onResponse handles backend responses, and removing hop-by-hop headers
func (p *Proxy) onResponse(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
	for _, h := range mubeng.HopHeaders {
		resp.Header.Del(h)
	}

	return resp
}
