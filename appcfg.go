package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// Define the namespace, pod name, and container name to check
	namespace := "default"
	podName := "example-pod"
	containerName := "example-container"

	// Parse kubeconfig file path
	kubeconfig := flag.String("kubeconfig", os.Getenv("HOME")+"/.kube/config", "Path to kubeconfig file")
	flag.Parse()

	// Build Kubernetes client configuration
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		log.Fatalf("Error building kubeconfig: %v", err)
	}

	// Create Kubernetes client
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Error creating Kubernetes client: %v", err)
	}

	// Get the specified pod
	pod, err := clientset.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		log.Fatalf("Error fetching pod: %v", err)
	}

	// Check container status
	found := false
	for _, containerStatus := range pod.Status.ContainerStatuses {
		if containerStatus.Name == containerName {
			found = true
			if containerStatus.State.Running != nil {
				fmt.Printf("Container '%s' in pod '%s' is running.\n", containerName, podName)
			} else if containerStatus.State.Waiting != nil {
				fmt.Printf("Container '%s' in pod '%s' is in waiting state: %s\n", containerName, podName, containerStatus.State.Waiting.Reason)
			} else if containerStatus.State.Terminated != nil {
				fmt.Printf("Container '%s' in pod '%s' is terminated: %s\n", containerName, podName, containerStatus.State.Terminated.Reason)
			} else {
				fmt.Printf("Container '%s' in pod '%s' is in an unknown state.\n", containerName, podName)
			}
			break
		}
	}

	if !found {
		fmt.Printf("Container '%s' not found in pod '%s'.\n", containerName, podName)
	}
}
