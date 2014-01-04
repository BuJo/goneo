package sgi

import (
	"fmt"
)

type vf2State struct {
	candidates    []vf2Match
	query, target Graph

	queryPath  []int
	targetPath []int

	// q ~> t
	mapping map[int]int

	// Fsem
	isSemanticallyFeasable SemFeasFunc
}
type vf2Match struct {
	query, target int
}

func (state *vf2State) GetGraph() Graph    { return state.target }
func (state *vf2State) GetSubgraph() Graph { return state.query }

func (state *vf2State) NextPair() (int, int) {
	if len(state.candidates) == 0 {
		return NULL_NODE, NULL_NODE
	}

	candidate := state.candidates[len(state.candidates)-1]

	state.candidates = state.candidates[0 : len(state.candidates)-1]

	return candidate.query, candidate.target
}

func (state *vf2State) BackTrack() {
	if len(state.queryPath) == 0 || state.IsGoal() {
		state.clearMapping()
		return
	}

	if state.lastQueryNodeMapped() {
		return
	}

	state.clearMapping()

	for i, q := range state.queryPath {
		state.mapping[q] = state.targetPath[i]
	}
}

func (state *vf2State) lastQueryNodeMapped() bool {
	queryNode := state.queryPath[len(state.queryPath)-1]
	queryNeighbours := state.query.Relations(queryNode)

	for _, n := range queryNeighbours {
		if _, ok := state.mapping[n]; !ok {
			return false
		}

	}

	return true
}

func (state *vf2State) clearMapping() {
	for k, _ := range state.mapping {
		delete(state.mapping, k)
	}
}

func (state *vf2State) IsFeasablePair(queryNode, targetNode int) bool {

	fSyn := state.isSyntacticallyFeasable(queryNode, targetNode)

	fSem := true

	if fSyn && len(state.queryPath) > 0 {
		fmt.Printf("(%d,%d)~>(%d,%d) is syntactically feasable, sem: ", state.queryPath[len(state.queryPath)-1], state.targetPath[len(state.targetPath)-1], queryNode, targetNode)

		fSem = state.isSemanticallyFeasable(state, state.queryPath[len(state.queryPath)-1], state.targetPath[len(state.targetPath)-1], queryNode, targetNode)
	}

	return fSyn && fSem
}

func (state *vf2State) isSyntacticallyFeasable(queryNode, targetNode int) bool {

	// Already mapped
	if _, ok := state.mapping[queryNode]; ok {
		return false
	}
	for _, t := range state.mapping {
		if t == targetNode {
			return false
		}
	}

	// match neighbour counts
	targetNeighbours := state.target.Relations(targetNode)
	queryNeighbours := state.query.Relations(queryNode)

	if len(queryNeighbours) > len(targetNeighbours) {
		return false
	}

	// TODO: more tests if queryNode matches targetNode

	// match edges
	if len(state.queryPath) == 0 {
		return true
	}

	for i, q := range state.queryPath {

		// match edges in query to target
		if state.query.Contains(q, queryNode) {
			if !state.target.Contains(state.targetPath[i], targetNode) {
				return false
			}
		}

		// TODO: more test for edge compatibility
	}

	return true
}

func (state *vf2State) isFeasableCandidate(queryNode, targetNode int) bool {
	// Test: not already visited
	for q, _ := range state.mapping {
		if q == queryNode || state.mapping[queryNode] == targetNode {
			return false
		}
	}
	return true
}

func (state *vf2State) IsGoal() bool {
	return len(state.mapping) == state.query.Order()
}
func (state *vf2State) IsDead() bool {
	return state.query.Order() > state.target.Order()
}

func (state *vf2State) GetMapping() map[int]int {
	return state.mapping
}

func (state *vf2State) NextState(queryNode, targetNode int) State {
	next := new(vf2State)

	next.mapping = make(map[int]int, state.query.Order())
	next.query = state.query
	next.target = state.target

	next.queryPath = make([]int, 0, state.query.Order())
	next.targetPath = make([]int, 0, state.target.Order())

	copy(next.queryPath, state.queryPath)
	copy(next.targetPath, state.targetPath)

	for q, t := range state.mapping {
		next.mapping[q] = t
	}

	next.mapping[queryNode] = targetNode
	next.queryPath = append(next.queryPath, queryNode)
	next.targetPath = append(next.targetPath, targetNode)

	next.candidates = make([]vf2Match, 0, next.query.Order())
	next.loadCandidates(queryNode, targetNode)

	next.isSemanticallyFeasable = state.isSemanticallyFeasable

	return next
}

func (state *vf2State) String() string {
	str := ""

	for i, q := range state.queryPath {
		str += fmt.Sprintf("->(%d~>%d)", q, state.targetPath[i])
	}

	return str + fmt.Sprintf(":%o", state.mapping)
}

func (state *vf2State) loadRootCandidates() {
	for q := 0; q < state.query.Order(); q += 1 {
		for t := 0; t < state.target.Order(); t += 1 {
			state.candidates = append(state.candidates, vf2Match{q, t})
		}
	}

	fmt.Println("loaded new candidates: ", state.candidates)
}
func (state *vf2State) loadCandidates(queryNode, targetNode int) {
	targetNeighbours := state.target.Relations(targetNode)
	queryNeighbours := state.query.Relations(queryNode)

	for _, q := range queryNeighbours {
		for _, t := range targetNeighbours {
			if state.isFeasableCandidate(q, t) {
				state.candidates = append(state.candidates, vf2Match{q, t})
			}
		}
	}

	fmt.Println("loaded new candidates: ", state.candidates)
}

func newVF2State(query, target Graph, fsem SemFeasFunc) State {
	state := new(vf2State)

	state.mapping = make(map[int]int, query.Order())
	state.query = query
	state.target = target

	state.candidates = make([]vf2Match, 0, query.Order())

	state.queryPath = make([]int, 0, query.Order())
	state.targetPath = make([]int, 0, target.Order())

	if fsem != nil {
		fmt.Println("custom Fsem")
		state.isSemanticallyFeasable = fsem
	} else {
		fmt.Println("always true Fsem")
		state.isSemanticallyFeasable = func(s State, a, b, c, d int) bool { return true }
	}

	state.loadRootCandidates()

	return state
}

func FindVF2SubgraphIsomorphism(query, target Graph, fsem SemFeasFunc) map[int]int {
	state := newVF2State(query, target, fsem)

	return FindIsomorphism(state)
}
