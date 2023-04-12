// Copyright (c) 2020-2023 ysicing(ysicing@ysicing.cloud) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Y PUBLIC LICENSE 1.0 (YPL 1.0)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package factory

import (
	"github.com/ysicing/ergo/internal/pkg/util/log"
)

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
