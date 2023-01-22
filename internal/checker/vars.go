package checker

import (
	"net/http"
)

var (
	client *http.Client
	myip   myIP

	endpoint = "https://api.myip.com/"
)
