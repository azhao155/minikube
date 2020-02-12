/*
Copyright 2020 The Kubernetes Authors All rights reserved.

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
	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/minikube/pkg/minikube/cluster"
	"k8s.io/minikube/pkg/minikube/config"
	"k8s.io/minikube/pkg/minikube/exit"
	"k8s.io/minikube/pkg/minikube/machine"
	"k8s.io/minikube/pkg/minikube/node"
	"k8s.io/minikube/pkg/minikube/out"
)

var nodeDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Deletes a node from a cluster.",
	Long:  "Deletes a node from a cluster.",
	Run: func(cmd *cobra.Command, args []string) {
		name := viper.GetString("name")
		if name == "" {
			exit.UsageT("name is required")
		}

		// Make sure it's not running
		api, err := machine.NewAPIClient()
		if err != nil {
			exit.WithError("creating api client", err)
		}

		cc, err := config.Load(viper.GetString(config.MachineProfile))
		if err != nil {
			exit.WithError("loading config", err)
		}

		n, _, err := node.Retrieve(cc, name)
		if err != nil {
			exit.WithError("retrieving node", err)
		}

		if cluster.IsHostRunning(api, name) {
			err := node.Stop(cc, n)
			if err != nil {
				glog.Warningf("Failed to stop node, will still try to delete")
			}
		}

		err = node.Delete(cc, name)
		if err != nil {
			out.FatalT("Failed to delete node {{.name}}", out.V{"name": name})
		}
	},
}

func init() {
	nodeDeleteCmd.Flags().String("name", "", "The name of the node to delete")
	nodeCmd.AddCommand(nodeDeleteCmd)
}
