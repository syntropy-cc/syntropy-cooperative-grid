package node_test

import (
	"testing"

	node "github.com/syntropy-grid/syntropy-cooperative-grid/manager/interfaces/cli/node/src"
)

func TestNodeManager_ListNodes_Empty(t *testing.T) {
	nm := node.NewNodeManager()
	if err := nm.Initialize(); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	list, err := nm.ListNodes()
	if err != nil {
		t.Fatalf("ListNodes failed: %v", err)
	}
	if list == nil {
		t.Fatalf("expected non-nil NodeList")
	}
	if list.Total != 0 {
		t.Fatalf("expected total 0, got %d", list.Total)
	}
}

func TestNodeManager_GetNodeStatus_NotFound(t *testing.T) {
	nm := node.NewNodeManager()
	if err := nm.Initialize(); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	if _, err := nm.GetNodeStatus("does-not-exist"); err == nil {
		t.Fatalf("expected error for unknown nodeID")
	}
}

func TestNodeManager_GetNodeLogs_NodeNotFound(t *testing.T) {
	nm := node.NewNodeManager()
	if err := nm.Initialize(); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	logs, err := nm.GetNodeLogs("missing-node", &node.LogOptions{Lines: 10})
	if err == nil {
		t.Fatalf("expected error when fetching logs for unknown node, got logs: %+v", logs)
	}
}
