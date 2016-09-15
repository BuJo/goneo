package goneo

import "fmt"

type PropertyContainer interface {
	Property(string) interface{}
	Id() int
	String() string
}

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

func (path *simplePath) Nodes() (nodes []Node) {
	nodes = make([]Node, 0)

	nodes = append(nodes, path.start)

	if path.relations == nil || len(path.relations) == 0 {
		return
	}

	for _, rel := range path.relations {
		nodes = append(nodes, rel.End())
	}

	return
}

func (path *simplePath) Relations() []Relation {
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
		items = append(items, rel.End())
	}

	return items
}

func (path *simplePath) String() (str string) {
	str = fmt.Sprintf("%s", path.start)
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
