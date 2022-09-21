package helper

import (
	"bytes"
	"text/template"

	"github.com/valyala/fastrand"
	"github.com/valyala/fasttemplate"
)

// Eval subtitute `i` strings
func Eval(i string) string {
	m = getEnviron()

	t := fasttemplate.New(i, "{{", "}}")
	s := t.ExecuteStringStd(m)

	return s
}

// EvalFunc subtitute `i` strings with functions
func EvalFunc(i string) string {
	var b bytes.Buffer

	f := template.FuncMap{
		"uint32":  fastrand.Uint32,
		"uint32n": fastrand.Uint32n,
	}

	t, err := template.New("eval").Funcs(f).Parse(i)
	if err != nil {
		return i
	}

	if err := t.Execute(&b, nil); err != nil {
		return i
	}

	return b.String()
}
