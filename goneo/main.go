package main

import (
	"flag"
	"goneo/db/mem"
	"goneo/web"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"syscall"
)

var (
	binding = flag.String("bind", ":7474", "Bind to ip/port")
	size    = flag.String("size", "small", "Size of generated graph")

	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
	memprofile = flag.String("memprofile", "", "write mem profile to file")
)

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		runtime.SetCPUProfileRate(1000)
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		sigChan := make(chan os.Signal)
		signal.Notify(sigChan, syscall.SIGINT)
		go func() {
			s := <-sigChan

			if *memprofile != "" {
				f, err := os.Create(*memprofile)
				if err != nil {
					log.Fatal("could not create memory profile: ", err)
				}
				runtime.GC() // get up-to-date statistics
				if err := pprof.WriteHeapProfile(f); err != nil {
					log.Fatal("could not write memory profile: ", err)
				}
				f.Close()
			}

			log.Printf("Recieved: %+v\n", s)
			pprof.StopCPUProfile()
			os.Exit(0)
		}()
	}

	db := mem.NewDb()

	nodeA := db.NewNode()
	nodeA.SetProperty("foo", "bar")

	nodeB := db.NewNode()
	nodeA.RelateTo(nodeB, "BELONGS_TO")

	nodeC := db.NewNode()

	nodeB.RelateTo(nodeC, "BELONGS_TO")

	if *size == "big" {
		maxNodes := 5000
		rand.Seed(42)

		for n := db.NewNode(); n.Id() < maxNodes; n = db.NewNode() {
			t, _ := db.GetNode(rand.Intn(n.Id()))
			n.RelateTo(t, "HAS")
		}
	}

	web.NewGoneoServer(db).Bind(*binding).Start()
}
