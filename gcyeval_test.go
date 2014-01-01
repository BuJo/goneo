package goneo

import (
	"testing"
)

func TestBasicStartQuery(t *testing.T) {
	db := NewUniverseGenerator().Generate()

	table, err := db.Evaluate("start n=node(*) return n as node")
	if err != nil {
		t.Error(err)
		return
	}
	if table.Len() < 1 {
		t.Error(table)
	}
}
func TestUniverse(t *testing.T) {
	db := NewUniverseGenerator().Generate()

	creators := db.FindNodeByProperty("creator", "Joss Whedon")
	if len(creators) != 1 {
		t.Fail()
	}

}

func TestTagged(t *testing.T) {
	db := NewUniverseGenerator().Generate()

	table, err := db.Evaluate("match (n:Tag)<-[:TAGGED_BY]-(v) return v")
	if err != nil {
		t.Error(err)
		return
	}
	if table.Len() != 2 {
		t.Error("Evaluation not implemented")
	}
}

type UniverseGenerator struct {
	db *DatabaseService
}

func NewUniverseGenerator() *UniverseGenerator {
	gen := new(UniverseGenerator)

	gen.db = NewTemporaryDb()

	gen.addMeta()
	gen.addShip()
	gen.addActors()

	return gen
}

func (gen *UniverseGenerator) Generate() *DatabaseService {
	return gen.db
}

func (gen *UniverseGenerator) addShip() {
	node := gen.db.NewNode()
	node.SetProperty("ship", "Firefly")

}
func (gen *UniverseGenerator) addMeta() {
	creator := gen.db.NewNode("Person")
	creator.SetProperty("creator", "Joss Whedon")

	series := gen.db.NewNode("Series")
	series.SetProperty("series", "Firefly")

	movie := gen.db.NewNode("Movie")
	movie.SetProperty("movie", "Serenity")

	creator.RelateTo(series, "CREATED")

	for _, tag := range []string{"Adventure", "Drama", "Sci-Fi"} {
		t := gen.db.NewNode("Tag")
		t.SetProperty("tag", tag)
		series.RelateTo(t, "IS_TAGGED")
		movie.RelateTo(t, "IS_TAGGED")
	}
}
func (gen *UniverseGenerator) addActors() {
	gen.actor("Nathan Fillion").played("Captain Malcolm 'Mal' Reynolds")
	gen.actor("Gina Torres").played("Zoë Washburne")
	gen.actor("Alan Tudyk").played("Hoban 'Wash' Washburne")
	gen.actor("Morena Baccarin").played("Inara Serra")
	gen.actor("Adam Baldwin").played("Jayne Cobb")
	gen.actor("Jewel Staite").played("Kaylee Frye")
	gen.actor("Sean Maher").played("Simon Tam")
	gen.actor("Summer Glau").played("River Tam")
	gen.actor("Ron Glass").played("Shepherd Derrial Book")
	gen.actor("Skylar Roberge").played("River Tam")
	gen.actor("Zac Efron").played("Simon Tam")
	gen.actor("Joss Whedon").played("Man at Funeral")
	gen.actor("Jeff Ricketts").played("Blue Glove Man #1")
	gen.actor("Dennis Cockrum").played("Blue Glove Man #2")
	gen.actor("Michael Fairman").played("Adelai Niska")
	gen.actor("Chiwetel Ejiofor").played("The Operative")
}
func (gen *UniverseGenerator) addEpisodes() {

}

type actorBuilder struct {
	actor *Node

	db *DatabaseService
}

func (gen *UniverseGenerator) actor(name string) *actorBuilder {
	b := new(actorBuilder)
	b.db = gen.db

	b.actor = b.db.NewNode("Actor")
	b.actor.SetProperty("actor", name)

	return b
}

func (b *actorBuilder) played(names ...string) *actorBuilder {

	for _, name := range names {
		var character *Node
		characters := b.db.FindNodeByProperty("character", name)

		if len(characters) == 0 {
			character = b.db.NewNode("Character")
			character.SetProperty("character", name)
		} else {
			character = characters[0]
		}

		b.actor.RelateTo(character, "PLAYED")
	}

	return b
}
