package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"

	"github.com/ghodss/yaml"
	"go.uber.org/zap"
)

type StateFile struct {
	Path string `json: path`
}

type StateTarget struct {
	Files      []StateFile `json:"files"`
}

type StateConfig struct {
	APIVersion string                 `json:"apiVersion"`
	Kind       string                 `json:"kind"`
	Targets     map[string]StateTarget `json:"target"`
}

func NewStateConfig() StateConfig {
	config := StateConfig{
		APIVersion: "simopolis.xyz/v1alpha1",
		Kind:       "StateConfig",
	}
	if config.Targets == nil {
		config.Targets = make(map[string]StateTarget)
	}

	return config
}

func NewStateConfigFromFile(path string) (StateConfig, error) {
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return StateConfig{}, fmt.Errorf("Reading state config '%s': %s", path, err)
	}

	return NewStateConfigFromBytes(bs)
}

func NewStateConfigFromBytes(bs []byte) (StateConfig, error) {
	var config StateConfig

	err := yaml.Unmarshal(bs, &config)
	if err != nil {
		return StateConfig{}, fmt.Errorf("Unmarshaling state config: %s", err)
	}

	err = config.Validate()
	if err != nil {
		return StateConfig{}, fmt.Errorf("Validating state config: %s", err)
	}

	if config.Targets == nil {
		config.Targets = make(map[string]StateTarget)
	}
	return config, nil
}

func (c StateConfig) WriteToFile(path string) error {
	c.Sort()
	bs, err := c.AsBytes()
	if err != nil {
		return fmt.Errorf("Marshaling state config: %s", err)
	}

	err = ioutil.WriteFile(path, bs, 0700)
	if err != nil {
		return fmt.Errorf("Writing state config: %s", err)
	}

	return nil
}

func (c StateConfig) AsBytes() ([]byte, error) {
	bs, err := yaml.Marshal(c)
	if err != nil {
		return nil, fmt.Errorf("Marshaling state config: %s", err)
	}

	return bs, nil
}

func (c StateConfig) Sort() {
	for _, state := range c.Targets {
		sort.Slice(state.Files, func(i, j int) bool {
			return state.Files[i].Path < state.Files[j].Path
		})
	}
}

func (c StateConfig) Validate() error {
	const (
		knownAPIVersion = "simopolis.xyz/v1alpha1"
		knownKind       = "StateConfig"
	)

	if c.APIVersion != knownAPIVersion {
		return fmt.Errorf("Validating apiVersion: Unknown version (known: %s)", knownAPIVersion)
	}
	if c.Kind != knownKind {
		return fmt.Errorf("Validating kind: Unknown kind (known: %s)", knownKind)
	}
	return nil
}

func Cleanup(stateFilePath string, targetName string, previousState StateConfig, newState StateConfig, logger *zap.SugaredLogger) error {
	// Calculate what should be deleted
	toDelete := make(map[string]bool)
	for _, s := range previousState.Targets[targetName].Files {
		toDelete[s.Path] = true
	}

	// Make everything in the new state as to be kept
	for _, s := range newState.Targets[targetName].Files {
		toDelete[s.Path] = false
	}

	cwd, err := os.Getwd()
	if err != nil {
		logger.Errorf("Current working directory has vanished: %s", cwd)
	}
	os.Chdir(filepath.Dir(stateFilePath))
	for file, delete := range toDelete {
		if delete == false {
			continue
		}
		logger.Debugf("Deleting: %s", file)
		if err := os.Remove(file); err != nil {
			logger.Infof("Couldn't delete: %s", file)
			// TODO should delete empty directories
		}
	}
	os.Chdir(cwd)
	return nil
}
