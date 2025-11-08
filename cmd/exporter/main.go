package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/bruli/pinger-exporter/internal/app"
	"github.com/bruli/pinger-exporter/internal/config"
	infrahttp "github.com/bruli/pinger-exporter/internal/infra/http"
	infranats "github.com/bruli/pinger-exporter/internal/infra/nats"
	infraprom "github.com/bruli/pinger-exporter/internal/infra/prometheus"
	"github.com/bruli/pinger/pkg/events"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	log := buildLogger()

	conf, err := config.New()
	if err != nil {
		log.Err(err).Msg("configuration error")
		os.Exit(1)
	}

	natsUrl := conf.NatsServerURL
	nc, err := nats.Connect(natsUrl)
	if err != nil {
		log.Err(err).Msg("failed to connect to nats")
		os.Exit(1)
	}
	log.Info().Msgf("connected to NATS server %s", natsUrl)
	defer func() {
		_ = nc.Drain()
		nc.Close()
	}()

	resourceMetricRepo := infraprom.NewResourceMetricRepository()
	logCHMdw := app.NewLogCommandHandlerMiddleware(log)
	createResourceMetricCH := logCHMdw(app.NewCreateResourceMetric(resourceMetricRepo))
	consumer := infranats.NewPingResultConsumer(ctx, createResourceMetricCH, log)

	subject := events.PingSubjet

	sub, err := nc.Subscribe(subject, consumer.Consume)
	if err != nil {
		log.Err(err).Msg("failed to subscribe to subject")
	}
	defer func() {
		_ = sub.Unsubscribe()
	}()

	log.Info().Msgf("subscribed to subject %s", subject)

	infrahttp.Run(nc, log)
}

func buildLogger() *zerolog.Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log := zerolog.New(os.Stdout).With().Timestamp().Logger()
	return &log
}
