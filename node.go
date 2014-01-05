package goneo

import (
	"fmt"
	"sort"
)

type Node struct {
	db *DatabaseService
	id int

	labels     []string
	relations  []*Relation
	properties map[string]string
}

func (node *Node) String() string {
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

func (node *Node) Property(prop string) interface{} {
	return node.properties[prop]
}
func (node *Node) SetProperty(name, val string) {
	if node.properties == nil {
		node.properties = make(map[string]string)
	}
	node.properties[name] = val
}
func (node *Node) Properties() map[string]string {
	return node.properties
}

func (node *Node) HasProperty(prop string) bool {
	_, ok := node.properties[prop]
	return ok
}

func (node *Node) HasLabel(labels ...string) bool {
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

func (node *Node) Labels() []string {
	return node.labels
}

func (node *Node) RelateTo(end *Node, relType string) *Relation {

	for _, rel := range node.Relations(Outgoing) {
		if rel.End.id == end.id && rel.typ == relType {
			return rel
		}
	}

	rel := node.db.createRelation(node, end)
	rel.typ = relType

	node.relations = append(node.relations, rel)
	end.relations = append(end.relations, rel)

	return rel
}

func (node *Node) Relations(dir Direction) []*Relation {
	if dir == Both {
		return node.relations
	}

	rels := make([]*Relation, 0)

	for _, rel := range node.relations {
		if dir == Incoming && rel.End.id == node.id || dir == Outgoing && rel.Start.id == node.id {
			rels = append(rels, rel)
		}
	}

	return rels
}

func (node *Node) Id() int {
	return node.id
}
