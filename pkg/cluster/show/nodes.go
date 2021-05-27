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

package show

import (
	kubekeyapiv1alpha1 "github.com/kubesphere/kubekey/apis/kubekey/v1alpha1"
	"github.com/kubesphere/kubekey/pkg/util/manager"
	"github.com/mitchellh/mapstructure"
	"github.com/modood/table"
	"strings"
)

// PingResults defines the items to be checked.
type NodeResults struct {
	Name    string `table:"name"`
	Address string `table:"address"`
	Time    string `table:"time"`
	Error   string `table:"error"`
	Etcd    bool   `table:"etcd"`
	Master  bool   `table:"master"`
	Worker  bool   `table:"worker"`
}

var (
	// CheckResults is used to save save check results.
	CheckResults = make(map[string]interface{})
)

func PingNode(mgr *manager.Manager, node *kubekeyapiv1alpha1.HostCfg) error {
	var results = make(map[string]interface{})
	results["name"] = node.Name
	results["address"] = node.Address
	results["etcd"] = node.IsEtcd
	results["master"] = node.IsMaster
	results["worker"] = node.IsWorker
	output, err := mgr.Runner.ExecuteCmd("echo \"health\"", 0, false)
	if err != nil {
		results["time"] = ""
		results["error"] = err.Error()
	} else {
		results["time"] = strings.TrimSpace(output)
		results["error"] = ""
	}
	CheckResults[node.Name] = results
	return nil
}

func PingNodeConfirm(mgr *manager.Manager) error {
	var results []NodeResults
	for node := range CheckResults {
		var result NodeResults
		_ = mapstructure.Decode(CheckResults[node], &result)
		results = append(results, result)
	}
	table.OutputA(results)
	return nil
}

func PingNodes(mgr *manager.Manager) error {
	if err := mgr.RunTaskOnAllNodes(PingNode, true); err != nil {
		return err
	}
	PingNodeConfirm(mgr)
	return nil
}
