package config

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config contains all runtime configuration for market-data-service.
type Config struct {
	Kafka  KafkaConfig  `yaml:"kafka"`
	Input  InputConfig  `yaml:"input"`
	Market MarketConfig `yaml:"market"`
	Run    RunConfig    `yaml:"run"`
	Replay ReplayConfig `yaml:"replay"`
}

type KafkaConfig struct {
	Broker string `yaml:"broker"`
}

type InputConfig struct {
	CSVPath string `yaml:"csv_path"`
}

type MarketConfig struct {
	Symbol    string `yaml:"symbol"`
	Timeframe string `yaml:"timeframe"`
}

type RunConfig struct {
	ID string `yaml:"id"`
}

type ReplayConfig struct {
	Mode string `yaml:"mode"`
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

// Validate ensures all required configuration fields have usable values.
func (c Config) Validate() error {
	if strings.TrimSpace(c.Kafka.Broker) == "" {
		return fmt.Errorf("kafka.broker is required")
	}

	if strings.TrimSpace(c.Input.CSVPath) == "" {
		return fmt.Errorf("input.csv_path is required")
	}

	if strings.TrimSpace(c.Market.Symbol) == "" {
		return fmt.Errorf("market.symbol is required")
	}

	if strings.TrimSpace(c.Market.Timeframe) == "" {
		return fmt.Errorf("market.timeframe is required")
	}

	if strings.TrimSpace(c.Run.ID) == "" {
		return fmt.Errorf("run.id is required")
	}

	switch c.Replay.Mode {
	case "fast":
		return nil
	default:
		return fmt.Errorf(
			"replay.mode must be %q, got %q",
			"fast",
			c.Replay.Mode,
		)
	}
}

func applyEnvironmentOverrides(cfg *Config) {
	if value := os.Getenv("KAFKA_BROKER"); value != "" {
		cfg.Kafka.Broker = value
	}

	if value := os.Getenv("INPUT_CSV_PATH"); value != "" {
		cfg.Input.CSVPath = value
	}

	if value := os.Getenv("MARKET_SYMBOL"); value != "" {
		cfg.Market.Symbol = value
	}

	if value := os.Getenv("MARKET_TIMEFRAME"); value != "" {
		cfg.Market.Timeframe = value
	}

	if value := os.Getenv("RUN_ID"); value != "" {
		cfg.Run.ID = value
	}

	if value := os.Getenv("REPLAY_MODE"); value != "" {
		cfg.Replay.Mode = value
	}
}
