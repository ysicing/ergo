package kube

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ergoapi/util/color"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

const (
	warningThreshold  = 80.00
	criticalThreshold = 90.00
)

// NewGpuResource returns the list of NewGpuResource
func NewGpuResource(name v1.ResourceName, rl *v1.ResourceList) *resource.Quantity {
	if val, ok := (*rl)[name]; ok {
		return &val
	}
	return rl.Name(name, resource.DecimalSI)
}

// calcPercentage
func calcPercentage(dividend, divisor int64) float64 {
	if divisor > 0 {
		value, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(dividend)/float64(divisor)*100), 64)
		return value
	}
	return float64(0)
}

type MemoryResource struct {
	*resource.Quantity
}

// NewMemoryResource
func NewMemoryResource(value int64) *MemoryResource {
	return &MemoryResource{resource.NewQuantity(value, resource.BinarySI)}
}

// calcPercentage
func (r *MemoryResource) calcPercentage(divisor *resource.Quantity) float64 {
	return calcPercentage(r.Value(), divisor.Value())
}

func (r *MemoryResource) String() string {
	// XXX: Support more units
	return fmt.Sprintf("%vMi", r.Value()/(1024*1024))
}

// ToQuantity
func (r *MemoryResource) ToQuantity() *resource.Quantity {
	return resource.NewQuantity(r.Value(), resource.BinarySI)
}

type CPUResource struct {
	*resource.Quantity
}

// NewCPUResource
func NewCPUResource(value int64) *CPUResource {
	r := resource.NewMilliQuantity(value, resource.DecimalSI)
	return &CPUResource{r}
}

// String
func (r *CPUResource) String() string {
	// XXX: Support more units
	return fmt.Sprintf("%vm", r.MilliValue())
}

// float64ToString float64转string
func float64ToString(s float64) string {
	//return strconv.FormatFloat(s, 'G', -1, 32)
	return fmt.Sprintf("%v%%", strconv.FormatFloat(s, 'G', -1, 64))
}

// StringTofloat64
func stringTofloat64(a string) float64 {
	value, _ := strconv.ParseFloat(a, 64)
	return value
}

// calcPercentage
func (r *CPUResource) calcPercentage(divisor *resource.Quantity) float64 {
	return calcPercentage(r.MilliValue(), divisor.MilliValue())
}

// ToQuantity
func (r *CPUResource) ToQuantity() *resource.Quantity {
	return resource.NewMilliQuantity(r.MilliValue(), resource.DecimalSI)
}

// FieldString
func FieldString(str string) float64 {
	switch {
	case strings.Contains(str, "%"):
		str1 := strings.Split(str, "%")
		value, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", stringTofloat64(str1[0])), 64)
		return value
	case strings.Contains(str, "Mi"):
		str1 := strings.Split(str, "Mi")
		value, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", stringTofloat64(str1[0])), 64)
		return value
	case strings.Contains(str, "m"):
		str1 := strings.Split(str, "m")
		value, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", stringTofloat64(str1[0])), 64)
		return value
	default:
		return float64(0)
	}
}

// Compare
func ExceedsCompare(a string) string {
	if FieldString(a) > float64(criticalThreshold) {
		return redColor(a)
	} else if FieldString(a) > float64(warningThreshold) {
		return yellowColor(a)
	} else {
		return a
	}
}

func redColor(s string) string {
	return color.SRed(s)
}

func yellowColor(s string) string {
	return color.SYellow(s)
}
