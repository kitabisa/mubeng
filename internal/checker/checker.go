package checker

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/logrusorgru/aurora"
	"ktbs.dev/mubeng/common"
	"ktbs.dev/mubeng/pkg/helper"
	"ktbs.dev/mubeng/pkg/mubeng"
)

// Do checks proxy from list.
//
// Displays proxies that have died if verbose mode is enabled,
// or save live proxies into user defined files.
func Do(opt *common.Options) {
	for _, proxy := range opt.ProxyManager.Proxies {
		wg.Add(1)

		go func(address string) {
			defer wg.Done()

			addr, err := check(address, opt.Timeout)
			if len(opt.Countries) > 0 && !isMatchCC(opt.Countries, addr.CC) {
				return
			}

			if err != nil {
				if opt.Verbose {
					fmt.Printf("[%s] %s\n", aurora.Red("DIED"), address)
				}
			} else {
				fmt.Printf("[%s] [%s] [%s] %s\n", aurora.Green("LIVE"), aurora.Magenta(addr.CC), aurora.Cyan(addr.IP), address)

				if opt.Output != "" {
					fmt.Fprintf(opt.Result, "%s\n", address)
				}
			}
		}(helper.EvalFunc(proxy))
	}

	wg.Wait()
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

func check(address string, timeout time.Duration) (myIP, error) {
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return myip, err
	}

	tr, err := mubeng.Transport(address)
	if err != nil {
		return myip, err
	}

	proxy := &mubeng.Proxy{
		Address:   address,
		Transport: tr,
	}

	client, req = proxy.New(req)
	client.Timeout = timeout
	req.Header.Add("Connection", "close")

	resp, err := client.Do(req)
	if err != nil {
		return myip, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return myip, err
	}

	err = json.Unmarshal([]byte(body), &myip)
	if err != nil {
		return myip, err
	}

	defer resp.Body.Close()
	defer tr.CloseIdleConnections()

	return myip, nil
}
