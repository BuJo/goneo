package goneo

import (
	"bytes"
	"errors"
	"fmt"
	"goneo/gcy"
	"text/tabwriter"
)

type (
	evalContext struct {
		vars  map[string][]PropertyContainer
		paths map[string][]Path

		db *DatabaseService
	}

	Stringable interface {
		String() string
	}

	TabularData struct {
		line []map[string]Stringable
	}

	searchQuery struct{ q *gcy.SearchQuery }
	root        struct{ r *gcy.Root }
	returns     struct{ r *gcy.Return }
)

func (t TabularData) String() string {
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
			fmt.Fprint(w, line[header].String()+"\t")
		}
		fmt.Fprintln(w, "")
	}
	w.Flush()

	return b.String()
}
func (t TabularData) Len() int {
	return len(t.line)
}

func (q *searchQuery) evaluate(ctx evalContext) *TabularData {

	ctx.vars = make(map[string][]PropertyContainer)
	ctx.paths = make(map[string][]Path)

	for _, r := range q.q.Roots {
		(&root{r}).evaluate(ctx)
	}

	table := &TabularData{}

	for _, r := range q.q.Returns {
		table = (&returns{r}).evaluate(ctx)
	}

	return table
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
		line := make(map[string]Stringable)

		line[r.Alias] = o

		table.line = append(table.line, line)
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
	query, err := gcy.Parse("goneo", qry)
	if err != nil {
		return nil, err
	}
	table := (&searchQuery{query.(*gcy.SearchQuery)}).evaluate(evalContext{db: db})
	if table == nil {
		return nil, errors.New("Could not evaluate query")
	}

	return table, nil
}
