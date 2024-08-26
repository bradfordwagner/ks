package kube

import (
	"github.com/bradfordwagner/go-util/log"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func Client(kubeconfig string) (clientset kubernetes.Interface, err error) {
	l := log.Log()
	config, err := config(kubeconfig)
	if err != nil {
		l.With("error", err).Error("failed to create kubernetes config")
		return
	}

	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		l.With("error", err).Error("failed to create kubernetes client")
	}

	return
}

func Dynamic(kubeconfig string) (d dynamic.Interface, err error) {
	l := log.Log()

	config, err := config(kubeconfig)
	if err != nil {
		l.With("error", err).Error("failed to create kubernetes config")
		return
	}

	d, err = dynamic.NewForConfig(config)
	if err != nil {
		l.With("error", err).Error("failed to create kubernetes dynamic client")
	}
	return
}

func config(kubeconfig string) (config *rest.Config, err error) {
	// in cluster
	config, err = rest.InClusterConfig()
	if err != nil {
		// kubeconfig / file based
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	}

	return
}

func SetNamespace(kubeconfig, namespace string) (err error) {
	load, err := clientcmd.LoadFromFile(kubeconfig)
	for _, context := range load.Contexts {
		context.Namespace = namespace
	}
	return clientcmd.WriteToFile(*load, kubeconfig)
}
