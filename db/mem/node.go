package mem

import (
	"fmt"
	"sort"

	. "github.com/BuJo/goneo/db"
)

type node struct {
	db *databaseService
	id int

	labels     []string
	relations  []Relation
	properties map[string]string
}

func (n *node) String() string {
	props := " {"
	for key, val := range n.properties {
		props += key + ":\"" + val + "\","
	}
	props += "}"
	if len(props) == 3 {
		props = ""
	}

	labels := ""
	for _, l := range n.labels {
		labels += ":" + l
	}

	return fmt.Sprintf("(%d%s%s)", n.id, labels, props)
}

func (n *node) Property(prop string) interface{} {
	return n.properties[prop]
}
func (n *node) SetProperty(name, val string) {
	if n.properties == nil {
		n.properties = make(map[string]string)
	}
	n.properties[name] = val
}
func (n *node) Properties() map[string]string {
	return n.properties
}

func (n *node) HasProperty(prop string) bool {
	_, ok := n.properties[prop]
	return ok
}

func (n *node) HasLabel(labels ...string) bool {
	for _, label := range labels {
		i := sort.SearchStrings(n.labels, label)
		if i < len(n.labels) && n.labels[i] == label {
			// x is present at data[i]
		} else {
			// x is not present in data,
			return false
		}
	}

	return true
}

func (n *node) Labels() []string {
	return n.labels
}

func (n *node) RelateTo(endI Node, relType string) Relation {
	end, ok := endI.(*node)

	if !ok {
		panic("Handling Node of a different DB implementation")
	}

	for _, rel := range n.Relations(Outgoing) {
		if rel.End().Id() == end.id && rel.Type() == relType {
			return rel
		}
	}

	rel := n.db.createRelation(n, end)
	rel.typ = relType

	n.relations = append(n.relations, rel)
	end.relations = append(end.relations, rel)

	return rel
}

func (n *node) Relations(dir Direction) []Relation {
	if dir == Both {
		return n.relations
	}

	rels := make([]Relation, 0)

	for _, rel := range n.relations {
		if dir == Incoming && rel.End().Id() == n.id || dir == Outgoing && rel.Start().Id() == n.id {
			rels = append(rels, rel)
		}
	}

	return rels
}

func (n *node) Id() int {
	return n.id
}
