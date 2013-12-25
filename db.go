package goneo

import (
	"errors"
	"fmt"
)

type DatabaseService struct {
	nodes         []*Node
	relationships []*Relation
}

func NewTemporaryDb() *DatabaseService {
	db := new(DatabaseService)

	db.nodes = make([]*Node, 0)
	db.relationships = make([]*Relation, 0)

	return db
}

func (db *DatabaseService) NewNode() *Node {
	n := new(Node)
	n.db = db

	db.nodes = append(db.nodes, n)
	n.id = len(db.nodes) - 1

	fmt.Println("New Node: ", n.id, ", nlen:", len(db.nodes))

	return n
}

func (db *DatabaseService) createRelation(a, b *Node) *Relation {
	r := new(Relation)
	r.Start = a
	r.End = b

	db.relationships = append(db.relationships, r)
	r.id = len(db.relationships) - 1

	return r
}

func (db *DatabaseService) GetNode(id int) (*Node, error) {
	if db.nodes == nil || len(db.nodes) < id+1 {
		return nil, errors.New(fmt.Sprintf("Node %d not found", id))
	}
	return db.nodes[id], nil
}
func (db *DatabaseService) GetAllNodes() []*Node {
	return db.nodes
}
func (db *DatabaseService) GetRelation(id int) (*Relation, error) {
	if db.nodes == nil || len(db.relationships) < id+1 {
		return nil, errors.New("Relationship not found")
	}
	return db.relationships[id], nil
}
func (db *DatabaseService) FindPath(start, end *Node) Path {

	builder := NewPathBuilder(start)

	for _, rel := range start.Relations(Outgoing) {
		builder, done := findPathRec(builder.Append(rel), end)
		if done {
			return builder.Build()
		}

	}

	return nil
}

func findPathRec(builder *PathBuilder, end *Node) (b *PathBuilder, done bool) {
	start := builder.Last()

	if start.id == end.id {
		return builder, true
	}

	for _, rel := range start.Relations(Outgoing) {
		builder, done := findPathRec(builder.Append(rel), end)
		if done {
			return builder, done

		}
	}

	return builder, false
}
