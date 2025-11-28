package ssa_callgraph

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCallSite(t *testing.T) {
	ctx := context.Background()
	prog := Program{}
	args := InitProgramArgs{
		Path:      "./../../",
		Algorithm: "vta",
	}
	err := prog.Load(ctx, args)
	assert.Nil(t, err)
	g := prog.Graph
	t.Logf("g = %v", g)
}
