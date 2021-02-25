package main

import (
	"github.com/projectdiscovery/gologger"
	"ktbs.dev/mubeng/internal/runner"
)

func main() {
	opt := runner.Options()

	if err := runner.New(opt); err != nil {
		gologger.Fatal().Msgf("Error! %s", err)
	}
}
