package http

import (
	"errors"
	"net/http"

	"github.com/nats-io/nats.go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
)

const PortAddr = ":8080"

func Run(nc *nats.Conn, log *zerolog.Logger) {
	http.Handle("/metrics", promhttp.Handler())

	http.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		if nc.IsConnected() {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("ok"))
			return
		}
		w.WriteHeader(http.StatusServiceUnavailable)
	})

	if err := http.ListenAndServe(PortAddr, nil); errors.Is(err, http.ErrServerClosed) {
		log.Err(err).Msg("server closed")
	}
}
