package install

import (
	"fmt"
	kubekeyapiv1alpha1 "github.com/kubesphere/kubekey/apis/kubekey/v1alpha1"
	"github.com/kubesphere/kubekey/pkg/util"
	"github.com/kubesphere/kubekey/pkg/util/dialer"
	"github.com/kubesphere/kubekey/pkg/util/executor"
	"runtime"
	"testing"
)

var logger = util.InitLogger(true)

type MockConnector struct {
}

func (dialer *MockConnector) Connect(host kubekeyapiv1alpha1.HostCfg) (dialer.Connection, error) {
	return &MockConnection{}, nil
}

type MockConnection struct {
}

func (mock *MockConnection) Exec(cmd string, host *kubekeyapiv1alpha1.HostCfg) (stdout string, err error) {
	logger.Info("run cmd: %s on host: %s", cmd, host)
	return "OK",nil
}

func (mock *MockConnection) Scp(src, dst string) error {
	return nil
}

func (mock *MockConnection) Close() {
}

func GenTestClusterCfg(name string, nodeNum int) (*kubekeyapiv1alpha1.Cluster, string) {
	cfg := kubekeyapiv1alpha1.Cluster{}
	cfg.Spec.Hosts = append(cfg.Spec.Hosts, kubekeyapiv1alpha1.HostCfg{
		Name:            name,
		Address:         util.LocalIP(),
		InternalAddress: util.LocalIP(),
		Port:            kubekeyapiv1alpha1.DefaultSSHPort,
		User:            "user",
		Password:        "",
		PrivateKeyPath:  fmt.Sprintf("%s/.ssh/id_rsa", "/home/user"),
		Arch:            runtime.GOARCH,
	})

	cfg.Spec.RoleGroups = kubekeyapiv1alpha1.RoleGroups{
		Etcd:   []string{name},
		Master: []string{name},
		Worker: []string{name},
	}
	cfg.Spec.Kubernetes = kubekeyapiv1alpha1.Kubernetes{
		Version: kubekeyapiv1alpha1.DefaultKubeVersion,
	}
	return &cfg, name
}

func Test_install(t *testing.T) {
	cfg, objName := GenTestCfg("Test_install")

	executor := executor.NewExecutorWithOptions(&cfg.Spec, objName, logger, "", nil,
		executor.WithDebug(true), executor.WithSkipCheck(true),
		executor.WithSkipPullImages(true), executor.WithSkipFailTask(true),
		executor.WithConnector(&MockConnector{}))

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
}
