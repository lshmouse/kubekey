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

package install

import (
	"fmt"
	kubekeyapiv1alpha1 "github.com/kubesphere/kubekey/apis/kubekey/v1alpha1"
	"github.com/kubesphere/kubekey/pkg/connector"
	"github.com/kubesphere/kubekey/pkg/test"
	"github.com/kubesphere/kubekey/pkg/util"
	"github.com/kubesphere/kubekey/pkg/util/executor"
	"testing"
)

var logger = util.InitLogger(true)

type MockConnector struct {
	connections map[string]*MockConnection
}

func (connector *MockConnector) Connect(host kubekeyapiv1alpha1.HostCfg) (connector.Connection, error) {
	if val, ok := connector.connections[host.Address]; ok {
		logger.Infof("Connection exist for host: %s", host.Address)
		return val, nil
	}

	logger.Infof("New connection on host: %s", host.Address)
	operations := make([]string, 0)
	connection := &MockConnection{operations: operations}
	connector.connections[host.Address] = connection
	return connection, nil
}

type MockConnection struct {
	operations []string
}

func (mock *MockConnection) Exec(cmd string, host *kubekeyapiv1alpha1.HostCfg) (stdout string, err error) {
	logger.Infof("run cmd: %s on host: %s", cmd, host.Address)
	mock.operations = append(mock.operations, cmd)
	return "OK", nil
}

func (mock *MockConnection) Scp(src, dst string) error {
	return nil
}

func (mock *MockConnection) Close() {
}

func Test_install(t *testing.T) {
	cfg, objName := test.GenTestClusterCfg("Test_install", util.LocalIP())

	executor := executor.NewExecutor(&cfg.Spec, objName, logger, "", true, true, true, false, false, nil)

	connections := make(map[string]*MockConnection)
	executor.Connector = &MockConnector{connections: connections}

	executor.DownloadCommand = func(path, url string) string {
		// this is an extension point for downloading tools, for example users can set the timeout, proxy or retry under
		// some poor network environment. Or users even can choose another cli, it might be wget.
		// perhaps we should have a build-in download function instead of totally rely on the external one
		return fmt.Sprintf("curl -L -o %s %s", path, url)
	}

	mgr, err := executor.CreateManager()
	if err != nil {
		t.Errorf("Create executor manager failure: %s", err)
	}
	ExecTasks(mgr)
	test.AssertEqual(t, len(connections), 1)
	operations := connections[util.LocalIP()].operations
	test.AssertEqual(t, len(operations), 50)

	// first command
	test.AssertEqual(t, operations[0][20:27], "useradd")
	// TODO
}
