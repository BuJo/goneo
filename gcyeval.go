package goneo

import (
	"bytes"
	"errors"
	"fmt"
	"goneo/gcy"
	"goneo/sgi"
	"text/tabwriter"
)

type (
	evalContext struct {
		vars map[string][]PropertyContainer

		subgraphNameMap    map[string]int
		subgraphRevNameMap map[int]string

		subgraph *dbGraph

		db *DatabaseService
	}

	Stringable interface {
		String() string
	}

	TabularData struct {
		line []map[string]Stringable
	}

	query   struct{ q *gcy.Query }
	match   struct{ m *gcy.Match }
	root    struct{ r *gcy.Root }
	returns struct{ r *gcy.Variable }
)

func (t *TabularData) String() string {
	if len(t.line) == 0 {
		return ""
	}
	b := new(bytes.Buffer)

	w := new(tabwriter.Writer)

	// Format in tab-separated columns with a tab stop of 8.
	w.Init(b, 0, 8, 0, '\t', 0)

	headers := make([]string, 0)

	for header, _ := range t.line[0] {
		headers = append(headers, header)
	}

	for _, header := range headers {
		fmt.Fprint(w, header+"\t")
	}
	fmt.Fprintln(w, "")

	for _, line := range t.line {
		for _, header := range headers {
			if _, ok := line[header]; ok && line[header] != nil {
				fmt.Fprint(w, line[header].String()+"\t")
			}
		}
		fmt.Fprintln(w, "")
	}
	w.Flush()

	return b.String()
}
func (t *TabularData) Len() int {
	return len(t.line)
}
func (t *TabularData) Columns() int {
	return len(t.line[0])
}

func (q *query) evaluate(ctx evalContext) *TabularData {

	ctx.vars = make(map[string][]PropertyContainer)
	ctx.subgraphNameMap = make(map[string]int)
	ctx.subgraphRevNameMap = make(map[int]string)

	for _, r := range q.q.Roots {
		(&root{r}).evaluate(ctx)
	}

	if q.q.Match != nil {
		(&match{q.q.Match}).evaluate(ctx)
	}

	table := &TabularData{}

	for _, r := range q.q.Returns {
		table = (&returns{r}).evaluate(ctx)
	}

	return table
}

// BUG(jo): db cannot ancode undirected graph
// BUG(jo): cannot encode more than one possible relation type

func (mm *match) evaluate(ctx evalContext) *TabularData {
	m := mm.m

	subgraph := NewTemporaryDb()

	for _, p := range m.Paths {
		var builder *PathBuilder
		for currentNode := p.Start; currentNode != nil; {
			var n *Node

			if currentNode.Name != "" {

				if id, ok := ctx.subgraphNameMap[currentNode.Name]; ok {
					n, _ = subgraph.GetNode(id)
					fmt.Println("tried to find node name ", currentNode.Name, ", found", n)
				}
			}

			if n == nil {
				n = subgraph.NewNode(currentNode.Labels...)
				ctx.subgraphNameMap[currentNode.Name] = n.Id()
				ctx.subgraphRevNameMap[n.Id()] = currentNode.Name

				fmt.Println("created new node: ", n, "(", currentNode, ")")

				for k, v := range currentNode.Props {
					n.SetProperty(k, v)
				}
			}

			if builder == nil {
				fmt.Println("first run, node: ", n, "(", currentNode, ")")
				builder = NewPathBuilder(n)
			} else {
				prevNode := builder.Last()
				fmt.Println("next run, ", prevNode, "->", n, "(", currentNode, ")")

				// TODO: utter crap, path is specific, has no "optional variants"
				// TODO: utter crap, db relations have no "optional variants"
				for _, typ := range currentNode.LeftRel.Types {
					var rel *Relation
					if currentNode.LeftRel.Direction == "->" {
						rel = prevNode.RelateTo(n, typ)
					} else {
						rel = n.RelateTo(prevNode, typ)
					}
					//fmt.Println("____ created new subgraph rel: ", rel, currentNode.LeftRel)

					builder = builder.Append(rel)
					break
				}
			}

			if currentNode.RightRel == nil {
				currentNode = nil
			} else {

				currentNode = currentNode.RightRel.RightNode
			}
		}
	}

	mappings := sgi.FindVF2SubgraphIsomorphism(&dbGraph{subgraph}, &dbGraph{ctx.db}, isSemanticallyFeasable)

	for _, mapping := range mappings {
		for q, t := range mapping {
			name := ctx.subgraphRevNameMap[q]
			node, _ := ctx.db.GetNode(t)

			if _, ok := ctx.vars[name]; !ok {
				ctx.vars[name] = make([]PropertyContainer, 1)
			}

			ctx.vars[name] = append(ctx.vars[name], node)
		}
	}

	return nil
}

func (rr *root) evaluate(ctx evalContext) *TabularData {
	r := rr.r
	_, ok := ctx.vars[r.Name]

	if !ok {
		ctx.vars[r.Name] = make([]PropertyContainer, 0)
	}

	if r.Typ == "node" {
		if len(r.IdVars) == 1 && r.IdVars[0] == -1 {
			for _, node := range ctx.db.GetAllNodes() {
				ctx.vars[r.Name] = append(ctx.vars[r.Name], node)
			}
		} else {
			for _, id := range r.IdVars {
				node, _ := ctx.db.GetNode(id)
				ctx.vars[r.Name] = append(ctx.vars[r.Name], node)
			}
		}
	} else {
		if len(r.IdVars) == 1 && r.IdVars[0] == -1 {
			for _, rel := range ctx.db.GetAllRelations() {
				ctx.vars[r.Name] = append(ctx.vars[r.Name], rel)
			}
		} else {
			for _, id := range r.IdVars {
				rel, _ := ctx.db.GetRelation(id)
				ctx.vars[r.Name] = append(ctx.vars[r.Name], rel)
			}
		}
	}

	return &TabularData{}
}
func (rr *returns) evaluate(ctx evalContext) *TabularData {
	r := rr.r
	table := &TabularData{}

	table.line = make([]map[string]Stringable, 0)

	for _, o := range ctx.vars[r.Name] {
		if o != nil {
			line := make(map[string]Stringable)

			line[r.Alias] = o

			table.line = append(table.line, line)
		}
	}

	fmt.Println("evaluating return, ", table.line)

	return table
}

//Evaluate a gcy query
//
// Example:
//
// 	start n=node(*) return n
func (db *DatabaseService) Evaluate(qry string) (*TabularData, error) {
	q, err := gcy.Parse("goneo", qry)
	if err != nil {
		return nil, err
	}
	table := (&query{q}).evaluate(evalContext{db: db})
	if table == nil {
		return nil, errors.New("Could not evaluate query")
	}

	return table, nil
}

// defines semantic feasability of the given state M(s) and n' m'
func isSemanticallyFeasable(state sgi.State, fromQueryNode, fromTargetNode, toQueryNode, toTargetNode int) bool {
	graph := state.GetGraph().(*dbGraph).db
	subgraph := state.GetSubgraph().(*dbGraph).db

	q2, _ := subgraph.GetNode(toQueryNode)
	t2, _ := graph.GetNode(toTargetNode)

	// Labels
	labelsOk := t2.HasLabel(q2.Labels()...)

	if fromQueryNode == sgi.NULL_NODE {
		fmt.Printf("queryNode: %s , targetNode: %s\n", q2, t2)
		return labelsOk
	}

	q1, _ := subgraph.GetNode(fromQueryNode)
	t1, _ := graph.GetNode(fromTargetNode)

	// query relation
	qRel, qDir := &Relation{}, Both
	for _, rel := range q1.Relations(Both) {
		if rel.Start.Id() == q1.Id() && rel.End.Id() == q2.Id() {
			qRel, qDir = rel, Outgoing
		} else if rel.End.Id() == q1.Id() && rel.Start.Id() == q2.Id() {
			qRel, qDir = rel, Incoming
		}
	}

	// target relation
	tRel, tDir := &Relation{}, Both
	for _, rel := range t1.Relations(Both) {
		if rel.Start.Id() == t1.Id() && rel.End.Id() == t2.Id() {
			tRel, tDir = rel, Outgoing
		} else if rel.End.Id() == t1.Id() && rel.Start.Id() == t2.Id() {
			tRel, tDir = rel, Incoming
		}
	}

	fmt.Printf("queryNodes: %s%s%s , targetNodes: %s%s%s | ", q1, qDir, q2, t1, tDir, t2)
	//fmt.Printf("rel: (%s~>%s), dir: (%s~>%s) \n", qRel.Type(), tRel.Type(), qDir, tDir)
	fmt.Println(qRel, tRel)

	return (qRel.Type() == "" || qRel.Type() == tRel.Type()) && (qDir == Both || qDir == tDir) && labelsOk
}

type dbGraph struct {
	db *DatabaseService
}

func (g *dbGraph) Order() int {
	return len(g.db.GetAllNodes())
}

func (g *dbGraph) Contains(a, b int) bool {
	//fmt.Println("gr:Contains:",a,b)
	node, _ := g.db.GetNode(a)
	for _, rel := range node.Relations(Both) {
		if rel.End.Id() == b || rel.Start.Id() == b {
			return true
		}
	}
	return false
}
func (g *dbGraph) Successors(a int) []int {

	node, _ := g.db.GetNode(a)
	ids := make([]int, 0)
	for _, rel := range node.Relations(Outgoing) {
		ids = append(ids, rel.End.Id())
	}
	//fmt.Println("gr:succ:",a,ids)
	return ids
}
func (g *dbGraph) Predecessors(a int) []int {

	node, _ := g.db.GetNode(a)
	ids := make([]int, 0)
	for _, rel := range node.Relations(Incoming) {
		ids = append(ids, rel.Start.Id())
	}
	//fmt.Println("gr:Pred:",a,ids)
	return ids
}
func (g *dbGraph) Relations(a int) []int {

	node, _ := g.db.GetNode(a)
	ids := make([]int, 0)
	for _, rel := range node.Relations(Both) {
		if rel.Start.Id() == a {
			ids = append(ids, rel.End.Id())
		} else {
			ids = append(ids, rel.Start.Id())
		}
	}
	//fmt.Println("gr:Rel:", a, ids)
	return ids
}
