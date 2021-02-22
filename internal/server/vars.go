package server

import (
	"net/http"

	"github.com/henvic/httpretty"
	"github.com/mbndr/logo"
)

var (
	rotate string
	server *http.Server
	client *http.Client
	dump   *httpretty.Logger
	mime   = "text/plain"
	log    *logo.Logger
	ok     = 1
)
