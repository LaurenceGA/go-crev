// package flow detremines more complex user interaction flows
package flow

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

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

func NewIDSetter(configManipulator ConfigManipulator, githubUser Github) *IDSetter {
	return &IDSetter{
		configManipulator: configManipulator,
		githubUser:        githubUser,
	}
}

// IDSetter is responsible for high level flow of setting of a user's ID
type IDSetter struct {
	configManipulator ConfigManipulator
	githubUser        Github
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

	repo, err := i.githubUser.GetRepository(ctx, usr.Login, standardCrevProofRepoName)
	if err != nil {
		if errors.Is(err, github.NotFoundError) {
			fmt.Println("No repo!")
		} else {
			return fmt.Errorf("getting repository: %w", err)
		}
	}

	_ = repo

	return i.configManipulator.SetCurrentID(&id.ID{
		ID:   strconv.Itoa(int(usr.ID)),
		Type: id.Github,
	})
}
