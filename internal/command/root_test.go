package command

import (
	"testing"

	"github.com/LaurenceGA/go-crev/internal/command/io"
	"github.com/stretchr/testify/suite"
)

type RootTestSuite struct {
	suite.Suite
}

func TestRootTestSuite(t *testing.T) {
	suite.Run(t, &RootTestSuite{})
}

func (s *RootTestSuite) TestRootCommands() {
	tests := []struct {
		name        string
		args        []string
		expectError bool
	}{
		{
			name:        "base command",
			args:        []string{},
			expectError: true,
		},
		{
			name:        "plain store",
			args:        []string{"store"},
			expectError: true,
		},
		{
			name:        "invalid command",
			args:        []string{"nonExistant"},
			expectError: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		s.Run(tt.name, func() {
			cmd := NewRootCommand(&io.IO{})
			cmd.SetArgs(tt.args)

			err := cmd.Execute()

			if tt.expectError {
				s.Error(err)
			} else {
				s.NoError(err)
			}
		})
	}
}
