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

package tmpl

import (
	"github.com/kubesphere/kubekey/pkg/test"
	"github.com/kubesphere/kubekey/pkg/util"
	"github.com/kubesphere/kubekey/pkg/util/executor"
	"strings"
	"testing"
)

func Test_backupScript(t *testing.T) {
	var logger = util.InitLogger(true)
	cfg, objName := test.GenTestClusterCfg("Test_backupScript", "192.168.0.1", "192.168.0.2", "192.168.0.3")

	executor := executor.NewExecutor(&cfg.Spec, objName, logger, "", true, true, true, false, false, nil)
	mgr, err := executor.CreateManager()
	if err != nil {
		t.Errorf("Create executor manager failure: %s", err)
	}
	output, err := EtcdBackupScript(mgr, &mgr.EtcdNodes[0])
	if err != nil {
		t.Errorf("Create executor manager failure: %s", err)
	}

	var flag = true
	flag = strings.Contains(output, "ETCD_BACKUP_ENDPOINT='https://192.168.0.1:2379'")
	if !flag {
		t.Errorf("output: %s contains wrong ETCD_BACKUP_ENDPOINT", output)
	}

	flag = strings.Contains(output, "ENDPOINTS='https://192.168.0.1:2379,https://192.168.0.2:2379,https://192.168.0.3:2379'")
	if !flag {
		t.Errorf("output: %s contains wrong ENDPOINTS", output)
	}
}