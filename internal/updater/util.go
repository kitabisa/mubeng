package updater

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/projectdiscovery/gologger"
	"github.com/tcnksm/go-latest"
	"github.com/kitabisa/mubeng/common"
)

// isLatest check if current version is latest
func isLatest() (bool, string) {
	current := common.Version
	if current == "" {
		current = "0"
	}

	res, err := latest.Check(&latest.GithubTag{
		Owner:      "kitabisa",
		Repository: common.App,
	}, current)

	if err != nil {
		println()
		gologger.Fatal().Msgf("Error! %s", err)
	}

	if strings.Contains(res.Current, "dev") {
		return false, res.Current
	}

	return res.Latest, res.Current
}

// getChanges will returning changelog of given tag
func getChanges(tag string) (string, error) {
	chg := &changes{}

	resp, err := http.Get("https://api.github.com/repos/kitabisa/mubeng/releases/tags/" + tag)
	if err != nil {
		return "", errors.New("check your internet connection")
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(chg); err != nil {
		return "", err
	}

	if chg.Name != tag {
		return chg.Body, errors.New("latest tag doesn't same")
	}

	return chg.Body, nil
}
