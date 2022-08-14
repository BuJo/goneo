package sgi

import (
	"testing"
)

type dbGraphMock struct {
	nodes []int    // nodes with integer values
	edges [][]bool // array marking edges
}

func (g *dbGraphMock) Order() int             { return len(g.nodes) }
func (g *dbGraphMock) Contains(a, b int) bool { return g.edges[a][b] }
func (g *dbGraphMock) Successors(a int) []int {
	ids := make([]int, 0)
	for i, ok := range g.edges[a] {
		if ok {
			ids = append(ids, i)
		}
	}
	return ids
}
func (g *dbGraphMock) Predecessors(a int) []int {
	ids := make([]int, 0)
	for rownum, row := range g.edges {
		if row[a] {
			ids = append(ids, rownum)
		}
	}
	return ids
}
func (g *dbGraphMock) Relations(a int) []int {
	return append(g.Successors(a), g.Predecessors(a)...)
}

func isSemanticallyFeasable(state State, fromQueryNode, fromTargetNode, toQueryNode, toTargetNode int) bool {
	return true
}

var testgraph = &dbGraphMock{
	nodes: []int{100, 102, 103},
	edges: [][]bool{
		{false, true, true},
		{true, false, false},
		{true, false, false},
	},
}

func ExampleFindVF2SubgraphIsomorphism() {
	var db, subgraph *dbGraphMock

	FindVF2SubgraphIsomorphism(subgraph, db, func(state State, fromQueryNode, fromTargetNode, toQueryNode, toTargetNode int) bool {

		return isSemanticallyFeasable(state, fromQueryNode, fromTargetNode, toQueryNode, toTargetNode)
	})
}

func TestFailingIsomorphism(t *testing.T) {
	subgraph := &dbGraphMock{
		nodes: []int{100, 102},
		edges: [][]bool{
			{false, false},
			{false, false},
		},
	}

	mappings := FindVF2SubgraphIsomorphism(subgraph, testgraph, func(state State, fromQueryNode, fromTargetNode, toQueryNode, toTargetNode int) bool {

		return true
	})

	found := false
	for _, mapping := range mappings {
		// Each mapping contains the query and target node ids
		for q, t := range mapping {
			if subgraph.nodes[q] == 102 && testgraph.nodes[t] == 102 {
				found = true
			}
		}
	}

	if found {
		t.Errorf("Should not have found isomorphism")
	}
}

func TestSemanticIsomorphism(t *testing.T) {
	subgraph := &dbGraphMock{
		nodes: []int{100, 102},
		edges: [][]bool{
			{false, true},
			{true, false},
		},
	}

	mappings := FindVF2SubgraphIsomorphism(subgraph, testgraph, func(state State, fromQueryNode, fromTargetNode, toQueryNode, toTargetNode int) bool {

		// Compare node values to determine equality
		nodeIdOk := subgraph.nodes[toQueryNode] == testgraph.nodes[toTargetNode]

		// And don't care about edge compatibility
		return nodeIdOk
	})

	found := false
	for _, mapping := range mappings {
		// Each mapping contains the query and target node ids
		for q := range mapping {
			if subgraph.nodes[q] == 100 {
				found = true
			}
		}
	}

	if !found {
		t.Errorf("Should find isomorphism")
	}
}

func TestSemanticallyFailingIsomorphism(t *testing.T) {
	subgraph := &dbGraphMock{
		nodes: []int{100, 200},
		edges: [][]bool{
			{false, true},
			{false, false},
		},
	}

	mappings := FindVF2SubgraphIsomorphism(subgraph, testgraph, func(state State, fromQueryNode, fromTargetNode, toQueryNode, toTargetNode int) bool {

		// Compare node values to determine equality
		nodeIdOk := subgraph.nodes[toQueryNode] == testgraph.nodes[toTargetNode]

		// And don't care about edge compatibility
		return nodeIdOk
	})

	found := false
	for _, mapping := range mappings {
		// Each mapping contains the query and target node ids
		for q := range mapping {
			if subgraph.nodes[q] == 100 {
				found = true
			}
		}
	}

	if found {
		t.Errorf("Should not find isomorphism")
	}
}
