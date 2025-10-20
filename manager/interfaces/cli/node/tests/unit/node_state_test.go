package node_test

import (
	"os"
	"testing"
	"time"

	node "github.com/syntropy-grid/syntropy-cooperative-grid/manager/interfaces/cli/node/src"
)

func TestNodeStateManager_Initialize_CreatesStateDir(t *testing.T) {
	nsm := node.NewNodeStateManager(node.NewLogger())
	if err := nsm.Initialize(); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}
	// Verifica que o diretório base foi criado
	cfgDir := os.ExpandEnv("~/.syntropy")
	if _, err := os.Stat(cfgDir); err != nil && os.IsNotExist(err) {
		t.Fatalf("expected state directory to exist: %s", cfgDir)
	}
}

func TestNodeStateManager_SaveAndLoadState_NoPanic(t *testing.T) {
	nsm := node.NewNodeStateManager(node.NewLogger())
	if err := nsm.Initialize(); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}
	// Apenas exercita caminhos de persistência sem depender de tipos internos
	if err := nsm.SaveState(); err != nil {
		t.Fatalf("SaveState failed: %v", err)
	}
	if err := nsm.LoadState(); err != nil {
		t.Fatalf("LoadState failed: %v", err)
	}
	_ = time.Now()
}
