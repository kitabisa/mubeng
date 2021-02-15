package server

import (
	"io"
	"math/rand"
	"net/http"
	"strings"

	"ktbs.dev/mubeng/pkg/mubeng"
)

// ServeHTTP as request handler.
func (p *Proxy) ServeHTTP(wr http.ResponseWriter, req *http.Request) {
	log.Debugf("%s %s %s", req.RemoteAddr, req.Method, req.URL)

	if !strings.HasPrefix(req.URL.Scheme, "http") {
		msg := "Unsupported protocol scheme: " + req.URL.Scheme
		http.Error(wr, msg, http.StatusBadRequest)
		log.Warnf("%s %s", req.RemoteAddr, msg)
		return
	}

	// Rotate proxy IP every AFTER request
	if (rotate == "") || (ok >= p.rotate) {
		rotate = p.list[rand.Intn(len(p.list))]
		if ok >= p.rotate {
			ok = 1
		}
	} else {
		ok++
	}

	tr, err := mubeng.Transport(rotate)
	if err != nil {
		http.Error(wr, "Transport Error", http.StatusInternalServerError)
		log.Errorf("%s %s", req.RemoteAddr, err)
		return
	}

	proxy := &mubeng.Proxy{
		Address:   rotate,
		Transport: tr,
		Verbose:   p.verbose,
		Color:     true,
	}

	client, req = proxy.New(req)
	client.Timeout = p.timeout

	resp, err := client.Do(req)
	if err != nil {
		http.Error(wr, "Server Error", http.StatusInternalServerError)
		log.Errorf("%s %s", req.RemoteAddr, err)
		return
	}
	defer resp.Body.Close()

	log.Debug(req.RemoteAddr, " ", resp.Status)

	// Removing Hop-by-hop headers
	for _, h := range mubeng.HopHeaders {
		resp.Header.Del(h)
	}

	// Copy headers
	for k, vv := range resp.Header {
		for _, v := range vv {
			wr.Header().Add(k, v)
		}
	}

	wr.WriteHeader(resp.StatusCode)
	if _, err := io.Copy(wr, resp.Body); err != nil {
		log.Fatalf("Error! %s", err)
	}
}
