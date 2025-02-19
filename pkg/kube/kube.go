package kube

import (
	"fmt"
	"log/slog"
	"os"

	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// GetKubeClient returns a kubernetes client. It try to load the kubeconfig file
// from the given path, if it's empty it will try to load the default config. If
// that doesn't exist it falls back to the in-cluster config.
func GetKubeClient(kubeconfigPath string) (*kubernetes.Clientset, error) {
	// Try to load the kubeconfig from the given path
	client, err := loadKubeConfig(kubeconfigPath)
	if err != nil {
		slog.Info("Failed to load kubeconfig from the given path", slog.Any("error", err))
		slog.Info("Trying to load the in-cluster config")
		client, err = restclient.InClusterConfig()
		if err != nil {
			err = fmt.Errorf("Failed to load the in-cluster config: %w", err)
			return nil, err
		}
	}
	return kubernetes.NewForConfig(client)
}

// Load the kubeconfig from a given path
func loadKubeConfig(path string) (*restclient.Config, error) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	loadingRules.ExplicitPath = path
	kubeconfigPath := loadingRules.GetDefaultFilename()
	// Strip the leading slash from the path
	slog.Info("Loading kubeconfig", slog.String("path", kubeconfigPath))

	// Check if the file exists
	if _, err := os.Stat(kubeconfigPath); err != nil {
		return nil, err
	}

	return clientcmd.BuildConfigFromFlags("", kubeconfigPath)
}
