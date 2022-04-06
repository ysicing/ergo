package cmd

import (
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/pkg/cli/kubectl"
)

// KubectlCommand kubectl command.
func KubectlCommand() *cobra.Command {
	return kubectl.EmbedCommand()
}
