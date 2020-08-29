// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package cmd

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	mid "github.com/ysicing/ginmid"
	"runtime"
)

var webport int

var dashboardCmd = &cobra.Command{
	Use:     "dashboard",
	Aliases: []string{"ui"},
	Short:   "Simple web UIs",
	Run: func(cmd *cobra.Command, args []string) {
		simpleweb := gin.Default()
		simpleweb.Use(mid.RequestID(), mid.PromMiddleware(nil))
		simpleweb.GET("/metrics", mid.PromHandler(promhttp.Handler()))
		simpleweb.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"os":           runtime.GOOS,
				"arch":         runtime.GOARCH,
				"go":           runtime.Version(),
				"x-request-id": mid.GetRequestID(c),
			})
		})
		if webport < 80 || webport > 65530 {
			webport = 12306
		}
		simpleweb.Run(fmt.Sprintf("0.0.0.0:%v", webport))
	},
}

func init() {
	rootCmd.AddCommand(dashboardCmd)
	dashboardCmd.PersistentFlags().IntVar(&webport, "port", 12306, "端口")
}
