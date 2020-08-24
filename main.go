package main

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
	"time"
)

func main() {

	// Location of kubeconfig file
	kubeconfig := os.Getenv("HOME") + "/.kube/config"

	// Create a Config (k8s.io/client-go/rest)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// Create an API Clientset (k8s.io/client-go/kubernetes)
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// Create a CoreV1Client (k8s.io/client-go/kubernetes/typed/core/v1)
	coreV1Client := clientset.CoreV1()
	// Create an AppsV1Client (k8s.io/client-go/kubernetes/typed/apps/v1)
	appsV1Client := clientset.AppsV1()

	//-------------------------------------------------------------------------//
	// List pods (all namespaces)
	//-------------------------------------------------------------------------//

	// Get a *PodList (k8s.io/api/core/v1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	pods, err := coreV1Client.Pods("").List(ctx, metaV1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	// List each Pod (k8s.io/api/core/v1)
	for i, pod := range pods.Items {
		fmt.Printf("Pod %d: %s, Namespace: %s\n", i+1, pod.ObjectMeta.Name, pod.ObjectMeta.Namespace)
	}

	watch, err := clientset.CoreV1().Pods("").Watch(ctx, metaV1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	go func() {
		for event := range watch.ResultChan() {
			fmt.Printf("Type: %v\n", event.Type)
			p, ok := event.Object.(*v1.Pod)
			if !ok {
				log.Fatal("unexpected type")
			}
			// fmt.Printf("Container status:\n %v \n",p.Status.ContainerStatuses)
			fmt.Printf("\n\nPod Name: %s\n",p.ObjectMeta.Name)
			fmt.Printf("Pod Namespace: %s\n",p.ObjectMeta.Namespace)
			fmt.Printf("Status.Phase: %s\n",p.Status.Phase)

			if p.Status.Phase == corev1.PodRunning {
				fmt.Printf("HostIP: %s\n", p.Status.HostIP)
				fmt.Printf("PodIP: %s\n",p.Status.PodIP)
				fmt.Printf("StartTime: %v\n",p.Status.StartTime.Time)
			}
			fmt.Printf("Status Phase:\n%v\n",p.Status.Phase)
		}
	}()
	time.Sleep(15 * time.Second)

	//-------------------------------------------------------------------------//
	// List nodes
	//-------------------------------------------------------------------------//

	// Get a *NodeList (k8s.io/api/core/v1)
	nodes, err := coreV1Client.Nodes().List(ctx, metaV1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	// For each Node (k8s.io/api/core/v1)
	for i, node := range nodes.Items {
		fmt.Printf("Node %d: %s\n", i+1, node.ObjectMeta.Name)
	}

	//-------------------------------------------------------------------------//
	// List deployments (all namespaces)
	//-------------------------------------------------------------------------//

	// Get a *DeploymentList (k8s.io/api/apps/v1)
	deployments, err := appsV1Client.Deployments("").List(ctx, metaV1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	// For each Deployment (k8s.io/api/apps/v1)
	for i, deployment := range deployments.Items {
		fmt.Printf("Deployment %d: %s\n", i+1, deployment.ObjectMeta.Name)
	}

}
