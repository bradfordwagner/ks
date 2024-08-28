package cmds

import (
	"context"
	"errors"
	"fmt"
	"github.com/bradfordwagner/go-util/log"
	"github.com/bradfordwagner/ks/internal/args"
	"github.com/bradfordwagner/ks/internal/choose"
	"github.com/bradfordwagner/ks/internal/kube"
	"github.com/bradfordwagner/ks/internal/tmux"
	"github.com/koki-develop/go-fzf"
	"io"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
)

func KubeNewNamespace(a args.Standard) (err error) {
	client, err := kube.Client(a.Kubeconfig)
	if err != nil {
		return
	}

	// follow kubeconfig symlink to get actual kubeconfig
	sourceKubeconfig, err := ResolveIfSymlink(a.Kubeconfig)
	if err != nil {
		return err
	}
	a.Kubeconfig = sourceKubeconfig

	// get namespaces
	ctx, _ := context.WithTimeout(context.Background(), a.Timeout)
	namespaceList, err := client.CoreV1().Namespaces().List(ctx, v1.ListOptions{})
	var namespaces []string
	for _, namespace := range namespaceList.Items {
		namespaces = append(namespaces, namespace.Name)
	}

	// select a namespace
	namespace, err := choose.One(namespaces)
	if errors.Is(fzf.ErrAbort, err) {
		return nil
	} else if err != nil {
		return
	}

	// copy kubeconfig to kubeconfig-namespace
	src, _ := os.Open(a.Kubeconfig)
	defer src.Close()
	dstName := fmt.Sprintf("%s.%s", a.Kubeconfig, namespace)
	dst, _ := os.Create(dstName)
	defer dst.Close()
	_, err = io.Copy(dst, src)

	// log
	log.Log().With("kubeconfig", dstName).With("namespace", namespace).Info("created kubeconfig with namespace")

	// copy kubeconfig to tmux buffer
	tmux.LoadBufferKubeconfig(dstName)

	// set namespace
	return kube.SetNamespace(dstName, namespace)
}

// ResolveIfSymlink checks if the given file path is a symlink and resolves it if true
func ResolveIfSymlink(path string) (string, error) {
	fileInfo, err := os.Lstat(path)
	if err != nil {
		return "", fmt.Errorf("failed to get file info: %w", err)
	}

	if fileInfo.Mode()&os.ModeSymlink != 0 {
		resolvedPath, err := os.Readlink(path)
		if err != nil {
			return "", fmt.Errorf("failed to resolve symlink: %w", err)
		}
		return resolvedPath, nil
	}

	return path, nil
}
