package proxymanager

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/kitabisa/mubeng/pkg/mubeng"
)

// ProxyManager defines the proxy list and current proxy position
type ProxyManager struct {
	CurrentIndex int
	filepath     string
	Length       int
	Proxies      []string
}

func init() {
	rand.Seed(time.Now().UnixNano())

	manager = &ProxyManager{CurrentIndex: -1}
}

// New initialize ProxyManager
func New(filename string) (*ProxyManager, error) {
	keys := make(map[string]bool)

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	manager.Proxies = []string{}
	manager.filepath = filename

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		proxy := scanner.Text()
		if _, value := keys[proxy]; !value {
			if _, err = mubeng.Transport(proxy); err == nil {
				keys[proxy] = true
				manager.Proxies = append(manager.Proxies, proxy)
			}
		}
	}

	manager.Length = len(manager.Proxies)
	if manager.Length < 1 {
		return manager, fmt.Errorf("open %s: has no valid proxy URLs", filename)
	}

	return manager, scanner.Err()
}
