package resource

import (
	"fmt"
	"os"
	"strings"

	"github.com/gosuri/uitable"
	"k8s.io/apimachinery/pkg/labels"
	"github.com/ysicing/ergo/internal/kube"
	"github.com/ysicing/ergo/pkg/util/output"
)

func (o *Option) NodeValidate() {
	// TODO
}

func (o *Option) RunResourceNode() error {
	selector := labels.Everything()
	var err error
	if len(o.Selector) > 0 {
		selector, err = labels.Parse(o.Selector)
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
	data, err := k.GetNodeResources(o.SortBy, selector)
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
		table.AddRow("Name", "IP", "CPU使用", "CPU分配", "CPU限制", "CPU容量", "内存使用", "内存分配", "内存限制", "内存容量", "pod数", "pod容量", "存活时间")
		for _, d := range data {
			table.AddRow(d.NodeName, d.NodeIP,
				d.CPUUsages, fmt.Sprintf("%v(%v)", d.CPURequests, d.CPURequestsFraction), fmt.Sprintf("%v(%v)", d.CPULimits, d.CPULimitsFraction), d.CPUCapacity,
				d.MemoryUsages, fmt.Sprintf("%v(%v)", d.MemoryRequests, d.MemoryRequestsFraction), fmt.Sprintf("%v(%v)", d.MemoryLimits, d.MemoryLimitsFraction), d.MemoryCapacity,
				fmt.Sprintf("%v(%v)", d.AllocatedPods, d.PodFraction), d.PodCapacity, d.Age)
		}
		return output.EncodeTable(os.Stdout, table)
	}
}
