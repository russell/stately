package config

import (
	"fmt"
	"sort"
	"io/ioutil"

	"github.com/ghodss/yaml"
)

type StateFile struct {
	Path string `json: path`
}

type StateConfig struct {
	APIVersion      string      `json:"apiVersion"`
	Kind            string      `json:"kind"`
	OutputDirectory string      `json: outputDirectory`
	Files           []StateFile `json:"directories"`
}

func NewStateConfig() StateConfig {
	return StateConfig{
		APIVersion: "simopolis.xyz/v1alpha1",
		Kind:       "StateConfig",
	}
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

func (c StateConfig) Sort() () {
	sort.Slice(c.Files, func(i, j int) bool {
		return c.Files[i].Path < c.Files[j].Path
	})
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
