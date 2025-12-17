package checks

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type AzureYaml struct {
	Name             string             `yaml:"name"`
	Services         map[string]Service `yaml:"services"`
	Hooks            Hooks              `yaml:"hooks"`
	Infra            Infra              `yaml:"infra"`
	RequiredVersions RequiredVersions   `yaml:"requiredVersions"`
}

type RequiredVersions struct {
	Azd        string            `yaml:"azd"`
	Extensions map[string]string `yaml:"extensions"`
}

type Infra struct {
	Provider string `yaml:"provider"`
}

type Service struct {
	Language string       `yaml:"language"`
	Host     string       `yaml:"host"`
	Project  string       `yaml:"project"`
	Image    string       `yaml:"image"`
	Hooks    Hooks        `yaml:"hooks"`
	Docker   DockerConfig `yaml:"docker"`
}

type DockerConfig struct {
	Remote bool `yaml:"remote"`
}

type Hooks map[string]HookConfig

type HookConfig struct {
	Shell string `yaml:"shell"`
	Run   string `yaml:"run"`
}

// UnmarshalYAML implements custom unmarshaling for HookConfig to handle both string and object formats
func (h *HookConfig) UnmarshalYAML(value *yaml.Node) error {
	// Case 1: hook is just a string (the command to run)
	if value.Kind == yaml.ScalarNode {
		h.Run = value.Value
		return nil
	}

	// Case 2: hook is an object with shell and run properties
	type plain HookConfig
	return value.Decode((*plain)(h))
}

func LoadProjectConfig(path string) (*AzureYaml, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config AzureYaml
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse azure.yaml: %w", err)
	}

	return &config, nil
}
