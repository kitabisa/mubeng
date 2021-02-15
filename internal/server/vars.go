package server

import (
	"net/http"

	"github.com/mbndr/logo"
)

var (
	rotate string
	client *http.Client
	log    *logo.Logger
	ok     = 1
)
