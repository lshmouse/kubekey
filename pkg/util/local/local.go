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

package local

import kubekeyapiv1alpha1 "github.com/kubesphere/kubekey/apis/kubekey/v1alpha1"

type LocalConnection struct {
}

func (c *LocalConnection) Exec(cmd string, host *kubekeyapiv1alpha1.HostCfg) (string, error) {
	return "OK", nil
}

func (c *LocalConnection) Scp(src, dst string) error {
	return nil
}

func (c *LocalConnection) Close() {
}