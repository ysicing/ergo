package incluster

import (
	"github.com/ergoapi/util/file"
	"github.com/ysicing/ergo/common"
	"github.com/ysicing/ergo/internal/pkg/cluster"
	"github.com/ysicing/ergo/internal/pkg/providers"
	"github.com/ysicing/ergo/internal/pkg/types"
	"github.com/ysicing/ergo/pkg/util/log"
)

// providerName is the name of this provider.
const providerName = "incluster"

const createUsageExample = `
	create default cluster:
		ergo kube k3s init
`

type InCluster struct {
	*cluster.Cluster
}

func init() {
	providers.RegisterProvider(providerName, func() (providers.Provider, error) {
		return newProvider(), nil
	})
}

func newProvider() *InCluster {
	c := cluster.NewCluster()
	c.Provider = providerName
	return &InCluster{
		Cluster: c,
	}
}

// GetUsageExample returns native usage example prompt.
func (p *InCluster) GetUsageExample(action string) string {
	switch action {
	case "create":
		return createUsageExample
	default:
		return "not support"
	}
}

// GetCreateFlags returns native create flags.
func (p *InCluster) GetCreateFlags() []types.Flag {
	return nil
}

func (p *InCluster) GetProviderName() string {
	return p.Provider
}

// InitCluster init cluster.
func (p *InCluster) InitCluster() (err error) {
	return nil
}

// JoinCluster join cluster.
func (p *InCluster) JoinCluster() (err error) {
	return nil
}

func (p *InCluster) InitSystem() error {
	return nil
}

func (p *InCluster) InitBigcat() error {
	log.Flog.Info("start init bigcat")
	if err := p.InstallBigCat(); err != nil {
		return err
	}
	file.Writefile(common.GetCustomConfig(common.InitModeCluster), "in cluster ok")
	return nil
}
