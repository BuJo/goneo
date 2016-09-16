package main

import (
	"flag"
	. "goneo"
	. "goneo/db"
	"goneo/db/mem"
	"goneo/web"
	"math/rand"
)

var (
	binding = flag.String("bind", ":7474", "Bind to ip/port")
	size    = flag.String("size", "small", "Size of generated graph")
)

func main() {
	flag.Parse()

	var db DatabaseService

	if *size == "universe" {
		db = NewUniverseGenerator().Generate()
	} else if *size == "big" {
		maxNodes := 5000
		rand.Seed(42)

		db = mem.NewDb()
		db.NewNode()
		for n := db.NewNode(); n.Id() < maxNodes; n = db.NewNode() {
			t, _ := db.GetNode(rand.Intn(n.Id()))
			n.RelateTo(t, "HAS")
		}
	} else {
		db = mem.NewDb()

		nodeA := db.NewNode()
		nodeA.SetProperty("foo", "bar")

		nodeB := db.NewNode()
		nodeA.RelateTo(nodeB, "BELONGS_TO")

		nodeC := db.NewNode()

		nodeB.RelateTo(nodeC, "BELONGS_TO")

	}

	web.NewGoneoServer(db).Bind(*binding).Start()
}
