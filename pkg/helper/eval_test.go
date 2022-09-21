package helper

import (
	"os"
	"testing"
)

func setupTestEval() {
	_ = os.Setenv("USERNAME", "FOO")
	_ = os.Setenv("PASSWORD", "BAR")
}

func TestEval(t *testing.T) {
	setupTestEval()

	got := Eval("http://{{USERNAME}}:{{PASSWORD}}@127.0.0.1:1337")
	want := "http://FOO:BAR@127.0.0.1:1337"

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestEvalFunc(t *testing.T) {
	input := "http://{{uint32}}:{{uint32n 10000}}@127.0.0.1:1337"

	got := EvalFunc(input)

	if input == got {
		t.Errorf("got %q, wanted pseudo-random", got)
	}
}
