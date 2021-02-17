package server

import (
	"bytes"
	"io/ioutil"
	"math/rand"
	"net/http"

	"github.com/elazarl/goproxy"
	"ktbs.dev/mubeng/pkg/mubeng"
)

// OnRequest handles client request
func (p *Proxy) OnRequest(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	log.Debugf("%s %s %s", req.RemoteAddr, req.Method, req.URL)

	// Rotate proxy IP every AFTER request
	if (rotate == "") || (ok >= p.Options.Rotate) {
		rotate = p.Options.List[rand.Intn(len(p.Options.List))]
		if ok >= p.Options.Rotate {
			ok = 1
		}
	} else {
		ok++
	}

	tr, err := mubeng.Transport(rotate)
	if err != nil {
		log.Errorf("%s %s", req.RemoteAddr, err)
		return req, goproxy.NewResponse(req, "text/plain", http.StatusInternalServerError, "Proxy transport error")
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
		log.Errorf("%s %s", req.RemoteAddr, err)
		return req, goproxy.NewResponse(req, "text/plain", http.StatusBadGateway, "Proxy error: "+err.Error())
	}
	defer resp.Body.Close()

	// Copying response body
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("%s %s", req.RemoteAddr, err)
		return req, goproxy.NewResponse(req, "text/plain", http.StatusInternalServerError, "Proxy Error: "+err.Error())
	}

	resp.Body = ioutil.NopCloser(bytes.NewBuffer(buf))

	// Removing hop-by-hop headers
	for _, h := range mubeng.HopHeaders {
		resp.Header.Del(h)
	}

	log.Debug(req.RemoteAddr, " ", resp.Status)
	return req, resp
}

// OnConnect handles CONNECT method
func (p *Proxy) OnConnect(host string, ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
	return goproxy.MitmConnect, host
}

// OnResponse handles backend responses, and removing hop-by-hop headers
func (p *Proxy) OnResponse(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
	for _, h := range mubeng.HopHeaders {
		resp.Header.Del(h)
	}

	return resp
}
