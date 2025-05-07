package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config is the top-level YAML structure. Mirrors LiteLLM's liteLLM_config.yaml but simplified.
//
// Example:
// model_list:
//   - model_name: gpt-3.5-turbo
//     provider: openai
//     api_base: https://api.openai.com/v1
//     api_key_env: OPENAI_API_KEY
// rate_limit:
//   requests_per_minute: 60

type Config struct {
	ModelList []ModelConfig `yaml:"model_list"`
	RateLimit *RateLimit    `yaml:"rate_limit,omitempty"`
}

type ModelConfig struct {
	ModelName  string `yaml:"model_name"`
	Provider   string `yaml:"provider"`
	APIBase    string `yaml:"api_base"`
	APIKeyEnv  string `yaml:"api_key_env"`
	APIVersion string `yaml:"api_version,omitempty"`     // azure
	Deployment string `yaml:"deployment_name,omitempty"` // azure
}

type RateLimit struct {
	RequestsPerMinute int `yaml:"requests_per_minute"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	if len(cfg.ModelList) == 0 {
		return nil, fmt.Errorf("model_list cannot be empty")
	}
	return &cfg, nil
}
