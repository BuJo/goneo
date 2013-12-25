package main

import (
	"bytes"
	"fmt"
	"strconv"
	"text/tabwriter"
)

type (
	Context struct {
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

func (q *searchQuery) evaluate(ctx Context) TabularData {

	ctx.vars = make(map[string][]PropertyContainer)
	ctx.paths = make(map[string][]Path)

	for _, r := range q.roots {
		r.evaluate(ctx)
	}

	table := TabularData{}

	for _, r := range q.returns {
		table = r.evaluate(ctx)
	}

	return table
}

func (r *root) evaluate(ctx Context) TabularData {
	_, ok := ctx.vars[r.name]

	if !ok {
		ctx.vars[r.name] = make([]PropertyContainer, 0)
	}

	if r.typ == "node" {

		if r.idRange == "*" {
			for _, node := range ctx.db.GetAllNodes() {
				ctx.vars[r.name] = append(ctx.vars[r.name], node)
			}
		} else {
			id, _ := strconv.Atoi(r.idRange)
			ctx.vars[r.name] = append(ctx.vars[r.name], ctx.db.GetNode(id))
		}
	} else {
		id, _ := strconv.Atoi(r.idRange)
		ctx.vars[r.name] = append(ctx.vars[r.name], ctx.db.GetRelation(id))
	}

	return TabularData{}
}
func (r *returns) evaluate(ctx Context) TabularData {
	table := TabularData{}

	table.line = make([]map[string]Stringable, 0)

	for _, o := range ctx.vars[r.name] {
		line := make(map[string]Stringable)

		line[r.alias] = o

		table.line = append(table.line, line)
	}

	fmt.Println("evaluating return, ", table.line)

	return table
}
