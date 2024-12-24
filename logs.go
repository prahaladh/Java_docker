package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// Parse kubeconfig
	kubeconfig := flag.String("kubeconfig", os.Getenv("HOME")+"/.kube/config", "path to the kubeconfig file")
	namespace := flag.String("namespace", "default", "namespace of the pod")
	podName := flag.String("pod", "", "name of the pod")
	container := flag.String("container", "", "container name (optional)")
	pattern := flag.String("pattern", `\w+Exception(:.*)?`, "regex pattern to match logs")
	flag.Parse()

	if *podName == "" {
		log.Fatalf("Pod name is required. Use the -pod flag.")
	}

	// Load the kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		log.Fatalf("Failed to load kubeconfig: %v", err)
	}

	// Create Kubernetes client
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed to create Kubernetes client: %v", err)
	}

	// Create regex
	re, err := regexp.Compile(*pattern)
	if err != nil {
		log.Fatalf("Invalid regex pattern: %v", err)
	}

	// Get logs
	logs, err := getPodLogs(clientset, *namespace, *podName, *container)
	if err != nil {
		log.Fatalf("Failed to fetch logs: %v", err)
	}

	// Check logs against regex
	matches := re.FindAllString(logs, -1)
	if len(matches) > 0 {
		fmt.Println("Regex matches found:")
		for _, match := range matches {
			fmt.Println(match)
		}
	} else {
		fmt.Println("No matches found in the logs.")
	}
}

// getPodLogs fetches logs from the specified pod and container
func getPodLogs(clientset *kubernetes.Clientset, namespace, podName, container string) (string, error) {
	podLogOpts := v1.PodLogOptions{
		Follow:    false,
		Container: container,
	}

	req := clientset.CoreV1().Pods(namespace).GetLogs(podName, &podLogOpts)
	podLogs, err := req.Stream(context.TODO())
	if err != nil {
		return "", fmt.Errorf("failed to open log stream: %v", err)
	}
	defer podLogs.Close()

	var logs string
	buf := make([]byte, 2048)
	for {
		n, err := podLogs.Read(buf)
		if n > 0 {
			logs += string(buf[:n])
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("error reading log stream: %v", err)
		}
	}

	return logs, nil
}
