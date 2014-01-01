package sgi

import (
	"fmt"
	"goneo"
	"testing"
)

func TestFoo(t *testing.T) {
	graph, subgraph := setUp(), setUp()

	n1, n2, n3 := graph.db.NewNode(), graph.db.NewNode(), graph.db.NewNode()
	n1.RelateTo(n2, "A")
	n2.RelateTo(n3, "A")
	n1.RelateTo(n3, "A")

	s1, s2 := subgraph.db.NewNode(), subgraph.db.NewNode()
	s1.RelateTo(s2, "A")

	isoMap := FindVF2SubgraphIsomorphism(subgraph, graph)
	if isoMap == nil || len(isoMap) == 0 {
		t.Fail()
	}
	fmt.Println("map:", isoMap)
}

type dbGraph struct {
	db *goneo.DatabaseService
}

func (g *dbGraph) Order() int {
	return len(g.db.GetAllNodes())
}

func (g *dbGraph) GetEdges() []Edge {
	rels := g.db.GetAllRelations()
	edges := make([]Edge, len(rels), len(rels))

	for i, r := range rels {
		edges[i] = r
	}
	return edges
}
func (g *dbGraph) Contains(a, b int) bool {
	//fmt.Println("gr:Contains:",a,b)
	node, _ := g.db.GetNode(a)
	for _, rel := range node.Relations(goneo.Outgoing) {
		if rel.End.Id() == b {
			return true
		}
	}
	return false
}
func (g *dbGraph) Successors(a int) []int {

	node, _ := g.db.GetNode(a)
	ids := make([]int, 0)
	for _, rel := range node.Relations(goneo.Outgoing) {
		ids = append(ids, rel.End.Id())
	}
	//fmt.Println("gr:succ:",a,ids)
	return ids
}
func (g *dbGraph) Predecessors(a int) []int {

	node, _ := g.db.GetNode(a)
	ids := make([]int, 0)
	for _, rel := range node.Relations(goneo.Incoming) {
		ids = append(ids, rel.Start.Id())
	}
	//fmt.Println("gr:Pred:",a,ids)
	return ids
}
func (g *dbGraph) Relations(a int) []int {

	node, _ := g.db.GetNode(a)
	ids := make([]int, 0)
	for _, rel := range node.Relations(goneo.Both) {
		if rel.Start.Id() == a {
			ids = append(ids, rel.End.Id())
		} else {
			ids = append(ids, rel.Start.Id())
		}
	}
	fmt.Println("gr:Rel:", a, ids)
	return ids
}

func setUp() *dbGraph {
	graph := new(dbGraph)
	graph.db = goneo.NewTemporaryDb()
	return graph
}
