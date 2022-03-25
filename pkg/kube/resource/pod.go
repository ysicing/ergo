package resource

import (
	"fmt"
	"os"
	"strings"

	"github.com/gosuri/uitable"
	"github.com/ysicing/ergo/internal/kube"
	"github.com/ysicing/ergo/pkg/util/output"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
)

func (o *Option) RunResourcePod() error {
	labelSelector := labels.Everything()
	var err error
	if len(o.LabelSelector) > 0 {
		labelSelector, err = labels.Parse(o.LabelSelector)
		if err != nil {
			return err
		}
	}
	fieldSelector := fields.Everything()
	if len(o.FieldSelector) > 0 {
		fieldSelector, err = fields.ParseSelector(o.FieldSelector)
		if err != nil {
			return err
		}
	}
	cfg := kube.ClientConfig{
		KubeCtx:    o.KubeCtx,
		KubeConfig: o.KubeConfig,
	}
	k, err := kube.NewKubeClient(&cfg)
	if err != nil {
		return err
	}
	metrics, err := k.GetPodMetricsFromMetricsAPI(o.Namespace, labelSelector, fieldSelector)
	if err != nil {
		return err
	}
	if len(metrics.Items) == 0 {
		return nil
	}
	data, err := k.GetPodResources(metrics.Items, o.Namespace, o.SortBy)
	if err != nil {
		return err
	}
	switch strings.ToLower(o.Output) {
	case "json":
		return output.EncodeJSON(os.Stdout, data)
	case "yaml":
		return output.EncodeYAML(os.Stdout, data)
	default:
		table := uitable.New()
		table.AddRow("Namespace", "Name", "CPU使用", "CPU分配", "CPU限制", "内存使用", "内存分配", "内存限制")
		for _, d := range data {
			table.AddRow(d.Namespace, d.Name,
				fmt.Sprintf("%v(%v)", d.CPUUsages, d.CPUUsagesFraction), d.CPURequests, d.CPULimits,
				fmt.Sprintf("%v(%v)", d.MemoryUsages, d.MemoryUsagesFraction), d.MemoryRequests, d.MemoryLimits)
		}
		return output.EncodeTable(os.Stdout, table)
	}
}
