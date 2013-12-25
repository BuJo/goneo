package goneo

import "fmt"

type PropertyContainer interface {
	Property(string) interface{}
	String() string
}

type Path interface {
	Nodes() []*Node
	Relations() []*Relation
	Items() []PropertyContainer

	String() string
}

type simplePath struct {
	start     *Node
	relations []*Relation
}

func (path *simplePath) Nodes() (nodes []*Node) {
	nodes = make([]*Node, 0)

	nodes = append(nodes, path.start)

	if path.relations == nil || len(path.relations) == 0 {
		return
	}

	for _, rel := range path.relations {
		nodes = append(nodes, rel.End)
	}

	return
}

func (path *simplePath) Relations() []*Relation {
	return path.relations
}

func (path *simplePath) Items() []PropertyContainer {
	items := make([]PropertyContainer, 0)
	items = append(items, path.start)

	if path.relations == nil || len(path.relations) == 0 {
		return items
	}

	for _, rel := range path.relations {
		items = append(items, rel)
		items = append(items, rel.End)
	}

	return items
}

func (path *simplePath) String() (str string) {
	str = fmt.Sprintf("(%d)", path.start.id)

	for _, rel := range path.relations {
		str = fmt.Sprintf("%s-[:%s]->(%d)", str, rel.Type, rel.End.id)
	}

	return
}

type PathBuilder struct {
	start     *Node
	relations []*Relation
}

func (builder *PathBuilder) Build() Path {
	return &simplePath{builder.start, builder.relations}
}

func (builder *PathBuilder) Append(rel *Relation) *PathBuilder {
	b := new(PathBuilder)
	b.start = builder.start
	b.relations = append(builder.relations, rel)

	return b
}
func (builder *PathBuilder) Last() *Node {
	if builder.relations == nil || len(builder.relations) == 0 {
		return builder.start
	}
	return builder.relations[len(builder.relations)-1].End
}

func NewPathBuilder(start *Node) *PathBuilder {
	builder := new(PathBuilder)
	builder.start = start
	builder.relations = make([]*Relation, 0)

	return builder
}
