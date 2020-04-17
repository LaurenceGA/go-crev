// package flow detremines more complex user interaction flows
package flow

import (
	"context"
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

type GithubUser interface {
	GetUser(context.Context, string) (*github.User, error)
}

func NewIDSetter(configManipulator ConfigManipulator, githubUser GithubUser) *IDSetter {
	return &IDSetter{
		configManipulator: configManipulator,
		githubUser:        githubUser,
	}
}

// IDSetter is responsible for high level flow of setting of a user's ID
type IDSetter struct {
	configManipulator ConfigManipulator
	githubUser        GithubUser
}

// SetFromUsername takes a username input and sets our local current ID to that by finding
// resolving the ID from the username.
// By default we assume it's a Github ID
func (i *IDSetter) SetFromUsername(ctx context.Context, usernameRaw string) error {
	username := strings.TrimPrefix(usernameRaw, "@")

	usr, err := i.githubUser.GetUser(ctx, username)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	return i.configManipulator.SetCurrentID(&id.ID{
		ID:   strconv.Itoa(int(usr.ID)),
		Type: id.Github,
	})
}
