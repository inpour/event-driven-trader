package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/inpour/event-driven-trader/services/order-execution/internal/app"
	"github.com/inpour/event-driven-trader/services/order-execution/internal/config"
	"github.com/inpour/event-driven-trader/services/order-execution/internal/kafka"
)

const defaultConfigPath = "/app/configs/config.yaml"

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

	kafkaConsumer, err := kafka.NewConsumer(
		cfg.Kafka.Broker,
		cfg.Kafka.GroupID,
		cfg.Consumer.StartOffset,
	)
	if err != nil {
		return fmt.Errorf("create Kafka consumer: %w", err)
	}
	defer kafkaConsumer.Close()

	service := app.NewService(kafkaConsumer)

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	err = service.Run(ctx)
	if errors.Is(err, context.Canceled) {
		return nil
	}

	return err
}
