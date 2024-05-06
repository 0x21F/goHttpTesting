package models

import "container/list"

type Graph struct {
	Nodes map[uint]*GraphNode
}

type GraphNode struct {
	Id    uint
	Edges list.List
}

type Edge struct {
	u, v uint
}

func (g Graph) HasEdge(u, v uint) bool {
	if !g.verify(u, v) {
		return false
	}

	return false
}

func (g *Graph) AddEdge(u, v uint) bool {
	if !g.verify(u, v) {
		return false
	}

	return true
}

func (g *Graph) DelEdge(u, v uint) bool {
	if !g.verify(u, v) {
		return false
	}

	return true
}

func (g Graph) verify(u, v uint) bool {
	return g.Nodes[u] != nil || g.Nodes[v] != nil
}
