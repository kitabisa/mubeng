package runner

import (
	"bufio"
	"errors"
	"os"

	"ktbs.dev/mubeng/common"
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
