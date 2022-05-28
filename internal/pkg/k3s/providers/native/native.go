package native

import (
	"fmt"

	"github.com/ergoapi/util/exnet"
	"github.com/ysicing/ergo/internal/pkg/k3s/cluster"
	"github.com/ysicing/ergo/internal/pkg/k3s/providers"
	"github.com/ysicing/ergo/internal/pkg/k3s/types"
	"github.com/ysicing/ergo/pkg/util/log"
	"github.com/ysicing/ergo/pkg/util/preflight"
	utilsexec "k8s.io/utils/exec"
)

// providerName is the name of this provider.
const providerName = "native"

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

const (
	createUsageExample = `
	create default cluster:
		ergo kube init

	create custom cluster
		ergo kube init --podsubnet "10.42.0.0/16" \
 			--svcsubnet "10.43.0.0/16" \
			--eip 1.1.1.1  \
			--san kubeapi.k8s.io
`
	joinUsageExample = `
	join node to cluster:

		# use k3s api & k3s nodetoken
		ergo kube join --api <api address> --token <api token>
`
)

// GetUsageExample returns native usage example prompt.
func (p *Native) GetUsageExample(action string) string {
	switch action {
	case "create":
		return createUsageExample
	case "join":
		return joinUsageExample
	default:
		return "not support"
	}
}

// GetCreateFlags returns native create flags.
func (p *Native) GetCreateFlags() []types.Flag {
	fs := p.GetCreateOptions()
	fs = append(fs, p.GetCreateExtOptions()...)
	return fs
}

// GetJoinFlags returns native join flags.
func (p *Native) GetJoinFlags() []types.Flag {
	return p.GetJoinOptions()
}

func (p *Native) GetProviderName() string {
	return p.Provider
}

// CreateCluster create cluster.
func (p *Native) CreateCluster() (err error) {
	log.Flog.Info("start init cluster")
	return p.InitCluster()
}

// JoinNode join node.
func (p *Native) JoinNode() (err error) {
	return p.JoinCluster()
}

func (p *Native) InitBigCat() error {
	return p.InstallBigCat()
}

func (p *Native) CreateCheck() error {
	log.Flog.Info("start pre-flight checks")
	if err := preflight.RunInitNodeChecks(utilsexec.New(), &p.Metadata, false); err != nil {
		return err
	}
	log.Flog.Done("pre-flight checks passed")
	return nil
}

func (p *Native) PreSystemInit() error {
	log.Flog.Info("start system init")
	if err := p.SystemInit(); err != nil {
		return err
	}
	log.Flog.Done("system init passed")
	return nil
}

// GenerateManifest generates manifest deploy command.
func (p *Native) GenerateManifest() []string {
	// no need to support.
	return nil
}

// Show show cluster info.
func (p *Native) Show() {
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
