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

func NewCreator(commandIO *io.IO,
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

func (c *Creator) CreateTrust(ctx context.Context, usernameRaw string, options CreatorOptions) error {
	conf, err := c.loadConfig()
	if err != nil {
		return err
	}

	username := strings.TrimPrefix(usernameRaw, "@")

	fmt.Fprintln(c.commandIO.Out(), "Looking for Github user", username)

	usr, err := c.githubClient.GetUser(ctx, username)
	if err != nil {
		return err
	}

	c.commandIO.VerbosePrintf("Found user %d (%s)\n", usr.ID, username)

	idURL := c.getUserIDURL(ctx, usr.Login)

	sshKeySigner, err := c.keyLoader.LoadKey(options.IdentityFile)
	if err != nil {
		return fmt.Errorf("loading SSH key: %w", err)
	}

	trusteeID := &id.ID{
		ID:    strconv.Itoa(int(usr.ID)),
		Type:  id.Github,
		URL:   idURL,
		Alias: usr.Login,
	}

	trustLevel, err := c.getTrustLevel()
	if err != nil {
		return err
	}

	trustComment, err := c.prompter.Prompt("Comment")
	if err != nil {
		return err
	}

	trustObj := trust.New(uuid.NewString(), *conf.CurrentID, trustLevel, trustComment, []*id.ID{trusteeID})

	c.commandIO.VerbosePrintln("Signing trust")

	if err := trustObj.Sign(sshKeySigner); err != nil {
		return err
	}

	userStore := &store.ProofStore{Dir: conf.CurrentStore}

	c.commandIO.VerbosePrintln("Saving trust to current store")

	if err := c.storeWriter.SaveTrust(userStore, trustObj); err != nil {
		return fmt.Errorf("saving trust: %w", err)
	}

	fmt.Fprintf(c.commandIO.Out(), "Created trust proof %s in %s\n", trustObj.Data.ID, userStore.Dir)

	return nil
}

func (c *Creator) loadConfig() (*config.Configuration, error) {
	c.commandIO.VerbosePrintln("Loading user config")

	conf, err := c.configReader.Load()
	if err != nil {
		return nil, err
	}

	if err := validateConfig(conf); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	c.commandIO.VerbosePrintf(
		"CurrentID: %s ID - %s (%s) - %s\n",
		conf.CurrentID.Type,
		conf.CurrentID.ID,
		conf.CurrentID.Alias,
		conf.CurrentID.URL,
	)

	c.commandIO.VerbosePrintf(
		"CurrentStore:: %s\n",
		conf.CurrentStore,
	)

	return conf, nil
}

func (c *Creator) getUserIDURL(ctx context.Context, username string) string {
	fmt.Fprintf(c.commandIO.Out(), "Checking for proof store %s/%s\n", username, store.StandardCrevProofRepoName)

	repo, err := c.githubClient.GetRepository(ctx, username, store.StandardCrevProofRepoName)
	if err != nil {
		if errors.Is(err, github.NotFoundError) {
			fmt.Fprintln(c.commandIO.Out(), "Not found")
		} else {
			// Non-fatal. Just print and move on...
			fmt.Fprintf(c.commandIO.Err(), "Failed trying to find repository with error: %v\n", err)
		}

		return "" // No known crev proof URL for ID
	}

	fmt.Fprintln(c.commandIO.Out(), "Found proof store", repo.HTMLurl)

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

func (c *Creator) getTrustLevel() (trust.Level, error) {
	levelResponse, err := c.prompter.Select("Trust level", trustPrompts())
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
