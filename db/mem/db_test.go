package mem

import (
	"fmt"
	. "github.com/BuJo/goneo/db"
	"math/rand"
	"testing"
)

func ExampleNewDb() {

	db := NewDb()

	nodeA := db.NewNode()
	nodeA.SetProperty("foo", "bar")

	nodeB := db.NewNode()

	fmt.Println("nodes: ", nodeA, nodeB)

	relAB := nodeA.RelateTo(nodeB, "BELONGS_TO")

	fmt.Println("relation: ", relAB)

	path := db.FindPath(nodeA, nodeB)

	fmt.Println("path: ", path)

	nodeC := db.NewNode()

	nodeB.RelateTo(nodeC, "BELONGS_TO")

	path = db.FindPath(nodeA, nodeC)

	fmt.Println("path: ", path)
	// Output:
	//nodes:  (0 {foo:"bar",}) (1)
	//relation:  (0 {foo:"bar",})-[:BELONGS_TO]->(1)
	//path:  (0 {foo:"bar",})-[:BELONGS_TO]->(1)
	//path:  (0 {foo:"bar",})-[:BELONGS_TO]->(1)-[:BELONGS_TO]->(2)
}

func TestDbCreation(t *testing.T) {
	db := NewDb()

	if nodes := db.GetAllNodes(); len(nodes) != 0 {
		t.Fatal("DB should be empty")
	}
}

func TestNodeCreation(t *testing.T) {
	db := NewDb()

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

func TestPathFinding(t *testing.T) {
	db := NewDb()

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

func TestNodeProperties(t *testing.T) {
	db := NewDb()

	node := db.NewNode()
	node.SetProperty("foo", "bar")
	db.NewNode()

	if nodes := db.FindNodeByProperty("foo", "bar"); len(nodes) != 1 {
		t.Fatal("DB should deliver one node for foo")
	}
}

func BenchmarkCreateNode(b *testing.B) {

	db := NewDb()

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

	db := NewDb()

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
