package incluster

import (
	"fmt"

	"github.com/ergoapi/util/exnet"
	"github.com/ergoapi/util/file"
	"github.com/ysicing/ergo/common"
	"github.com/ysicing/ergo/internal/pkg/k3s/cluster"
	"github.com/ysicing/ergo/internal/pkg/k3s/providers"
	"github.com/ysicing/ergo/internal/pkg/k3s/types"
	"github.com/ysicing/ergo/pkg/util/log"
)

// providerName is the name of this provider.
const providerName = "incluster"

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

const createUsageExample = `
	create BigCat cluster:
		ergo kube init
`

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
	fs := p.GetCreateExtOptions()
	return fs
}

func (p *InCluster) GetProviderName() string {
	return p.Provider
}

// CreateCluster create cluster.
func (p *InCluster) CreateCluster() (err error) {
	log.Flog.Warn("exists cluster, check cluster status")
	return nil
}

// JoinNode join node.
func (p *InCluster) JoinNode() (err error) {
	return nil
}

func (p *InCluster) InitBigCat() error {
	log.Flog.Info("start init BigCat")
	if err := p.InstallBigCat(); err != nil {
		return err
	}
	file.Writefile(common.GetCustomConfig(common.InitModeCluster), "in cluster ok")
	return nil
}

func (p *InCluster) CreateCheck() error {
	// no need to support.
	return nil
}

func (p *InCluster) PreSystemInit() error {
	// no need to support.
	return nil
}

// GenerateManifest generates manifest deploy command.
func (p *InCluster) GenerateManifest() []string {
	// no need to support.
	return nil
}

// Show show cluster info.
func (p *InCluster) Show() {
	loginip := p.Metadata.EIP
	if len(loginip) <= 0 {
		loginip = exnet.LocalIPs()[0]
	}
	// cfg, _ := config.LoadConfig()
	// if cfg != nil {
	// 	cfg.DB = "sqlite"
	// 	cfg.Token = kutil.GetNodeToken()
	// 	cfg.Master = []config.Node{
	// 		{
	// 			Name: zos.GetHostname(),
	// 			Host: loginip,
	// 			Init: true,
	// 		},
	// 	}
	// 	cfg.SaveConfig()
	// }

	log.Flog.Info("----------------------------")
	log.Flog.Donef("web:: %s", fmt.Sprintf("http://%s:32379", loginip))
	log.Flog.Donef("docs: %s", "https://github.com/ysicing/ergo")
}
