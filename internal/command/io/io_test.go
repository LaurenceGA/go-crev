package io

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestDefaulIO(t *testing.T) {
	zeroIO := &IO{}

	if zeroIO.In() != os.Stdin {
		t.Error("Expected default in to be os.Stdin")
	}

	if zeroIO.Out() != os.Stdout {
		t.Error("Expected default out to be os.Stdout")
	}

	if zeroIO.Err() != os.Stderr {
		t.Error("Expected default err to be os.Stderr")
	}

	if zeroIO.Verbose != false {
		t.Error("Expected not to be verbose by default")
	}
}

func TestCustomIO(t *testing.T) {
	testIn := strings.NewReader("")
	testOut := &bytes.Buffer{}
	testErr := &bytes.Buffer{}
	testVerbose := true
	customIO := New(testIn, testOut, testErr, testVerbose)

	if customIO.In() != testIn {
		t.Error("Expected in to be as provided through constructor")
	}

	if customIO.Out() != testOut {
		t.Error("Expected out to be as provided through constructor")
	}

	if customIO.Err() != testErr {
		t.Error("Expected err to be as provided through constructor")
	}

	if customIO.Verbose != testVerbose {
		t.Error("Expected verbosity to be as provided through constructor")
	}
}
