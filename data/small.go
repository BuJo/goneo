package data

import (
	. "github.com/BuJo/goneo/db"
)

type small struct{ DatabaseService }

func NewSmallGenerator(db DatabaseService) *small {

	return &small{db}
}

func (gen *small) Generate() DatabaseService {
	nodeA := gen.NewNode()
	nodeA.SetProperty("foo", "bar")

	nodeB := gen.NewNode()
	nodeA.RelateTo(nodeB, "BELONGS_TO")

	nodeC := gen.NewNode()

	nodeB.RelateTo(nodeC, "BELONGS_TO")

	return gen
}
