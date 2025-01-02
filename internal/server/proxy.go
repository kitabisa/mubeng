package server

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/elazarl/goproxy"
	"github.com/kitabisa/mubeng/common"
	"github.com/kitabisa/mubeng/internal/proxygateway"
)

// Proxy as ServeMux in proxy server handler.
type Proxy struct {
	HTTPProxy *goproxy.ProxyHttpServer
	Options   *common.Options
	Gateways  map[string]*proxygateway.ProxyGateway
	mu        sync.RWMutex
}

func (p *Proxy) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	p.mu.Lock()
	defer p.mu.Unlock()

	for _, gateway := range p.Gateways {
		if err := gateway.Close(ctx); err != nil {
			log.Error(fmt.Errorf("could not close gateway: %w", err))
		}
	}
}
