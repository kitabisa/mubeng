package helper

import (
	"os"
	"strings"
)

func getEnviron() map[string]interface{} {
	for _, env := range os.Environ() {
		pair := strings.SplitN(env, "=", 2)
		key := pair[0]
		val := pair[1]

		m[key] = val
	}

	return m
}
