package verifier

import (
	"github.com/LaurenceGA/go-crev/internal/command/io"
	"github.com/LaurenceGA/go-crev/internal/mod"

	"github.com/jedib0t/go-pretty/table"
)

type ModLister interface {
	List() ([]*mod.Module, error)
}

type LineCounter interface {
	CountLines(string) (int, error)
}

func New(commandIO *io.IO, modLister ModLister, lineCounter LineCounter) *Verifier {
	return &Verifier{
		modLister:   modLister,
		commandIO:   commandIO,
		lineCounter: lineCounter,
	}
}

type Verifier struct {
	commandIO   *io.IO
	modLister   ModLister
	lineCounter LineCounter
}

type Verification struct {
	Package, Version string
	LineOfCode       int
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

		linesOfCode, err := v.lineCounter.CountLines(m.Dir)
		if err != nil {
			// Swallow error and move on to allow partial errors
			// Should log warning for error here...
			linesOfCode = -1
		}

		verifications = append(verifications, &Verification{
			Package:    m.Path,
			Version:    m.Version,
			LineOfCode: linesOfCode,
		})
	}

	return verifications
}

func (v *Verifier) writeVerifications(verifications []*Verification) {
	t := table.NewWriter()
	t.AppendHeader(table.Row{"Package", "Version", "Lines of code"})
	t.SetOutputMirror(v.commandIO.Out())
	t.SetStyle(table.StyleLight)

	for _, ver := range verifications {
		t.AppendRow(table.Row{ver.Package, ver.Version, ver.LineOfCode})
	}

	t.Render()
}
