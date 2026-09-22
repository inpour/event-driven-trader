package main

import (
	"context"
	"fmt"
	"os"

	"github.com/inpour/event-driven-trader/services/market-data/internal/app"
	"github.com/inpour/event-driven-trader/services/market-data/internal/config"
	"github.com/inpour/event-driven-trader/services/market-data/internal/kafka"
	csvfeed "github.com/inpour/event-driven-trader/services/market-data/internal/provider/csv"
)

const defaultConfigPath = "/app/config/config.yaml"

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = defaultConfigPath
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("load configuration: %w", err)
	}

	provider := csvfeed.NewProvider(cfg.Input.CSVPath)
	producer := kafka.NewProducer(cfg.Kafka.Broker)
	defer producer.Close()

	service := app.NewService(
		provider,
		producer,
		cfg.Run.ID,
		cfg.Market.Symbol,
		cfg.Market.Timeframe,
	)

	if err := service.Run(context.Background()); err != nil {
		return fmt.Errorf("run market-data service: %w", err)
	}

	return nil
}
