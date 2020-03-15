package command

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type FetchTestSuite struct {
	suite.Suite
}

func TestFetchTestSuite(t *testing.T) {
	suite.Run(t, &FetchTestSuite{})
}

func (s *FetchTestSuite) TestFetchArgs() {
	tests := []struct {
		name        string
		args        []string
		expectError bool
	}{
		{
			name:        "nil arguments",
			expectError: true,
		},
		{
			name:        "no arguments",
			args:        []string{},
			expectError: true,
		},
		{
			name:        "one argument",
			args:        []string{"an argument"},
			expectError: false,
		},
		{
			name:        "too many arguments",
			args:        []string{"one", "two"},
			expectError: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		s.Run(tt.name, func() {
			err := fetchStore(nil, tt.args)
			if tt.expectError {
				s.Error(err)
			} else {
				s.NoError(err)
			}
		})
	}
}
