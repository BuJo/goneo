// Simple memory based implementation of DatabaseService.
package mem

import (
	"errors"
	"fmt"
	. "github.com/BuJo/goneo/db"
	"sort"
)

type databaseService struct {
	nodes         []Node
	relationships []Relation
}

// Create a DB instance of a simple memory backed graph DB
func NewDb() DatabaseService {
	db := new(databaseService)

	db.nodes = make([]Node, 0)
	db.relationships = make([]Relation, 0)

	return db
}

func (db *databaseService) NewNode(labels ...string) Node {
	n := new(node)
	n.db = db

	db.nodes = append(db.nodes, n)
	n.id = len(db.nodes) - 1

	sort.Strings(labels)

	n.labels = make([]string, 0, 1)
	for _, l := range labels {
		n.labels = append(n.labels, l)
	}

	return n
}

func (db *databaseService) createRelation(a, b Node) *relation {
	r := new(relation)
	r.start = a
	r.end = b

	db.relationships = append(db.relationships, r)
	r.id = len(db.relationships) - 1

	return r
}

func (db *databaseService) GetNode(id int) (Node, error) {
	if db.nodes == nil || len(db.nodes) < id+1 || id < 0 {
		return nil, errors.New(fmt.Sprintf("Node %d not found", id))
	}
	return db.nodes[id], nil
}

func (db *databaseService) GetAllNodes() []Node {
	return db.nodes
}

func (db *databaseService) GetRelation(id int) (Relation, error) {
	if db.nodes == nil || len(db.relationships) < id+1 {
		return nil, errors.New("Relationship not found")
	}
	return db.relationships[id], nil
}

func (db *databaseService) GetAllRelations() []Relation {
	return db.relationships
}

func (db *databaseService) FindPath(start, end Node) Path {

	builder := NewPathBuilder(start)

	for _, rel := range start.Relations(Outgoing) {
		builder, done := findPathRec(builder.Append(rel), end)
		if done {
			return builder.Build()
		}

	}

	return nil
}

func findPathRec(builder *PathBuilder, end Node) (b *PathBuilder, done bool) {
	start := builder.Last()

	if start.Id() == end.Id() {
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

func (db *databaseService) FindNodeByProperty(prop, value string) []Node {
	found := make([]Node, 0)

	for _, node := range db.nodes {
		if node.HasProperty(prop) && node.Property(prop) == value {
			found = append(found, node)
		}
	}

	return found
}
