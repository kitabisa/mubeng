package server

import (
	"context"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"time"

	// "crypto/tls"
	// "net"
	// "bufio"
	// "regexp"

	"github.com/elazarl/goproxy"
	"github.com/mbndr/logo"
	"ktbs.dev/mubeng/common"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Run proxy server with a user defined listener.
//
// An active log have 2 receivers, especially stdout and into file if opt.Output isn't empty.
// Then close the proxy server if it receives a signal that interrupts the program.
func Run(opt *common.Options) {
	cli := logo.NewReceiver(os.Stderr, "")
	cli.Color = true
	cli.Level = logo.DEBUG

	file, _ := logo.Open(opt.Output)
	out := logo.NewReceiver(file, "")
	out.Format = "%s: %s"

	log = logo.NewLogger(cli, out)

	// cer, err := tls.LoadX509KeyPair("mubeng.crt", "mubeng.key")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	handler := &Proxy{}
	handler.Options = opt
	handler.HTTPProxy = goproxy.NewProxyHttpServer()
	// goproxy.OkConnect = &goproxy.ConnectAction{Action: goproxy.ConnectAccept, TLSConfig: goproxy.TLSConfigFromCA(&cer)}
	// handler.HTTPProxy.OnRequest(goproxy.ReqHostMatches(regexp.MustCompile("^.*$"))).HandleConnect(goproxy.AlwaysMitm)
	// handler.HTTPProxy.OnRequest(goproxy.ReqHostMatches(regexp.MustCompile("^.*:80$"))).HijackConnect(func(req *http.Request, client net.Conn, ctx *goproxy.ProxyCtx) {
	// 	defer func() {
	// 		if e := recover(); e != nil {
	// 			ctx.Logf("error connecting to remote: %v", e)
	// 			client.Write([]byte("HTTP/1.1 500 Cannot reach destination\r\n\r\n"))
	// 		}
	// 		client.Close()
	// 	}()
	// 	clientBuf := bufio.NewReadWriter(bufio.NewReader(client), bufio.NewWriter(client))
	// 	remote, err := net.Dial("tcp", req.URL.Host)
	// 	log.Fatal(err)
	// 	client.Write([]byte("HTTP/1.1 200 Ok\r\n\r\n"))
	// 	remoteBuf := bufio.NewReadWriter(bufio.NewReader(remote), bufio.NewWriter(remote))
	// 	for {
	// 		req, err := http.ReadRequest(clientBuf.Reader)
	// 		log.Fatal(err)
	// 		log.Fatal(req.Write(remoteBuf))
	// 		log.Fatal(remoteBuf.Flush())
	// 		resp, err := http.ReadResponse(remoteBuf.Reader, req)
	// 		log.Fatal(err)
	// 		log.Fatal(resp.Write(clientBuf.Writer))
	// 		log.Fatal(clientBuf.Flush())
	// 	}
	// })
	handler.HTTPProxy.OnRequest().DoFunc(handler.OnRequest)
	handler.HTTPProxy.OnRequest().HandleConnectFunc(handler.OnConnect)
	handler.HTTPProxy.OnResponse().DoFunc(handler.OnResponse)

	server := &http.Server{
		Addr:    opt.Address,
		Handler: handler.HTTPProxy,
	}

	go func() {
		log.Info("Starting proxy server on ", opt.Address)
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = server.Shutdown(ctx)
}
