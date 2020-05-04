package idset

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/LaurenceGA/go-crev/internal/command/io"
	"github.com/LaurenceGA/go-crev/internal/git"
	"github.com/LaurenceGA/go-crev/internal/github"
	"github.com/LaurenceGA/go-crev/internal/id"
	"github.com/LaurenceGA/go-crev/internal/store"
)

type ConfigManipulator interface {
	SetCurrentID(*id.ID) error
	SetCurrentStore(string) error
}

type Github interface {
	GetUser(context.Context, string) (*github.User, error)
	GetRepository(context.Context, string, string) (*github.Repository, error)
}

type RepoFetcher interface {
	Fetch(context.Context, string) (*store.ProofStore, error)
}

func NewIDSetter(commandIO *io.IO,
	configManipulator ConfigManipulator,
	githubUser Github,
	repoFetcher RepoFetcher) *IDSetter {
	return &IDSetter{
		commandIO:         commandIO,
		configManipulator: configManipulator,
		githubUser:        githubUser,
		repoFetcher:       repoFetcher,
	}
}

// IDSetter is responsible for high level flow of setting of a user's ID
type IDSetter struct {
	commandIO         *io.IO
	configManipulator ConfigManipulator
	githubUser        Github
	repoFetcher       RepoFetcher
}

// SetFromUsername takes a username input and sets our local current ID to that by finding
// resolving the ID from the username.
// By default we assume it's a Github ID
func (i *IDSetter) SetFromUsername(ctx context.Context, usernameRaw string) error {
	username := strings.TrimPrefix(usernameRaw, "@")

	usr, err := i.githubUser.GetUser(ctx, username)
	if err != nil {
		return fmt.Errorf("getting user: %w", err)
	}

	idStoreURL := i.loadExistingStandardRepo(ctx, usr.Login)

	return i.configManipulator.SetCurrentID(&id.ID{
		ID:    strconv.Itoa(int(usr.ID)),
		Type:  id.Github,
		URL:   idStoreURL,
		Alias: usr.Login,
	})
}

func (i *IDSetter) loadExistingStandardRepo(ctx context.Context, owner string) string {
	repo, err := i.githubUser.GetRepository(ctx, owner, store.StandardCrevProofRepoName)
	if err != nil {
		if errors.Is(err, github.NotFoundError) {
			fmt.Fprintf(i.commandIO.Out(),
				"Couldn't find proof repo in Github for %s/%s. You should make one here...\n",
				owner,
				store.StandardCrevProofRepoName)
		} else {
			// Non-fatal. Just print and move on...
			fmt.Fprintf(i.commandIO.Err(), "Failed trying to find repository with error: %v\n", err)
		}

		return "" // No known crev proof URL for ID
	}

	fmt.Fprintln(i.commandIO.Out(), "Found existing proof repo!")

	i.loadRepoAsCurrentStore(ctx, repo.CloneURL)

	return repo.HTMLurl
}

func (i *IDSetter) loadRepoAsCurrentStore(ctx context.Context, cloneURL string) {
	store, err := i.repoFetcher.Fetch(ctx, cloneURL)
	if err != nil {
		if !errors.Is(err, git.ErrRepositoryAlreadyExists) {
			fmt.Fprintf(i.commandIO.Err(), "Failed trying to clone proof repo: %v\n", err)

			return
		}

		fmt.Fprintln(i.commandIO.Out(), "It's already there!")
	}

	if err := i.configManipulator.SetCurrentStore(store.Dir); err != nil {
		fmt.Fprintf(i.commandIO.Err(), "Failed to set current store to %s: %v\n", store.Dir, err)
	}
}
