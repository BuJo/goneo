package goneo

import (
	"fmt"
)

type Label string

type Node struct {
	db *DatabaseService
	id int

	Labels     []Label
	relations  []*Relation
	properties map[string]string
}

func (node *Node) String() string {
	props := " {"
	for key, val := range node.properties {
		props += key + ":" + val + ","
	}
	props += "}"
	if len(props) == 3 {
		props = ""
	}
	return fmt.Sprintf("(%d%s)", node.id, props)
}

func (node *Node) Property(prop string) interface{} {
	return node.properties[prop]
}
func (node *Node) SetProperty(name, val string) {
	if (node.properties == nil ) {
		node.properties = make(map[string]string)
	}
	node.properties[name] = val
}
func (node *Node) GetProperties()  map[string]string {
	return node.properties
}

func (node *Node) RelateTo(end *Node, relType string) *Relation {
	rel := node.db.createRelation(node, end)
	rel.Type = relType

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
