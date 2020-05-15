package trust

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/LaurenceGA/go-crev/internal/command/io"
	"github.com/LaurenceGA/go-crev/internal/config"
	"github.com/LaurenceGA/go-crev/internal/github"
	"github.com/LaurenceGA/go-crev/internal/id"
	"github.com/LaurenceGA/go-crev/internal/store"
	"github.com/LaurenceGA/go-crev/proof/trust"
	"github.com/google/uuid"
	"golang.org/x/crypto/ssh"
)

type ConfigReader interface {
	Load() (*config.Configuration, error)
}

type Github interface {
	GetUser(context.Context, string) (*github.User, error)
	GetRepository(context.Context, string, string) (*github.Repository, error)
}

type Prompter interface {
	Select(string, []string) (string, error)
	Prompt(string) (string, error)
}

type KeyLoader interface {
	LoadKey(string) (ssh.Signer, error)
}

type StoreWriter interface {
	SaveTrust(*store.ProofStore, *trust.Trust) error
}

func NewTrustCreator(commandIO *io.IO,
	configReader ConfigReader,
	githubClient Github,
	prompter Prompter,
	keyLoader KeyLoader,
	storeWriter StoreWriter) *Creator {
	return &Creator{
		commandIO:    commandIO,
		configReader: configReader,
		githubClient: githubClient,
		prompter:     prompter,
		keyLoader:    keyLoader,
		storeWriter:  storeWriter,
	}
}

type Creator struct {
	commandIO    *io.IO
	configReader ConfigReader
	githubClient Github
	prompter     Prompter
	keyLoader    KeyLoader
	storeWriter  StoreWriter
}

type CreatorOptions struct {
	IdentityFile string
}

func (t *Creator) CreateTrust(ctx context.Context, usernameRaw string, options CreatorOptions) error {
	config, err := t.loadConfig()
	if err != nil {
		return err
	}

	username := strings.TrimPrefix(usernameRaw, "@")

	usr, err := t.githubClient.GetUser(ctx, username)
	if err != nil {
		return err
	}

	idURL := t.getUserIDURL(ctx, usr.Login)

	sshKeySigner, err := t.keyLoader.LoadKey(options.IdentityFile)
	if err != nil {
		return fmt.Errorf("loading SSH key: %w", err)
	}

	trusteeID := &id.ID{
		ID:    strconv.Itoa(int(usr.ID)),
		Type:  id.Github,
		URL:   idURL,
		Alias: usr.Login,
	}

	trustLevel, err := t.getTrustLevel()
	if err != nil {
		return err
	}

	trustComment, err := t.prompter.Prompt("Comment")
	if err != nil {
		return err
	}

	trustObj := trust.New(uuid.New().String(), *config.CurrentID, trustLevel, trustComment, []*id.ID{trusteeID})

	if err := trustObj.Sign(sshKeySigner); err != nil {
		return err
	}

	userStore := &store.ProofStore{Dir: config.CurrentStore}

	if err := t.storeWriter.SaveTrust(userStore, trustObj); err != nil {
		return fmt.Errorf("saving trust: %w", err)
	}

	return nil
}

func (t *Creator) loadConfig() (*config.Configuration, error) {
	conf, err := t.configReader.Load()
	if err != nil {
		return nil, err
	}

	if err := validateConfig(conf); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return conf, nil
}

func (t *Creator) getUserIDURL(ctx context.Context, username string) string {
	repo, err := t.githubClient.GetRepository(ctx, username, store.StandardCrevProofRepoName)
	if err != nil {
		if errors.Is(err, github.NotFoundError) {
			fmt.Fprintf(t.commandIO.Out(),
				"Couldn't find proof repo in Github for %s/%s\n",
				username,
				store.StandardCrevProofRepoName)
		} else {
			// Non-fatal. Just print and move on...
			fmt.Fprintf(t.commandIO.Err(), "Failed trying to find repository with error: %v\n", err)
		}

		return "" // No known crev proof URL for ID
	}

	return repo.HTMLurl
}

type constError string

func (e constError) Error() string {
	return string(e)
}

const (
	ErrCurrentStoreEmpty constError = "user current store is empty"
	ErrCurrentIDNotSet   constError = "user current ID not set"
)

func validateConfig(c *config.Configuration) error {
	// Should check if location exists in filesystem?
	if c.CurrentStore == "" {
		return ErrCurrentStoreEmpty
	}

	if c.CurrentID == nil {
		return ErrCurrentIDNotSet
	}

	return nil
}

type InvalidLevelError string

func (e InvalidLevelError) Error() string {
	return fmt.Sprintf("invalid level: %s", string(e))
}

func (t *Creator) getTrustLevel() (trust.Level, error) {
	levelResponse, err := t.prompter.Select("Trust level", trustPrompts())
	if err != nil {
		return "", err
	}

	level, ok := trust.ToLevel(levelResponse)
	if !ok {
		return "", InvalidLevelError(levelResponse)
	}

	return level, nil
}

func trustPrompts() []string {
	levels := trust.Levels()
	promptLevels := make([]string, 0, len(levels))

	for _, l := range levels {
		promptLevels = append(promptLevels, strings.Title(string(l)))
	}

	return promptLevels
}
