// Copyright (c) 2020-2023 ysicing(ysicing@ysicing.cloud) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Y PUBLIC LICENSE 1.0 (YPL 1.0)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package kubectl

import (
	"flag"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/kubectl/pkg/cmd"
	"os"
)

// Main kubectl main function.
// Borrowed from https://github.com/kubernetes/kubernetes/blob/master/cmd/kubectl/kubectl.go.
func Main() {

	pflag.CommandLine.SetNormalizeFunc(cliflag.WordSepNormalizeFunc)
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)

	if err := EmbedCommand().Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

// EmbedCommand Used to embed the kubectl command.
func EmbedCommand() *cobra.Command {
	c := cmd.NewDefaultKubectlCommand()
	c.Short = "Kubectl controls the Kubernetes cluster manager"

	return c
}
