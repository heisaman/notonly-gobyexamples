package main

import (
	"context"
	"flag"
	"fmt"
	"os"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/tools/clientcmd"
    "k8s.io/client-go/kubernetes"
)


func main() {
	kubeconfig := flag.String("kubeconfig", "config_east1_test", "kubeconfig file")
	flag.Parse()
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
        fmt.Printf("The kubeconfig cannot be loaded: %v\n", err)
        os.Exit(1)
    }
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
        fmt.Printf("The clientset cannot be created: %v\n", err)
        os.Exit(1)
    }

	pod, err := clientset.CoreV1().ConfigMaps("monitoring").Get(context.TODO(), "blackbox-exporter-configuration", metav1.GetOptions{})
	if err != nil {
        fmt.Printf("The cm cannot be got: %v\n", err)
        os.Exit(1)
    }
	fmt.Println(pod)
}
