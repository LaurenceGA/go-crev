package mod

import (
	"bytes"
	"io"
	"os/exec"
)

type ModulesWrapper interface {
	List() (io.Reader, error)
}

func NewWrapper() *Wrapper {
	return &Wrapper{}
}

type Wrapper struct {
}

func (m *Wrapper) List() (io.Reader, error) {
	// Consider adding back in "-u". It will slow things down...
	cmd := exec.Command("go", "list", "-m", "-json", "all")

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	return &out, nil
}
