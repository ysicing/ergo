package kube

import (
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/pkg/kube/resource"
	"github.com/ysicing/ergo/pkg/util/factory"
	"k8s.io/kubectl/pkg/util/templates"
)

var (
	KRNodeExample = templates.Examples(`
	ergo kube res node
	`)
	KRPodExample = templates.Examples(`
	ergo kube res pod
	ergo kube res pod -l app=my-nginx
	ergo kube res pod -l app=my-nginx -o json
	ergo kube res pod -l app=my-nginx -n default -o yaml
	`)
)

func KRCmd(f factory.Factory) *cobra.Command {
	o := resource.Option{}
	kr := &cobra.Command{
		Use:   "res",
		Short: "resource",
		Args:  cobra.NoArgs,
	}
	kr.PersistentFlags().StringVarP(&o.KubeCtx, "context", "", "", "context to use for Kubernetes config")
	kr.PersistentFlags().StringVarP(&o.KubeConfig, "kubeconfig", "", "", "kubeconfig file to use for Kubernetes config")
	kr.PersistentFlags().StringVarP(&o.Output, "output", "o", "", "prints the output in the specified format. Allowed values: table, json, yaml (default table)")
	kr.PersistentFlags().StringVarP(&o.Selector, "selector", "l", "", "Selector (label query) to filter on, supports '=', '==', and '!='.(e.g. -l key1=value1,key2=value2)")
	kr.PersistentFlags().StringVarP(&o.SortBy, "sortBy", "s", "cpu", "sort by cpu or memory")
	podCmd := &cobra.Command{
		Use:                   "pod [NAME | -l label]",
		Short:                 "pod provides an overview of the pod",
		DisableFlagsInUseLine: true,
		Example:               KRPodExample,
		Aliases:               []string{"pods", "po"},
		Run: func(cmd *cobra.Command, args []string) {
			o.Validate()
			o.RunResourcePod()
		},
	}
	podCmd.PersistentFlags().StringVarP(&o.LabelSelector, "label", "l", "", "Selector (label query) to filter on, supports '=', '==', and '!='.(e.g. -l key1=value1,key2=value2)")
	podCmd.PersistentFlags().StringVarP(&o.FieldSelector, "field", "f", "", "Selector (field query) to filter on, supports '=', '==', and '!='.(e.g. -f key1=value1,key2=value2)")
	podCmd.PersistentFlags().StringVarP(&o.Namespace, "namespace", "n", "", "only include rosource from this namespace")
	kr.AddCommand(podCmd)
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
	kr.AddCommand(nodeCmd)
	return kr
}
