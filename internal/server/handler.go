package server

import (
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
		return req, goproxy.NewResponse(req, "text/plain", http.StatusInternalServerError, "Server Error")
	}

	proxy := &mubeng.Proxy{
		Address:   rotate,
		Transport: tr,
		Verbose:   p.Options.Verbose,
		Color:     true,
	}

	client, req = proxy.New(req)
	client.Timeout = p.Options.Timeout

	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("%s %s", req.RemoteAddr, err)
		return req, goproxy.NewResponse(req, "text/plain", http.StatusInternalServerError, "Server Error")
	}
	defer resp.Body.Close()

	log.Debug(req.RemoteAddr, " ", resp.Status)
	return req, resp
}

// OnResponse handles from backend response, and removing hop-by-hop headers
func (p *Proxy) OnResponse(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
	for _, h := range mubeng.HopHeaders {
		resp.Header.Del(h)
	}

	return resp
}

// OnConnect handles CONNECT method
func (p *Proxy) OnConnect(host string, ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
	return goproxy.MitmConnect, host
}

// // ServeHTTP as request handler.
// func (p *Proxy) ServeHTTP(wr http.ResponseWriter, req *http.Request) {

// }

// ServeHTTP as request handler.
// func (p *Proxy) ServeHTTP(wr http.ResponseWriter, req *http.Request) {
// 	log.Debugf("%s %s %s", req.RemoteAddr, req.Method, req.URL)

// 	if !strings.HasPrefix(req.URL.Scheme, "http") {
// 		msg := "Unsupported protocol scheme: " + req.URL.Scheme
// 		http.Error(wr, msg, http.StatusBadRequest)
// 		log.Warnf("%s %s", req.RemoteAddr, msg)
// 		return
// 	}

// 	// Rotate proxy IP every AFTER request
// 	if (rotate == "") || (ok >= p.Options.Rotate) {
// 		rotate = p.Options.List[rand.Intn(len(p.Options.List))]
// 		if ok >= p.Options.Rotate {
// 			ok = 1
// 		}
// 	} else {
// 		ok++
// 	}

// 	tr, err := mubeng.Transport(rotate)
// 	if err != nil {
// 		http.Error(wr, "Transport Error", http.StatusInternalServerError)
// 		log.Errorf("%s %s", req.RemoteAddr, err)
// 		return
// 	}

// 	proxy := &mubeng.Proxy{
// 		Address:   rotate,
// 		Transport: tr,
// 		Verbose:   p.Options.Verbose,
// 		Color:     true,
// 	}

// 	client, req = proxy.New(req)
// 	client.Timeout = p.Options.Timeout

// 	resp, err := client.Do(req)
// 	if err != nil {
// 		http.Error(wr, "Server Error", http.StatusInternalServerError)
// 		log.Errorf("%s %s", req.RemoteAddr, err)
// 		return
// 	}
// 	defer resp.Body.Close()

// 	log.Debug(req.RemoteAddr, " ", resp.Status)

// 	// Removing Hop-by-hop headers
// 	for _, h := range mubeng.HopHeaders {
// 		resp.Header.Del(h)
// 	}

// 	// Copy headers
// 	for k, vv := range resp.Header {
// 		for _, v := range vv {
// 			wr.Header().Add(k, v)
// 		}
// 	}

// 	wr.WriteHeader(resp.StatusCode)
// 	if _, err := io.Copy(wr, resp.Body); err != nil {
// 		log.Fatalf("Error! %s", err)
// 	}
// }
