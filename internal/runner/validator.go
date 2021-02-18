package runner

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"ktbs.dev/mubeng/common"
	"ktbs.dev/mubeng/pkg/mubeng"
)

// validate user-supplied option values before Runner.
func validate(opt *common.Options) error {
	var err error

	if opt.File == "" {
		return errors.New("no proxy file provided")
	}

	opt.List, err = readFile(opt.File)
	if err != nil {
		return err
	}

	opt.List = uniq(opt.List)
	if len(opt.List) < 1 {
		return fmt.Errorf("open %s: has no valid proxy URLs", opt.File)
	}

	if opt.Output != "" {
		opt.Result, err = os.OpenFile(opt.Output, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			return err
		}
	}

	return nil
}

// readFile which is returned as a string array.
func readFile(path string) ([]string, error) {
	var lines []string

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func uniq(list []string) []string {
	keys := make(map[string]bool)
	uniq := []string{}

	for _, proxy := range list {
		if _, value := keys[proxy]; !value {
			_, err := mubeng.Transport(proxy)
			if err == nil {
				keys[proxy] = true
				uniq = append(uniq, proxy)
			}
		}
	}

	return uniq
}
