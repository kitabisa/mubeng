package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/mbndr/logo"
	"ktbs.dev/mubeng/common"
)

// Run proxy server with a user defined listener.
//
// An active log have 2 receivers, especially stdout and into file if opt.Output isn't empty.
// Then close the proxy server if it receives a signal that interrupts the program.
func Run(opt *common.Options) {
	cli := logo.NewReceiver(os.Stdout, "")
	cli.Color = true
	cli.Level = logo.DEBUG

	file, _ := logo.Open(opt.Output)
	out := logo.NewReceiver(file, "")
	out.Format = "%s: %s"

	log = logo.NewLogger(cli, out)

	handler := &Proxy{
		list:    opt.List,
		verbose: opt.Verbose,
	}

	server := &http.Server{
		Addr:    opt.Address,
		Handler: handler,
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
