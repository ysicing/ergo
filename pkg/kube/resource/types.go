package resource

type Option struct {
	Namespace     string
	LabelSelector string
	FieldSelector string
	SortBy        string
	QPS           float32
	Burst         int
	KubeCtx       string
	KubeConfig    string
	Output        string
	Selector      string
}

func (o *Option) Validate() {
	if len(o.SortBy) > 0 {
		if o.SortBy != "cpu" {
			o.SortBy = "memory"
		}
	}
}
