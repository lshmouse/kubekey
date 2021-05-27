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

package ping

import (
	"fmt"
	kubekeyapiv1alpha1 "github.com/kubesphere/kubekey/apis/kubekey/v1alpha1"
	"github.com/kubesphere/kubekey/pkg/config"
	"github.com/kubesphere/kubekey/pkg/util"
	"github.com/kubesphere/kubekey/pkg/util/executor"
	"github.com/kubesphere/kubekey/pkg/util/manager"
	"github.com/mitchellh/mapstructure"
	"github.com/modood/table"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
)

// PingResults defines the items to be checked.
type Results struct {
	Name  string `table:"name"`
	Time  string `table:"time"`
	Error string `table:"error"`
}

var (
	// CheckResults is used to save save check results.
	CheckResults = make(map[string]interface{})
)

func PingNode(mgr *manager.Manager, node *kubekeyapiv1alpha1.HostCfg) error {
	var results = make(map[string]interface{})
	results["name"] = node.Name
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

func PingConfirm(mgr *manager.Manager) error {
	var results []Results
	for node := range CheckResults {
		var result Results
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
	PingConfirm(mgr)
	return nil
}

func PingCluster(clusterCfgFile string, logger *log.Logger, verbose bool) error {
	currentDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return errors.Wrap(err, "Failed to get current dir")
	}
	if err := util.CreateDir(fmt.Sprintf("%s/kubekey", currentDir)); err != nil {
		return errors.Wrap(err, "Failed to create work dir")
	}

	cfg, objName, err := config.ParseClusterCfg(clusterCfgFile, "", "", false, logger)
	if err != nil {
		return errors.Wrap(err, "Failed to download cluster config")
	}

	executorInstance := executor.NewExecutor(&cfg.Spec, objName, logger, "", verbose, true, true, false, false, nil)
	executorInstance.DownloadCommand = func(path, url string) string {
		// this is an extension point for downloading tools, for example users can set the timeout, proxy or retry under
		// some poor network environment. Or users even can choose another cli, it might be wget.
		// perhaps we should have a build-in download function instead of totally rely on the external one
		return fmt.Sprintf("curl -L -o %s %s", path, url)
	}
	return Execute(executorInstance)
	//return Execute(executor.NewExecutor(&cfg.Spec, objName, logger, "", verbose, true, skipPullImages, false, false, nil))
}

func ExecTasks(mgr *manager.Manager) error {
	pingTasks := []manager.Task{
		{Task: PingNodes, ErrMsg: "Failed to ping all nodes of the clusters"},
	}
	for _, step := range pingTasks {
		if err := step.Run(mgr); err != nil {
			return errors.Wrap(err, step.ErrMsg)
		}
	}
	return nil
}

func Execute(executor *executor.Executor) error {
	mgr, err := executor.CreateManager()
	if err != nil {
		return err
	}
	return ExecTasks(mgr)
}
