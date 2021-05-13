package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	instanceIds = strings.Split(os.Getenv("DB_INSTANCE_IDS"), ",")

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
	fmt.Println("Scraping metrics data...")

	client, err := rds.NewClientWithAccessKey(os.Getenv("REGION_ID"), os.Getenv("ACCESS_KEY_ID"), os.Getenv("ACCESS_KEY_SECRET"))
	if err != nil {
        panic(err)
    }
	now := time.Now()
	loc, _ := time.LoadLocation("UTC")
	endTime := now.In(loc)
	startTime := endTime.Add(time.Duration(-3) * time.Minute)
	for _, insId := range instanceIds {
		// Fetch db instance attributes
		fmt.Println(insId)
		attrReq := rds.CreateDescribeDBInstanceAttributeRequest()
		attrReq.DBInstanceId = insId
		response, err := client.DescribeDBInstanceAttribute(attrReq)
		if err != nil {
			panic(err)
		}
		dbInstanceAttr := response.Items.DBInstanceAttribute[0]
		dbInstanceStorage := dbInstanceAttr.DBInstanceStorage
		// dbInstanceMemory := dbInstanceAttr.DBInstanceMemory
		dbInstanceDescription := dbInstanceAttr.DBInstanceDescription

		// Fetch db instance performance data
		perfReq := rds.CreateDescribeDBInstancePerformanceRequest()
		perfReq.DBInstanceId = insId
		perfReq.StartTime = startTime.Format("2006-01-02T15:04Z")
		perfReq.EndTime = endTime.Format("2006-01-02T15:04Z")
		perfReq.Key = "MySQL_MemCpuUsage,MySQL_SpaceUsage,MySQL_IOPS"
		resp, err := client.DescribeDBInstancePerformance(perfReq)
		if err != nil {
			panic(err)
		}
		keys := resp.PerformanceKeys.PerformanceKey
		v_disk := 0.0
		v_io := 0.0
		v_cpu := 0.0
		v_mem := 0.0
		for _, key := range keys {
			fmt.Println(key)
			keyName := key.Key
			valueFormat := key.ValueFormat
			values := key.Values.PerformanceValue
			lastValue := values[len(values)-1]
			value := lastValue.Value
			if keyName == "MySQL_SpaceUsage" {
				intValue, err := strconv.Atoi(value)
				if err != nil {
					panic(err)
				}
				diskUse := intValue * 100 / 1024 / 1024 / 1024
				preValue := strconv.Itoa(diskUse/ dbInstanceStorage)
				v_disk, err = strconv.ParseFloat(preValue, 64)
			} else if keyName == "MySQL_MemCpuUsage" {
				vs := strings.Split(value, "&")
				vfs := strings.Split(valueFormat, "&")
				for i, v := range vs {
					if vfs[i] == "cpuusage" {
						v_cpu, err = strconv.ParseFloat(v, 64)
					} else {
						v_mem, err = strconv.ParseFloat(v, 64)
					}
				}
			} else if keyName == "MySQL_IOPS" {
				v_io, err = strconv.ParseFloat(value, 64)
			} else {
				fmt.Println("Unexpected key: " + keyName)
			}
			if err != nil {
				panic(err)
			}
		}

		ch <- prometheus.MustNewConstMetric(diskDesc, prometheus.GaugeValue, v_disk, dbInstanceDescription, insId)
		ch <- prometheus.MustNewConstMetric(cpuDesc, prometheus.GaugeValue, v_cpu, dbInstanceDescription, insId)
		ch <- prometheus.MustNewConstMetric(memoryDesc, prometheus.GaugeValue, v_mem, dbInstanceDescription, insId)
		ch <- prometheus.MustNewConstMetric(ioDesc, prometheus.GaugeValue, v_io, dbInstanceDescription, insId)
	}
}

func main() {
	arc := AliyunRdsCollector{}
	prometheus.MustRegister(arc)
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8000", nil))
}
