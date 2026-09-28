package config

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Kafka    KafkaConfig    `yaml:"kafka"`
	Consumer ConsumerConfig `yaml:"consumer"`
}

type KafkaConfig struct {
	Broker  string `yaml:"broker"`
	GroupID string `yaml:"group_id"`
}

type ConsumerConfig struct {
	StartOffset string `yaml:"start_offset"`
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

	if strings.TrimSpace(c.Kafka.GroupID) == "" {
		return fmt.Errorf("kafka.group_id is required")
	}

	switch c.Consumer.StartOffset {
	case "earliest", "latest":
		return nil
	default:
		return fmt.Errorf(
			"consumer.start_offset must be %q or %q, got %q",
			"earliest",
			"latest",
			c.Consumer.StartOffset,
		)
	}
}

func applyEnvironmentOverrides(cfg *Config) {
	if value := os.Getenv("KAFKA_BROKER"); value != "" {
		cfg.Kafka.Broker = value
	}

	if value := os.Getenv("KAFKA_GROUP_ID"); value != "" {
		cfg.Kafka.GroupID = value
	}

	if value := os.Getenv("CONSUMER_START_OFFSET"); value != "" {
		cfg.Consumer.StartOffset = value
	}
}
