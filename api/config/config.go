package config

// Config holds all configuration for our application
type Config struct {
	Env  *EnvConfig
	YAML *YAMLConfig
}

// LoadConfig loads all configuration
func LoadConfig() (*Config, error) {
	// Load environment configuration
	envConfig := LoadEnvConfig()

	// Load YAML configuration
	yamlConfig, err := LoadYAMLConfig("")
	if err != nil {
		return nil, err
	}

	// Combine configurations
	config := &Config{
		Env:  envConfig,
		YAML: yamlConfig,
	}

	return config, nil
}

// GetDSN returns the Data Source Name for the database
func (c *Config) GetDSN() string {
	return c.Env.GetDSN()
}
