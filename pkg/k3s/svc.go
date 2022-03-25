package k3s

import "github.com/ysicing/ergo/internal/kube"

func (o *Option) RebuildService() error {
	cfg := kube.ClientConfig{
		KubeCtx:    o.KubeCtx,
		KubeConfig: o.KubeConfig,
	}
	// TODO
	_, err := kube.NewKubeClient(&cfg)
	if err != nil {
		return err
	}
	return nil
}
