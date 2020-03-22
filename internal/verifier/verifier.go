package verifier

import (
	"github.com/LaurenceGA/go-crev/internal/command/io"
	"github.com/LaurenceGA/go-crev/internal/mod"

	"github.com/jedib0t/go-pretty/table"
)

type ModLister interface {
	List() ([]*mod.Module, error)
}

func New(modLister ModLister, commandIO *io.IO) *Verifier {
	return &Verifier{
		modLister: modLister,
		commandIO: commandIO,
	}
}

type Verifier struct {
	modLister ModLister
	commandIO *io.IO
}

type Verification struct {
	Package, Version string
}

func (v *Verifier) Verify() error {
	modules, err := v.modLister.List()
	if err != nil {
		return err
	}

	verifications := v.createVerifications(modules)

	v.writeVerifications(verifications)

	return nil
}

func (v *Verifier) createVerifications(allModules []*mod.Module) []*Verification {
	verifications := make([]*Verification, 0, len(allModules))

	for _, m := range allModules {
		if m.Main {
			continue
		}

		verifications = append(verifications, &Verification{
			Package: m.Path,
			Version: m.Version,
		})
	}

	return verifications
}

func (v *Verifier) writeVerifications(verifications []*Verification) {
	t := table.NewWriter()
	t.AppendHeader(table.Row{"Package", "Version"})
	t.SetOutputMirror(v.commandIO.Out())
	t.SetStyle(table.StyleLight)

	for _, ver := range verifications {
		t.AppendRow(table.Row{ver.Package, ver.Version})
	}

	t.Render()
}
