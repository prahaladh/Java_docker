package main

import (
    "context"
    "flag"
    "fmt"
    "log"
    "os"
    "path/filepath"

    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/clientcmd"
    "k8s.io/client-go/util/homedir"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    corev1 "k8s.io/api/core/v1"
)

// GetKubeClient creates a Kubernetes clientset
func GetKubeClient() (*kubernetes.Clientset, error) {
    var kubeconfig string
    if home := homedir.HomeDir(); home != "" {
        kubeconfig = filepath.Join(home, ".kube", "config")
    } else {
        kubeconfig = os.Getenv("KUBECONFIG")
    }

    // Build Kubernetes client configuration
    config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
    if err != nil {
        return nil, err
    }

    // Create a new clientset for interacting with the Kubernetes cluster
    return kubernetes.NewForConfig(config)
}

// GetPodsByService retrieves the pods associated with a given service
func GetPodsByService(clientset *kubernetes.Clientset, namespace, serviceName string) ([]corev1.Pod, error) {
    // Retrieve the service by name
    service, err := clientset.CoreV1().Services(namespace).Get(context.TODO(), serviceName, metav1.GetOptions{})
    if err != nil {
        return nil, fmt.Errorf("failed to get service: %v", err)
    }

    // Get the label selector from the service's spec
    selector := service.Spec.Selector
    if len(selector) == 0 {
        return nil, fmt.Errorf("service has no selectors")
    }

    // Convert the label selector into a string that can be used to query pods
    labelSelector := metav1.FormatLabelSelector(&metav1.LabelSelector{MatchLabels: selector})

    // List the pods in the same namespace that match the service's label selector
    pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
        LabelSelector: labelSelector,
    })
    if err != nil {
        return nil, fmt.Errorf("failed to list pods: %v", err)
    }

    return pods.Items, nil
}

// GetDeployedImages extracts the container image details for each pod
func GetDeployedImages(pods []corev1.Pod) {
    for _, pod := range pods {
        fmt.Printf("Pod Name: %s\n", pod.Name)
        for _, container := range pod.Spec.Containers {
            fmt.Printf("Container Name: %s, Image: %s\n", container.Name, container.Image)
        }
    }
}

func main() {
    // Command-line arguments for namespace and service name
    namespace := flag.String("namespace", "default", "Kubernetes namespace")
    serviceName := flag.String("service", "", "Kubernetes service name")
    flag.Parse()

    if *serviceName == "" {
        log.Fatal("Service name must be provided")
    }

    // Initialize Kubernetes client
    clientset, err := GetKubeClient()
    if err != nil {
        log.Fatalf("Failed to create Kubernetes client: %v", err)
    }

    // Get pods associated with the service
    pods, err := GetPodsByService(clientset, *namespace, *serviceName)
    if err != nil {
        log.Fatalf("Error retrieving pods for service %v: %v", *serviceName, err)
    }

    if len(pods) == 0 {
        log.Fatalf("No pods found for service %s", *serviceName)
    }

    // Get deployed container images from the pods
    GetDeployedImages(pods)
}
