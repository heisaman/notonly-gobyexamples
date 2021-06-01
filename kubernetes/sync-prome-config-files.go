package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

var (
	logfmtLogger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
	logger       = log.With(logfmtLogger, "operation", "syncCmdbEcsToFiles")
)

type EcsTarget struct {
	Targets []string          `json:"targets,omitempty"`
	Labels  map[string]string `json:"labels,omitempty"`
}

func getEcsAddresses(ctx context.Context) ([]EcsTarget, error) {

	cmdbHost := ""
	cmdbPort := ""
	user := ""
	password := ""
	cmdbName := ""
	db, err := sql.Open("mysql", fmt.Sprint(user, ":", password, "@tcp(", cmdbHost, ":", cmdbPort, ")/", cmdbName))
	if err != nil {
		level.Error(logger).Log("err", err)
		return nil, errors.Wrapf(err, "failed to create mysql database connection to DB (%s)", cmdbName)
	}
	defer db.Close()

	results, err := db.QueryContext(ctx, `SELECT COALESCE(public_ip_address, '') as ip, COALESCE(applications, '') as applications, COALESCE(cluster, '') as cluster,
						COALESCE(services, '') as services, os_type FROM create_ecs_instance where instance_network_type='classic' and 
						instance_name NOT LIKE '%k8s%' and instance_name != '医美生产'
						union 
						SELECT COALESCE(private_ip_address, '') as ip, COALESCE(applications, '') as applications, COALESCE(cluster, '') as cluster,
						COALESCE(services, '') as services, os_type FROM create_ecs_instance where instance_network_type='vpc' and 
						instance_name NOT LIKE '%k8s%' and instance_name != '医美生产'`)
	if err != nil {
		level.Error(logger).Log("err", err)
		return nil, errors.Wrap(err, "failed to create mysql database connection to DB")
	}
	defer results.Close()

	targets := []EcsTarget{}
	for results.Next() {
		var ip, applications, cluster, services, osType string
		err = results.Scan(&ip, &applications, &cluster, &services, &osType)
		if err != nil {
			level.Error(logger).Log("err", err)
			continue
		}
		if osType == "linux" {
			ip = ip + ":9100"
		} else {
			ip = ip + ":9182"
		}
		targets = append(targets, EcsTarget{
			Targets: []string{ip},
			Labels:  map[string]string{"application": applications, "cluster": cluster, "service": services},
		})
	}

	return targets, nil
}

func syncCmdbEcsToFiles(ctx context.Context) error {

	targets, err := getEcsAddresses(ctx)
	if err != nil {
		level.Error(logger).Log("err", err)
	}
	level.Info(logger).Log("msg", "Successful fetched all ecs targets", "num_targets", len(targets))

	t, err := json.Marshal(targets)
	if err != nil {
		panic(err)
	}

	// create or update prometheus config file
	f, err := os.OpenFile("/opt/prometheus/ecs_targets.json", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = f.WriteString(string(t))
	if err != nil {
		panic(err)
	}
	f.Sync()

	return nil
}

func main() {

	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)
	ctx := context.Background()
	err := syncCmdbEcsToFiles(ctx)
	if err != nil {
		fmt.Printf("error sync cmdb ecs to files: %w\n", err)
	}
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			err = syncCmdbEcsToFiles(ctx)
			if err != nil {
				fmt.Printf("error sync cmdb ecs to files: %w\n", err)
			}
		}
	}
}
