package data

import (
	"math/rand"

	. "github.com/BuJo/goneo/db"
)

type random struct{ DatabaseService }

func NewLargeRandomGenerator(db DatabaseService) *random {

	return &random{db}
}

func (gen *random) Generate() DatabaseService {
	maxNodes := 5000
	rand.Seed(42)

	gen.NewNode()
	for n := gen.NewNode(); n.Id() < maxNodes; n = gen.NewNode() {
		t, _ := gen.GetNode(rand.Intn(n.Id()))
		n.RelateTo(t, "HAS")
	}

	return gen
}
