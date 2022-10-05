package tencent

import (
	"github.com/ysicing/ergo/pkg/providers"
)

const providerName = "tencent"

func init() {
	providers.RegisterProvider(providerName, func() (providers.Provider, error) {
		return newProvider(), nil
	})
}

type Tencent struct {
}

func newProvider() *Tencent {
	return &Tencent{}
}

// GetProviderName returns provider name.
func (p *Tencent) GetProviderName() string {
	return providerName
}
