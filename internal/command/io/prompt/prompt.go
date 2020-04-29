package prompt

import (
	stdio "io"
	"io/ioutil"

	"github.com/LaurenceGA/go-crev/internal/command/io"
	"github.com/manifoldco/promptui"
)

func NewPrompter(commandIO *io.IO) *Prompter {
	return &Prompter{
		commandIO: commandIO,
	}
}

type Prompter struct {
	commandIO *io.IO
}

func (p *Prompter) Select(label string, options []string) (string, error) {
	prompt := promptui.Select{
		Label:  label,
		Items:  options,
		Stdin:  ioutil.NopCloser(p.commandIO.In()),
		Stdout: noopCloser(p.commandIO.Out()),
	}

	_, result, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return result, nil
}

func (p *Prompter) Prompt(label string) (string, error) {
	prompt := promptui.Prompt{
		Label:  label,
		Stdin:  ioutil.NopCloser(p.commandIO.In()),
		Stdout: noopCloser(p.commandIO.Out()),
	}

	return prompt.Run()
}

func noopCloser(w stdio.Writer) stdio.WriteCloser {
	return &noopWriteCloser{w}
}

type noopWriteCloser struct {
	stdio.Writer
}

func (n *noopWriteCloser) Close() error {
	// Noop
	return nil
}
