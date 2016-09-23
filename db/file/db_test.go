package file

import (
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
