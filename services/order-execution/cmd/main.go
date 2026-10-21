package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

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

	// Start Prometheus metrics HTTP server in the background
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	metricsServer := &http.Server{
		Addr:    ":" + cfg.Prometheus.MetricsServerPort,
		Handler: mux,
	}

	go func() {
		if err := metricsServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Fprintf(os.Stderr, "metrics server error: %v\n", err)
		}
	}()
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = metricsServer.Shutdown(shutdownCtx)
	}()

	err = service.Run(ctx)
	if errors.Is(err, context.Canceled) {
		return nil
	}

	return err
}
