package runner

import (
	"fmt"
	"net/http"
	"os"

	"github.com/logrusorgru/aurora"
	"ktbs.dev/mubeng/common"
)

// showBanner to stderr
func showBanner() {
	fmt.Fprintf(os.Stderr, "%s\n\n", aurora.Cyan(common.Banner))
}

// showUsage to stderr
func showUsage() {
	fmt.Fprint(os.Stderr, "Usage:", common.Usage)
}

// isConnected to the internet?
func isConnected() bool {
	if _, err := http.Get(google204); err != nil {
		return false
	}

	return true
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
