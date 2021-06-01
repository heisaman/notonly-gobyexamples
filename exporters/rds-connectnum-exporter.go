package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	threadsConnected := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "rds_connectnum_count",
			Help:      "RDS threads connected count.",
		},
		[]string{"endpoint", "dbinstanceid", "department", "cluster"},
	)
	maxConnections := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "rds_connectnum_maxcount",
			Help:      "RDS max connections count.",
		},
		[]string{"endpoint", "dbinstanceid", "department", "cluster"},
	)
)

func init() {
	// Metrics have to be registered to be exposed:
	prometheus.MustRegister(threadsConnected)
	prometheus.MustRegister(maxConnections)
}

func main() {
	threadsConnected.WithLabelValues("bob", "put", "bob", "put").Set(4)
	maxConnections.WithLabelValues("bob", "put", "bob", "put").Set(4)

	// The Handler function provides a default handler to expose metrics
	// via an HTTP server. "/metrics" is the usual endpoint for that.
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}