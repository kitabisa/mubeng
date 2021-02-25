package runner

import (
	"errors"
	"net/http"
	"runtime"
	"strings"

	"github.com/inconshreveable/go-update"
	"ktbs.dev/mubeng/common"
)

var doUpdate bool

func updateNow(ver string) error {
	opt := update.Options{}

	if err := opt.CheckPermissions(); err != nil {
		return err
	}

	rep := strings.NewReplacer(
		"APP", common.App,
		"VVER", "v"+ver,
		"VER", ver,
		"OS", runtime.GOOS,
		"ARCH", runtime.GOARCH,
	)

	binaryURL = rep.Replace(binaryURL)
	if runtime.GOOS == "windows" {
		binaryURL += ".exe"
	}

	resp, err := http.Get(binaryURL)
	if err != nil {
		return errors.New("check your internet connection")
	}
	defer resp.Body.Close()

	err = update.Apply(resp.Body, opt)
	if err != nil {
		return err
	}

	return nil
}
