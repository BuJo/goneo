package file

import (
	. "github.com/BuJo/goneo/db"
)

type edge struct {
	db *filedb

	id         int
	typ        string
	start, end *node
}

func (e *edge) Id() int        { return e.id }
func (e *edge) String() string { return "" }
func (e *edge) Start() Node    { return e.start }
func (e *edge) End() Node      { return e.end }
func (e *edge) Type() string   { return e.typ }

func (e *edge) Property(prop string) interface{}        { return nil }
func (e *edge) Properties() map[string]interface{}      { return nil }
func (e *edge) SetProperty(nam string, val interface{}) {}
