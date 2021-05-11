package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	// "github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	// 	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	instanceIds []string

	diskDesc = prometheus.NewDesc(
		"rds_baseinf_disk",
		"Aliyun Rds disk space usage.",
		[]string{"endpoint", "dbinstanceid"}, nil,
	)
	ioDesc = prometheus.NewDesc(
		"rds_baseinf_io",
		"Aliyun Rds IOPS.",
		[]string{"endpoint", "dbinstanceid"}, nil,
	)
	cpuDesc = prometheus.NewDesc(
		"rds_baseinf_cpu",
		"Aliyun Rds cpu usage.",
		[]string{"endpoint", "dbinstanceid"}, nil,
	)
	memoryDesc = prometheus.NewDesc(
		"rds_baseinf_memory",
		"Aliyun Rds memory usage.",
		[]string{"endpoint", "dbinstanceid"}, nil,
	)
)

type AliyunRdsCollector struct {
}

func (arc AliyunRdsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- diskDesc
	ch <- ioDesc
	ch <- cpuDesc
	ch <- memoryDesc
}

func (arc AliyunRdsCollector) Collect(ch chan<- prometheus.Metric)  {
	fmt.Println(instanceIds)
	for _, ins := range instanceIds {
		ch <- prometheus.MustNewConstMetric(
			diskDesc, prometheus.GaugeValue, 0.0, "localhost", ins)
	}
}

func main() {
	flag.Parse()
	instanceIds = flag.Args()
	arc := AliyunRdsCollector{}
	prometheus.MustRegister(arc)
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8000", nil))
}
