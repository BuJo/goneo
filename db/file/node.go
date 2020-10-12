package file

import (
	. "github.com/BuJo/goneo/db"
)

type node struct {
	db *filedb

	id     int
	labels []string
	edges  []*edge
}

func (n *node) Id() int                          { return n.id }
func (n *node) String() string                   { return "" }
func (n *node) Property(prop string) interface{} { return nil }
func (n *node) Properties() map[string]string    { return nil }
func (n *node) SetProperty(name, val string)     {}
func (n *node) HasProperty(prop string) bool     { return false }
func (n *node) HasLabel(labels ...string) bool   { return false }
func (n *node) Labels() []string                 { return n.labels }

func (n *node) RelateTo(end Node, relType string) Relation {
	edge, _ := n.db.createEdge(n, end.(*node))
	return edge
}

func (n *node) Relations(dir Direction) []Relation { return nil }
