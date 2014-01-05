package goneo

import "fmt"

type Relation struct {
	db *DatabaseService
	id int

	typ        string
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
	return fmt.Sprintf("%s-[:%s]->%s", rel.Start, rel.typ, rel.End)
}

func (rel *Relation) Property(prop string) interface{} {
	return rel.Properties[prop]
}

func DirectionFromString(str string) Direction {
	if str == "out" || str == "outgoing" {
		return Outgoing
	} else if str == "in" || str == "incoming" {
		return Incoming
	}
	return Both
}


func (d Direction) String() string {
	switch d {
	case Both:
		return "-"
	case Incoming:
		return "<-"
	case Outgoing:
		return "->"
	default:
		return ""
	}
}

func (rel *Relation) Type() string {
	return rel.typ
}
