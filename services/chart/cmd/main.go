package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/inpour/event-driven-trader/services/chart/internal/app"
	"github.com/inpour/event-driven-trader/services/chart/internal/config"
	"github.com/inpour/event-driven-trader/services/chart/internal/kafka"
	"github.com/inpour/event-driven-trader/services/chart/internal/render"
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

	RunCompletedConsumer := kafka.NewRunCompletedConsumer(
		cfg.Kafka.Broker,
		cfg.Kafka.Consumer.GroupID,
		kafka.TopicRunCompleted,
	)
	defer func() {
		if err := RunCompletedConsumer.Close(); err != nil {
			fmt.Fprintln(os.Stderr, fmt.Errorf("close run-completed consumer: %w", err))
		}
	}()

	batchReader := kafka.NewBatchReader(
		cfg.Kafka.Broker,
		kafka.TopicCandleClosed,
		kafka.TopicMarkerCreated,
	)

	candleLoader := app.NewCandleLoader(batchReader)
	markerLoader := app.NewMarkerLoader(batchReader)

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	chartRenderer := render.NewChartRenderer(cfg.Result.OutputPath)

	handler := app.NewRunCompletedHandler(
		candleLoader,
		markerLoader,
		chartRenderer,
	)

	consumer := app.NewRunCompleteConsumer(
		RunCompletedConsumer,
		handler,
	)

	return consumer.Consume(ctx)
}
