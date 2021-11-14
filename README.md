# goneo

Author: Jonathan Buch <jonathan.buch@gmail.com>

This software contains:

* a simple in-memory node/edge database for a directed graph
* a query mechanism on the db for single nodes and a simple depth first path
* a simple cypher-like language
* an http server for access

### Hacking

[![Build Status](https://travis-ci.org/BuJo/goneo.svg?branch=master)](https://travis-ci.org/BuJo/goneo)

This project uses go modules to handle dependencies.

Testing:

	go test ./...

Building:

```
go build github.com/BuJo/goneo/cmd/goneo
```

Running the program with a few generated nodes in the in-memory-db:

```
./goneo -size universe
```

Running cypher to search for some labelled nodes:

```
curl -F 'gocy=match (t:Tag) RETURN t.tag AS tag' localhost:7474/table
```

The web interface can also render out the db as graphviz format via `/graphviz` (which also understands form field gocy with a search query).

#### Release

	govendor sync
	VERSION=v1.0
	go install -ldflags "-X main.buildversion=$VERSION"

### Literature

* OPTIMIZED BACKTRACKING FOR SUBGRAPH ISOMORPHISM
  Lixin Fu and Shruthi Chandra
* A (Sub)Graph Isomorphism Algorithm for Matching Large Graphs
  Luigi P. Cordella, Pasquale Foggia, Carlo Sansone, and Mario Vento
* Labeled Subgraph Matching Using Degree Filtering
  Lixin Fu and Surya Prakash R Kommireddy
* An Improved Algorithm for Matching Large Graphs
  L. P. Cordella, P. Foggia, C. Sansone, M. Vento
* Read, R. C. and Corneil, D. G. (1977). The graph isomorphism disease. Journal of Graph Theory 1, 339–363.
* Gati, G. (1979). Further annotated bibliography on the isomorphism disease. Journal of Graph Theory 3, 95–109.
* An Algorithm for Subgraph Isomorphism
  J. R. ULLMANN
