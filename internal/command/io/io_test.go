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
}

func TestCustomIO(t *testing.T) {
	testIn := strings.NewReader("")
	testOut := &bytes.Buffer{}
	testErr := &bytes.Buffer{}
	zeroIO := New(testIn, testOut, testErr)

	if zeroIO.In() != testIn {
		t.Error("Expected default in to be as provided through constructor")
	}

	if zeroIO.Out() != testOut {
		t.Error("Expected default out to be as provided through constructor")
	}

	if zeroIO.Err() != testErr {
		t.Error("Expected default err to be as provided through constructor")
	}
}
