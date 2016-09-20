package file

import (
	. "github.com/BuJo/goneo/db"
)

type filedb struct {
	nodes []*node
}

func NewDb(name string, options map[string][]string) DatabaseService {
	db := &filedb{}

	db.nodes = make([]*node, 0)

	return db
}

func (db *filedb) NewNode(labels ...string) Node {
	n := &node{}

	db.nodes = append(db.nodes, n)
	n.id = len(db.nodes) - 1

	return n
}
func (db *filedb) GetNode(id int) (Node, error) {
	if id > len(db.nodes) || i < 0 {
		return nil, errors.New("Did not find id")
	}

	return db.nodes[id], nil
}

func (db *filedb) GetAllNodes() []Node                          { return nil }
func (db *filedb) GetRelation(id int) (Relation, error)         { return nil, nil }
func (db *filedb) GetAllRelations() []Relation                  { return nil }
func (db *filedb) FindPath(start, end Node) Path                { return nil }
func (db *filedb) FindNodeByProperty(prop, value string) []Node { return nil }

func (db *filedb) Close() {}
