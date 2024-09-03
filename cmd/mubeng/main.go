package main

import (
	"github.com/projectdiscovery/gologger"
	"github.com/kitabisa/mubeng/internal/runner"
)

func main() {
	opt := runner.Options()

	if err := runner.New(opt); err != nil {
		gologger.Fatal().Msgf("Error! %s", err)
	}
}
