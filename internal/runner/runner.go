package runner

import (
	"errors"

	"github.com/kitabisa/mubeng/common"
	"github.com/kitabisa/mubeng/internal/checker"
	"github.com/kitabisa/mubeng/internal/daemon"
	"github.com/kitabisa/mubeng/internal/server"
)

// New to switch an action, whether to check or run a proxy server.
func New(opt *common.Options) error {
	if opt.Address != "" {
		if opt.Daemon {
			return daemon.New(opt)
		}

		server.Run(opt)
	} else if opt.Check {
		checker.Do(opt)

		if opt.Output != "" {
			defer opt.Result.Close()
		}
	} else {
		return errors.New("no action to run")
	}

	return nil
}
