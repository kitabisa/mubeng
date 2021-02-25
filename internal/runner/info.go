package runner

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/logrusorgru/aurora"
	"github.com/tcnksm/go-latest"
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

// isLatest check if current version is latest
func isLatest() (bool, string) {
	if common.Version == "" {
		return false, ""
	}

	res, _ := latest.Check(&latest.GithubTag{
		Owner:      "kitabisa",
		Repository: common.App,
	}, common.Version)

	if strings.Contains(res.Current, "dev") {
		return false, res.Current
	}

	return res.Latest, res.Current
}
