package io

import (
	"io"
	"os"
)

func New(in io.Reader, out, err io.Writer) *IO {
	return &IO{
		in:  in,
		out: out,
		err: err,
	}
}

type IO struct {
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
