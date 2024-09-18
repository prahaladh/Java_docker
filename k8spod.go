package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// Parse flags for Kubernetes config
	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// Build the config from the kubeconfig file
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// Create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// Service name as a parameter
	if len(os.Args) < 2 {
		fmt.Println("Please provide a service name as a parameter.")
		os.Exit(1)
	}
	serviceName := os.Args[1]

	// Get the pods in the same namespace as the service
	namespace := "default" // Change to your namespace
	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: fmt.Sprintf("app=%s", serviceName),
	})
	if err != nil {
		panic(err.Error())
	}

	// Check the overall status
	overallStatus := "Running" // Assume all are running initially
	for _, pod := range pods.Items {
		switch pod.Status.Phase {
		case "Failed":
			overallStatus = "Failed"
			break // No need to continue if we found a failure
		case "Pending":
			overallStatus = "Pending"
		}
	}

	fmt.Printf("Overall Pod Status for Service '%s': %s\n", serviceName, overallStatus)
}

// Helper function to get the home directory
func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
