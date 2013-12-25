package main

import "fmt"

type Relation struct {
	db *DatabaseService
	id int

	Type       string
	Start, End *Node
	Properties map[string]interface{}
}

type Direction int

const (
	Both Direction = iota
	Incoming
	Outgoing
)

func (rel *Relation) String() string {
	return fmt.Sprintf("(%d)-[:%s]->(%d)", rel.Start.id, rel.Type, rel.End.id)
}

func (rel *Relation) Property(prop string) interface{} {
	return rel.Properties[prop]
}
