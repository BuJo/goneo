package sgi

import (
	"fmt"
)

var NULL_NODE int = -1

type Graph interface {
	Order() int

	// contains link from node a to b
	Contains(a, b int) bool

	Successors(n int) []int
	Predecessors(n int) []int
	Relations(n int) []int
}

type Edge interface {
	Type() string
}

type State interface {
	GetGraph() Graph
	GetSubgraph() Graph

	NextPair() (int, int)

	BackTrack()

	IsFeasablePair(queryNode, targetNode int) bool

	IsGoal() bool
	IsDead() bool

	GetMapping() map[int]int // (partial) mapping of the graphs

	NextState(queryNode, targetNode int) State

	String() string
}

type StateFunc func(graph, graph2 Graph) State

func FindIsomorphism(graph, graph2 Graph, makeInitialState StateFunc) map[int]int {

	state := makeInitialState(graph, graph2)

	isoMapping := make(map[int]int, 0)

	if match(state, isoMapping) {
		return isoMapping
	} else {
		return make(map[int]int, 0)
	}
}

func match(state State, isoMapping map[int]int) bool {
	fmt.Println("Start Match")

	if state.IsDead() {
		fmt.Println("Match is dead")
		return false
	}

	if state.IsGoal() {
		fmt.Println("Match Goal reached")
		for k, v := range state.GetMapping() {
			isoMapping[k] = v
		}

		return true
	}

	n1, n2 := state.NextPair()
	found := false

	fmt.Println("Match starting pair: ", n1, n2)
	fmt.Println("State:", state)

	for ; !found && n1 != NULL_NODE; n1, n2 = state.NextPair() {

		if state.IsFeasablePair(n1, n2) {
			fmt.Println("are feasable: ", n1, n2)

			next := state.NextState(n1, n2)

			found = match(next, isoMapping)

			next.BackTrack()
		}
	}

	fmt.Println("End Match, found? ", found)

	return found
}
