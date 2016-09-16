package goneo

import (
	"fmt"
	. "goneo/db"
)

// Format the given DB in Graphviz format.
func DumpDot(db DatabaseService) string {

	str := "digraph G {\n"

	for _, n := range db.GetAllNodes() {
		label := ""
		if len(n.Labels()) > 0 {
			label = n.Labels()[0]
		}
		str += fmt.Sprintf("\tn%d [label=\"%d\\n%s\"]\n", n.Id(), n.Id(), label)
	}
	str += "\n\n"

	for _, r := range db.GetAllRelations() {
		str += fmt.Sprintf("\tn%d -> n%d [label=\"%s\"]\n", r.Start().Id(), r.End().Id(), r.Type())
	}

	str += "}\n"

	return str
}
