package node_test

import (
	"testing"

	node "github.com/syntropy-grid/syntropy-cooperative-grid/manager/interfaces/cli/node/src"
)

func TestNodeManager_Initialize(t *testing.T) {
	nm := node.NewNodeManager()
	if err := nm.Initialize(); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}
	if !nm.IsRunning() {
		t.Fatalf("expected IsRunning true")
	}
}
