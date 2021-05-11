package main

import (
	"context"
	"flag"
	"fmt"
	"os"

    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
    "k8s.io/client-go/tools/clientcmd"
    "k8s.io/client-go/dynamic"
)


func main() {
	kubeconfig := flag.String("kubeconfig", "config_east1_test", "kubeconfig file")
	flag.Parse()
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
        fmt.Printf("The kubeconfig cannot be loaded: %v\n", err)
        os.Exit(1)
    }
	client, err := dynamic.NewForConfig(config)
	if err != nil {
        fmt.Printf("The client cannot be created: %v\n", err)
        os.Exit(1)
    }

	gvr := schema.GroupVersionResource{
		Group: "monitoring.coreos.com",
		Version: "v1",
		Resource: "prometheusrules",
	}

	rule, err := client.Resource(gvr).Namespace("monitoring").Get(context.TODO(), "prometheus-prometheus-oper-general.rules", metav1.GetOptions{})

	if err != nil {
        fmt.Printf("The prometheusrule cannot be got: %v\n", err)
        os.Exit(1)
    }
	// &{map[apiVersion:monitoring.coreos.com/v1 kind:PrometheusRule metadata:map[creationTimestamp:2020-07-20T11:06:19Z generation:1 labels:map[app:prometheus-operator chart:prometheus-operator-8.13.2 heritage:Tiller release:prometheus] name:prometheus-prometheus-oper-general.rules namespace:monitoring resourceVersion:306713194 selfLink:/apis/monitoring.coreos.com/v1/namespaces/monitoring/prometheusrules/prometheus-prometheus-oper-general.rules uid:d7dd1f6b-b0c2-469b-9585-cd6d9d5b382a] spec:map[groups:[map[name:general.rules rules:[map[alert:TargetDown annotations:map[message:{{ printf "%.4g" $value }}% of the {{ $labels.job }}/{{ $labels.service }} targets in {{ $labels.namespace }} namespace are down.] expr:100 * (count(up == 0) BY (job, namespace, service) / count(up) BY (job, namespace, service)) > 10 for:10m labels:map[severity:warning]] map[alert:Watchdog annotations:map[message:This is an alert meant to ensure that the entire alerting pipeline is functional.
    // This alert is always firing, therefore it should always be firing in Alertmanager and always fire against a receiver. There are integrations with various notification mechanisms that send a notification when this alert is not firing. For example the "DeadMansSnitch" integration in PagerDuty. ] expr:vector(1) labels:map[severity:none]]]]]]]}
	// fmt.Printf("%#v\n", rule)

	json, err := rule.MarshalJSON()
	fmt.Println(string(json))

	name, found, err := unstructured.NestedString(rule.Object, "metadata", "name")
	fmt.Println(name, found)

	name, found, err := unstructured.NestedFieldNoCopy(rule.Object, "metadata", "name")
	

}
