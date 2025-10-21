package node_test

import (
	"testing"

	node "github.com/syntropy-grid/syntropy-cooperative-grid/manager/interfaces/cli/node/src"
)

func TestNodeManager_CreateNode_InvalidUbuntuVersion(t *testing.T) {
	nm := node.NewNodeManager()
	if err := nm.Initialize(); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	_, err := nm.CreateNode(&node.CreateOptions{UbuntuVersion: "19.04"})
	if err == nil {
		t.Fatalf("expected error for invalid Ubuntu version")
	}
}

func TestNodeManager_UpdateConfiguration_NotRunning(t *testing.T) {
	nm := node.NewNodeManager()
	cfg := nm.GetConfiguration()
	if err := nm.UpdateConfiguration(cfg); err == nil {
		t.Fatalf("expected error when updating config while not running")
	}
}

