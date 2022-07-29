package status

import (
	"context"
	"fmt"
	"time"

	"github.com/ergoapi/log"
	"github.com/ergoapi/util/file"
	"github.com/ysicing/ergo/common"
	"github.com/ysicing/ergo/internal/kube"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type K8sStatusOption struct {
	Namespace      string
	KubeConfig     string
	Wait           bool
	WaitDuration   time.Duration
	IgnoreWarnings bool
	ListOutput     string
}

type K8sStatusCollector struct {
	log    log.Logger
	client *kube.Client
	option K8sStatusOption
}

func NewK8sStatusCollector(option K8sStatusOption) (*K8sStatusCollector, error) {
	client, err := kube.NewClient("", option.KubeConfig)
	return &K8sStatusCollector{
		client: client,
		option: option,
	}, err
}

func (s K8sStatusOption) waitTimeout() time.Duration {
	if s.WaitDuration != time.Duration(0) {
		return s.WaitDuration
	}

	return common.StatusWaitDuration
}

func (k *K8sStatusCollector) Status(ctx context.Context) (*Status, error) {
	var mostRecentStatus *Status
	ctx, cancel := context.WithTimeout(ctx, k.option.waitTimeout())
	defer cancel()

retry:
	select {
	case <-ctx.Done():
		return mostRecentStatus, fmt.Errorf("timeout while waiting for status to become successful: %w", ctx.Err())
	default:
	}
	s := k.status(ctx)
	if s != nil {
		mostRecentStatus = s
	}
	if !k.statusIsReady(s) && k.option.Wait {
		time.Sleep(common.WaitRetryInterval)
		goto retry
	}

	return mostRecentStatus, nil
}

func (k *K8sStatusCollector) statusIsReady(s *Status) bool {
	// TODO
	return true
}

func (k *K8sStatusCollector) status(ctx context.Context) *Status {
	status := newStatus(k.option.ListOutput)
	if err := k.nodesStatus(ctx, status); err != nil {
		k.log.Errorf("failed to get nodes status: %v", err)
	}
	if err := k.bigcatStatus(ctx, status); err != nil {
		k.log.Errorf("failed to get bigcat status: %v", err)
	}
	return status
}

func (k *K8sStatusCollector) deploymentStatus(ctx context.Context, ns, name, t string, status *Status) (bool, error) {
	k.log.Debugf("check cm %s status", name)
	stateCount := PodStateCount{Type: "Deployment"}
	d, err := k.client.GetDeployment(ctx, ns, name, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		stateCount.Disabled = true
		if t == "k8s" {
			status.KubeStatus.PodState[name] = stateCount
		}
		return true, nil
	}

	if err != nil {
		stateCount.Disabled = false
		if t == "k8s" {
			status.KubeStatus.PodState[name] = stateCount
		}
		return false, err
	}

	if d == nil {
		stateCount.Disabled = false
		// TODO type: k8s, api
		if t == "k8s" {
			status.KubeStatus.PodState[name] = stateCount
		}
		return false, fmt.Errorf("component %s is not available", name)
	}

	stateCount.Desired = int(d.Status.Replicas)
	stateCount.Ready = int(d.Status.ReadyReplicas)
	stateCount.Available = int(d.Status.AvailableReplicas)
	stateCount.Unavailable = int(d.Status.UnavailableReplicas)
	stateCount.Disabled = false
	if t == "k8s" {
		status.KubeStatus.PodState[name] = stateCount
	}
	notReady := stateCount.Desired - stateCount.Ready
	if notReady > 0 {
		k.log.Warnf("%d pods of Deployment %s are not ready", notReady, name)
	}
	if unavailable := stateCount.Unavailable - notReady; unavailable > 0 {
		k.log.Warnf("%d pods of Deployment %s are not available", unavailable, name)
	}
	return false, nil
}

func (k *K8sStatusCollector) bigcatStatus(ctx context.Context, status *Status) error {
	// 集群
	k.deploymentStatus(ctx, "kube-system", "coredns", "k8s", status)
	k.deploymentStatus(ctx, "kube-system", "metrics-server", "k8s", status)
	k.deploymentStatus(ctx, "kube-system", "local-path-provisioner", "k8s", status)
	return nil
}

func (k *K8sStatusCollector) nodesStatus(ctx context.Context, status *Status) error {
	nodes, err := k.client.ListNodes(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}
	nodetotal := len(nodes.Items)
	status.KubeStatus.NodeCount["total"] = nodetotal
	readynode := 0
	master := 0
	worker := 0
	for _, node := range nodes.Items {
		for _, nc := range node.Status.Conditions {
			if nc.Type == "Ready" && nc.Status == corev1.ConditionTrue {
				readynode++
			}
		}
		if node.Labels["node-role.kubernetes.io/master"] == "true" || node.Labels["node-role.kubernetes.io/control-plane"] == "true" {
			master++
		} else {
			worker++
		}
	}
	status.KubeStatus.NodeCount["ready"] = readynode
	status.KubeStatus.NodeCount["master"] = master
	status.KubeStatus.NodeCount["worker"] = worker
	if master == 0 {
		status.KubeStatus.Type = "cloud-managed"
	} else {
		if file.CheckFileExists(common.GetCustomConfig(common.InitFileName)) {
			status.KubeStatus.Type = "self-managed"
		} else {
			status.KubeStatus.Type = "owner-self-managed"
		}
	}
	version, _ := k.client.GetGitVersion(ctx)
	if len(version) > 0 {
		status.KubeStatus.Version = version
	}
	return nil
}
