package proxymanager

import (
	"math/rand"

	"github.com/fsnotify/fsnotify"
	"github.com/kitabisa/mubeng/pkg/helper"
)

// NextProxy will navigate the next proxy to use
func (p *ProxyManager) NextProxy() string {
	p.CurrentIndex++
	if p.CurrentIndex > len(p.Proxies)-1 {
		p.CurrentIndex = 0
	}

	proxy := p.Proxies[p.CurrentIndex]

	return proxy
}

// RandomProxy will choose a proxy randomly from the list
func (p *ProxyManager) RandomProxy() string {
	return p.Proxies[rand.Intn(len(p.Proxies))]
}

// Rotate proxy based on method
//
// Valid methods are "sequent" and "random", default return empty string.
func (p *ProxyManager) Rotate(method string) string {
	var proxy string

	switch method {
	case "sequent":
		proxy = p.NextProxy()
	case "random":
		proxy = p.RandomProxy()
	}

	if proxy != "" {
		proxy = helper.EvalFunc(proxy)
	}

	return proxy
}

// Watch proxy file from events
func (p *ProxyManager) Watch() (*fsnotify.Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return watcher, err
	}

	if err := watcher.Add(p.filepath); err != nil {
		return watcher, err
	}

	return watcher, nil
}

// Reload proxy pool
func (p *ProxyManager) Reload() error {
	i := p.CurrentIndex

	p, err := New(p.filepath)
	if err != nil {
		return err
	}
	p.CurrentIndex = i

	return nil
}
