package providers

import (
	"fmt"
	"sync"

	"github.com/ysicing/ergo/internal/pkg/k3s/types"
)

// Factory is a function that returns a Provider.Interface.
type Factory func() (Provider, error)

var (
	providersMutex sync.Mutex
	providers      = make(map[string]Factory)
)

type Provider interface {
	GetProviderName() string
	// Get command usage example.
	GetUsageExample(action string) string
	CreateCluster() error
	JoinNode() error
	InitBigCat() error
	GetCreateFlags() []types.Flag
	CreateCheck() error
	PreSystemInit() error
	Show()
}

// RegisterProvider registers a provider.Factory by name.
func RegisterProvider(name string, p Factory) {
	providersMutex.Lock()
	defer providersMutex.Unlock()
	if _, found := providers[name]; !found {
		// log.Flog.Infof("registered provider %s", name)
		providers[name] = p
	}
}

// GetProvider creates an instance of the named provider, or nil if
// the name is unknown.  The error return is only used if the named provider
// was known but failed to initialize.
func GetProvider(name string) (Provider, error) {
	if name == "" {
		name = "native"
	}
	providersMutex.Lock()
	defer providersMutex.Unlock()
	f, found := providers[name]
	if !found {
		return nil, fmt.Errorf("provider %s is not registered", name)
	}
	return f()
}

// ListProviders returns current supported providers.
func ListProviders() []types.Provider {
	providersMutex.Lock()
	defer providersMutex.Unlock()
	list := make([]types.Provider, 0)
	for p := range providers {
		list = append(list, types.Provider{
			Name: p,
		})
	}
	return list
}
