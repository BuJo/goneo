package sgi

import (
	"log"
	"reflect"
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

type SemFeasFunc func(state State, fromQueryNode, fromTargetNode, toQueryNode, toTargetNode int) bool

func FindIsomorphism(initialState State) []map[int]int {

	isoMappings := make([]map[int]int, 0, 0)

	match(initialState, &isoMappings)

	return isoMappings

}

func match(state State, isoMappings *[]map[int]int) {
	log.Print("Start Match")

	if state.IsGoal() {
		log.Print("Match Goal reached, len:", len(*isoMappings))
		isoMapping := make(map[int]int)

		for k, v := range state.GetMapping() {
			isoMapping[k] = v
		}

		if !alreadyInMappings(*isoMappings, isoMapping) {
			log.Print("not in mapping")
			*isoMappings = append(*isoMappings, isoMapping)
		} else {
			log.Print("in mapping")
		}

		return
	}

	if state.IsDead() {
		log.Print("Match is dead")
		return
	}

	n1, n2 := state.NextPair()

	for ; n1 != NULL_NODE; n1, n2 = state.NextPair() {
		log.Print("State:", state, " next pair: ", n1, n2)

		if state.IsFeasablePair(n1, n2) {
			log.Print("are feasable: ", n1, n2)

			next := state.NextState(n1, n2)

			match(next, isoMappings)

			next.BackTrack()
		}
	}
}

func alreadyInMappings(mappings []map[int]int, mapping map[int]int) bool {
	log.Print("in mapping? ", mappings, mapping)

	for _, m0 := range mappings {
		eq := reflect.DeepEqual(m0, mapping)
		if eq {
			return true
		}
	}

	return false
}
