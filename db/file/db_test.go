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
