package providers

import (
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
)

type Factory func() (Provider, error)

var (
	providersMutex sync.Mutex
	providers      = make(map[string]Factory)
)

type Provider interface {
	GetProviderName() string
}

func RegisterProvider(name string, p Factory) {
	providersMutex.Lock()
	defer providersMutex.Unlock()
	if _, found := providers[name]; !found {
		logrus.Debugf("registered provider %s", name)
		providers[name] = p
	}
}

func GetProvider(name string) (Provider, error) {
	providersMutex.Lock()
	defer providersMutex.Unlock()
	f, found := providers[name]
	if !found {
		return nil, fmt.Errorf("provider %s is not registered", name)
	}
	return f()
}
