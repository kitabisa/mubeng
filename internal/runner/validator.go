package runner

import (
	"errors"
	"ktbs.dev/mubeng/common"
	"ktbs.dev/mubeng/internal/proxymanager"
	"os"
	"path/filepath"
)

// validate user-supplied option values before Runner.
func validate(opt *common.Options) error {
	var err error

	if opt.File == "" {
		return errors.New("no proxy file provided")
	}

	opt.File, err = filepath.Abs(opt.File)
	if err != nil {
		return err
	}

	opt.ProxyManager, err = proxymanager.NewManager(opt.File)
	if err != nil {
		return err
	}

	if opt.Output != "" {
		opt.Output, err = filepath.Abs(opt.Output)
		if err != nil {
			return err
		}

		opt.Result, err = os.OpenFile(opt.Output, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			return err
		}
	}

	return nil
}
