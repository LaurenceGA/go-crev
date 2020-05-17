package io

import (
	"fmt"
	"io"
	"os"
)

func New(in io.Reader, out, err io.Writer, verbose bool) *IO {
	return &IO{
		Verbose: verbose,
		in:      in,
		out:     out,
		err:     err,
	}
}

type IO struct {
	Verbose  bool
	in       io.Reader
	out, err io.Writer
}

func (i *IO) In() io.Reader {
	if i.in == nil {
		return os.Stdin
	}

	return i.in
}

func (i *IO) Out() io.Writer {
	if i.out == nil {
		return os.Stdout
	}

	return i.out
}

func (i *IO) Err() io.Writer {
	if i.err == nil {
		return os.Stderr
	}

	return i.err
}

// VerbosePrintf will print out the items with the given format if set to verbose.
func (i *IO) VerbosePrintf(format string, a ...interface{}) {
	if i.Verbose {
		fmt.Fprintf(i.Out(), format, a...)
	}
}

// VerbosePrintln will print out the given items and a newline if set to verbose.
func (i *IO) VerbosePrintln(a ...interface{}) {
	if i.Verbose {
		fmt.Fprintln(i.Out(), a...)
	}
}
