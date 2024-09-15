package proxymanager

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/kitabisa/mubeng/pkg/helper"
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
		proxy := helper.Eval(scanner.Text())
		if _, value := keys[proxy]; !value {
			_, err = mubeng.Transport(placeholder.ReplaceAllString(proxy, ""))
			if err == nil {
				keys[proxy] = true
				manager.Proxies = append(manager.Proxies, proxy)
			}
		}
	}

	manager.Count()

	if manager.Length < 1 {
		return manager, fmt.Errorf("open %s: has no valid proxy URLs", filename)
	}

	return manager, scanner.Err()
}
