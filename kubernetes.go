package main

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func kubernetesClient(kubeconfig string) (*kubernetes.Clientset, error) {
	var clientset *kubernetes.Clientset
	var config *rest.Config
	var err error

	if kubeconfig == "" {
		fmt.Println("kubernetesClient: using in-cluster authentication")

		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
	} else {
		fmt.Printf("kubernetesClient: using kubeconfig %s with current context\n", kubeconfig)

		// BuildConfigFromFlags will fallback to InClusterConfig as well so we could spare us the if block,
		// but we don't want the potential warning messages mentioning different application flags
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, err
		}
	}

	clientset, err = kubernetes.NewForConfig(config)
	return clientset, err
}

func listExternalIPs(client *kubernetes.Clientset) ([]string, error) {
	ipList := make([]string, 0)

	nodes, err := client.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, node := range nodes.Items {
		for _, address := range node.Status.Addresses {
			if address.Type == v1.NodeExternalIP {
				ipList = append(ipList, address.Address)
			}
		}
	}

	return ipList, nil
}
