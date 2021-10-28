// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package factory

import "github.com/ergoapi/log"

// Factory is the main interface for various client creations
type Factory interface {
	// GetLog retrieves the log instance
	GetLog() log.Logger
}

// DefaultFactoryImpl is the default factory implementation
type DefaultFactoryImpl struct{}

// DefaultFactory returns the default factory implementation
func DefaultFactory() Factory {
	return &DefaultFactoryImpl{}
}

// GetLog implements interface
func (f *DefaultFactoryImpl) GetLog() log.Logger {
	return log.GetInstance()
}
