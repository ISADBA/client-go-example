package main

import (
	"context"
	"log"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

func main() {
	conf, err := config.GetConfig()
	if err != nil {
		log.Fatal("Get config file Fail", err)
	}
	clientSet, err := kubernetes.NewForConfig(conf)
	if err != nil {
		log.Fatal("Get ClientSet Fail", err)
	}

	for {
		pods, err := clientSet.CoreV1().Pods("default").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("There are %d pods in the cluster\n", len(pods.Items))
		for i, pod := range pods.Items {
			log.Printf("%d -> %s/%s", i+1, pod.Namespace, pod.Name)
		}
		<-time.Tick(5 * time.Second)
	}
}
