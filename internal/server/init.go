package server

import (
	"math/rand"
	"time"
)

func init() {
	// TODO(dwisiswant0): deprecated, update this later.
	// nolint: staticcheck
	rand.Seed(time.Now().UnixNano())
}
