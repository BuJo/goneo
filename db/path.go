package db

import "fmt"

// PropertyContainer encapsulates graph nodes and edges so they can be
// handled interchangeably for output.
type PropertyContainer interface {
	Property(string) interface{}
	Id() int
	String() string
}

// Path encapsulates a simple path in a graph
type Path interface {
	Nodes() []Node
	Relations() []Relation
	Items() []PropertyContainer

	String() string
}

type simplePath struct {
	start     Node
	relations []Relation
}

// Nodes returns all nodes from a path
func (path *simplePath) Nodes() (nodes []Node) {
	nodes = make([]Node, 0)

	nodes = append(nodes, path.start)

	if len(path.relations) == 0 {
		return
	}

	for _, rel := range path.relations {
		nodes = append(nodes, rel.End())
	}

	return
}

// Relations returns all relations from a path
func (path *simplePath) Relations() []Relation {
	return path.relations
}

// Items returns all elements from a path
func (path *simplePath) Items() []PropertyContainer {
	items := make([]PropertyContainer, 0)
	items = append(items, path.start)

	if len(path.relations) == 0 {
		return items
	}

	for _, rel := range path.relations {
		items = append(items, rel)
		items = append(items, rel.End())
	}

	return items
}

func (path *simplePath) String() (str string) {
	str = path.start.String()
	left := path.start

	for _, rel := range path.relations {
		relstr := ""
		if rel.Type() != "" {
			relstr = fmt.Sprintf("[:%s]", rel.Type())
		}

		direction := Both
		if rel.Start() == left {
			direction = Outgoing
			left = rel.End()
		} else {
			direction = Incoming
			left = rel.Start()
		}

		switch direction {
		case Both:
			relstr = "-" + relstr + "-"
		case Incoming:
			relstr = "<-" + relstr + "-"
		case Outgoing:
			relstr = "-" + relstr + "->"
		}
		str = fmt.Sprintf("%s%s%s", str, relstr, left)
	}

	return
}

// PathBuilder provides a way to build a path.
type PathBuilder struct {
	start, end Node
	relations  []Relation
}

func (builder *PathBuilder) Build() Path {
	return &simplePath{builder.start, builder.relations}
}

func (builder *PathBuilder) Append(rel Relation) *PathBuilder {
	b := new(PathBuilder)
	b.start, b.end = builder.start, builder.end
	b.relations = append(builder.relations, rel)

	if rel.End() == b.end {
		b.end = rel.Start()
	} else {
		b.end = rel.End()
	}
	//fmt.Println("end of path ", b.Build(), " is now: ", b.end, " added rel ", rel)
	return b
}
func (builder *PathBuilder) Last() Node {
	return builder.end
}

func NewPathBuilder(start Node) *PathBuilder {
	builder := new(PathBuilder)
	builder.start, builder.end = start, start
	builder.relations = make([]Relation, 0)

	return builder
}
