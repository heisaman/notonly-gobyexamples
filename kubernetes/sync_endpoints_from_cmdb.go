package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	kubeconfig   *string
	logfmtLogger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
	logger       = log.With(logfmtLogger, "operation", "syncCmdbEcsEndpoints")
)

func getEcsAddresses(ctx context.Context) ([]v1.EndpointAddress, []error) {
	addresses := make([]v1.EndpointAddress, 0)
	errs := make([]error, 0)

	cmdbHost := os.Getenv("DB_HOST")
	cmdbPort := os.Getenv("DB_PORT")
	user := os.Getenv("USER")
	password := os.Getenv("PASSWORD")
	cmdbName := os.Getenv("DB_NAME")
	db, err := sql.Open("mysql", fmt.Sprint(user, ":", password, "@tcp(", cmdbHost, ":", cmdbPort, ")/", cmdbName))
	if err != nil {
		errs = append(errs, errors.Wrapf(err, "failed to create mysql database connection to DB (%s)", cmdbName))
		return addresses, errs
	}
	defer db.Close()

	osType := os.Getenv("OS_TYPE")
	environment := os.Getenv("ENVIRONMENT")
	networkType := os.Getenv("NETWORK_TYPE")
	results, err := db.QueryContext(ctx, "SELECT public_ip_address FROM create_ecs_instance where os_type=? and environment=? and instance_network_type=?", osType, environment, networkType)
	if err != nil {
		errs = append(errs, errors.Wrap(err, "failed to query create_ecs_instance table"))
		return addresses, errs
	}
	defer results.Close()

	for results.Next() {
		var ip string
		err = results.Scan(&ip)
		if err != nil {
			errs = append(errs, errors.Wrap(err, "failed to read ip column"))
			continue
		}
		addresses = append(addresses, v1.EndpointAddress{
			IP: ip,
		})
	}

	return addresses, errs
}

func syncCmdbEcsEndpoints(ctx context.Context) error {
	cfg, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		return fmt.Errorf("error creating config from %s: %w", *kubeconfig, err)
	}
	client, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return fmt.Errorf("The client cannot be created: %w", err)
	}

	eps := &v1.Endpoints{
		ObjectMeta: metav1.ObjectMeta{
			Name:   "windows-exporter",
			Labels: map[string]string{"app": "windows-exporter"},
		},
		Subsets: []v1.EndpointSubset{
			{
				Ports: []v1.EndpointPort{
					{
						Name: "windows-metrics",
						Port: 9182,
					},
				},
			},
		},
	}

	addresses, errs := getEcsAddresses(ctx)
	if len(errs) > 0 {
		for _, err := range errs {
			level.Warn(logger).Log("err", err)
		}
	}
	level.Debug(logger).Log("msg", "CMDB servers converted to endpoint addresses", "num_addresses", len(addresses))

	eps.Subsets[0].Addresses = addresses

	svc := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:   "windows-exporter",
			Labels: map[string]string{"app": "windows-exporter"},
		},
		Spec: v1.ServiceSpec{
			Type:      v1.ServiceTypeClusterIP,
			ClusterIP: "None",
			Ports: []v1.ServicePort{
				{
					Name: "windows-metrics",
					Port: 9182,
				},
			},
		},
	}

	level.Debug(logger).Log("msg", "Updating Kubernetes service", "service", "windows-exporter", "ns", "monitoring")
	err = CreateOrUpdateService(ctx, client.CoreV1().Services("monitoring"), svc)
	if err != nil {
		return errors.Wrap(err, "synchronizing windows-exporter service object failed")
	}

	level.Debug(logger).Log("msg", "Updating Kubernetes endpoint", "endpoint", "windows-exporter", "ns", "monitoring")
	err = CreateOrUpdateEndpoints(ctx, client.CoreV1().Endpoints("monitoring"), eps)
	if err != nil {
		return errors.Wrap(err, "synchronizing windows-exporter endpoints object failed")
	}

	return nil
}

func main() {
	kubeconfig = flag.String("kubeconfig", "config_east1_test", "kubeconfig file")
	flag.Parse()

	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)
	ctx := context.Background()
	err := syncCmdbEcsEndpoints(ctx)
	if err != nil {
		fmt.Printf("error sync cmdb ecs endpoints: %w\n", err)
	}
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			err = syncCmdbEcsEndpoints(ctx)
			if err != nil {
				fmt.Printf("error sync cmdb ecs endpoints: %w\n", err)
			}
		}
	}
}
