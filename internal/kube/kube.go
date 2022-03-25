package kube

import (
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
)

type ClientConfig struct {
	QPS        float32
	Burst      int
	KubeCtx    string
	KubeConfig string
}

// New returns a kubernetes client.
// It tries first with in-cluster config, if it fails it will try with out-of-cluster config.
func New(cc *ClientConfig) (client kubernetes.Interface, metricsClient *metrics.Clientset, err error) {
	client, metricsClient, err = NewInCluster(cc)
	if err == nil {
		return
	}

	client, metricsClient, err = NewFromConfig(cc)
	if err != nil {
		return
	}

	return
}

// NewFromConfig returns a new out-of-cluster kubernetes client.
func NewFromConfig(cc *ClientConfig) (client kubernetes.Interface, metricsClient *metrics.Clientset, err error) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	if cc.KubeConfig == "" {
		cc.KubeConfig = os.Getenv("KUBECONFIG")
		if cc.KubeConfig == "" {
			dir, err := os.UserHomeDir()
			if err != nil {
				return nil, nil, err
			}
			cc.KubeConfig = filepath.Join(dir, ".kube", "config")
		}
	}

	loadingRules.ExplicitPath = cc.KubeConfig

	var config *rest.Config

	if cc.KubeCtx != "" {
		config, err = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			loadingRules,
			&clientcmd.ConfigOverrides{CurrentContext: cc.KubeCtx},
		).ClientConfig()
	} else {
		// use the current context in kubeconfig
		config, err = clientcmd.BuildConfigFromFlags("", cc.KubeConfig)
	}

	if err != nil {
		return
	}

	cc.apply(config)

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return
	}

	client = clientset

	metricsClient, err = metrics.NewForConfig(config)
	if err != nil {
		return
	}

	return client, metricsClient, nil
}

// NewInCluster returns a new in-cluster kubernetes client.
func NewInCluster(cc *ClientConfig) (client kubernetes.Interface, metricsClient *metrics.Clientset, err error) {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		return
	}

	cc.apply(config)

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return
	}

	client = clientset

	metricsClient, err = metrics.NewForConfig(config)
	if err != nil {
		return
	}

	return
}

func (cc *ClientConfig) apply(config *rest.Config) {
	if cc.QPS > 0.0 {
		config.QPS = cc.QPS // the default is rest.DefaultQPS which is 5.0
	}

	if cc.Burst > 0 {
		config.Burst = cc.Burst // the default is rest.DefaultBurst which is 10
	}
}
