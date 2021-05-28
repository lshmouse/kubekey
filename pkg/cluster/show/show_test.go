package show

import (
	kubekeyapiv1alpha1 "github.com/kubesphere/kubekey/apis/kubekey/v1alpha1"
	"github.com/kubesphere/kubekey/pkg/config"
	"github.com/kubesphere/kubekey/pkg/util"
	"github.com/kubesphere/kubekey/pkg/util/dialer"
	"github.com/kubesphere/kubekey/pkg/util/executor"
	"testing"
)

type MockConnector struct {
}

func (dialer *MockConnector) Connect(host kubekeyapiv1alpha1.HostCfg) (dialer.Connection, error) {
	return &MockConnection{}, nil
}

type MockConnection struct {
}

func (mock *MockConnection) Exec(cmd string, host *kubekeyapiv1alpha1.HostCfg) (stdout string, err error) {
	return "OK",nil
}

func (mock *MockConnection) Scp(src, dst string) error {
	return nil
}

func (mock *MockConnection) Close() {
}

func Test_showNodes(t *testing.T) {
	logger := util.InitLogger(true)

	cfg, objName, err := config.ParseClusterCfg("", "", "", false, logger)
	if err != nil {
		t.Errorf("Failed to download cluster config: %s", err)
	}

	executor := executor.NewExecutorWithOptions(&cfg.Spec, objName, logger, "", nil,
		executor.WithDebug(true), executor.WithSkipCheck(true),
		executor.WithSkipPullImages(true), executor.WithSkipFailTask(true),
		executor.WithConnector(&MockConnector{}))

	mgr, err := executor.CreateManager()
	if err != nil {
		t.Errorf("Create executor manager failure: %s", err)
	}
	PingNodes(mgr)
}
