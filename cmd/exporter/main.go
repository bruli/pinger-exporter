package main

import (
	"log"
	"net/http"
	"os"

	"github.com/bruli/pinger/pkg/events"
	"github.com/nats-io/nats.go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/protobuf/proto"
)

var (
	eventsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "pinger_events_total",
			Help: "Nombre total d'events pinger processats",
		},
		[]string{"resource"},
	)
	latencyHist = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "pinger_latency_seconds",
			Help:    "Latència observada pels recursos (en segons)",
			Buckets: prometheus.ExponentialBuckets(0.01, 2, 150),
		},
		[]string{"resource"},

	)
	lastEventTs = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "pinger_last_event_timestamp_seconds",
			Help: "Epoch del darrer event rebut per recurs",
		},
		[]string{"resource"},
	)
)

func parseSeconds(e *events.PingResult) (float64, bool) {
	return float64(e.GetLatency() / 1000), true
}

func mustEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func main() {
	prometheus.MustRegister(eventsTotal, latencyHist, lastEventTs)

	natsURL := mustEnv("NATS_URL", "nats://nats:4222")
	subject := mustEnv("NATS_SUBJECT", events.PingSubjet)
	httpAddr := mustEnv("LISTEN_ADDR", ":8080")

	// Connexió a NATS
	nc, err := nats.Connect(natsURL, nats.Name("pinger-prom-exporter"))
	if err != nil {
		log.Fatalf("Error connectant a NATS: %v", err)
	}
	defer nc.Drain()

	// Subscriure'ns als events
	_, err = nc.Subscribe(subject, func(msg *nats.Msg) {
		var e events.PingResult
		if err = proto.Unmarshal(msg.Data, &e); err != nil {
			// Si vols admetre format "resource=...,latency_ms=..." molt simple:
			// intenta un parse fallback aquí.
			log.Printf("PROTOBUFF invàlid: %v", err)
			return
		}
		if e.Resource == "" {
			log.Printf("Event sense resource, ignorat")
			return
		}
		if sec, ok := parseSeconds(&e); ok {
			eventsTotal.WithLabelValues(e.Resource).Inc()
			latencyHist.WithLabelValues(e.Resource).Observe(sec)
			lastEventTs.WithLabelValues(e.Resource).Set(float64(e.GetCreatedAt().AsTime().Unix()))

			eventsTotal.WithLabelValues(e.Status).Inc()
			latencyHist.WithLabelValues(e.Status).Observe(sec)
			lastEventTs.WithLabelValues(e.Status).Set(float64(e.GetCreatedAt().AsTime().Unix()))
		} else {
			log.Printf("Event sense latència: %+v", &e)
		}
	})
	if err != nil {
		log.Fatalf("Error subscrivint-se a %s: %v", subject, err)
	}

	// /metrics
	http.Handle("/metrics", promhttp.Handler())

	// /ready per readinessProbe (minim)
	http.HandleFunc("/ready", func(w http.ResponseWriter, _ *http.Request) {
		if nc.IsConnected() {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("ok"))
			return
		}
		w.WriteHeader(http.StatusServiceUnavailable)
	})

	log.Printf("Escoltant %s (raspatge Prometheus a /metrics). NATS=%s subject=%s",
		httpAddr, natsURL, subject)
	log.Fatal(http.ListenAndServe(httpAddr, nil))
}
