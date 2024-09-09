package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type LogRequest struct {
	Namespace     string `json:"namespace"`
	PodName       string `json:"podName"`
	ContainerName string `json:"containerName"`
	Follow        bool   `json:"follow"`
	Pattern       string `json:"pattern"` // Regex pattern to search for
}

type LogResponse struct {
	Logs string `json:"logs"`
}

func main() {
	http.HandleFunc("/logs", getPodLogsHandler)

	port := "8080"
	fmt.Printf("Server listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func getPodLogsHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the JSON request body
	var logReq LogRequest
	if err := json.NewDecoder(r.Body).Decode(&logReq); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Fetch logs with regex matching
	logs, err := getPodLogs(logReq.Namespace, logReq.PodName, logReq.ContainerName, logReq.Follow, logReq.Pattern)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching logs: %v", err), http.StatusInternalServerError)
		return
	}

	// Respond with logs
	response := LogResponse{Logs: logs}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getPodLogs(namespace, podName, containerName string, follow bool, pattern string) (string, error) {
	// Load the kubeconfig
	kubeconfig := os.Getenv("HOME") + "/.kube/config"
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return "", fmt.Errorf("error building kubeconfig: %v", err)
	}

	// Create Kubernetes client
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return "", fmt.Errorf("error creating kubernetes clientset: %v", err)
	}

	// Compile regex pattern
	var regex *regexp.Regexp
	if pattern != "" {
		regex, err = regexp.Compile(pattern)
		if err != nil {
			return "", fmt.Errorf("invalid regex pattern: %v", err)
		}
	}

	// Request pod logs
	req := clientset.CoreV1().Pods(namespace).GetLogs(podName, &v1.PodLogOptions{
		Container: containerName,
		Follow:    follow,
	})

	// Stream logs
	logStream, err := req.Stream(context.Background())
	if err != nil {
		return "", fmt.Errorf("error opening log stream: %v", err)
	}
	defer logStream.Close()

	// Read logs and filter by regex
	var logs string
	reader := bufio.NewReader(logStream)
	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return "", fmt.Errorf("error reading log stream: %v", err)
		}
		if err == io.EOF {
			break
		}

		// If a regex pattern is provided, only append matching lines
		if regex == nil || regex.MatchString(line) {
			logs += line
		}
	}

	return logs, nil
}
