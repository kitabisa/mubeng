package daemon

import (
	"path/filepath"
	"strconv"

	"github.com/kardianos/service"
	"github.com/projectdiscovery/gologger"
	"ktbs.dev/mubeng/common"
)

// New to initialize mubeng in daemon
func New(opt *common.Options) error {
	file, err := filepath.Abs(opt.File)
	if err != nil {
		return err
	}

	// Copying user-supplied arguments
	args := []string{
		"-f", file,
		"-a", opt.Address,
		"-t", opt.Timeout.String(),
		"-r", strconv.Itoa(opt.Rotate),
		"-o", opt.Output,
	}

	o := make(service.KeyValue)
	o["Restart"] = "on-success"
	o["SuccessExitStatus"] = "1 2 8 SIGKILL"

	cfg := &service.Config{
		Name:        "mubeng",
		DisplayName: "mubeng",
		Description: "An incredibly fast proxy checker & IP rotator with ease.",
		Arguments:   args,
		Option:      o,
	}

	p := &program{opt: opt}
	s, err := service.New(p, cfg)
	if err != nil {
		return err
	}

	// Stop & uninstall current mubeng service, then re-installing & start
	_ = service.Control(s, "stop")
	_ = service.Control(s, "uninstall")
	err = service.Control(s, "install")
	if err != nil {
		return err
	}

	gologger.Infof("Running as daemon...")

	err = service.Control(s, "start")
	if err != nil {
		return err
	}

	return nil
}
