package goneo

import (
	"fmt"

	. "github.com/BuJo/goneo/db"
)

// DumpDot formats the given DB in Graphviz format.
func DumpDot(db DatabaseService) string {

	str := "digraph G {\n"

	for _, n := range db.GetAllNodes() {
		label := ""
		if len(n.Properties()) > 0 {
			// Select a property which is long'ish
			for _, p := range n.Properties() {
				if len(p) > 2 {
					label = p
					break
				}
			}
		} else if len(n.Labels()) > 0 {
			// Otherwise try to use a label
			label = n.Labels()[0]
		}
		str += fmt.Sprintf("\tn%d [label=\"[%d]\\n%s\"]\n", n.Id(), n.Id(), label)
	}
	str += "\n\n"

	for _, r := range db.GetAllRelations() {
		str += fmt.Sprintf("\tn%d -> n%d [label=\"%s\"]\n", r.Start().Id(), r.End().Id(), r.Type())
	}

	str += "}\n"

	return str
}
