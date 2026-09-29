package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Service ServiceConfig `yaml:"service"`
	Kafka   KafkaConfig   `yaml:"kafka"`
}

type ServiceConfig struct {
	Name     string `yaml:"name"`
	LogLevel string `yaml:"log_level"`
}

type KafkaConfig struct {
	Broker   string              `yaml:"broker"`
	Consumer KafkaConsumerConfig `yaml:"consumer"`
	Producer KafkaProducerConfig `yaml:"producer"`
}

type KafkaConsumerConfig struct {
	GroupID             string `yaml:"group_id"`
	InputCompletedTopic string `yaml:"input_completed_topic"`
	CandleTopic         string `yaml:"candle_topic"`
}

type KafkaProducerConfig struct {
	MarkerCreatedTopic string `yaml:"marker_created_topic"`
	RunCompletedTopic  string `yaml:"run_completed_topic"`
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
	var result error

	if strings.TrimSpace(c.Service.Name) == "" {
		result = errors.Join(
			result,
			errors.New("service.name is required"),
		)
	}

	if strings.TrimSpace(c.Service.LogLevel) == "" {
		result = errors.Join(
			result,
			errors.New("service.log_level is required"),
		)
	}

	if strings.TrimSpace(c.Kafka.Broker) == "" {
		result = errors.Join(
			result,
			errors.New("kafka.broker is required"),
		)
	}

	if strings.TrimSpace(c.Kafka.Consumer.GroupID) == "" {
		result = errors.Join(
			result,
			errors.New("kafka.consumer.group_id is required"),
		)
	}

	if strings.TrimSpace(c.Kafka.Consumer.InputCompletedTopic) == "" {
		result = errors.Join(
			result,
			errors.New(
				"kafka.consumer.input_completed_topic is required",
			),
		)
	}

	if strings.TrimSpace(c.Kafka.Consumer.CandleTopic) == "" {
		result = errors.Join(
			result,
			errors.New("kafka.consumer.candle_topic is required"),
		)
	}

	if strings.TrimSpace(c.Kafka.Producer.MarkerCreatedTopic) == "" {
		result = errors.Join(
			result,
			errors.New(
				"kafka.producer.marker_created_topic is required",
			),
		)
	}

	if strings.TrimSpace(c.Kafka.Producer.RunCompletedTopic) == "" {
		result = errors.Join(
			result,
			errors.New(
				"kafka.producer.run_completed_topic is required",
			),
		)
	}

	return result
}

func applyEnvironmentOverrides(cfg *Config) {
	if value := strings.TrimSpace(os.Getenv("SERVICE_NAME")); value != "" {
		cfg.Service.Name = value
	}

	if value := strings.TrimSpace(os.Getenv("LOG_LEVEL")); value != "" {
		cfg.Service.LogLevel = value
	}

	if value := strings.TrimSpace(os.Getenv("KAFKA_BROKER")); value != "" {
		cfg.Kafka.Broker = value
	}

	if value := strings.TrimSpace(os.Getenv("KAFKA_GROUP_ID")); value != "" {
		cfg.Kafka.Consumer.GroupID = value
	}

	if value := strings.TrimSpace(os.Getenv("KAFKA_INPUT_COMPLETED_TOPIC")); value != "" {
		cfg.Kafka.Consumer.InputCompletedTopic = value
	}

	if value := strings.TrimSpace(os.Getenv("KAFKA_CANDLE_TOPIC")); value != "" {
		cfg.Kafka.Consumer.CandleTopic = value
	}

	if value := strings.TrimSpace(os.Getenv("KAFKA_MARKER_CREATED_TOPIC")); value != "" {
		cfg.Kafka.Producer.MarkerCreatedTopic = value
	}

	if value := strings.TrimSpace(os.Getenv("KAFKA_RUN_COMPLETED_TOPIC")); value != "" {
		cfg.Kafka.Producer.RunCompletedTopic = value
	}
}
