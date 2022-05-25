package kube

import (
	"context"
	"os"

	"github.com/ergoapi/util/color"
	"github.com/ergoapi/util/file"
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/common"
	"github.com/ysicing/ergo/internal/pkg/k3s/status"
	"github.com/ysicing/ergo/internal/pkg/k3s/status/top"
	"github.com/ysicing/ergo/pkg/util/log"
	"k8s.io/kubectl/pkg/util/templates"
)

func StatusCmd() *cobra.Command {
	var params = status.K8sStatusOption{}
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Display status",
		Long:  ``,
		PreRun: func(cmd *cobra.Command, args []string) {
			defaultArgs := os.Args
			if !file.CheckFileExists(common.GetDefaultKubeConfig()) {
				log.Flog.Warnf("not found cluster. just run %s init cluster", color.SGreen("%s init", defaultArgs[0]))
				os.Exit(0)
			}
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			collector, err := status.NewK8sStatusCollector(params)
			if err != nil {
				return err
			}
			s, err := collector.Status(context.Background())
			// Report the most recent status even if an error occurred.
			s.Format()
			if err != nil {
				log.Flog.Fatalf("Unable to determine status:  %s", err)
			}
			return err
		},
	}
	cmd.Flags().StringVarP(&params.KubeConfig, "kubeconfig", "c", "", "Kubernetes configuration file")
	cmd.Flags().BoolVar(&params.Wait, "wait", false, "Wait for status to report success (no errors and warnings)")
	cmd.Flags().DurationVar(&params.WaitDuration, "wait-duration", common.StatusWaitDuration, "Maximum time to wait for status")
	cmd.Flags().BoolVar(&params.IgnoreWarnings, "ignore-warnings", false, "Ignore warnings when waiting for status to report success")
	cmd.Flags().StringVarP(&params.ListOutput, "output", "o", "", "prints the output in the specified format. Allowed values: table, json, yaml (default table)")
	cmd.AddCommand(topNodeCmd())
	return cmd
}

var (
	KRNodeExample = templates.Examples(`
	ergo kube status node
	`)
)

func topNodeCmd() *cobra.Command {
	o := top.NodeOption{}
	nodeCmd := &cobra.Command{
		Use:                   "node",
		DisableFlagsInUseLine: true,
		Short:                 "node provides an overview of the node",
		Aliases:               []string{"nodes", "no"},
		Example:               KRNodeExample,
		Run: func(cmd *cobra.Command, args []string) {
			o.Validate()
			o.RunResourceNode()
		},
	}
	nodeCmd.PersistentFlags().StringVarP(&o.KubeCtx, "context", "", "", "context to use for Kubernetes config")
	nodeCmd.PersistentFlags().StringVarP(&o.KubeConfig, "kubeconfig", "", "", "kubeconfig file to use for Kubernetes config")
	nodeCmd.PersistentFlags().StringVarP(&o.Output, "output", "o", "", "prints the output in the specified format. Allowed values: table, json, yaml (default table)")
	nodeCmd.PersistentFlags().StringVarP(&o.Selector, "selector", "l", "", "Selector (label query) to filter on, supports '=', '==', and '!='.(e.g. -l key1=value1,key2=value2)")
	nodeCmd.PersistentFlags().StringVarP(&o.SortBy, "sortBy", "s", "cpu", "sort by cpu or memory")
	return nodeCmd
}
