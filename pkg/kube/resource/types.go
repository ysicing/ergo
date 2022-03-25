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

func (p *Option) Validate() {
	if len(p.SortBy) > 0 {
		if p.SortBy != "cpu" {
			p.SortBy = "memory"
		}
	}
}
