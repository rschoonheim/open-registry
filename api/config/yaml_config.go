package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// YAMLConfig holds the configuration loaded from open-registry.yaml
type YAMLConfig struct {
	Environment struct {
		File string `yaml:"file"`
	} `yaml:"environment"`
	Authentication struct {
		Register struct {
			Enabled bool `yaml:"enabled"`
		} `yaml:"register"`
	} `yaml:"authentication"`
}

// LoadYAMLConfig loads configuration from the open-registry.yaml file
func LoadYAMLConfig(configPath string) (*YAMLConfig, error) {
	// If configPath is empty, use default path
	if configPath == "" {
		// Get the current working directory
		workDir, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("failed to get working directory: %w", err)
		}

		// Check if we're in the cmd directory and need to go up one level
		if filepath.Base(workDir) == "cmd" {
			workDir = filepath.Dir(workDir)
		}

		configPath = filepath.Join(workDir, "open-registry.yaml")
	}

	// Check if the file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Printf("Warning: YAML config file not found at %s, using defaults", configPath)
		return &YAMLConfig{}, nil // Initialize with default values
	}

	// Read the file
	yamlFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read YAML config file: %w", err)
	}

	// Parse YAML
	yamlConfig := &YAMLConfig{}
	err = yaml.Unmarshal(yamlFile, yamlConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to parse YAML config: %w", err)
	}

	log.Printf("Loaded YAML configuration from %s", configPath)
	return yamlConfig, nil
}
