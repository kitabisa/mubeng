package runner

import (
	"fmt"
	"os"

	"github.com/logrusorgru/aurora"
	"github.com/kitabisa/mubeng/common"
)

// showBanner to stderr
func showBanner() {
	fmt.Fprintf(os.Stderr, "%s\n\n", aurora.Cyan(common.Banner))
}

// showUsage to stderr
func showUsage() {
	fmt.Fprint(os.Stderr, "Usage:", common.Usage)
}

// showVersion and exit
func showVersion() {
	version := common.Version
	if version == "" {
		version = "unknown (go-get)"
	}

	fmt.Println(common.App, "version", version)
	os.Exit(1)
}
