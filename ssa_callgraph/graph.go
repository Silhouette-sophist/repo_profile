package ssa_callgraph

import (
	"golang.org/x/tools/go/ssa"
)

type Node struct {
	Func *ssa.Function
	ID   string
	In   map[string]*Edge
	Out  map[string]*Edge
}

type Edge struct {
	CallerID string
	CalleeID string
	Site     ssa.CallInstruction
}

type Graph struct {
	NodeMap map[string]*Node
}

func (g *Graph) Delete(ID string) {
	node := g.NodeMap[ID]
	for _, edge := range node.Out {
		if _, ok := g.NodeMap[edge.CalleeID]; !ok {
			continue
		}
		delete(g.NodeMap[edge.CalleeID].In, node.ID)
	}
	for _, edge := range node.In {
		if _, ok := g.NodeMap[edge.CallerID]; !ok {
			continue
		}
		delete(g.NodeMap[edge.CallerID].Out, node.ID)
	}
	delete(g.NodeMap, ID)
}
