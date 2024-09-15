package server

import (
	"net/http"
	"os"
	"os/signal"

	"github.com/elazarl/goproxy"
	"github.com/henvic/httpretty"
	"github.com/kitabisa/mubeng/common"
	"github.com/mbndr/logo"
)

// Run proxy server with a user defined listener.
//
// An active log have 2 receivers, especially stdout and into file if opt.Output isn't empty.
// Then close the proxy server if it receives a signal that interrupts the program.
func Run(opt *common.Options) {
	var recs []*logo.Receiver

	cli := logo.NewReceiver(os.Stderr, "")
	cli.Color = true
	cli.Level = logo.DEBUG
	recs = append(recs, cli)

	file, err := logo.Open(opt.Output)
	if err == nil {
		out := logo.NewReceiver(file, "")
		out.Format = "%s: %s"
		recs = append(recs, out)
	}

	dump = &httpretty.Logger{
		RequestHeader:  true,
		ResponseHeader: true,
		Colors:         true,
	}

	handler = &Proxy{}
	handler.Options = opt
	handler.HTTPProxy = goproxy.NewProxyHttpServer()
	handler.HTTPProxy.OnRequest().DoFunc(handler.onRequest)
	handler.HTTPProxy.OnRequest().HandleConnectFunc(handler.onConnect)
	handler.HTTPProxy.OnResponse().DoFunc(handler.onResponse)
	handler.HTTPProxy.NonproxyHandler = http.HandlerFunc(nonProxy)

	server = &http.Server{
		Addr:    opt.Address,
		Handler: handler.HTTPProxy,
	}

	log = logo.NewLogger(recs...)

	if opt.Watch {
		watcher, err := opt.ProxyManager.Watch()
		if err != nil {
			log.Fatal(err)
		}
		defer watcher.Close()

		go watch(watcher)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	go interrupt(stop)

	log.Infof("%d proxies loaded", opt.ProxyManager.Count())

	log.Infof("[PID: %d] Starting proxy server on %s", os.Getpid(), opt.Address)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
