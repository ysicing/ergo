/*
 * Copyright (c) 2021 ysicing <i@ysicing.me>
 */

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/cmd/cloud"
	"github.com/ysicing/ergo/pkg/util/factory"
)

// NewCloudCommand 云服务商支持
func newCloudCommand(f factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cloud [flags]",
		Short: "cloud tools",
	}
	cmd.AddCommand(cloud.CvmCmd(f))
	cmd.AddCommand(cloud.DomainCmd(f))
	cmd.AddCommand(cloud.CRCmd(f))
	return cmd
}
