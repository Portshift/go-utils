package monitor

import (
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Factory struct {
	KubeConfig string
	Context    string
}

func (f *Factory) toRawKubeConfigLoader() clientcmd.ClientConfig {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	loadingRules.DefaultClientConfig = &clientcmd.DefaultClientConfig
	loadingRules.ExplicitPath = f.KubeConfig
	configOverrides := &clientcmd.ConfigOverrides{
		ClusterDefaults: clientcmd.ClusterDefaults,
		CurrentContext:  f.Context,
	}

	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
}

type KubeOptions struct {
	// QPS indicates the maximum QPS to the master from this client.
	// If it's zero, the created RESTClient will use DefaultQPS: 5
	KubernetesAPIQPS float32

	// Maximum burst for throttle.
	// If it's zero, the created RESTClient will use DefaultBurst: 10.
	KubernetesAPIBurst int
}

func CreateK8sClientset(f *Factory, ops KubeOptions) (*kubernetes.Clientset, *rest.Config, error) {
	// Create Kubernetes go-client clientset
	var config *rest.Config
	var err error

	if f != nil {
		config, err = f.toRawKubeConfigLoader().ClientConfig()
	} else {
		config, err = rest.InClusterConfig()
	}
	if err != nil {
		return nil, nil, fmt.Errorf("failed to build config: %v", err)
	}

	config.QPS = ops.KubernetesAPIQPS
	config.Burst = ops.KubernetesAPIBurst

	// Create a rest client not targeting specific API version
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create a rest client: %v", err)
	}

	return clientset, config, nil
}

