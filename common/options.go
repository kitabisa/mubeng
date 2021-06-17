package common

import (
	"ktbs.dev/mubeng/internal/proxymanager"
	"os"
	"time"
)

// Options consists of the configuration required.
type Options struct {
	File         string
	Address      string
	Check        bool
	Timeout      time.Duration
	Rotate       int
	Verbose      bool
	Output       string
	Result       *os.File
	ProxyManager *proxymanager.ProxyManager
	Daemon       bool
}
