package server

import (
	"github.com/elazarl/goproxy"
	"github.com/kitabisa/mubeng/common"
)

// Proxy as ServeMux in proxy server handler.
type Proxy struct {
	HTTPProxy *goproxy.ProxyHttpServer
	Options   *common.Options
}
