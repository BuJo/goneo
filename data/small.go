package data

import (
	. "github.com/BuJo/goneo/db"
)

type smallgen struct{ DatabaseService }

func NewSmallGenerator(db DatabaseService) DatabaseGenerator {

	return smallgen{db}
}

func (db smallgen) Generate() DatabaseService {
	nodeA := db.NewNode()
	nodeA.SetProperty("foo", "bar")

	nodeB := db.NewNode()
	nodeA.RelateTo(nodeB, "BELONGS_TO")

	nodeC := db.NewNode()

	nodeB.RelateTo(nodeC, "BELONGS_TO")

	return db
}
