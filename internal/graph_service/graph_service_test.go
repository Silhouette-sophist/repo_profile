package graph_service

import (
	"context"
	"path/filepath"
	"testing"
)

func TestTransferGraph(t *testing.T) {
	ctx := context.Background()
	curRepoPath, err := filepath.Abs("./../..")
	if err != nil {
		t.Errorf("abs path err: %v", err)
	}
	TransferGraph(ctx, curRepoPath)
}
