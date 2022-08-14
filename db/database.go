// Package db is the basic contract package for a database handled by goneo.
package db

type Node interface {
	Id() int
	String() string
	Property(prop string) interface{}
	Properties() map[string]string
	SetProperty(name, val string)
	HasProperty(prop string) bool
	HasLabel(labels ...string) bool
	Labels() []string
	RelateTo(end Node, relType string) Relation
	Relations(dir Direction) []Relation
}

type Relation interface {
	Id() int
	String() string
	Start() Node
	End() Node
	Type() string

	Property(prop string) interface{}
	Properties() map[string]interface{}
	SetProperty(nam string, val interface{})
}

type DatabaseService interface {
	NewNode(labels ...string) Node

	GetNode(id int) (Node, error)
	GetAllNodes() []Node

	GetRelation(id int) (Relation, error)
	GetAllRelations() []Relation

	FindPath(start, end Node) Path

	FindNodeByProperty(prop, value string) []Node

	Close()
}
