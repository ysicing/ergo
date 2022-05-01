package incluster

import (
	"github.com/ysicing/ergo/internal/pkg/cluster"
	"github.com/ysicing/ergo/internal/pkg/providers"
	"github.com/ysicing/ergo/internal/pkg/types"
)

// providerName is the name of this provider.
const providerName = "incluster"

const createUsageExample = `
	create default cluster:
		ergo kube k3s init
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
	return nil
}

func (p *Native) GetProviderName() string {
	return p.Provider
}

// InitCluster init cluster.
func (p *Native) InitCluster() (err error) {
	return nil
}

// JoinCluster join cluster.
func (p *Native) JoinCluster() (err error) {
	return nil
}

func (p *Native) InitSystem() error {
	return nil
}
