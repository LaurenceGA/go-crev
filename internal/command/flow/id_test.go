package flow

import (
	"context"
	"errors"
	"testing"

	"github.com/LaurenceGA/go-crev/internal/command/io"
	"github.com/LaurenceGA/go-crev/internal/git"
	"github.com/LaurenceGA/go-crev/internal/github"
	"github.com/LaurenceGA/go-crev/internal/id"
	"github.com/LaurenceGA/go-crev/internal/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type getUserMock struct {
	expectedInput string
	usr           *github.User
	err           error
}

type getRepoInput struct {
	owner, repo string
}

type getRepoMock struct {
	expectedInput getRepoInput
	repo          *github.Repository
	err           error
}

type fetchRepoMock struct {
	expectedFetchURL string
	err              error
}

type setIDMock struct {
	expectedSetID *id.ID
	err           error
}

type idFlowTestCase struct {
	name          string
	inputUsername string
	getUserMock   getUserMock
	getRepoMock   *getRepoMock
	fetchRepoMock *fetchRepoMock
	setIDMock     *setIDMock
	// mockSetIDError error
	// expectedIDSet  *id.ID
	expectedError bool
}

func TestIDFlow(t *testing.T) {
	for _, testCase := range []idFlowTestCase{
		{
			name:          "simple success",
			inputUsername: "123",
			getUserMock: getUserMock{
				expectedInput: "123",
				usr:           &github.User{ID: 502, Login: "123"},
			},
			getRepoMock: &getRepoMock{
				expectedInput: expectedDefaultRepo("123"),
				repo:          &github.Repository{CloneURL: "cloneURL", HTMLurl: "htmlURL"},
			},
			fetchRepoMock: &fetchRepoMock{expectedFetchURL: "cloneURL"},
			setIDMock:     &setIDMock{expectedSetID: githubIDWithURL("502", "htmlURL")},
		},
		{
			name:          "username with @",
			inputUsername: "@user",
			getUserMock: getUserMock{
				expectedInput: "user",
				usr:           &github.User{ID: 18, Login: "user"},
			},
			getRepoMock: &getRepoMock{
				expectedInput: expectedDefaultRepo("user"),
				err:           github.NotFoundError,
			},
			setIDMock: &setIDMock{expectedSetID: githubID("18")},
		},
		{
			name:          "username with double @",
			inputUsername: "@@user",
			getUserMock: getUserMock{
				expectedInput: "@user",
				usr:           &github.User{ID: 5, Login: "@user"},
			},
			getRepoMock: &getRepoMock{
				expectedInput: expectedDefaultRepo("@user"),
				err:           github.NotFoundError,
			},
			setIDMock: &setIDMock{expectedSetID: githubID("5")},
		},
		{
			name:          "Get user error",
			inputUsername: "user",
			getUserMock: getUserMock{
				expectedInput: "user",
				err:           errors.New("not found"),
			},
			expectedError: true,
		},
		{
			name:          "Fail to set ID",
			inputUsername: "user",
			getUserMock: getUserMock{
				expectedInput: "user",
				usr:           &github.User{ID: 5, Login: "user"},
			},
			getRepoMock: &getRepoMock{
				expectedInput: expectedDefaultRepo("user"),
				err:           github.NotFoundError,
			},
			setIDMock: &setIDMock{
				expectedSetID: githubID("5"),
				err:           errors.New("failed to update config"),
			},
			expectedError: true,
		},
		{
			name:          "Fail to get repo from github",
			inputUsername: "user",
			getUserMock: getUserMock{
				expectedInput: "user",
				usr:           &github.User{ID: 5, Login: "user"},
			},
			getRepoMock: &getRepoMock{
				expectedInput: expectedDefaultRepo("user"),
				err:           errors.New("that's a fail"),
			},
			setIDMock:     &setIDMock{expectedSetID: githubID("5")},
			expectedError: false, // Failing to find repo is non-fatal
		},
		{
			name:          "Repo already exists",
			inputUsername: "user",
			getUserMock: getUserMock{
				expectedInput: "user",
				usr:           &github.User{ID: 5, Login: "user"},
			},
			getRepoMock: &getRepoMock{
				expectedInput: expectedDefaultRepo("user"),
				repo:          &github.Repository{CloneURL: "cloneURL", HTMLurl: "htmlURL"},
			},
			fetchRepoMock: &fetchRepoMock{
				expectedFetchURL: "cloneURL",
				err:              git.ErrRepositoryAlreadyExists,
			},
			setIDMock:     &setIDMock{expectedSetID: githubIDWithURL("5", "htmlURL")},
			expectedError: false, // Failing to find repo is non-fatal
		},
		{
			name:          "Fail to clone repo",
			inputUsername: "user",
			getUserMock: getUserMock{
				expectedInput: "user",
				usr:           &github.User{ID: 5, Login: "user"},
			},
			getRepoMock: &getRepoMock{
				expectedInput: expectedDefaultRepo("user"),
				repo:          &github.Repository{CloneURL: "cloneURL", HTMLurl: "htmlURL"},
			},
			fetchRepoMock: &fetchRepoMock{
				expectedFetchURL: "cloneURL",
				err:              errors.New("That's a fail"),
			},
			setIDMock:     &setIDMock{expectedSetID: githubIDWithURL("5", "htmlURL")},
			expectedError: false, // Failing to find repo is non-fatal
		},
	} {
		t.Run(testCase.name, runIDFlowTestCase(testCase))
	}
}

func runIDFlowTestCase(testCase idFlowTestCase) func(*testing.T) {
	return func(t *testing.T) {
		controller := gomock.NewController(t)
		mockConfigManipulator := mocks.NewMockConfigManipulator(controller)
		mockGithub := mocks.NewMockGithub(controller)
		mockRepoFetcher := mocks.NewMockRepoFetcher(controller)

		mockGithub.EXPECT().
			GetUser(gomock.Any(), testCase.getUserMock.expectedInput).
			Return(testCase.getUserMock.usr, testCase.getUserMock.err)

		if testCase.getRepoMock != nil {
			mockGithub.EXPECT().
				GetRepository(gomock.Any(), testCase.getRepoMock.expectedInput.owner, testCase.getRepoMock.expectedInput.repo).
				Return(testCase.getRepoMock.repo, testCase.getRepoMock.err)
		}

		if testCase.setIDMock != nil {
			mockConfigManipulator.EXPECT().
				SetCurrentID(testCase.setIDMock.expectedSetID).
				Return(testCase.setIDMock.err)
		}

		if testCase.fetchRepoMock != nil {
			mockRepoFetcher.EXPECT().
				Fetch(gomock.Any(), testCase.fetchRepoMock.expectedFetchURL).
				Return(testCase.fetchRepoMock.err)
		}

		idSetterFlow := NewIDSetter(&io.IO{}, mockConfigManipulator, mockGithub, mockRepoFetcher)

		err := idSetterFlow.SetFromUsername(context.Background(), testCase.inputUsername)

		if testCase.expectedError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}

		controller.Finish()
	}
}

func expectedDefaultRepo(owner string) getRepoInput {
	return getRepoInput{
		owner: owner,
		repo:  standardCrevProofRepoName,
	}
}

func githubID(githubID string) *id.ID {
	return &id.ID{
		ID:   githubID,
		Type: id.Github,
	}
}

func githubIDWithURL(githubID, url string) *id.ID {
	return &id.ID{
		ID:   githubID,
		Type: id.Github,
		URL:  url,
	}
}
