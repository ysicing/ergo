// Copyright (c) 2020-2023 ysicing(ysicing@ysicing.cloud) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Y PUBLIC LICENSE 1.0 (YPL 1.0)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/pkg/cli/kubectl"
)

// KubectlCommand kubectl command.
func KubectlCommand() *cobra.Command {
	return kubectl.EmbedCommand()
}
