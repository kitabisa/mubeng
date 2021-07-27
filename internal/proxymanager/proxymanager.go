package proxymanager

import (
	"bufio"
	"errors"
	"fmt"
	"ktbs.dev/mubeng/pkg/mubeng"
	"math/rand"
	"os"
	"time"
)

type ProxyManager struct {
	Proxies      []string
	CurrentIndex int
}

func NewManager(filename string) (*ProxyManager, error) {
	file, openErr := os.Open(filename)
	if openErr != nil {
		return nil, openErr
	}
	defer file.Close()
	manager := &ProxyManager{Proxies: []string{}, CurrentIndex: 0}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		proxy := scanner.Text()
		_, err := mubeng.Transport(proxy)
		if err == nil {
			manager.Proxies = append(manager.Proxies, proxy)
		}
	}
	if len(manager.Proxies) < 1 {
		return manager, fmt.Errorf("open %s: has no valid proxy URLs", filename)
	}
	return manager, scanner.Err()
}

func (p *ProxyManager) NextProxy() (string, error) {
	if len(p.Proxies) == 0 {
		return "", errors.New("ProxyManager.Proxies is empty, load proxies")
	}
	p.CurrentIndex++
	if p.CurrentIndex > len(p.Proxies)-1 {
		p.CurrentIndex = 0
	}
	proxy := p.Proxies[p.CurrentIndex]
	return proxy, nil
}

func (p *ProxyManager) RandomProxy() (string, error) {
	if len(p.Proxies) == 0 {
		return "", errors.New("ProxyManager.Proxies is empty, load proxies")
	}
	rand.Seed(time.Now().UnixNano())
	return p.Proxies[rand.Intn(len(p.Proxies))], nil
}
