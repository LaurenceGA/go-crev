package mod

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"path/filepath"
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

func TestMod(t *testing.T) {
	suite.Run(t, &ModSuite{})
}

type modTestCase struct {
	name          string
	inputFile     string
	modulesErr    error
	wantErr       bool
	wantedModules []*Module
}

func (s *ModSuite) TestLoadModules() {
	tests := []modTestCase{
		{
			name:       "modules command error",
			modulesErr: errors.New("failed to run command"),
			wantErr:    true,
		},
		{
			name:      "Bad JSON",
			inputFile: "badJSON",
			wantErr:   true,
		},
		{
			name:      "JSON is a list",
			inputFile: "JSONList",
			wantErr:   true,
		},
		{
			name:          "Single empty JSON",
			inputFile:     "singleEmptyJSON",
			wantedModules: []*Module{{}},
		},
		{
			name:      "Single full struct",
			inputFile: "singleFullStruct",
			wantedModules: []*Module{
				{
					Path:      "mod/path",
					Main:      false,
					Dir:       "some/directory",
					GoMod:     "/some/mod/file/go.mod",
					GoVersion: "1.14",
				},
			},
		},
		{
			name:      "Multiple modules",
			inputFile: "multipleModules",
			wantedModules: []*Module{
				{Path: "mod1/path"},
				{Path: "mod2/path"},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		s.Run(tt.name, func() {
			s.runModulesTest(tt)
		})
	}
}

func (s *ModSuite) runModulesTest(tt modTestCase) {
	modsCmd := mocks.NewMockModulesWrapper(s.controller)

	var modulesOutput io.Reader
	if tt.modulesErr == nil {
		modulesOutput = s.loadFile(tt.inputFile)
	}

	modsCmd.EXPECT().List().Return(modulesOutput, tt.modulesErr)

	lister := NewLister(modsCmd)

	modules, err := lister.List()

	if tt.wantErr {
		s.Error(err)
	} else {
		s.NoError(err)
	}

	s.assertModulesArraysEqual(tt.wantedModules, modules)
}

func (s *ModSuite) assertModulesArraysEqual(expected []*Module, actual []*Module) {
	s.Len(actual, len(expected))

	for i := range actual {
		s.Equal(expected[i], actual[i])
	}
}

func (s *ModSuite) loadFile(filename string) io.Reader {
	fileData, err := ioutil.ReadFile(filepath.Join("testdata", filename))
	s.NoError(err)

	return bytes.NewReader(fileData)
}
