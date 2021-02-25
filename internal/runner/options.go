package runner

import (
	"flag"
	"os"
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

	flag.DurationVar(&opt.Timeout, "t", 30*time.Second, "")
	flag.DurationVar(&opt.Timeout, "timeout", 30*time.Second, "")

	flag.IntVar(&opt.Rotate, "r", 1, "")
	flag.IntVar(&opt.Rotate, "rotate", 1, "")

	flag.BoolVar(&opt.Verbose, "v", false, "")
	flag.BoolVar(&opt.Verbose, "verbose", false, "")

	flag.BoolVar(&opt.Daemon, "d", false, "")
	flag.BoolVar(&opt.Daemon, "daemon", false, "")

	flag.StringVar(&opt.Output, "o", "", "")
	flag.StringVar(&opt.Output, "output", "", "")

	flag.BoolVar(&doUpdate, "u", false, "")
	flag.BoolVar(&doUpdate, "update", false, "")

	flag.Usage = func() {
		showBanner()
		showUsage()
	}

	flag.Parse()
	showBanner()

	if isConnected() {
		lat, ver := isLatest()

		if !lat && ver != "" {
			gologger.Info().Msgf("New version v%s is available!", ver)

			if doUpdate {
				gologger.Info().Msgf("Updating...")
				if err := updateNow(ver); err != nil {
					gologger.Fatal().Msgf("Error while update! %s.", err)
				}

				gologger.Info().Msgf("Successfully update!")
				os.Exit(1)
			}
		}
	}

	if err := validate(opt); err != nil {
		gologger.Fatal().Msgf("Error! %s.", err)
	}

	return opt
}
