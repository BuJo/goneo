package mem

import (
	"fmt"
	. "github.com/BuJo/goneo/db"
)

type relation struct {
	db *databaseService
	id int

	typ        string
	start, end Node
	properties map[string]interface{}
}

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
func (rel *relation) Properties() map[string]interface{} {
	return rel.properties
}
func (rel *relation) SetProperty(name string, val interface{}) {
	if rel.properties == nil {
		rel.properties = make(map[string]interface{})
	}
	rel.properties[name] = val
}

func (rel *relation) Type() string { return rel.typ }
func (rel *relation) Id() int      { return rel.id }
func (rel *relation) Start() Node  { return rel.start }
func (rel *relation) End() Node    { return rel.end }
