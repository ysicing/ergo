package kube

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/ergoapi/util/exmap"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	storagev1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
)

type Client struct {
	Clientset        kubernetes.Interface
	DynamicClientset dynamic.Interface
	MetricsClientset *metrics.Clientset
	Config           *rest.Config
	RawConfig        clientcmdapi.Config
	restClientGetter genericclioptions.RESTClientGetter
	contextName      string
}

func NewClient(contextName, kubeconfig string) (*Client, error) {
	if kubeconfig == "" {
		kubeconfig = os.Getenv("KUBECONFIG")
		if kubeconfig == "" {
			dir, err := os.UserHomeDir()
			if err != nil {
				return nil, err
			}
			kubeconfig = filepath.Join(dir, ".kube", "config")
		}
	}
	restClientGetter := genericclioptions.ConfigFlags{
		Context:    &contextName,
		KubeConfig: &kubeconfig,
	}
	rawKubeConfigLoader := restClientGetter.ToRawKubeConfigLoader()

	config, err := rawKubeConfigLoader.ClientConfig()
	if err != nil {
		return nil, err
	}

	rawConfig, err := rawKubeConfigLoader.RawConfig()
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	dynamicClientset, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	metricsClient, err := metrics.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	if contextName == "" {
		contextName = rawConfig.CurrentContext
	}

	return &Client{
		Clientset:        clientset,
		Config:           config,
		DynamicClientset: dynamicClientset,
		MetricsClientset: metricsClient,
		RawConfig:        rawConfig,
		restClientGetter: &restClientGetter,
		contextName:      contextName,
	}, nil
}

func (c *Client) ContextName() (name string) {
	return c.contextName
}

func (c *Client) GetVersion(_ context.Context) (string, error) {
	v, err := c.Clientset.Discovery().ServerVersion()
	if err != nil {
		return "", fmt.Errorf("failed to get Kubernetes version: %w", err)
	}
	return fmt.Sprintf("%#v", *v), nil
}

func (c *Client) GetGitVersion(_ context.Context) (string, error) {
	v, err := c.Clientset.Discovery().ServerVersion()
	if err != nil {
		return "", fmt.Errorf("failed to get Kubernetes version: %w", err)
	}
	return v.GitVersion, nil
}

func (c *Client) ListSC(ctx context.Context, opts metav1.ListOptions) (*storagev1.StorageClassList, error) {
	return c.Clientset.StorageV1().StorageClasses().List(ctx, opts)
}

func (c *Client) GetDefaultSC(ctx context.Context) (*storagev1.StorageClass, error) {
	scs, err := c.ListSC(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, sc := range scs.Items {
		defaultclass := exmap.GetLabelValue(sc.Annotations, "storageclass.kubernetes.io/is-default-class")
		if defaultclass == "true" {
			return &sc, nil
		}
	}
	return nil, fmt.Errorf("no default storage class found")
}

func (c *Client) PatchDefaultSC(ctx context.Context) (*storagev1.StorageClass, error) {
	sc, err := c.GetDefaultSC(ctx)
	if err == nil {
		return sc, nil
	}
	scs, err := c.ListSC(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	scnew := scs.Items[0]
	scan := scnew.Annotations
	scan = exmap.MergeLabels(scan, map[string]string{
		"storageclass.kubernetes.io/is-default-class": "true",
	})
	scnew.Annotations = scan
	return c.Clientset.StorageV1().StorageClasses().Update(ctx, &scnew, metav1.UpdateOptions{})
}

func (c *Client) ListNodes(ctx context.Context, opts metav1.ListOptions) (*corev1.NodeList, error) {
	return c.Clientset.CoreV1().Nodes().List(ctx, opts)
}

func (c *Client) CreateSecret(ctx context.Context, namespace string, secret *corev1.Secret, opts metav1.CreateOptions) (*corev1.Secret, error) {
	return c.Clientset.CoreV1().Secrets(namespace).Create(ctx, secret, opts)
}

func (c *Client) UpdateSecret(ctx context.Context, namespace string, secret *corev1.Secret, opts metav1.UpdateOptions) (*corev1.Secret, error) {
	return c.Clientset.CoreV1().Secrets(namespace).Update(ctx, secret, opts)
}

func (c *Client) PatchSecret(ctx context.Context, namespace, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions) (*corev1.Secret, error) {
	return c.Clientset.CoreV1().Secrets(namespace).Patch(ctx, name, pt, data, opts)
}

func (c *Client) DeleteSecret(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	return c.Clientset.CoreV1().Secrets(namespace).Delete(ctx, name, opts)
}

func (c *Client) GetSecret(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*corev1.Secret, error) {
	return c.Clientset.CoreV1().Secrets(namespace).Get(ctx, name, opts)
}

func (c *Client) CreateServiceAccount(ctx context.Context, namespace string, account *corev1.ServiceAccount, opts metav1.CreateOptions) (*corev1.ServiceAccount, error) {
	return c.Clientset.CoreV1().ServiceAccounts(namespace).Create(ctx, account, opts)
}

func (c *Client) DeleteServiceAccount(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	return c.Clientset.CoreV1().ServiceAccounts(namespace).Delete(ctx, name, opts)
}

func (c *Client) CreateClusterRole(ctx context.Context, role *rbacv1.ClusterRole, opts metav1.CreateOptions) (*rbacv1.ClusterRole, error) {
	return c.Clientset.RbacV1().ClusterRoles().Create(ctx, role, opts)
}

func (c *Client) DeleteClusterRole(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.Clientset.RbacV1().ClusterRoles().Delete(ctx, name, opts)
}

func (c *Client) CreateClusterRoleBinding(ctx context.Context, role *rbacv1.ClusterRoleBinding, opts metav1.CreateOptions) (*rbacv1.ClusterRoleBinding, error) {
	return c.Clientset.RbacV1().ClusterRoleBindings().Create(ctx, role, opts)
}

func (c *Client) DeleteClusterRoleBinding(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.Clientset.RbacV1().ClusterRoleBindings().Delete(ctx, name, opts)
}

func (c *Client) GetConfigMap(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*corev1.ConfigMap, error) {
	return c.Clientset.CoreV1().ConfigMaps(namespace).Get(ctx, name, opts)
}

func (c *Client) PatchConfigMap(ctx context.Context, namespace, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions) (*corev1.ConfigMap, error) {
	return c.Clientset.CoreV1().ConfigMaps(namespace).Patch(ctx, name, pt, data, opts)
}

func (c *Client) UpdateConfigMap(ctx context.Context, configMap *corev1.ConfigMap, opts metav1.UpdateOptions) (*corev1.ConfigMap, error) {
	return c.Clientset.CoreV1().ConfigMaps(configMap.Namespace).Update(ctx, configMap, opts)
}

func (c *Client) CreateConfigMap(ctx context.Context, namespace string, config *corev1.ConfigMap, opts metav1.CreateOptions) (*corev1.ConfigMap, error) {
	return c.Clientset.CoreV1().ConfigMaps(namespace).Create(ctx, config, opts)
}

func (c *Client) DeleteConfigMap(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	return c.Clientset.CoreV1().ConfigMaps(namespace).Delete(ctx, name, opts)
}

func (c *Client) CreateService(ctx context.Context, namespace string, service *corev1.Service, opts metav1.CreateOptions) (*corev1.Service, error) {
	return c.Clientset.CoreV1().Services(namespace).Create(ctx, service, opts)
}

func (c *Client) DeleteService(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	return c.Clientset.CoreV1().Services(namespace).Delete(ctx, name, opts)
}

func (c *Client) GetService(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*corev1.Service, error) {
	return c.Clientset.CoreV1().Services(namespace).Get(ctx, name, opts)
}

func (c *Client) ListServices(ctx context.Context, namespace string, options metav1.ListOptions) (*corev1.ServiceList, error) {
	return c.Clientset.CoreV1().Services(namespace).List(ctx, options)
}

func (c *Client) CreateDeployment(ctx context.Context, namespace string, deployment *appsv1.Deployment, opts metav1.CreateOptions) (*appsv1.Deployment, error) {
	return c.Clientset.AppsV1().Deployments(namespace).Create(ctx, deployment, opts)
}

func (c *Client) GetDeployment(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*appsv1.Deployment, error) {
	return c.Clientset.AppsV1().Deployments(namespace).Get(ctx, name, opts)
}

func (c *Client) DeleteDeployment(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	return c.Clientset.AppsV1().Deployments(namespace).Delete(ctx, name, opts)
}

func (c *Client) PatchDeployment(ctx context.Context, namespace, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions) (*appsv1.Deployment, error) {
	return c.Clientset.AppsV1().Deployments(namespace).Patch(ctx, name, pt, data, opts)
}

func (c *Client) CheckDeploymentStatus(ctx context.Context, namespace, deployment string) error {
	d, err := c.GetDeployment(ctx, namespace, deployment, metav1.GetOptions{})
	if err != nil {
		return err
	}

	if d == nil {
		return fmt.Errorf("deployment is not available")
	}

	if d.Status.ObservedGeneration != d.Generation {
		return fmt.Errorf("observed generation (%d) is older than generation of the desired state (%d)",
			d.Status.ObservedGeneration, d.Generation)
	}

	if d.Status.Replicas == 0 {
		return fmt.Errorf("replicas count is zero")
	}

	if d.Status.AvailableReplicas != d.Status.Replicas {
		return fmt.Errorf("only %d of %d replicas are available", d.Status.AvailableReplicas, d.Status.Replicas)
	}

	if d.Status.ReadyReplicas != d.Status.Replicas {
		return fmt.Errorf("only %d of %d replicas are ready", d.Status.ReadyReplicas, d.Status.Replicas)
	}

	if d.Status.UpdatedReplicas != d.Status.Replicas {
		return fmt.Errorf("only %d of %d replicas are up-to-date", d.Status.UpdatedReplicas, d.Status.Replicas)
	}

	return nil
}

func (c *Client) CreateDaemonSet(ctx context.Context, namespace string, ds *appsv1.DaemonSet, opts metav1.CreateOptions) (*appsv1.DaemonSet, error) {
	return c.Clientset.AppsV1().DaemonSets(namespace).Create(ctx, ds, opts)
}

func (c *Client) PatchDaemonSet(ctx context.Context, namespace, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions) (*appsv1.DaemonSet, error) {
	return c.Clientset.AppsV1().DaemonSets(namespace).Patch(ctx, name, pt, data, opts)
}

func (c *Client) GetDaemonSet(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*appsv1.DaemonSet, error) {
	return c.Clientset.AppsV1().DaemonSets(namespace).Get(ctx, name, opts)
}

func (c *Client) ListDaemonSet(ctx context.Context, namespace string, o metav1.ListOptions) (*appsv1.DaemonSetList, error) {
	return c.Clientset.AppsV1().DaemonSets(namespace).List(ctx, o)
}

func (c *Client) DeleteDaemonSet(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	return c.Clientset.AppsV1().DaemonSets(namespace).Delete(ctx, name, opts)
}

func (c *Client) CreateNamespace(ctx context.Context, namespace string, opts metav1.CreateOptions) (*corev1.Namespace, error) {
	return c.Clientset.CoreV1().Namespaces().Create(ctx, &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: namespace}}, opts)
}

func (c *Client) GetNamespace(ctx context.Context, namespace string, options metav1.GetOptions) (*corev1.Namespace, error) {
	return c.Clientset.CoreV1().Namespaces().Get(ctx, namespace, options)
}

func (c *Client) DeleteNamespace(ctx context.Context, namespace string, opts metav1.DeleteOptions) error {
	return c.Clientset.CoreV1().Namespaces().Delete(ctx, namespace, opts)
}

func (c *Client) ListNamespaces(ctx context.Context, o metav1.ListOptions) (*corev1.NamespaceList, error) {
	return c.Clientset.CoreV1().Namespaces().List(ctx, o)
}

func (c *Client) ListEvents(ctx context.Context, o metav1.ListOptions) (*corev1.EventList, error) {
	return c.Clientset.CoreV1().Events(corev1.NamespaceAll).List(ctx, o)
}

func (c *Client) DeletePod(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	return c.Clientset.CoreV1().Pods(namespace).Delete(ctx, name, opts)
}

func (c *Client) DeletePodCollection(ctx context.Context, namespace string, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	return c.Clientset.CoreV1().Pods(namespace).DeleteCollection(ctx, opts, listOpts)
}

func (c *Client) ListPods(ctx context.Context, namespace string, options metav1.ListOptions) (*corev1.PodList, error) {
	return c.Clientset.CoreV1().Pods(namespace).List(ctx, options)
}

func (c *Client) PodLogs(namespace, name string, opts *corev1.PodLogOptions) *rest.Request {
	return c.Clientset.CoreV1().Pods(namespace).GetLogs(name, opts)
}

func (c *Client) GetLogs(ctx context.Context, namespace, name, container string, sinceTime time.Time, limitBytes int64, previous bool) (string, error) {
	t := metav1.NewTime(sinceTime)
	o := corev1.PodLogOptions{
		Container:  container,
		Follow:     false,
		LimitBytes: &limitBytes,
		Previous:   previous,
		SinceTime:  &t,
		Timestamps: true,
	}
	r := c.Clientset.CoreV1().Pods(namespace).GetLogs(name, &o)
	s, err := r.Stream(ctx)
	if err != nil {
		return "", err
	}
	defer s.Close()
	var b bytes.Buffer
	if _, err = io.Copy(&b, s); err != nil {
		return "", err
	}
	return b.String(), nil
}

func (c *Client) ExecInPodWithStderr(ctx context.Context, namespace, pod, container string, command []string) (bytes.Buffer, bytes.Buffer, error) {
	result, err := c.execInPod(ctx, ExecParameters{
		Namespace: namespace,
		Pod:       pod,
		Container: container,
		Command:   command,
	})
	return result.Stdout, result.Stderr, err
}

func (c *Client) ExecInPodWithTTY(ctx context.Context, namespace, pod, container string, command []string) (bytes.Buffer, error) {
	result, err := c.execInPod(ctx, ExecParameters{
		Namespace: namespace,
		Pod:       pod,
		Container: container,
		Command:   command,
		TTY:       true,
	})
	// Using TTY (for context cancellation support) fuses stderr into stdout
	return result.Stdout, err
}

func (c *Client) ExecInPod(ctx context.Context, namespace, pod, container string, command []string) (bytes.Buffer, error) {
	result, err := c.execInPod(ctx, ExecParameters{
		Namespace: namespace,
		Pod:       pod,
		Container: container,
		Command:   command,
	})
	if err != nil {
		return bytes.Buffer{}, err
	}

	if errString := result.Stderr.String(); errString != "" {
		return bytes.Buffer{}, fmt.Errorf("command failed: %s", errString)
	}

	return result.Stdout, nil
}

func (c *Client) CreateIngressClass(ctx context.Context, ingressClass *networkingv1.IngressClass, opts metav1.CreateOptions) (*networkingv1.IngressClass, error) {
	return c.Clientset.NetworkingV1().IngressClasses().Create(ctx, ingressClass, opts)
}

func (c *Client) DeleteIngressClass(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.Clientset.NetworkingV1().IngressClasses().Delete(ctx, name, opts)
}
