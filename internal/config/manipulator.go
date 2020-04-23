package config

import (
	"fmt"
	"io/ioutil"

	"github.com/LaurenceGA/go-crev/internal/files"
	"github.com/LaurenceGA/go-crev/internal/id"
	"gopkg.in/yaml.v3"
)

func NewManipulator(appDirs files.AppDirs) *Manipulator {
	return &Manipulator{
		appDirs: appDirs,
	}
}

type Manipulator struct {
	appDirs files.AppDirs
}

func (m *Manipulator) Load() (*Configuration, error) {
	configFilepath, err := m.appDirs.ConfigFile()
	if err != nil {
		return nil, fmt.Errorf("getting config file path: %w", err)
	}

	var conf Configuration

	data, err := ioutil.ReadFile(configFilepath)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	if err := yaml.Unmarshal(data, &conf); err != nil {
		return nil, fmt.Errorf("umarshalling config YAML: %w", err)
	}

	return &conf, nil
}

func (m *Manipulator) Save(conf *Configuration) error {
	configFilepath, err := m.appDirs.ConfigFile()
	if err != nil {
		return fmt.Errorf("getting config file path: %w", err)
	}

	data, err := yaml.Marshal(&conf)
	if err != nil {
		return fmt.Errorf("marshalling configuration: %w", err)
	}

	if err := ioutil.WriteFile(configFilepath, data, 0666); err != nil {
		return fmt.Errorf("writing config file: %w", err)
	}

	return nil
}

func (m *Manipulator) CurrentStore() (string, error) {
	config, err := m.Load()
	if err != nil {
		return "", err
	}

	return config.CurrentStore, nil
}

func (m *Manipulator) SetCurrentStore(path string) error {
	config, err := m.Load()
	if err != nil {
		return err
	}

	fmt.Printf("Setting current store to %s\n", path)

	config.CurrentStore = path

	return m.Save(config)
}

func (m *Manipulator) CurrentID() (*id.ID, error) {
	config, err := m.Load()
	if err != nil {
		return nil, err
	}

	return config.CurrentID, nil
}

func (m *Manipulator) SetCurrentID(id *id.ID) error {
	config, err := m.Load()
	if err != nil {
		return err
	}

	fmt.Printf("Setting current ID to %s\n", id)

	config.CurrentID = id

	return m.Save(config)
}
