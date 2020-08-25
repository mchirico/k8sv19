package main

import (
	"fmt"
	"log"
	"os"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/runtime"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	// SOME_LABEL is the key name to retrieve
	SOME_LABEL = "mytag"
)

func main() {
	log.Print("Shared Informer app started")
	// Location of kubeconfig file
	kubeconfig := os.Getenv("HOME") + "/.kube/config"

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Panic(err.Error())
	}

	factory := informers.NewSharedInformerFactory(clientset, 0)
	informer := factory.Core().V1().Pods().Informer()
	stopper := make(chan struct{})
	defer close(stopper)
	defer runtime.HandleCrash()
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    onAdd,
		DeleteFunc: onDelete,
		UpdateFunc: onUpdate,
	})
	go informer.Run(stopper)
	if !cache.WaitForCacheSync(stopper, informer.HasSynced) {
		runtime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
		return
	}
	<-stopper
}

// onCheckStatus is the function executed when the kubernetes informer notified the
// presence of a new kubernetes pod in the cluster
func onAdd(obj interface{}) {
	// Cast the obj as node
	pod := obj.(*corev1.Pod)
	label, ok := pod.GetLabels()[SOME_LABEL]
	if ok {
		fmt.Printf("onAdd\n")
		fmt.Printf("It has the label: %s\n", label)
	}
}

func onDelete(obj interface{}) {
	// Cast the obj as node
	pod := obj.(*corev1.Pod)
	label, ok := pod.GetLabels()[SOME_LABEL]
	if ok {
		fmt.Printf("onDelete\n")
		fmt.Printf("It has the label: %s\n", label)
	}
}

func onUpdate(oldObj, newObj interface{}) {
	// Cast the obj as node
	pod := oldObj.(*corev1.Pod)
	label, ok := pod.GetLabels()[SOME_LABEL]
	if ok {
		fmt.Printf("onUpdate\n")
		fmt.Printf("old label: %s\n", label)
	}

	pod = newObj.(*corev1.Pod)
	label, ok = pod.GetLabels()[SOME_LABEL]
	if ok {
		fmt.Printf("onUpdate\n")
		fmt.Printf("new label: %s\n", label)
	}
}
