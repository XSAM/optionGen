package main

import (
	"bytes"
	"fmt"
)

const (
	optionDeclarationSuffix = "OptionDeclaration"
	optionGen               = "optionGen"
)

type BufWrite struct {
	buf *bytes.Buffer
}

func (b BufWrite) wf(format string, vals ...interface{}) {
	fmt.Fprintf(b.buf, format, vals...)
}

func (b BufWrite) wln(vals ...interface{}) {
	fmt.Fprintln(b.buf, vals...)
}
