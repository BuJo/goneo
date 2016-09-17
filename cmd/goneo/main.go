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
	"github.com/BuJo/goneo"
	"github.com/BuJo/goneo/data"
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

	db, _ := goneo.OpenDb("mem:testdb")

	if *size == "universe" {
		db = data.NewUniverseGenerator(db).Generate()
	} else if *size == "big" {
		maxNodes := 5000
		rand.Seed(42)

		db.NewNode()
		for n := db.NewNode(); n.Id() < maxNodes; n = db.NewNode() {
			t, _ := db.GetNode(rand.Intn(n.Id()))
			n.RelateTo(t, "HAS")
		}
	} else {
		nodeA := db.NewNode()
		nodeA.SetProperty("foo", "bar")

		nodeB := db.NewNode()
		nodeA.RelateTo(nodeB, "BELONGS_TO")

		nodeC := db.NewNode()

		nodeB.RelateTo(nodeC, "BELONGS_TO")

	}

	server := web.NewGoneoServer(db)

	if port := os.Getenv("PORT"); port != "" {
		*binding = ":" + port
	}

	server.Bind(*binding)

	if apikey := os.Getenv("HOSTEDGRAPHITE_APIKEY"); apikey != "" {
		host, port := "carbon.hostedgraphite.com", 2003
		server.EnableGraphite(host, port, apikey)
	}

	server.Start()
}
