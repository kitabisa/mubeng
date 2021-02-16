package checker

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/logrusorgru/aurora"
	"ktbs.dev/mubeng/common"
	"ktbs.dev/mubeng/pkg/mubeng"
)

// Do checks proxy from list.
//
// Displays proxies that have died if verbose mode is enabled,
// or save live proxies into user defined files.
func Do(opt *common.Options) {
	for _, proxy := range opt.List {
		wg.Add(1)

		go func(address string) {
			cc, err := check(address, opt.Timeout)
			if err != nil {
				if opt.Verbose {
					fmt.Printf("[%s] %s\n", aurora.Red("DIED"), address)
				}
			} else {
				fmt.Printf("[%s] [%s] %s\n", aurora.Green("LIVE"), aurora.Magenta(cc), address)

				if opt.Output != "" {
					fmt.Fprintf(opt.Result, "%s\n", address)
				}
			}

			defer wg.Done()
		}(proxy)
	}

	wg.Wait()
}

func check(address string, timeout time.Duration) (string, error) {
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return "", err
	}

	tr, err := mubeng.Transport(address)
	if err != nil {
		return "", err
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
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal([]byte(body), &myip)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	return myip.CC, nil
}
