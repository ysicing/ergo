package native

import (
	"github.com/ysicing/ergo/internal/pkg/cluster"
	"github.com/ysicing/ergo/internal/pkg/providers"
	"github.com/ysicing/ergo/internal/pkg/types"
	"github.com/ysicing/ergo/pkg/util/log"
)

// providerName is the name of this provider.
const providerName = "native"

const createUsageExample = `
	create default cluster:
		ergo kube k3s init
	create custom cluster
		ergo kube k3s --podsubnet "10.42.0.0/16" \
 			--svcsubnet "10.43.0.0/16" \
			--plugins lb \
			--plugins ingress \
			--eip 1.1.1.1  \
			--san kubeapi.k8s.io
`

type Native struct {
	*cluster.Cluster
}

func init() {
	providers.RegisterProvider(providerName, func() (providers.Provider, error) {
		return newProvider(), nil
	})
}

func newProvider() *Native {
	c := cluster.NewCluster()
	c.Provider = providerName
	return &Native{
		Cluster: c,
	}
}

// GetUsageExample returns native usage example prompt.
func (p *Native) GetUsageExample(action string) string {
	switch action {
	case "create":
		return createUsageExample
	default:
		return "not support"
	}
}

// GetCreateFlags returns native create flags.
func (p *Native) GetCreateFlags() []types.Flag {
	fs := p.GetCreateOptions()
	fs = append(fs, p.GetCreateOptions()...)
	return fs
}

func (p *Native) GetProviderName() string {
	return p.Provider
}

// InitCluster init cluster.
func (p *Native) InitCluster() (err error) {
	log.Flog.Info("start init cluster")
	return p.InitKubeCluster()
}

// JoinCluster join cluster.
func (p *Native) JoinCluster() (err error) {
	return nil
}

func (p *Native) InitSystem() error {
	log.Flog.Info("start system init")
	if err := p.SystemInit(); err != nil {
		return err
	}
	log.Flog.Donef("system init passed")
	return nil
}

func (p *Native) InitBigcat() error {
	return p.InstallBigCat()
}
