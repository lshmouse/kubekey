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

package test

import (
	"fmt"
	kubekeyapiv1alpha1 "github.com/kubesphere/kubekey/apis/kubekey/v1alpha1"
	"github.com/kubesphere/kubekey/pkg/util"
	"runtime"
	"testing"
)

func AssertEqual(t *testing.T, actual interface{}, expected interface{}) {
	if actual != expected {
		t.Fatalf("actual: %s != expected: %s", actual, expected)
	}
}

func GenTestClusterCfg(name string, hosts ... string) (*kubekeyapiv1alpha1.Cluster, string) {
	cfg := kubekeyapiv1alpha1.Cluster{}
	for _, host := range hosts {
		cfg.Spec.Hosts = append(cfg.Spec.Hosts, kubekeyapiv1alpha1.HostCfg{
			Name:            host,
			Address:         host,
			InternalAddress: host,
			Port:            kubekeyapiv1alpha1.DefaultSSHPort,
			User:            "user",
			Password:        "",
			PrivateKeyPath:  fmt.Sprintf("%s/.ssh/id_rsa", "/home/user"),
			Arch:            runtime.GOARCH,
		})
	}
	master := hosts[:util.MinOf(3, len(hosts))]
	cfg.Spec.RoleGroups = kubekeyapiv1alpha1.RoleGroups{
		Etcd:   master,
		Master: master,
		Worker: hosts,
	}
	if len(master) >= 3 {
		cfg.Spec.ControlPlaneEndpoint.Address = name + ".lb.kubekey.com"
	}

	cfg.Spec.Kubernetes = kubekeyapiv1alpha1.Kubernetes{
		Version: kubekeyapiv1alpha1.DefaultKubeVersion,
	}
	return &cfg, name
}
