package mem

import (
	"fmt"
	. "goneo/db"
	"sort"
)

type node struct {
	db *databaseService
	id int

	labels     []string
	relations  []Relation
	properties map[string]string
}

func (node *node) String() string {
	props := " {"
	for key, val := range node.properties {
		props += key + ":\"" + val + "\","
	}
	props += "}"
	if len(props) == 3 {
		props = ""
	}

	labels := ""
	for _, l := range node.labels {
		labels += ":" + l
	}

	return fmt.Sprintf("(%d%s%s)", node.id, labels, props)
}

func (node *node) Property(prop string) interface{} {
	return node.properties[prop]
}
func (node *node) SetProperty(name, val string) {
	if node.properties == nil {
		node.properties = make(map[string]string)
	}
	node.properties[name] = val
}
func (node *node) Properties() map[string]string {
	return node.properties
}

func (node *node) HasProperty(prop string) bool {
	_, ok := node.properties[prop]
	return ok
}

func (node *node) HasLabel(labels ...string) bool {
	for _, label := range labels {
		i := sort.SearchStrings([]string(node.labels), label)
		if i < len(node.labels) && node.labels[i] == label {
			// x is present at data[i]
		} else {
			// x is not present in data,
			return false
		}
	}

	return true
}

func (node *node) Labels() []string {
	return node.labels
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

func (node *node) Relations(dir Direction) []Relation {
	if dir == Both {
		return node.relations
	}

	rels := make([]Relation, 0)

	for _, rel := range node.relations {
		if dir == Incoming && rel.End().Id() == node.id || dir == Outgoing && rel.Start().Id() == node.id {
			rels = append(rels, rel)
		}
	}

	return rels
}

func (node *node) Id() int {
	return node.id
}
