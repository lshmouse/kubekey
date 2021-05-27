/*
Copyright 2020 The KubeSphere Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"github.com/kubesphere/kubekey/pkg/cluster/ping"
	"github.com/kubesphere/kubekey/pkg/util"
	"github.com/spf13/cobra"
)

// pingCmd represents the ping command
var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "ping all nodes of your cluster with this command",
	RunE: func(cmd *cobra.Command, args []string) error {
		logger := util.InitLogger(opt.Verbose)
		return ping.PingCluster(opt.ClusterCfgFile, logger, opt.Verbose)
	},
}

func init() {
	rootCmd.AddCommand(pingCmd)
	pingCmd.Flags().StringVarP(&opt.ClusterCfgFile, "filename", "f", "", "Path to a cluster configuration file")
}
