package server

import (
	"fmt"
	"net/http"

	"github.com/mbndr/logo"
)

// repo:hashicorp/go-retryablehttp /v\.(Error|Info|Debug|Warn)\("/
type ReleveledLogo struct {
	*logo.Logger
	*http.Request

	Verbose bool
}

func (l ReleveledLogo) toArgs(msg string, keysAndValues ...interface{}) []interface{} {
	args := []interface{}{msg}

	for i := 0; i < len(keysAndValues); i += 2 {
		if i+1 < len(keysAndValues) {
			key := keysAndValues[i]
			value := keysAndValues[i+1]

			switch key {
			case "request", "timeout":
				continue
			case "error":
				return []interface{}{value}
			case "remaining":
				value = fmt.Sprintf("%d", value)
			}

			args = append(args, fmt.Sprintf(" [%v=%q]", key, value))
		}
	}

	return args
}

func (l ReleveledLogo) Error(msg string, keysAndValues ...interface{}) {
	if l.Verbose {
		return
	}

	args := l.toArgs(msg, keysAndValues...)
	l.Logger.Error(args...)
}

func (l ReleveledLogo) Info(msg string, keysAndValues ...interface{}) {
	args := l.toArgs(msg, keysAndValues...)
	l.Logger.Info(args...)
}

func (l ReleveledLogo) Debug(msg string, keysAndValues ...interface{}) {
	switch msg {
	case "performing request":
		return
	case "retrying request":
		msg = fmt.Sprintf("%s Retrying %s %s", l.Request.RemoteAddr, l.Request.Method, l.Request.URL)
	}

	args := l.toArgs(msg, keysAndValues...)
	l.Logger.Debug(args...)
}

func (l ReleveledLogo) Warn(msg string, keysAndValues ...interface{}) {
	args := l.toArgs(msg, keysAndValues...)
	l.Logger.Warn(args...)
}
