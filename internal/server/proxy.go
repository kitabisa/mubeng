package server

import (
	"github.com/elazarl/goproxy"
	"ktbs.dev/mubeng/common"
)

// Proxy as ServeMux in proxy server handler.
type Proxy struct {
	HTTPProxy *goproxy.ProxyHttpServer
	Options   *common.Options
}
