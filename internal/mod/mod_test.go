package mod

import (
	"errors"
	"testing"

	"github.com/LaurenceGA/go-crev/internal/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type ModSuite struct {
	suite.Suite

	controller *gomock.Controller
}

func (s *ModSuite) SetupTest() {
	s.controller = gomock.NewController(s.T())
}

func (s *ModSuite) TearDownTest() {
	s.controller.Finish()
}

func TestMod(t *testing.T) {
	suite.Run(t, &ModSuite{})
}

func (s *ModSuite) TestLoadModules() {
	tests := []struct {
		name          string
		modulesErr    error
		wantErr       bool
		wantedModules []*Module
	}{
		{
			name:       "modules command error",
			modulesErr: errors.New("failed to run command"),
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		tt := tt
		s.Run(tt.name, func() {
			modsCmd := mocks.NewMockModulesWrapper(s.controller)
			modsCmd.EXPECT().List().Return(nil, tt.modulesErr)

			lister := NewLister(modsCmd)

			modules, err := lister.List()
			if tt.wantErr {
				s.Error(err)
			} else {
				s.NoError(err)
			}

			s.ElementsMatch(modules, tt.wantedModules)
		})
	}
}
