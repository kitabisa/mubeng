package server

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/elazarl/goproxy"
	"ktbs.dev/mubeng/common"
	"ktbs.dev/mubeng/pkg/helper"
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

	rotate = helper.EvalFunc(rotate)
	resChan := make(chan interface{})

	go func(r *http.Request) {
		if (r.URL.Scheme != "http") && (r.URL.Scheme != "https") {
			resChan <- fmt.Errorf("Unsupported protocol scheme: %s", r.URL.Scheme)
			return
		}

		log.Debugf("%s %s %s", r.RemoteAddr, r.Method, r.URL)

		tr, err := mubeng.Transport(rotate)
		if err != nil {
			resChan <- err
			return
		}

		proxy := &mubeng.Proxy{
			Address:   rotate,
			Transport: tr,
		}

		client, err := proxy.New(req)
		if err != nil {
			resChan <- err
			return
		}

		client.Timeout = p.Options.Timeout
		if p.Options.Verbose {
			client.Transport = dump.RoundTripper(tr)
		}

		resp, err := client.Do(req)
		if err != nil {
			resChan <- err
			return
		}
		defer resp.Body.Close()

		buf, err := io.ReadAll(resp.Body)
		if err != nil {
			resChan <- err
			return
		}

		resp.Body = io.NopCloser(bytes.NewBuffer(buf))

		resChan <- resp
	}(req)

	var resp *http.Response

	res := <-resChan
	switch res := res.(type) {
	case *http.Response:
		resp = res
		log.Debug(req.RemoteAddr, " ", resp.Status)
	case error:
		err := res
		log.Errorf("%s %s", req.RemoteAddr, err)
		resp = goproxy.NewResponse(req, mime, http.StatusBadGateway, "Proxy server error")
	}

	return req, resp
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

// nonProxy handles non-proxy requests
func nonProxy(w http.ResponseWriter, req *http.Request) {
	if common.Version != "" {
		w.Header().Add("X-Mubeng-Version", common.Version)
	}

	if req.URL.Path == "/cert" {
		w.Header().Add("Content-Type", "application/octet-stream")
		w.Header().Add("Content-Disposition", fmt.Sprint("attachment; filename=", "goproxy-cacert.der"))
		w.WriteHeader(http.StatusOK)

		if _, err := w.Write(goproxy.GoproxyCa.Certificate[0]); err != nil {
			http.Error(w, "Failed to get proxy certificate authority.", 500)
			log.Errorf("%s %s %s %s", req.RemoteAddr, req.Method, req.URL, err.Error())
		}

		return
	}

	http.Error(w, "This is a mubeng proxy server. Does not respond to non-proxy requests.", 500)
}
