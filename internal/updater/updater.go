package updater

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/briandowns/spinner"
	"github.com/inconshreveable/go-update"
	"github.com/logrusorgru/aurora"
	"github.com/projectdiscovery/gologger"
	"github.com/kitabisa/mubeng/common"
)

// New to initialize updater, such as:
// Get latest & changes, download & update binary version.
func New() error {
	s := spinner.New(spinner.CharSets[11], 90*time.Millisecond, spinner.WithWriter(os.Stderr))
	if err := s.Color("cyan"); err != nil {
		return err
	}

	s.Start()
	s.Suffix = " Getting latest version..."
	if lat, ver := isLatest(); !lat {
		s.Stop()
		gologger.Info().Msgf("New version v%s is available!", ver)
		s.Restart()

		s.Suffix = " Get changes..."
		chg, err := getChanges("v" + ver)
		if err != nil {
			s.Stop()
			return err
		}

		s.Stop()
		fmt.Fprintf(os.Stderr, "\n%s", aurora.Magenta(chg))

		if err := survey.AskOne(&survey.Confirm{Message: "Are you sure want to update?"}, &sure); err != nil {
			return err
		}

		if !sure {
			gologger.Info().Msgf("Cancelled.")
			os.Exit(1)
		}

		s.Restart()
		s.Suffix = " Downloading..."
		if err := doUpdate(ver); err != nil {
			s.Stop()
			return err
		}

		s.Stop()
		gologger.Info().Msgf("Successfully update!")
	} else {
		s.Stop()
		gologger.Info().Msgf("It's the latest stable version, no need to update.")
	}
	os.Exit(1)

	return nil
}

// doUpdate will replace mubeng binary with latest downlaoded version
func doUpdate(ver string) error {
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

	return update.Apply(resp.Body, opt)
}
