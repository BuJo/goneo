package main

import (
	"goneo/db/mem"
	"goneo/web"
)

func main() {

	db := mem.NewDb()

	nodeA := db.NewNode()
	nodeA.SetProperty("foo", "bar")

	nodeB := db.NewNode()
	nodeA.RelateTo(nodeB, "BELONGS_TO")

	nodeC := db.NewNode()

	nodeB.RelateTo(nodeC, "BELONGS_TO")

	web.NewGoneoServer(db).Bind(":7474").Start()
}
