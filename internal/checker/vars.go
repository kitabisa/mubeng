package checker

import (
	"net/http"
)

var (
	client *http.Client
	ipinfo IPInfo

	endpoint = "https://ipinfo.io/json"
)
