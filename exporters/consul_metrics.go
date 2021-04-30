package main

import (
	"log"
	"net/http"
	"regexp"

	"github.com/hashicorp/consul/api"
	"github.com/prometheus/client_go/prometheus"
	"github.com/prometheus/client_go/prometheus/promhttp"
)

var (
	up = prometheus.NewDesc(
		"consul_up",
		"Was talking to Consul successful.",
		nil, nil,
	)
	invalidChars = regexp.MustCompile("[^a-zA-Z0-9:_]")
)

type ConsulCollector struct {
}

func (c ConsulCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- up	
}

func (c ConsulCollector) Collect(ch chan<- prometheus.Metric)  {
	
}

func main() {
	c := ConsulCollector{}
	prometheus.MustRegister(c)
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8000", nil))
}
