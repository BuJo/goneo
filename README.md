# goneo

Author: Jonathan Buch <jonathan.buch@gmail.com>

This software contains:

* a simple in-memory node/edge database for a directed graph
* a query mechanism on the db for single nodes and a simple depth first path
* a simple cypher-like language
* an http server for access

### Hacking

[![Build Status](https://travis-ci.org/BuJo/goneo.svg?branch=master)](https://travis-ci.org/BuJo/goneo)

This project uses `govendor` to handle dependencies. See the [govendor quickstart][govquick] for more information.

Testing:

	govendor sync
	govendor test +local

#### Release

	govendor sync
	VERSION=v1.0
	go install -ldflags "-X main.buildversion=$VERSION -X main.buildtime=$(date +%FT%X%z)"

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

[govquick]: https://github.com/kardianos/govendor#quick-start-also-see-the-faq
