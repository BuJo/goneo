package goneo

import "fmt"

type Relation interface {
	Id() int
	Start() Node
	End() Node
	Type() string

	Property(prop string) interface{}
	String() string

	setStart(start Node)
	setEnd(end Node)
	setId(id int)
	setType(typ string)
}

type relation struct {
	db *DatabaseService
	id int

	typ        string
	start, end Node
	properties map[string]interface{}
}

type Direction int

const (
	Both Direction = iota
	Incoming
	Outgoing
)

func (rel *relation) String() string {
	relstr := ""
	if rel.typ != "" {
		relstr = fmt.Sprintf("[:%s]", rel.typ)
	}
	return fmt.Sprintf("%s-%s->%s", rel.start, relstr, rel.end)
}

func (rel *relation) Property(prop string) interface{} {
	return rel.properties[prop]
}

func (rel *relation) Type() string { return rel.typ }
func (rel *relation) Id() int      { return rel.id }
func (rel *relation) Start() Node  { return rel.start }
func (rel *relation) End() Node    { return rel.end }

func (rel *relation) setStart(start Node) { rel.start = start }
func (rel *relation) setEnd(end Node)     { rel.end = end }
func (rel *relation) setId(id int)        { rel.id = id }
func (rel *relation) setType(typ string)  { rel.typ = typ }

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
