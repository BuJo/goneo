package data

import . "github.com/BuJo/goneo/db"

func ExampleNewUniverseGenerator() {
	var db DatabaseService

	NewUniverseGenerator(db).Generate()
}
