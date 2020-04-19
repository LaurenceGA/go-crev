// package flow detremines more complex user interaction flows
package flow

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
	Fetch(context.Context, string) error
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

// This is expected to be a well known name
// Repo doesn't have to be this name, but if it is we can automatically find it
const standardCrevProofRepoName = "crev-proofs"

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
		ID:   strconv.Itoa(int(usr.ID)),
		Type: id.Github,
		URL:  idStoreURL,
	})
}

func (i *IDSetter) loadExistingStandardRepo(ctx context.Context, owner string) string {
	repo, err := i.githubUser.GetRepository(ctx, owner, standardCrevProofRepoName)
	if err != nil {
		if errors.Is(err, github.NotFoundError) {
			fmt.Fprintf(i.commandIO.Out(),
				"Couldn't find proof repo in Github for %s/%s. You should make one here...\n",
				owner,
				standardCrevProofRepoName)
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
	if err := i.repoFetcher.Fetch(ctx, cloneURL); err != nil {
		if !errors.Is(err, git.ErrRepositoryAlreadyExists) {
			fmt.Fprintf(i.commandIO.Err(), "Failed trying to clone proof repo: %v\n", err)
		}

		fmt.Fprintln(i.commandIO.Out(), "It's already there!")
	}

	// update config with new repo
}
