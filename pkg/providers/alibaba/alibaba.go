package alibaba

import (
	"github.com/ysicing/ergo/pkg/providers"
)

const providerName = "alibaba"

func init() {
	providers.RegisterProvider(providerName, func() (providers.Provider, error) {
		return newProvider(), nil
	})
}

type Alibaba struct {
}

func newProvider() *Alibaba {
	return &Alibaba{}
}

// GetProviderName returns provider name.
func (a *Alibaba) GetProviderName() string {
	return providerName
}

func (a *Alibaba) ListLighthouse() string {
	return ""
}

func (a *Alibaba) RebootLighthouse(region string, ids []string) error {
	return nil
}
