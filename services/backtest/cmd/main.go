package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/inpour/event-driven-trader/services/backtest/internal/app"
	"github.com/inpour/event-driven-trader/services/backtest/internal/config"
	"github.com/inpour/event-driven-trader/services/backtest/internal/kafka"
	"github.com/inpour/event-driven-trader/services/backtest/internal/strategy/sp2l"
	"github.com/inpour/event-driven-trader/shared/events"
)

const defaultConfigPath = "/app/configs/config.yaml"

func main() {
	if err := run(); err != nil {
		slog.Error("backtest service stopped", "error", err)
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

	logger := newLogger(cfg.Service.LogLevel)

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	inputConsumer := kafka.NewInputCompletedConsumer(
		cfg.Kafka.Broker,
		cfg.Kafka.Consumer.GroupID,
		cfg.Kafka.Consumer.InputCompletedTopic,
	)
	defer func() {
		if err := inputConsumer.Close(); err != nil {
			logger.Error(
				"close input-completed consumer",
				"error", err,
			)
		}
	}()

	batchReader := kafka.NewBatchReader(
		cfg.Kafka.Broker,
		cfg.Kafka.Consumer.CandleTopic,
	)

	producer, err := kafka.NewProducer(
		cfg.Kafka.Broker,
		cfg.Kafka.Producer.MarkerCreatedTopic,
		cfg.Kafka.Producer.RunCompletedTopic,
	)
	if err != nil {
		return fmt.Errorf("create Kafka producer: %w", err)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			logger.Error("close Kafka producer", "error", err)
		}
	}()

	strategy, err := sp2l.New(sp2l.DefaultConfig())

	candleLoader := app.NewCandleLoader(batchReader)
	runner := app.NewRunner(strategy)
	markerPublisher := app.NewMarkerPublisher(producer)

	handler := app.NewInputCompletedHandler(
		candleLoader,
		runner,
		markerPublisher,
		logger,
	)

	logger.Info(
		"backtest service started",
		"service", cfg.Service.Name,
		"group_id", cfg.Kafka.Consumer.GroupID,
		"input_completed_topic",
		cfg.Kafka.Consumer.InputCompletedTopic,
		"candle_topic",
		cfg.Kafka.Consumer.CandleTopic,
		"marker_created_topic",
		cfg.Kafka.Producer.MarkerCreatedTopic,
		"run_completed_topic",
		cfg.Kafka.Producer.RunCompletedTopic,
	)

	return consumeInputCompletedEvents(
		ctx,
		inputConsumer,
		handler,
		logger,
	)
}

func consumeInputCompletedEvents(
	ctx context.Context,
	consumer *kafka.InputCompletedConsumer,
	handler *app.InputCompletedHandler,
	logger *slog.Logger,
) error {
	for {
		inputCompleted, message, err := consumer.Fetch(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return nil
			}

			logger.Error(
				"fetch backtest input-completed event",
				"error", err,
			)
		}

		if inputCompleted.EventType != events.EventTypeBacktestInputCompleted {
			logger.Warn(
				"skipping unexpected event type",
				"event_type", inputCompleted.EventType,
				"topic", message.Topic,
				"partition", message.Partition,
				"offset", message.Offset,
			)

			if err := consumer.Commit(ctx, message); err != nil {
				return fmt.Errorf(
					"commit unexpected event at offset %d: %w",
					message.Offset,
					err,
				)
			}

			continue
		}

		logger.Info(
			"received backtest input-completed event",
			"run_id", inputCompleted.RunID,
			"candle_count", inputCompleted.Payload.CandleCount,
			"topic", message.Topic,
			"partition", message.Partition,
			"offset", message.Offset,
		)

		if err := handler.Handle(ctx, inputCompleted); err != nil {
			if errors.Is(err, context.Canceled) {
				return nil
			}

			logger.Error(
				"backtest run failed; input event will remain uncommitted",
				"run_id", inputCompleted.RunID,
				"topic", message.Topic,
				"partition", message.Partition,
				"offset", message.Offset,
				"error", err,
			)
		}

		if err := consumer.Commit(ctx, message); err != nil {
			return fmt.Errorf(
				"commit successful input event for run %q at offset %d: %w",
				inputCompleted.RunID,
				message.Offset,
				err,
			)
		}

		logger.Info(
			"backtest input event committed",
			"run_id", inputCompleted.RunID,
			"topic", message.Topic,
			"partition", message.Partition,
			"offset", message.Offset,
		)
	}
}

func waitForRetry(
	ctx context.Context,
	delay time.Duration,
) error {
	timer := time.NewTimer(delay)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func newLogger(level string) *slog.Logger {
	return slog.New(
		slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level: parseLogLevel(level),
			},
		),
	)
}

func parseLogLevel(value string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
