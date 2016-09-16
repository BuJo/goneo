// Commandline test application serving an in memory DB via http.
/*
Flags:

	-size=[small|big|universe]
	-bind=:7474
	-version

Sizes:

	* small: Three-node cluster
	* big: randomly generated tree
	* universe: sci-fi tv series information
*/
package main

import (
	"flag"
	"fmt"
	. "github.com/BuJo/goneo"
	. "github.com/BuJo/goneo/db"
	"github.com/BuJo/goneo/db/mem"
	"github.com/BuJo/goneo/web"
	"log"
	"math/rand"
	"os"
)

var (
	binding = flag.String("bind", ":7474", "Bind to ip/port")
	size    = flag.String("size", "small", "Size of generated graph")
	version = flag.Bool("version", false, "Print version information")

	buildversion, buildtime string = "SNAPSHOT", ""
)

func main() {
	flag.Parse()

	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile | log.LUTC)

	if *version {
		fmt.Println("goneo version", buildversion, "(", buildtime, ")")
		return
	}

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

	port := os.Getenv("PORT")
	if port != "" {
		*binding = ":" + port
	}

	web.NewGoneoServer(db).Bind(*binding).Start()
}
