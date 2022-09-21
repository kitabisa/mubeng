package proxymanager

import "regexp"

var (
	manager     *ProxyManager
	placeholder = regexp.MustCompile(`\{\{.*?\}\}`)
)
