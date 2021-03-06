package mod

import (
	"encoding/json"
	"errors"
	"io"
	"time"
)

func NewLister(modWrapper ModulesWrapper) *Lister {
	return &Lister{
		modWrapper: modWrapper,
	}
}

type Lister struct {
	modWrapper ModulesWrapper
}

func (l *Lister) List() ([]*Module, error) {
	cmdOutput, err := l.modWrapper.List()
	if err != nil {
		return nil, err
	}

	var modules []*Module

	dec := json.NewDecoder(cmdOutput)

	for {
		var v Module

		err := dec.Decode(&v)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return modules, nil
			}

			return nil, err
		}

		modules = append(modules, &v)
	}
}

// Module holds information about a specific module listed by go list.
// Copied from documentation from "go help list".
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

// ModuleError represents the error when a module cannot be loaded.
type ModuleError struct {
	Err string // the error itself
}
