package trust

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/LaurenceGA/go-crev/internal/command/io"
	"github.com/LaurenceGA/go-crev/internal/config"
	"github.com/LaurenceGA/go-crev/internal/github"
	"github.com/LaurenceGA/go-crev/internal/store"
)

type ConfigReader interface {
	Load() (*config.Configuration, error)
}

type Github interface {
	GetUser(context.Context, string) (*github.User, error)
	GetRepository(context.Context, string, string) (*github.Repository, error)
}

func NewTrustCreator(commandIO *io.IO,
	configReader ConfigReader,
	githubClient Github) *Creator {
	return &Creator{
		commandIO:    commandIO,
		configReader: configReader,
		githubClient: githubClient,
	}
}

type Creator struct {
	commandIO    *io.IO
	configReader ConfigReader
	githubClient Github
}

type CreatorOptions struct {
	IdentityFile string
}

func (t *Creator) CreateTrust(ctx context.Context, usernameRaw string, options CreatorOptions) error {
	_, err := t.loadConfig()
	if err != nil {
		return err
	}

	// Load local SSH key (verify?)

	// Get user ID
	username := strings.TrimPrefix(usernameRaw, "@")

	usr, err := t.githubClient.GetUser(ctx, username)
	if err != nil {
		return err
	}

	// Test usr.ID == config.ID.ID

	idURL := t.getUserIDURL(ctx, usr.Login)

	fmt.Println(idURL)

	// Look for standard crev-proofs repo
	// Present UI for rating
	// Present UI for comment
	// Sign
	// Write file
	// Commit

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

func validateConfig(c *config.Configuration) error {
	// Should check if location exists in filesystem?
	if c.CurrentStore == "" {
		return errors.New("user current store is empty")
	}

	if c.CurrentID == nil {
		return errors.New("user current ID not set")
	}

	return nil
}
