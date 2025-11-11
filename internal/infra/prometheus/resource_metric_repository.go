package prometheus

import (
	"context"

	"github.com/bruli/pinger-exporter/internal/domain/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

type ResourceMetricRepository struct {
	eventsTotal         *prometheus.CounterVec
	eventsByStatusTotal *prometheus.CounterVec
	latencyHist         *prometheus.HistogramVec
	lastEventTs         *prometheus.GaugeVec
}

func (rr ResourceMetricRepository) Create(ctx context.Context, r *metrics.Resource) {
	select {
	case <-ctx.Done():
		return
	default:
		rr.eventsTotal.WithLabelValues(r.Name()).Inc()
		rr.eventsByStatusTotal.WithLabelValues(r.Name(), r.Status()).Inc()
		rr.latencyHist.WithLabelValues(r.Name()).Observe(r.Seconds())
		rr.lastEventTs.WithLabelValues(r.Name()).Set(float64(r.CreatedAt().Unix()))
	}
}

func NewResourceMetricRepository() *ResourceMetricRepository {
	repo := ResourceMetricRepository{
		eventsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "pinger_events_total",
				Help: "Total number from Pinger events processed",
			},
			[]string{"resource"}),
		latencyHist: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "pinger_latency_seconds",
				Help:    "Observed latency by resource (in seconds)",
				Buckets: prometheus.ExponentialBuckets(0.01, 2, 150),
			},
			[]string{"resource"}),
		lastEventTs: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "pinger_last_event_timestamp_seconds",
				Help: "Epoch from last event received by resource",
			},
			[]string{"resource"}),
		eventsByStatusTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "pinger_events_total_by_status",
				Help: "Total status number from Pinger events processed",
			},
			[]string{"resource", "status"}),
	}
	prometheus.MustRegister(repo.eventsTotal, repo.latencyHist, repo.lastEventTs, repo.eventsByStatusTotal)

	return &repo
}
