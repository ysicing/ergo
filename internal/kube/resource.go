package kube

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metricsapi "k8s.io/metrics/pkg/apis/metrics"
)

// CPUResources describes node allocated resources.
type CPUResources struct {
	// CPUUsages is number of allocated milicores.
	CPUUsages *CPUResource

	// CPURequests is number of allocated milicores.
	CPURequests *CPUResource

	// CPURequestsFraction is a fraction of CPU, that is allocated.
	CPURequestsFraction float64 `json:"cpuRequestsFraction"`

	// CPULimits is defined CPU limit.
	CPULimits *CPUResource

	// CPULimitsFraction is a fraction of defined CPU limit, can be over 100%, i.e.
	// overcommitted.
	CPULimitsFraction float64 `json:"cpuLimitsFraction"`

	// CPUCapacity is specified node CPU capacity in milicores.
	CPUCapacity *CPUResource
}

// MemoryResources describes node allocated resources.
type MemoryResources struct {
	// MemoryUsages is a fraction of memory, that is allocated.
	MemoryUsages *MemoryResource

	// MemoryRequests is a fraction of memory, that is allocated.
	MemoryRequests *MemoryResource

	// MemoryRequestsFraction is a fraction of memory, that is allocated.
	MemoryRequestsFraction float64 `json:"memoryRequestsFraction"`

	// MemoryLimits is defined memory limit.
	MemoryLimits *MemoryResource

	// MemoryLimitsFraction is a fraction of defined memory limit, can be over 100%, i.e.
	// overcommitted.
	MemoryLimitsFraction float64 `json:"memoryLimitsFraction"`

	// MemoryCapacity is specified node memory capacity in bytes.
	MemoryCapacity *MemoryResource
}

// PodResources describes node allocated resources.
type PodResources struct {
	// AllocatedPods in number of currently allocated pods on the node.
	AllocatedPods int `json:"allocatedPods"`

	// PodCapacity is maximum number of pods, that can be allocated on the node.
	PodCapacity int64 `json:"podCapacity"`

	// PodFraction is a fraction of pods, that can be allocated on given node.
	PodFraction float64 `json:"podFraction"`
}

// NodeAllocatedResources describes node allocated resources.
type NodeAllocatedResources struct {
	CPUResources
	MemoryResources
	PodResources
}

// PodAllocatedResources describes node allocated resources.
type PodAllocatedResources struct {
	// CPUUsages is number of allocated milicores.
	CPUUsages *CPUResource

	// CPURequestsFraction is a fraction of CPU, that is allocated.
	CPUUsagesFraction float64 `json:"cpuUsagesFraction"`

	// CPURequests is number of allocated milicores.
	CPURequests *CPUResource

	// CPULimits is defined CPU limit.
	CPULimits *CPUResource

	// MemoryUsages is a fraction of memory, that is allocated.
	MemoryUsages *MemoryResource

	// MemoryRequestsFraction is a fraction of memory, that is allocated.
	MemoryUsagesFraction float64 `json:"memoryUsagesFraction"`

	// MemoryRequests is a fraction of memory, that is allocated.
	MemoryRequests *MemoryResource

	// MemoryLimits is defined memory limit.
	MemoryLimits *MemoryResource
}

// NodeCapacity
func NodeCapacity(node *v1.Node) v1.ResourceList {
	allocatable := node.Status.Capacity
	if len(node.Status.Allocatable) > 0 {
		allocatable = node.Status.Allocatable
	}
	return allocatable
}

// getNodeMetricsByNodeName returns a map of node metrics where the keys are the particular node names
func getNodeMetricsByNodeName(nodeMetricsList *metricsapi.NodeMetricsList) map[string]metricsapi.NodeMetrics {
	nodeMetricsByName := make(map[string]metricsapi.NodeMetrics)
	for _, metrics := range nodeMetricsList.Items {
		nodeMetricsByName[metrics.Name] = metrics
	}

	return nodeMetricsByName
}

// func getPodMetrics(m *metricsapi.PodMetrics) v1.ResourceList {
// 	podMetrics := make(v1.ResourceList)
// 	for _, res := range metricsutil.MeasuredResources {
// 		podMetrics[res], _ = resource.ParseQuantity("0")
// 	}

// 	for _, c := range m.Containers {
// 		for _, res := range metricsutil.MeasuredResources {
// 			quantity := podMetrics[res]
// 			quantity.Add(c.Usage[res])
// 			podMetrics[res] = quantity
// 		}
// 	}
// 	return podMetrics
// }

// getNodeAllocatedResources https://github.com/kubernetes/dashboard/blob/d386ff60597b6eab0222f2c3c4aecf8e49b3014e/src/app/backend/resource/node/detail.go\#L171
func getNodeAllocatedResources(node v1.Node, podList *v1.PodList, nodeMetricsList *metricsapi.NodeMetricsList) (NodeAllocatedResources, error) {
	reqs, limits := map[v1.ResourceName]resource.Quantity{}, map[v1.ResourceName]resource.Quantity{}

	for _, pod := range podList.Items {
		podReqs, podLimits, err := PodRequestsAndLimits(pod)
		if err != nil {
			return NodeAllocatedResources{}, err
		}
		for podReqName, podReqValue := range podReqs {
			if value, ok := reqs[podReqName]; !ok {
				reqs[podReqName] = podReqValue.DeepCopy()
			} else {
				value.Add(podReqValue)
				reqs[podReqName] = value
			}
		}
		for podLimitName, podLimitValue := range podLimits {
			if value, ok := limits[podLimitName]; !ok {
				limits[podLimitName] = podLimitValue.DeepCopy()
			} else {
				value.Add(podLimitValue)
				limits[podLimitName] = value
			}
		}
	}

	nodeMetricsByNodeName := getNodeMetricsByNodeName(nodeMetricsList)
	usageMetrics := nodeMetricsByNodeName[node.Name]

	capacity := NodeCapacity(&node)
	_cpuRequests, _cpuLimits, _memoryRequests, _memoryLimits := reqs[v1.ResourceCPU], limits[v1.ResourceCPU],
		reqs[v1.ResourceMemory], limits[v1.ResourceMemory]
	_cpuUsages, _memoryUsages := usageMetrics.Usage.Cpu().MilliValue(), usageMetrics.Usage.Memory().Value()

	cpuUsages := NewCPUResource(_cpuUsages)
	cpuRequests := NewCPUResource(_cpuRequests.MilliValue())
	cpuLimits := NewCPUResource(_cpuLimits.MilliValue())

	memoryUsages := NewMemoryResource(_memoryUsages)
	memoryRequests := NewMemoryResource(_memoryRequests.Value())
	memoryLimits := NewMemoryResource(_memoryLimits.Value())
	podCapacity := capacity.Pods().Value()
	podFraction := calcPercentage(int64(len(podList.Items)), podCapacity)

	nodeAllocatedResources := NodeAllocatedResources{
		CPUResources{
			CPUUsages:           cpuUsages,
			CPURequests:         cpuRequests,
			CPURequestsFraction: cpuRequests.calcPercentage(capacity.Cpu()),
			CPULimits:           cpuLimits,
			CPULimitsFraction:   cpuLimits.calcPercentage(capacity.Cpu()),
			CPUCapacity:         NewCPUResource(capacity.Cpu().MilliValue()),
		},
		MemoryResources{
			MemoryUsages:           memoryUsages,
			MemoryRequests:         memoryRequests,
			MemoryRequestsFraction: memoryRequests.calcPercentage(capacity.Memory()),
			MemoryLimits:           memoryLimits,
			MemoryLimitsFraction:   memoryLimits.calcPercentage(capacity.Memory()),
			MemoryCapacity:         NewMemoryResource(capacity.Memory().Value()),
		},
		PodResources{
			AllocatedPods: len(podList.Items),
			PodCapacity:   podCapacity,
			PodFraction:   podFraction,
		},
	}
	return nodeAllocatedResources, nil
}

// getPodAllocatedResources
func getPodAllocatedResources(pod v1.Pod, podmetric *metricsapi.PodMetrics) (PodAllocatedResources, error) {
	reqs, limits := map[v1.ResourceName]resource.Quantity{}, map[v1.ResourceName]resource.Quantity{}

	podReqs, podLimits, err := PodRequestsAndLimits(pod)
	if err != nil {
		return PodAllocatedResources{}, err
	}
	for podReqName, podReqValue := range podReqs {
		if value, ok := reqs[podReqName]; !ok {
			reqs[podReqName] = podReqValue.DeepCopy()
		} else {
			value.Add(podReqValue)
			reqs[podReqName] = value
		}
	}
	for podLimitName, podLimitValue := range podLimits {
		if value, ok := limits[podLimitName]; !ok {
			limits[podLimitName] = podLimitValue.DeepCopy()
		} else {
			value.Add(podLimitValue)
			limits[podLimitName] = value
		}
	}
	//podMetricsByPodName := getPodMetricsByPodName(podMetricsList)
	//usageMetrics := podMetricsByPodName[pod.Name]
	usageMetrics := getPodMetrics(podmetric)

	_cpuRequests, _cpuLimits, _memoryRequests, _memoryLimits := reqs[v1.ResourceCPU], limits[v1.ResourceCPU],
		reqs[v1.ResourceMemory], limits[v1.ResourceMemory]
	_cpuUsages, _memoryUsages := usageMetrics[v1.ResourceCPU], usageMetrics[v1.ResourceMemory]

	cpuUsages := NewCPUResource(_cpuUsages.MilliValue())
	cpuRequests := NewCPUResource(_cpuRequests.MilliValue())
	cpuLimits := NewCPUResource(_cpuLimits.MilliValue())

	memoryUsages := NewMemoryResource(_memoryUsages.Value())
	memoryRequests := NewMemoryResource(_memoryRequests.Value())
	memoryLimits := NewMemoryResource(_memoryLimits.Value())

	podAllocatedResources := PodAllocatedResources{
		CPUUsages:            cpuUsages,
		CPUUsagesFraction:    cpuUsages.calcPercentage(&_cpuLimits),
		CPURequests:          cpuRequests,
		CPULimits:            cpuLimits,
		MemoryUsages:         memoryUsages,
		MemoryUsagesFraction: memoryUsages.calcPercentage(&_memoryLimits),
		MemoryRequests:       memoryRequests,
		MemoryLimits:         memoryLimits,
	}
	return podAllocatedResources, nil
}

// PodRequestsAndLimits returns a dictionary of all defined resources summed up for all
// containers of the pod. If pod overhead is non-nil, the pod overhead is added to the
// total container resource requests and to the total container limits which have a
// non-zero quantity.
func PodRequestsAndLimits(pod v1.Pod) (reqs, limits v1.ResourceList, err error) {
	reqs, limits = v1.ResourceList{}, v1.ResourceList{}
	for _, container := range pod.Spec.Containers {
		addResourceList(reqs, container.Resources.Requests)
		addResourceList(limits, container.Resources.Limits)
	}
	// init containers define the minimum of any resource
	for _, container := range pod.Spec.InitContainers {
		maxResourceList(reqs, container.Resources.Requests)
		maxResourceList(limits, container.Resources.Limits)
	}

	// Add overhead for running a pod to the sum of requests and to non-zero limits:
	if pod.Spec.Overhead != nil {
		addResourceList(reqs, pod.Spec.Overhead)

		for name, quantity := range pod.Spec.Overhead {
			if value, ok := limits[name]; ok && !value.IsZero() {
				value.Add(quantity)
				limits[name] = value
			}
		}
	}
	return
}

// addResourceList adds the resources in newList to list
func addResourceList(list, rl v1.ResourceList) {
	for name, quantity := range rl {
		if value, ok := list[name]; !ok {
			list[name] = quantity.DeepCopy()
		} else {
			value.Add(quantity)
			list[name] = value
		}
	}
}

// maxResourceList sets list to the greater of list/newList for every resource
// either list
func maxResourceList(list, rl v1.ResourceList) {
	for name, quantity := range rl {
		if value, ok := list[name]; !ok {
			list[name] = quantity.DeepCopy()
			continue
		} else if quantity.Cmp(value) > 0 {
			list[name] = quantity.DeepCopy()
		}
	}
}
