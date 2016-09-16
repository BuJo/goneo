package goneo

import (
	"errors"
	. "github.com/BuJo/goneo/db"
	"github.com/BuJo/goneo/db/mem"
	"github.com/BuJo/goneo/gcy"
	"github.com/BuJo/goneo/sgi"
	"log"
)

type (
	evalContext struct {
		vars map[string][]PropertyContainer

		subgraphNameMap    map[string]int
		subgraphRevNameMap map[int]string

		subgraph *dbGraph

		db DatabaseService
	}

	query   struct{ q *gcy.Query }
	match   struct{ m *gcy.Match }
	root    struct{ r *gcy.Root }
	returns struct{ r *gcy.Returnable }
)

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
		table = table.Merge((&returns{r}).evaluate(ctx))
	}

	return table
}

// BUG(jo): db cannot encode undirected graph
// BUG(jo): cannot encode more than one possible relation type

func (mm *match) evaluate(ctx evalContext) *TabularData {
	m := mm.m

	subgraph := mem.NewDb()

	for _, p := range m.Paths {
		var builder *PathBuilder
		for currentNode := p.Start; currentNode != nil; {
			var n Node

			if currentNode.Name != "" {
				if id, ok := ctx.subgraphNameMap[currentNode.Name]; ok {
					n, _ = subgraph.GetNode(id)
					log.Print("tried to find node name ", currentNode.Name, ", found", n)
				}
			}

			if n == nil {
				n = subgraph.NewNode(currentNode.Labels...)
				ctx.subgraphNameMap[currentNode.Name] = n.Id()
				ctx.subgraphRevNameMap[n.Id()] = currentNode.Name

				log.Print("created new node: ", n, "(", currentNode, ")")

				for k, v := range currentNode.Props {
					n.SetProperty(k, v)
				}
			}

			if builder == nil {
				log.Print("first run, node: ", n, "(", currentNode, ")")
				builder = NewPathBuilder(n)
			} else {
				prevNode := builder.Last()
				log.Print("next run, ", prevNode, "->", n, "(", currentNode, ")")

				// TODO: utter crap, path is specific, has no "optional variants"
				// TODO: utter crap, db relations have no "optional variants"
				typ := ""
				if ok := len(currentNode.LeftRel.Types) > 0; ok {
					typ = currentNode.LeftRel.Types[0]
				}
				var rel Relation
				if currentNode.LeftRel.Direction == "->" {
					rel = prevNode.RelateTo(n, typ)
				} else {
					rel = n.RelateTo(prevNode, typ)
				}
				//log.Print("created new subgraph rel: ", rel, currentNode.LeftRel)

				builder = builder.Append(rel)
			}

			if currentNode.RightRel == nil {
				currentNode = nil
			} else {
				currentNode = currentNode.RightRel.RightNode
			}
		}
	}

	knownMappings := map[string]int{}
	for k, v := range ctx.vars {
		for _, n := range v {
			knownMappings[k] = n.Id()
		}
	}

	mappings := sgi.FindVF2SubgraphIsomorphism(&dbGraph{subgraph}, &dbGraph{ctx.db}, func(state sgi.State, fromQueryNode, fromTargetNode, toQueryNode, toTargetNode int) bool {
		//log.Print("tyring to find mapping for subgraph id ", toQueryNode, " in ", ctx.subgraphRevNameMap)
		if name, hasName := ctx.subgraphRevNameMap[toQueryNode]; hasName {

			if t2, hasMapping := knownMappings[name]; hasMapping {
				//log.Print(name, " mapped in ", ctx.vars, "targetId should be ", t2, " is ", toTargetNode)
				return t2 == toTargetNode
			}
			//log.Print(name, " not mapped in ", knownMappings, ", trying normal mapping")
		}

		return isSemanticallyFeasable(state, fromQueryNode, fromTargetNode, toQueryNode, toTargetNode)
	})

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

	log.Print("handled root: ", r, ", vars: ", ctx.vars)

	return &TabularData{}
}
func (rr *returns) evaluate(ctx evalContext) *TabularData {
	r := rr.r
	table := &TabularData{}

	table.line = make([]map[string]interface{}, 0)

	for _, o := range evaluateReturnable(ctx, r) {
		line := make(map[string]interface{})

		switch o.(type) {
		case Node:
			if r.Field != "" {
				line[r.Alias] = o.(Node).Property(r.Field)
			} else {
				line[r.Alias] = o
			}
		default:
			line[r.Alias] = o
		}
		table.line = append(table.line, line)

	}

	log.Print("evaluating return, ", table.line)

	return table
}

func evaluateReturnable(ctx evalContext, r *gcy.Returnable) []interface{} {
	objs := make([]interface{}, 0)

	switch r.Type {
	case "variable":
		for _, o := range ctx.vars[r.Object] {
			if o != nil {
				objs = append(objs, o)
			}
		}
	case "function":
		subobjs := make([][]interface{}, 0)
		for _, item := range r.Vars {
			subobjs = append(subobjs, evaluateReturnable(ctx, item))
		}
		switch r.Object {
		case "count":
			if len(subobjs) != 1 {
				log.Fatal("ERROR evaluating variable ", r.Name, ", need 1 variable")
			}
			objs = append(objs, len(subobjs[0]))
		}
	}

	return objs
}

//Evaluate a gcy query
//
// Example:
//
// 	start n=node(*) return n
func Evaluate(db DatabaseService, qry string) (*TabularData, error) {
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

	// Properties
	propsOk := true
	for k, v := range q2.Properties() {
		if ok := t2.HasProperty(k); !ok || t2.Property(k) != v {
			propsOk = false
			break
		}
	}

	if fromQueryNode == sgi.NULL_NODE {
		log.Printf("queryNode: %s , targetNode: %s\n", q2, t2)
		return labelsOk && propsOk
	}

	q1, _ := subgraph.GetNode(fromQueryNode)
	t1, _ := graph.GetNode(fromTargetNode)

	// query relation
	var qRel Relation
	var qDir Direction = Both
	for _, rel := range q1.Relations(Both) {
		if rel.Start().Id() == q1.Id() && rel.End().Id() == q2.Id() {
			qRel, qDir = rel, Outgoing
		} else if rel.End().Id() == q1.Id() && rel.Start().Id() == q2.Id() {
			qRel, qDir = rel, Incoming
		}
	}

	// target relation
	var tRel Relation
	var tDir Direction = Both
	for _, rel := range t1.Relations(Both) {
		if rel.Start().Id() == t1.Id() && rel.End().Id() == t2.Id() {
			tRel, tDir = rel, Outgoing
		} else if rel.End().Id() == t1.Id() && rel.Start().Id() == t2.Id() {
			tRel, tDir = rel, Incoming
		}
	}

	log.Printf("queryNodes: %s%s%s , targetNodes: %s%s%s | ", q1, qDir, q2, t1, tDir, t2)
	//log.Printf("rel: (%s~>%s), dir: (%s~>%s) \n", qRel.Type(), tRel.Type(), qDir, tDir)
	log.Print(qRel, tRel)

	return (qRel.Type() == "" || qRel.Type() == tRel.Type()) && (qDir == Both || qDir == tDir) && labelsOk && propsOk
}

type dbGraph struct {
	db DatabaseService
}

func (g *dbGraph) Order() int {
	return len(g.db.GetAllNodes())
}

func (g *dbGraph) Contains(a, b int) bool {
	//log.Print("gr:Contains:",a,b)
	node, _ := g.db.GetNode(a)
	for _, rel := range node.Relations(Both) {
		if rel.End().Id() == b || rel.Start().Id() == b {
			return true
		}
	}
	return false
}
func (g *dbGraph) Successors(a int) []int {

	node, _ := g.db.GetNode(a)
	ids := make([]int, 0)
	for _, rel := range node.Relations(Outgoing) {
		ids = append(ids, rel.End().Id())
	}
	//log.Print("gr:succ:",a,ids)
	return ids
}
func (g *dbGraph) Predecessors(a int) []int {

	node, _ := g.db.GetNode(a)
	ids := make([]int, 0)
	for _, rel := range node.Relations(Incoming) {
		ids = append(ids, rel.Start().Id())
	}
	//log.Print("gr:Pred:",a,ids)
	return ids
}
func (g *dbGraph) Relations(a int) []int {

	node, _ := g.db.GetNode(a)
	ids := make([]int, 0)
	for _, rel := range node.Relations(Both) {
		if rel.Start().Id() == a {
			ids = append(ids, rel.End().Id())
		} else {
			ids = append(ids, rel.Start().Id())
		}
	}
	//log.Print("gr:Rel:", a, ids)
	return ids
}
