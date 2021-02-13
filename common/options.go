package common

import (
	"os"
	"time"
)

// Options consists of the configuration required.
type Options struct {
	File    string
	Address string
	Check   bool
	Timeout time.Duration
	Rotate  int
	Verbose bool
	Daemon  bool
	Output  string
	Result  *os.File
	List    []string
	// After int
}
