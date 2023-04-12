// Copyright (c) 2020-2023 ysicing(ysicing@ysicing.cloud) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Y PUBLIC LICENSE 1.0 (YPL 1.0)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package flags

import (
	flag "github.com/spf13/pflag"
	"github.com/ysicing/ergo/common"
)

// GlobalFlags is the flags that contains the global flags
type GlobalFlags struct {
	Debug      bool
	Silent     bool
	ConfigPath string
	Vars       []string
	Flags      *flag.FlagSet
}

// SetGlobalFlags applies the global flags
func SetGlobalFlags(flags *flag.FlagSet) *GlobalFlags {
	globalFlags := &GlobalFlags{
		Vars:  []string{},
		Flags: flags,
	}
	flags.BoolVar(&globalFlags.Debug, "debug", false, "Prints the stack trace if an error occurs")
	flags.BoolVar(&globalFlags.Silent, "silent", false, "Run in silent mode and prevents any ergo log output except panics & fatals")
	flags.StringVar(&globalFlags.ConfigPath, "config", common.GetDefaultErgoCfg(), "The ergo config file to use")

	return globalFlags
}
