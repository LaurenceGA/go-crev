package ssh

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"golang.org/x/crypto/ssh"
)

type Prompter interface {
	Password(string) (string, error)
}

func NewLoader(prompter Prompter) *Loader {
	return &Loader{
		prompter: prompter,
	}
}

type Loader struct {
	prompter Prompter
}

type constError string

func (e constError) Error() string {
	return string(e)
}

const SSHKeyNotFound constError = "ssh key not found"

func (l *Loader) LoadKey(path string) (ssh.Signer, error) {
	keyPath, err := findValidKeyPath(path)
	if err != nil {
		return nil, fmt.Errorf("finding key file: %w", err)
	}

	key, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("reading key file: %w", err)
	}

	signer, err := l.parsePrivateKey(key, keyPath)
	if err != nil {
		return nil, fmt.Errorf("parsing key: %w", err)
	}

	return signer, nil
}

func (l *Loader) parsePrivateKey(key []byte, keyPath string) (ssh.Signer, error) {
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		var perr *ssh.PassphraseMissingError
		if errors.As(err, &perr) {
			return l.parsePrivateKeyWithPassphrase(key, keyPath)
		}

		return nil, err
	}

	return signer, nil
}

func (l *Loader) parsePrivateKeyWithPassphrase(key []byte, keyPath string) (ssh.Signer, error) {
	password, err := l.prompter.Password(fmt.Sprintf("Passphrase for %s", keyPath))
	if err != nil {
		return nil, err
	}

	return ssh.ParsePrivateKeyWithPassphrase(key, []byte(password))
}

func findValidKeyPath(path string) (string, error) {
	if path != "" {
		cleanPath := filepath.Clean(path)
		_, err := os.Stat(cleanPath)

		return cleanPath, err
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("finding home directory: %w", err)
	}

	for _, keyFileName := range keyFileNames() {
		keyPath := keySSHPath(home, keyFileName)
		if _, err := os.Stat(keyPath); err == nil {
			return keyPath, nil
		}
	}

	return "", SSHKeyNotFound
}

// Well known .ssh directory for all things ssh.
const sshDir = ".ssh"

func keySSHPath(home, keyFileName string) string {
	return filepath.Join(home, sshDir, keyFileName)
}

func keyFileNames() []string {
	return []string{"id_dsa", "id_ecdsa", "id_ecdsa_sk", "id_ed25519", "id_ed25519_sk", "id_rsa"}
}
