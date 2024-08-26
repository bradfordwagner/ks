package cmds

import (
	"context"
	"errors"
	"github.com/bradfordwagner/go-util/log"
	"github.com/bradfordwagner/ks/internal/args"
	"github.com/bradfordwagner/ks/internal/choose"
	"github.com/bradfordwagner/ks/internal/kube"
	"github.com/koki-develop/go-fzf"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func SetNamespace(a args.Standard) (err error) {
	l := log.Log()
	// get the current context
	client, err := kube.Client(a.Kubeconfig)
	if err != nil {
		return err
	}

	// list namespaces
	ctx, _ := context.WithTimeout(context.Background(), a.Timeout)
	res, err := client.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}
	var namespaces []string
	for _, ns := range res.Items {
		namespaces = append(namespaces, ns.Name)
	}

	// choose a namespace
	selectedNamespace, err := choose.One(namespaces)
	if errors.Is(fzf.ErrAbort, err) {
		return nil
	} else if err != nil {
		return err
	}

	// set the namespace
	err = kube.SetNamespace(a.Kubeconfig, selectedNamespace)
	if err != nil {
		return err
	}
	l.With("namespace", selectedNamespace).With("kubeconfig", a.Kubeconfig).Info("namespace set")

	return
}
