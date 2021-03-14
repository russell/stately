package config

import (
	"fmt"
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

func (c StateConfig) WriteToFile(path string) error {
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
