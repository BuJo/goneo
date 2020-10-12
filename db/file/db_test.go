package file

import (
	. "github.com/BuJo/goneo/db"
	"os"
	"testing"
)

func TestOpeningAndClosingDb(t *testing.T) {
	db, _ := NewDb("file.db", nil)
	if db == nil {
		t.Fatal("Db should have been created")
	}
	defer db.Close()
	defer os.Remove("file.db")

	if nodes := db.GetAllNodes(); len(nodes) != 0 {
		t.Fatal("DB should be empty")
	}
}

func TestSavingAndGettingNodes(t *testing.T) {
	db, _ := NewDb("file.db", nil)
	defer db.Close()
	defer os.Remove("file.db")

	node := db.NewNode()
	if node == nil {
		t.Fatal("Node creation should work")
	}

	db.NewNode()
	db.NewNode()
	db.NewNode()
	db.NewNode()
	db.NewNode()
	db.NewNode()
	node = db.NewNode()

	_, err := db.GetNode(node.Id())
	if err != nil {
		t.Fatalf("Getting node %d should work", node.Id())
	}

	if nodes := db.GetAllNodes(); len(nodes) == 0 {
		t.Fatal("DB should not be empty")
	}
}

func TestGettingInvalidNode(t *testing.T) {
	db, _ := NewDb("file.db", nil)
	defer db.Close()
	defer os.Remove("file.db")

	node, err := db.GetNode(77)
	if err == nil {
		t.Fatal("Expected error")
	}
	if node != nil {
		t.Fatal("Expected nil node on error")
	}
}

func TestGettingNodesAfterReOpenDb(t *testing.T) {
	db, _ := NewDb("file.db", nil)
	defer os.Remove("file.db")

	node := db.NewNode()
	if node == nil {
		t.Fatal("Node creation should work")
	}

	db.Close()
	db, _ = NewDb("file.db", nil)
	defer db.Close()

	node, _ = db.GetNode(0)
	if node == nil {
		t.Fatal("Getting node after re-opening the database should work")
	}
}

func TestRelationCreation(t *testing.T) {
	db, _ := NewDb("file.db", nil)
	defer os.Remove("file.db")

	nodeA := db.NewNode()
	nodeB := db.NewNode()
	rel := nodeA.RelateTo(nodeB, "HAS")

	if rel.Id() < 0 {
		t.Error("Relationship should have ok id")
	}

	db.Close()
	db, _ = NewDb("file.db", nil)
	defer db.Close()

	node, _ := db.GetNode(0)
	rels := node.Relations(Both)
	if len(rels) != 1 {
		t.Error("There should be one relationship")
	}
}
