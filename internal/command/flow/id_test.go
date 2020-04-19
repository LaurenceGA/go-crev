package flow

import (
	"context"
	"errors"
	"testing"

	"github.com/LaurenceGA/go-crev/internal/command/io"
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

type idFlowTestCase struct {
	name           string
	inputUsername  string
	getUserMock    getUserMock
	getRepoMock    getRepoMock
	mockSetIDError error
	expectedIDSet  *id.ID
	expectedError  bool
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
			getRepoMock: getRepoMock{
				expectedInput: expectedDefaultRepo("123"),
				repo:          &github.Repository{},
			},
			expectedIDSet: &id.ID{ID: "502", Type: id.Github},
		},
		{
			name:          "username with @",
			inputUsername: "@user",
			getUserMock: getUserMock{
				expectedInput: "user",
				usr:           &github.User{ID: 18, Login: "user"},
			},
			getRepoMock: getRepoMock{
				expectedInput: expectedDefaultRepo("user"),
				err:           github.NotFoundError,
			},
			expectedIDSet: &id.ID{ID: "18", Type: id.Github},
		},
		{
			name:          "username with double @",
			inputUsername: "@@user",
			getUserMock: getUserMock{
				expectedInput: "@user",
				usr:           &github.User{ID: 5, Login: "@user"},
			},
			getRepoMock: getRepoMock{
				expectedInput: expectedDefaultRepo("@user"),
				err:           github.NotFoundError,
			},
			expectedIDSet: &id.ID{ID: "5", Type: id.Github},
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
			getRepoMock: getRepoMock{
				expectedInput: expectedDefaultRepo("user"),
				err:           github.NotFoundError,
			},
			mockSetIDError: errors.New("failed to update config"),
			expectedIDSet:  &id.ID{ID: "5", Type: id.Github},
			expectedError:  true,
		},
		{
			name:          "Fail to get repo",
			inputUsername: "user",
			getUserMock: getUserMock{
				expectedInput: "user",
				usr:           &github.User{ID: 5, Login: "user"},
			},
			getRepoMock: getRepoMock{
				expectedInput: expectedDefaultRepo("user"),
				err:           errors.New("that's a fail"),
			},
			expectedIDSet: &id.ID{ID: "5", Type: id.Github},
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

		mockGithub.EXPECT().
			GetUser(gomock.Any(), testCase.getUserMock.expectedInput).
			Return(testCase.getUserMock.usr, testCase.getUserMock.err)

		if testCase.getUserMock.err == nil {
			mockGithub.EXPECT().
				GetRepository(gomock.Any(), testCase.getRepoMock.expectedInput.owner, testCase.getRepoMock.expectedInput.repo).
				Return(testCase.getRepoMock.repo, testCase.getRepoMock.err)

			mockConfigManipulator.EXPECT().
				SetCurrentID(testCase.expectedIDSet).
				Return(testCase.mockSetIDError)
		}

		idSetterFlow := NewIDSetter(&io.IO{}, mockConfigManipulator, mockGithub)

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
