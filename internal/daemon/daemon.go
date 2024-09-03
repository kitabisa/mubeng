package daemon

import (
	"runtime"
	"strconv"

	"github.com/kardianos/service"
	"github.com/projectdiscovery/gologger"
	"github.com/kitabisa/mubeng/common"
)

// New to initialize mubeng in daemon
func New(opt *common.Options) error {
	// Copying user-supplied arguments
	args := []string{
		"-f", opt.File,
		"-a", opt.Address,
		"-A", opt.Auth,
		"-t", opt.Timeout.String(),
		"-r", strconv.Itoa(opt.Rotate),
		"-m", opt.Method,
		"-o", opt.Output,
	}

	if opt.Sync {
		args = append(args, "-s")
	}

	if opt.Verbose {
		args = append(args, "-v")
	}

	if opt.Watch {
		args = append(args, "-w")
	}

	o := make(service.KeyValue)
	o["Restart"] = "on-success"
	o["SuccessExitStatus"] = "1 2 8 SIGKILL"

	cfg.Arguments = args
	cfg.Option = o

	p := &program{opt: opt}
	s, err := service.New(p, &cfg)
	if err != nil {
		return err
	}

	// Stop & uninstall current mubeng service, then re-installing & start
	_ = service.Control(s, "stop")
	_ = service.Control(s, "uninstall")

	if runtime.GOOS == "windows" {
		err = service.Control(s, "install")
		if err != nil {
			return err
		}

		gologger.Info().Msg("Service installed!")
		gologger.Info().Msg("Type 'net start mubeng' to start it up in daemon.")
	} else {
		_ = service.Control(s, "install")

		gologger.Info().Msg("Running as daemon...")
		return service.Control(s, "start")
	}

	return nil
}
