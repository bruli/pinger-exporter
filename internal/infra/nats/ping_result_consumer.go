package nats

import (
	"context"

	"github.com/bruli/pinger-exporter/internal/app"
	"github.com/bruli/pinger-exporter/internal/domain/metrics"
	"github.com/bruli/pinger/pkg/events"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
	"google.golang.org/protobuf/proto"
)

type PingResultConsumer struct {
	log        *zerolog.Logger
	resourceCH app.CommandHandler
	ctx        context.Context
}

func (c PingResultConsumer) Consume(msg *nats.Msg) {
	select {
	case <-c.ctx.Done():
		return
	default:
		var m events.PingResult
		if err := proto.Unmarshal(msg.Data, &m); err != nil {
			c.log.Error().Err(err).Msg("error while unmarshalling message")
			return
		}
		re, err := metrics.NewResource(m.GetResource(), float64(m.GetLatency()/1000), m.GetCreatedAt().AsTime())
		if err != nil {
			c.log.Error().Err(err).Msg("error while unmarshalling message")
			return
		}
		_ = c.resourceCH.Handle(context.Background(), app.CreateResourceMetricCommand{Resource: re})
	}
}

func NewPingResultConsumer(ctx context.Context, resourceCH app.CommandHandler, log *zerolog.Logger) *PingResultConsumer {
	return &PingResultConsumer{log: log, resourceCH: resourceCH, ctx: ctx}
}
