package mod

import (
	"bytes"
	"encoding/json"
	"io"
	"os/exec"
	"time"
)

func NewLister() *Lister {
	return &Lister{}
}

type Lister struct {
}

func (l *Lister) List() ([]*Module, error) {
	cmdOutput, err := l.modList()
	if err != nil {
		return nil, err
	}

	var modules []*Module

	dec := json.NewDecoder(cmdOutput)

	for {
		var v Module

		err := dec.Decode(&v)
		if err != nil {
			return modules, nil
		}

		modules = append(modules, &v)
	}
}

func (l *Lister) modList() (io.Reader, error) {
	// Consider adding back in "-u". It will slow things down.
	cmd := exec.Command("go", "list", "-m", "-json", "all")

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	return &out, nil
}

// Module holds information about a specific module listed by go list
// Copied from documentation from "go help list"
type Module struct {
	Path      string       `json:",omitempty"` // module path
	Version   string       `json:",omitempty"` // module version
	Versions  []string     `json:",omitempty"` // available module versions (with -versions)
	Replace   *Module      `json:",omitempty"` // replaced by this module
	Time      *time.Time   `json:",omitempty"` // time version was created
	Update    *Module      `json:",omitempty"` // available update, if any (with -u)
	Main      bool         `json:",omitempty"` // is this the main module?
	Indirect  bool         `json:",omitempty"` // is this module only an indirect dependency of main module?
	Dir       string       `json:",omitempty"` // directory holding files for this module, if any
	GoMod     string       `json:",omitempty"` // path to go.mod file describing module, if any
	GoVersion string       `json:",omitempty"` // go version used in module
	Error     *ModuleError `json:",omitempty"` // error loading module
}

// ModuleError represents the error when a module cannot be loaded
type ModuleError struct {
	Err string // the error itself
}
