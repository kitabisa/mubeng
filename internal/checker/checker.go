package checker

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/logrusorgru/aurora"
	"github.com/sourcegraph/conc/pool"
	"github.com/kitabisa/mubeng/common"
	"github.com/kitabisa/mubeng/pkg/helper"
	"github.com/kitabisa/mubeng/pkg/mubeng"
)

// Do checks proxy from list.
//
// Displays proxies that have died if verbose mode is enabled,
// or save live proxies into user defined files.
func Do(opt *common.Options) {
	p := pool.New().WithMaxGoroutines(opt.Goroutine)

	for _, proxy := range opt.ProxyManager.Proxies {
		address := helper.EvalFunc(proxy)

		p.Go(func() {
			addr, err := check(address, opt.Timeout)
			if len(opt.Countries) > 0 && !isMatchCC(opt.Countries, addr.Country) {
				return
			}

			if err != nil {
				if opt.Verbose {
					fmt.Printf("[%s] %s\n", aurora.Red("DIED"), address)
				}
			} else {
				fmt.Printf("[%s] [%s] [%s] %s\n", aurora.Green("LIVE"), aurora.Magenta(addr.Country), aurora.Cyan(addr.IP), address)

				if opt.Output != "" {
					fmt.Fprintf(opt.Result, "%s\n", address)
				}
			}
		})
	}

	p.Wait()
}

func isMatchCC(cc []string, code string) bool {
	if code == "" {
		return false
	}

	for _, c := range cc {
		if code == strings.ToUpper(strings.TrimSpace(c)) {
			return true
		}
	}

	return false
}

func check(address string, timeout time.Duration) (IPInfo, error) {
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return ipinfo, err
	}

	tr, err := mubeng.Transport(address)
	if err != nil {
		return ipinfo, err
	}

	proxy := &mubeng.Proxy{
		Address:   address,
		Transport: tr,
	}

	client, err := proxy.New(req)
	if err != nil {
		return ipinfo, err
	}

	client.Timeout = timeout
	req.Header.Add("Connection", "close")

	resp, err := client.Do(req)
	if err != nil {
		return ipinfo, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ipinfo, err
	}

	err = json.Unmarshal([]byte(body), &ipinfo)
	if err != nil {
		return ipinfo, err
	}

	defer resp.Body.Close()
	defer tr.CloseIdleConnections()

	return ipinfo, nil
}
