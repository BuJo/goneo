package db

import (
	"fmt"
	"testing"
)

type mocknode struct{ name string }
type mockrel struct{ start, end Node }

func (*mocknode) Id() int          { return 0 }
func (m *mocknode) String() string { return "(" + m.name + ")" }

func (*mocknode) Property(prop string) interface{}             { return "" }
func (*mocknode) Properties() map[string]string                { return nil }
func (*mocknode) SetProperty(name, val string)                 {}
func (*mocknode) HasProperty(prop string) bool                 { return false }
func (*mocknode) HasLabel(labels ...string) bool               { return false }
func (*mocknode) Labels() []string                             { return nil }
func (m *mocknode) RelateTo(end Node, relType string) Relation { return &mockrel{m, end} }
func (*mocknode) Relations(dir Direction) []Relation           { return nil }

func (*mockrel) Id() int                          { return 0 }
func (m *mockrel) Start() Node                    { return m.start }
func (m *mockrel) End() Node                      { return m.end }
func (*mockrel) Type() string                     { return "HAS" }
func (*mockrel) Property(prop string) interface{} { return nil }
func (*mockrel) String() string                   { return "HAS" }

func ExampleNewPathBuilder() {
	var start Node = &mocknode{"a"}
	var end Node = &mocknode{"b"}
	var rel Relation = start.RelateTo(end, "HAS")

	builder := NewPathBuilder(start)
	builder = builder.Append(rel)
	path := builder.Build()
	fmt.Println(path.String())
	// Output: (a)-[:HAS]->(b)
}

func TestPathBuilding(t *testing.T) {
	var start Node = &mocknode{"a"}
	var end Node = &mocknode{"b"}
	var rel Relation = start.RelateTo(end, "HAS")

	builder := NewPathBuilder(start)
	builder = builder.Append(rel)
	path := builder.Build()

	if len(path.Items()) != 3 {
		t.Fatal("Should contain 3 elements, two nodes one relation")
	}
}
