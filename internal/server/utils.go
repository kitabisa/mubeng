package server

import (
	"context"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
)

// Stop will terminate proxy server
func Stop(ctx context.Context) {
	_ = server.Shutdown(ctx)
}

func interrupt(sig chan os.Signal) {
	<-sig
	log.Warn("Interuppted. Exiting...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	Stop(ctx)
}

func watch(w *fsnotify.Watcher) {
	for {
		select {
		case event := <-w.Events:
			if event.Op == 2 {
				log.Info("Proxy file has changed, reloading...")

				err := handler.Options.ProxyManager.Reload()
				if err != nil {
					log.Fatal(err)
				}
			}
		case err := <-w.Errors:
			log.Fatal(err)
		}
	}
}
