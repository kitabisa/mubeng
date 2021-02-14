package server

import (
	"net/http"

	"github.com/mbndr/logo"
)

var (
	client *http.Client
	log    *logo.Logger
)
