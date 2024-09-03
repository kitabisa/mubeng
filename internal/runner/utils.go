package runner

import "os"

func hasStdin() bool {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return false
	}

	// Check if the input is from a pipe or redirected input
	return (stat.Mode() & os.ModeCharDevice) == 0
}
