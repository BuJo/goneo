package mem

import (
	"fmt"
	"testing"
)

func TestBasic(t *testing.T) {

	db := NewTemporaryDb()

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
}
