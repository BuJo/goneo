package file

import (
	. "github.com/BuJo/goneo/db"
)

type filedb struct{}

func NewDb(db string, options map[string][]string) DatabaseService {
	return &filedb{}
}

func (db *filedb) NewNode(labels ...string) Node                { return nil }
func (db *filedb) Close()                                       {}
func (db *filedb) GetNode(id int) (Node, error)                 { return nil, nil }
func (db *filedb) GetAllNodes() []Node                          { return nil }
func (db *filedb) GetRelation(id int) (Relation, error)         { return nil, nil }
func (db *filedb) GetAllRelations() []Relation                  { return nil }
func (db *filedb) FindPath(start, end Node) Path                { return nil }
func (db *filedb) FindNodeByProperty(prop, value string) []Node { return nil }
