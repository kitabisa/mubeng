package runner

import (
	"flag"
	"time"

	"github.com/projectdiscovery/gologger"
	"ktbs.dev/mubeng/common"
)

// Options defines the values needed to execute the Runner.
func Options() *common.Options {
	opt := &common.Options{}

	flag.StringVar(&opt.File, "f", "", "")
	flag.StringVar(&opt.File, "file", "", "")

	flag.StringVar(&opt.Address, "a", "", "")
	flag.StringVar(&opt.Address, "address", "", "")

	flag.BoolVar(&opt.Check, "c", false, "")
	flag.BoolVar(&opt.Check, "check", false, "")

	flag.DurationVar(&opt.Timeout, "t", 5*time.Second, "")
	flag.DurationVar(&opt.Timeout, "timeout", 5*time.Second, "")

	flag.IntVar(&opt.Rotate, "r", 10, "")
	flag.IntVar(&opt.Rotate, "rotate", 10, "")

	flag.BoolVar(&opt.Verbose, "v", false, "")
	flag.BoolVar(&opt.Verbose, "verbose", false, "")

	flag.BoolVar(&opt.Daemon, "d", false, "")
	flag.BoolVar(&opt.Daemon, "daemon", false, "")

	flag.StringVar(&opt.Output, "o", "", "")
	flag.StringVar(&opt.Output, "output", "", "")

	flag.Usage = func() {
		showBanner()
		showUsage()
	}

	flag.Parse()
	showBanner()

	if err := validate(opt); err != nil {
		gologger.Fatalf("Error! %s", err)
	}

	return opt
}
