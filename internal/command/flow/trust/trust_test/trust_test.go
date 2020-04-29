package trust_test

import (
	"context"
	"errors"
	"testing"

	"github.com/LaurenceGA/go-crev/internal/command/flow/trust"
	"github.com/LaurenceGA/go-crev/internal/command/flow/trust/mock"
	"github.com/LaurenceGA/go-crev/internal/command/io"
	"github.com/LaurenceGA/go-crev/internal/config"
	"github.com/LaurenceGA/go-crev/internal/github"
	"github.com/LaurenceGA/go-crev/internal/id"
	"github.com/LaurenceGA/go-crev/internal/store"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type mockConfigResponse struct {
	config *config.Configuration
	err    error
}

func (m *mockConfigResponse) getMock(controller *gomock.Controller) *mock.MockConfigReader {
	mck := mock.NewMockConfigReader(controller)
	mck.EXPECT().Load().Return(m.config, m.err)

	return mck
}

type trustPrompt struct {
	selection string
	err       error
}

type commentPrompt struct {
	comment string
	err     error
}

type mockPromptResponse struct {
	trustPrompt   *trustPrompt
	commentPrompt *commentPrompt
}

func (m *mockPromptResponse) getMock(controller *gomock.Controller) *mock.MockPrompter {
	mck := mock.NewMockPrompter(controller)

	if m.trustPrompt != nil {
		mck.EXPECT().
			Select(gomock.Any(), gomock.Any()).
			Return(m.trustPrompt.selection, m.trustPrompt.err)
	}

	if m.commentPrompt != nil {
		mck.EXPECT().
			Prompt(gomock.Any()).
			Return(m.commentPrompt.comment, m.commentPrompt.err)
	}

	return mck
}

type mockGetUser struct {
	expectedUsername string
	usr              *github.User
	err              error
}

type mockGetRepo struct {
	expectedOwner, expectedRepo string
	repo                        *github.Repository
	err                         error
}

type mockGithubResponse struct {
	mockGetUser *mockGetUser
	mockGetRepo *mockGetRepo
}

func (m *mockGithubResponse) getMock(controller *gomock.Controller) *mock.MockGithub {
	mck := mock.NewMockGithub(controller)

	if m.mockGetUser != nil {
		mck.EXPECT().
			GetUser(gomock.Any(), m.mockGetUser.expectedUsername).
			Return(m.mockGetUser.usr, m.mockGetUser.err)
	}

	if m.mockGetRepo != nil {
		mck.EXPECT().
			GetRepository(gomock.Any(), m.mockGetRepo.expectedOwner, m.mockGetRepo.expectedRepo).
			Return(m.mockGetRepo.repo, m.mockGetRepo.err)
	}

	return mck
}

func TestCannotReadConfig(t *testing.T) {
	const (
		testStore          = "/my/store"
		testUsername       = "user"
		testID             = 123
		testLevelSelection = "None"
	)

	var (
		testMockConfig = mockConfigResponse{
			config: &config.Configuration{
				CurrentStore: testStore,
				CurrentID:    &id.ID{},
			},
		}
		testMockGetUser = &mockGetUser{
			expectedUsername: testUsername,
			usr: &github.User{
				ID:    testID,
				Login: testUsername,
			},
		}
		testMockGetRepo = &mockGetRepo{
			expectedOwner: testUsername,
			expectedRepo:  store.StandardCrevProofRepoName,
			repo:          &github.Repository{},
		}
		testMockGithub = mockGithubResponse{
			mockGetUser: testMockGetUser,
			mockGetRepo: testMockGetRepo,
		}
		testPrompt = mockPromptResponse{
			&trustPrompt{
				selection: "None",
			},
			&commentPrompt{},
		}
	)

	for _, testCase := range []struct {
		name               string
		usernameInput      string
		mockConfigResponse mockConfigResponse
		mockGithubResponse mockGithubResponse
		mockPromptResponse mockPromptResponse
		expectError        bool
	}{
		{
			name: "Cannot read config",
			mockConfigResponse: mockConfigResponse{
				err: errors.New("can't read config"),
			},
			expectError: true,
		},
		{
			name: "No store set",
			mockConfigResponse: mockConfigResponse{
				config: &config.Configuration{
					CurrentStore: "",
					CurrentID:    &id.ID{},
				},
			},
			expectError: true,
		},
		{
			name: "No ID set",
			mockConfigResponse: mockConfigResponse{
				config: &config.Configuration{
					CurrentStore: testStore,
					CurrentID:    nil,
				},
			},
			expectError: true,
		},
		{
			name:               "Error getting user",
			usernameInput:      "user",
			mockConfigResponse: testMockConfig,
			mockGithubResponse: mockGithubResponse{
				mockGetUser: &mockGetUser{
					expectedUsername: "user",
					err:              errors.New("failed to talk to Github"),
				},
			},
			expectError: true,
		},
		{
			name:               "Fail trying to get repo",
			usernameInput:      testUsername,
			mockConfigResponse: testMockConfig,
			mockGithubResponse: mockGithubResponse{
				mockGetUser: testMockGetUser,
				mockGetRepo: &mockGetRepo{
					expectedOwner: testUsername,
					expectedRepo:  store.StandardCrevProofRepoName,
					err:           errors.New("can't talk to Github"),
				},
			},
			mockPromptResponse: testPrompt,
			expectError:        false, // Error is non-fatal
		},
		{
			name:               "Invalid level selection",
			usernameInput:      testUsername,
			mockConfigResponse: testMockConfig,
			mockGithubResponse: testMockGithub,
			mockPromptResponse: mockPromptResponse{
				trustPrompt: &trustPrompt{
					selection: "Not a level",
				},
			},
			expectError: true,
		},
		{
			name:               "Failed trust level prompt",
			usernameInput:      testUsername,
			mockConfigResponse: testMockConfig,
			mockGithubResponse: testMockGithub,
			mockPromptResponse: mockPromptResponse{
				trustPrompt: &trustPrompt{
					err: errors.New("nope"),
				},
			},
			expectError: true,
		},
	} {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			controller := gomock.NewController(t)

			trustCreator := trust.NewTrustCreator(
				&io.IO{},
				testCase.mockConfigResponse.getMock(controller),
				testCase.mockGithubResponse.getMock(controller),
				testCase.mockPromptResponse.getMock(controller),
			)

			err := trustCreator.CreateTrust(context.Background(), testCase.usernameInput, trust.CreatorOptions{})

			if testCase.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			controller.Finish()
		})
	}
}
