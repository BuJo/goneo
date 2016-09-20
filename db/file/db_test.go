package file

import (
	"testing"
)

func TestOpeningAndClosingDb(t *testing.T) {
	db := NewDb("file.db", nil)
	defer db.Close()
	if db == nil {
		t.Fatal("Db should have been created")
	}
}

func TestSavingAndGettingNodes(t *testing.T) {
	db := NewDb("file.db", nil)
	defer db.Close()

	node := db.NewNode()
	if node == nil {
		t.Fatal("Node creation should work")
	}

	node, _ = db.GetNode(node.Id())
	if node == nil {
		t.Fatal("Node getting should work")
	}
}

func TestGettingInvalidNode(t *testing.T) {
	db := NewDb("file.db", nil)
	defer db.Close()

	node, err := db.GetNode(77)
	if err == nil {
		t.Fatal("Expected error")
	}
	if node != nil {
		t.Fatal("Expected nil node on error")
	}
}

func TestGettingNodesAfterReOpenDb(t *testing.T) {
	db := NewDb("file.db", nil)
	db.Close()

	node := db.NewNode()
	if node == nil {
		t.Fatal("Node creation should work")
	}

	id := node.Id()

	db.Close()
	db = NewDb("file.db", nil)
	defer db.Close()

	node, _ = db.GetNode(id)
	if node == nil {
		t.Fatal("Node getting should work")
	}
}
