package main

import (
	"fmt"
	"goneo/db/mem"
	"goneo/web"
)

func main() {

	db := mem.NewDb()

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

	web.NewGoneoServer(db).Bind(":7474").Start()
}
