package data

import (
	. "github.com/BuJo/goneo/db"
	"math/rand"
)

type randgenerator struct{ DatabaseService }

func NewLargeRandomGenerator(db DatabaseService) DatabaseGenerator {

	return randgenerator{db}
}

func (db randgenerator) Generate() DatabaseService {
	maxNodes := 5000
	rand.Seed(42)

	db.NewNode()
	for n := db.NewNode(); n.Id() < maxNodes; n = db.NewNode() {
		t, _ := db.GetNode(rand.Intn(n.Id()))
		n.RelateTo(t, "HAS")
	}

	return db
}
