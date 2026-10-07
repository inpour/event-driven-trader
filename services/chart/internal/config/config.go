package config

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Kafka  KafkaConfig `yaml:"kafka"`
	Result Result      `yaml:"result"`
}

type KafkaConfig struct {
	Broker   string   `yaml:"broker"`
	Consumer Consumer `yaml:"consumer"`
}

type Consumer struct {
	GroupID string `yaml:"group_id"`
}

type Result struct {
	OutputPath string `yaml:"output_path"`
}

// Load reads YAML configuration from path and applies supported environment
// variable overrides.
func Load(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("read config file %q: %w", path, err)
	}

	var cfg Config

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("decode config file %q: %w", path, err)
	}

	applyEnvironmentOverrides(&cfg)

	if err := cfg.Validate(); err != nil {
		return Config{}, fmt.Errorf("validate config: %w", err)
	}

	return cfg, nil
}

func (c Config) Validate() error {
	if strings.TrimSpace(c.Kafka.Broker) == "" {
		return fmt.Errorf("kafka.broker is required")
	}

	if strings.TrimSpace(c.Kafka.Consumer.GroupID) == "" {
		return fmt.Errorf("kafka.consumer.group_id is required")
	}

	if strings.TrimSpace(c.Result.OutputPath) == "" {
		return fmt.Errorf("result.output_path is required")
	}

	return nil
}

func applyEnvironmentOverrides(cfg *Config) {
	if value := os.Getenv("KAFKA_BROKER"); value != "" {
		cfg.Kafka.Broker = value
	}

	if value := os.Getenv("KAFKA_CONSUMER_GROUP_ID"); value != "" {
		cfg.Kafka.Consumer.GroupID = value
	}

	if value := os.Getenv("RESULT_OUTPUT_PATH"); value != "" {
		cfg.Result.OutputPath = value
	}
}
