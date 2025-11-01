package ssa_callgraph

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCallSite(t *testing.T) {
	prog := Program{}
	err := prog.Load(InitProgramArgs{
		Path:      "./",
		Algorithm: "vta",
	})
	assert.Nil(t, err)
	g := prog.Graph
	t.Logf("g = %v", g)
}
