package mem

import (
	"fmt"
	. "github.com/BuJo/goneo/db"
	"math/rand"
	"testing"
)

func ExampleNewDb() {

	db, _ := NewDb("test", nil)

	nodeA := db.NewNode()
	nodeA.SetProperty("foo", "bar")

	nodeB := db.NewNode()

	fmt.Println("nodes: ", nodeA, nodeB)

	relAB := nodeA.RelateTo(nodeB, "BELONGS_TO")

	fmt.Println("relation: ", relAB)

	path := db.FindPath(nodeA, nodeB)

	fmt.Println("path: ", path)

	nodeC := db.NewNode("A")

	nodeB.RelateTo(nodeC, "BELONGS_TO")

	path = db.FindPath(nodeA, nodeC)

	fmt.Println("path: ", path)
	// Output:
	//nodes:  (0 {foo:"bar",}) (1)
	//relation:  (0 {foo:"bar",})-[:BELONGS_TO]->(1)
	//path:  (0 {foo:"bar",})-[:BELONGS_TO]->(1)
	//path:  (0 {foo:"bar",})-[:BELONGS_TO]->(1)-[:BELONGS_TO]->(2:A)
}

func TestDbCreation(t *testing.T) {
	db, _ := NewDb("test", nil)

	if nodes := db.GetAllNodes(); len(nodes) != 0 {
		t.Fatal("DB should be empty")
	}
}

func TestNodeCreation(t *testing.T) {
	db, _ := NewDb("test", nil)

	db.NewNode()
	db.NewNode()

	if nodes := db.GetAllNodes(); len(nodes) != 2 {
		t.Fatal("DB should be filled")
	}

	if n, err := db.GetNode(0); err != nil || n.Id() != 0 {
		t.Fatal("First node should be id 0")
	}

	if n, err := db.GetNode(1); err != nil || n.Id() != 1 {
		t.Fatal("Second node should be id 1")
	}
}

func TestRelationCreation(t *testing.T) {
	db, _ := NewDb("test", nil)

	nodeA := db.NewNode()
	nodeB := db.NewNode()
	rel := nodeA.RelateTo(nodeB, "HAS")
	rel.SetProperty("foo", "bar")

	if rel.Id() < 0 {
		t.Error("Relationship should have ok id")
	}

	if len(rel.Properties()) != 1 {
		t.Error("Rel should have one property")
	}

	if rel.Property("foo") != "bar" {
		t.Error("Should be able to retrieve property")
	}
}

func TestRelationRetrieval(t *testing.T) {
	db, _ := NewDb("test", nil)

	nodeA := db.NewNode()
	nodeB := db.NewNode()
	rel := nodeA.RelateTo(nodeB, "HAS")

	if o, err := db.GetRelation(rel.Id()); err != nil || rel.Id() != o.Id() {
		t.Error("Should retrieve saved relation")
	}

	if len(db.GetAllRelations()) != 1 {
		t.Error("Should contain one relation")
	}
}

func TestFailingRelationRetrieval(t *testing.T) {
	db, _ := NewDb("test", nil)

	nodeA := db.NewNode()
	nodeB := db.NewNode()
	nodeA.RelateTo(nodeB, "HAS")

	if _, err := db.GetRelation(6); err == nil {
		t.Error("Should fail to retrieve relation")
	}
}

func TestNoNode(t *testing.T) {
	db, _ := NewDb("test", nil)

	newNode := db.NewNode()

	if node, err := db.GetNode(newNode.Id() + 1); err == nil || node != nil {
		t.Fatal("Expected error")
	}
}

func TestPathFinding(t *testing.T) {
	db, _ := NewDb("test", nil)

	nodeA := db.NewNode()
	nodeB := db.NewNode()
	nodeC := db.NewNode()

	nodeA.RelateTo(nodeB, "BELONGS_TO")
	nodeB.RelateTo(nodeC, "BELONGS_TO")

	path := db.FindPath(nodeA, nodeC)

	if len(path.Nodes()) != 3 {
		t.Error("Should have 3 nodes in path")
	}
	if len(path.Relations()) != 2 {
		t.Error("Should have 2 relationships in path")
	}
}

func TestFailedPathFinding(t *testing.T) {
	db, _ := NewDb("test", nil)

	nodeA := db.NewNode()
	nodeB := db.NewNode()
	nodeC := db.NewNode()

	nodeA.RelateTo(nodeB, "BELONGS_TO")
	nodeA.RelateTo(nodeC, "BELONGS_TO")

	path := db.FindPath(nodeB, nodeC)

	if path != nil {
		t.Error("Should not have found path")
	}
}

func TestNodeProperties(t *testing.T) {
	db, _ := NewDb("test", nil)

	node := db.NewNode()
	node.SetProperty("foo", "bar")
	db.NewNode()

	nodes := db.FindNodeByProperty("foo", "bar")
	if len(nodes) != 1 {
		t.Fatal("DB should deliver one node for foo")
	}

	if node := nodes[0]; len(node.Properties()) != 1 {
		t.Fatal("Node should have properties")
	}
}

func TestNodeLabels(t *testing.T) {
	db, _ := NewDb("test", nil)

	node := db.NewNode("Human")

	if node.HasLabel("Human") == false {
		t.Error("Should have label")
	}

	if node.HasLabel("Robot") == true {
		t.Error("Should not have label")
	}

	if len(node.Labels()) != 1 {
		t.Error("Should have only one label")
	}
}

func TestNodeRelating(t *testing.T) {
	db, _ := NewDb("test", nil)

	nodeA := db.NewNode()
	nodeB := db.NewNode()

	nodeA.RelateTo(nodeB, "HAS")

	if rels := nodeA.Relations(Both); len(rels) != 1 {
		t.Error("There should be one relation")
	}

	if rels := nodeA.Relations(Outgoing); len(rels) != 1 {
		t.Error("There should be one outgoing relation")
	}

	if rels := nodeA.Relations(Incoming); len(rels) != 0 {
		t.Error("There should be no incoming relation")
	}

	if rels := nodeB.Relations(Incoming); len(rels) != 1 {
		t.Error("There should be one incoming relation")
	}

	nodeA.RelateTo(nodeB, "HAS")

	if rels := nodeA.Relations(Both); len(rels) != 1 {
		t.Error("There should be one relation")
	}
}

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

func (*mockrel) Id() int                                 { return 0 }
func (m *mockrel) Start() Node                           { return m.start }
func (m *mockrel) End() Node                             { return m.end }
func (*mockrel) Type() string                            { return "HAS" }
func (*mockrel) Property(prop string) interface{}        { return nil }
func (*mockrel) Properties() map[string]interface{}      { return nil }
func (*mockrel) SetProperty(nam string, val interface{}) {}

func (*mockrel) String() string { return "HAS" }

func TestPanicOnMixingDBImplementations(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("Should have panicked")
		}
	}()

	db, _ := NewDb("test", nil)

	nodeA := db.NewNode()
	nodeA.RelateTo(&mocknode{}, "HAS")
}

func BenchmarkCreateNode(b *testing.B) {

	db, _ := NewDb("test", nil)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		db.NewNode()
	}
}

var DB_CACHE DatabaseService

func createRandomDb(maxNodes int) DatabaseService {
	if DB_CACHE != nil {
		return DB_CACHE
	}

	rand.Seed(42)

	db, _ := NewDb("test", nil)

	db.NewNode()

	for n := db.NewNode(); n.Id() < maxNodes; n = db.NewNode() {
		t, _ := db.GetNode(rand.Intn(n.Id()))
		n.RelateTo(t, "HAS")
	}

	DB_CACHE = db

	return db
}

func BenchmarkPathFinding(b *testing.B) {
	maxNodes := 500000
	rand.Seed(42)

	db := createRandomDb(maxNodes)

	b.ResetTimer()

	var x, y Node

	for i := 0; i < b.N; i++ {
		x, _ = db.GetNode(rand.Intn(maxNodes))
		y, _ = db.GetNode(rand.Intn(maxNodes))

		db.FindPath(x, y)
	}
}
