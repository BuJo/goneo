package goneo

import (
	"fmt"
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

	table, err := db.Evaluate("match (n:Tag)<-[:IS_TAGGED]-(v) return v")
	if err != nil {
		t.Error(err)
		return
	}
	if table == nil {
		t.Error("nil table")
		return
	}
	if table.Len() < 2 {
		t.Error("Evaluation not implemented")
	}
	if table.Len() < 5 {
		t.Error("Multiple matches not implemented")
	}
	if table.Len() > 12 {
		fmt.Println(table)
		t.Error("Too much matches")
	}
}

func TestTwoReturns(t *testing.T) {
	db := NewUniverseGenerator().Generate()

	table, err := db.Evaluate("match (n:Tag)<-[r:IS_TAGGED]-(v) return v, n")
	if err != nil {
		t.Error(err)
		return
	}
	if table == nil {
		t.Error("nil table")
		return
	}
	if table.Len() < 2 {
		t.Error("Evaluation not implemented")
		return
	}
	if table.Columns() < 2 {
		t.Error("Columns missing")
	}
	if table.Len() < 5 {
		t.Error("Multiple matches not implemented")
	}
	if table.Len() > 12 {
		fmt.Println(table)
		t.Error("Too much matches")
	}
}

func TestPropertyRetMatch(t *testing.T) {
	db := NewUniverseGenerator().Generate()

	table, err := db.Evaluate("match (n:Person {actor: \"Joss Whedon\"}) return n.actor as actor")
	if err != nil {
		t.Error(err)
		return
	}
	if table.Len() < 1 {
		t.Error("Evaluation not implemented")
		return
	}
}

func TestStartMatch(t *testing.T) {
	db := NewUniverseGenerator().Generate()

	creator, _ := db.GetNode(0)
	fmt.Println(creator, creator.Relations(Both))

	table, err := db.Evaluate("start joss=node(0) match (joss)-->(o) return o.series")
	if err != nil {
		t.Error(err)
		return
	}
	if table.Len() < 1 {
		t.Error("Evaluation not implemented")
		return
	}
}

func TestLongPathMatch(t *testing.T) {
	db := NewUniverseGenerator().Generate()

	table, err := db.Evaluate("match (e1:Episode)<-[:APPEARED_IN]-(niska {character: \"Adelai Niska\"})-[:APPEARED_IN]->(e2:Episode) return e1, e2")
	if err != nil {
		t.Error(err)
		return
	}
	if table.Len() < 1 {
		t.Error("Evaluation not implemented")
		return
	}
}

func TestMultiMatch(t *testing.T) {
	db := NewUniverseGenerator().Generate()

	table, err := db.Evaluate("match (e1)-[:ARCS_TO]->(e2), (e1)<-[:APPEARED_IN]-(niska {character: \"Adelai Niska\"})-[:APPEARED_IN]->(e2) return e1, e2")
	if err != nil {
		t.Error(err)
		return
	}
	if table.Len() < 1 {
		t.Skip("Evaluation not implemented")
		return
	}
}

func TestPathVariable(t *testing.T) {
	db := NewUniverseGenerator().Generate()

	table, err := db.Evaluate("start joss=node(0) match path = (joss)-->(o) return path")
	if err != nil {
		t.Error(err)
		return
	}
	if table.Len() < 1 {
		t.Skip("Saving paths not implemented")
		return
	}
}

type UniverseGenerator struct {
	db *DatabaseService

	crew     []*Node
	episodes []*Node
	enemies  []*Node
}

func NewUniverseGenerator() *UniverseGenerator {
	gen := new(UniverseGenerator)

	gen.db = NewTemporaryDb()

	gen.addMeta()
	gen.addCharacters()
	gen.addEpisodes()

	return gen
}

func (gen *UniverseGenerator) Generate() *DatabaseService {
	return gen.db
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
func (gen *UniverseGenerator) addCharacters() {
	ship := gen.db.NewNode("Ship")
	ship.SetProperty("character", "Firefly")

	mal := gen.actor("Nathan Fillion").played("Captain Malcolm 'Mal' Reynolds")
	zoe := gen.actor("Gina Torres").played("Zoë Washburne")
	wash := gen.actor("Alan Tudyk").played("Hoban 'Wash' Washburne")
	inara := gen.actor("Morena Baccarin").played("Inara Serra")
	jayne := gen.actor("Adam Baldwin").played("Jayne Cobb")
	kaylee := gen.actor("Jewel Staite").played("Kaylee Frye")
	simon := gen.actor("Sean Maher").played("Simon Tam")
	river := gen.actor("Summer Glau").played("River Tam")
	sheperd := gen.actor("Ron Glass").played("Shepherd Derrial Book")
	gen.actor("Skylar Roberge").played("River Tam")
	gen.actor("Zac Efron").played("Simon Tam")
	gen.actor("Joss Whedon").played("Man at Funeral")
	blue1 := gen.actor("Jeff Ricketts").played("Blue Glove Man #1")
	blue2 := gen.actor("Dennis Cockrum").played("Blue Glove Man #2")
	niska := gen.actor("Michael Fairman").played("Adelai Niska")
	operative := gen.actor("Chiwetel Ejiofor").played("The Operative")

	gen.character(river).is("SISTER").of(simon)
	gen.character(simon).is("BROTHER").of(river)

	gen.crew = []*Node{mal, zoe, wash, inara, jayne, kaylee, simon, river, sheperd}
	gen.enemies = []*Node{blue1, blue2, niska, operative}

	gen.characters(gen.enemies...).are("ENEMY").of(gen.crew...)

	gen.characters(gen.crew...).are("CREW").of(ship)
	gen.character(mal).is("CAPTAIN").of(ship)
}
func (gen *UniverseGenerator) addEpisodes() {

	ep01 := gen.createEpisode(1, "")
	ep02 := gen.createEpisode(2, "Train Job")
	ep13 := gen.createEpisode(13, "War Stories")

	niska := gen.enemies[2]

	niska.RelateTo(ep02, "APPEARED_IN")
	ep02.RelateTo(ep13, "ARCS_TO")
	niska.RelateTo(ep13, "APPEARED_IN")

	gen.episodes = []*Node{ep01, ep02, ep13}
}

func (gen *UniverseGenerator) createEpisode(nr int, title string) *Node {
	ep := gen.db.NewNode("Episode")
	ep.SetProperty("episode", fmt.Sprintf("%d", nr))
	ep.SetProperty("title", title)
	return ep
}

type actorBuilder struct {
	actor *Node

	db *DatabaseService
}

func (gen *UniverseGenerator) actor(name string) *actorBuilder {
	b := new(actorBuilder)
	b.db = gen.db

	actors := b.db.FindNodeByProperty("actor", name)
	if len(actors) == 1 {
		b.actor = actors[0]
	} else {
		b.actor = b.db.NewNode("Actor", "Person")
		b.actor.SetProperty("actor", name)
	}

	return b
}

func (b *actorBuilder) played(names ...string) *Node {

	var character *Node

	for _, name := range names {
		characters := b.db.FindNodeByProperty("character", name)

		if len(characters) == 0 {
			character = b.db.NewNode("Character")
			character.SetProperty("character", name)
		} else {
			character = characters[0]
		}

		b.actor.RelateTo(character, "PLAYED")
	}

	return character
}

type characterBuilder struct {
	characters []*Node
	relType    string

	db *DatabaseService
}

func (gen *UniverseGenerator) character(char *Node) *characterBuilder {
	b := new(characterBuilder)
	b.db = gen.db

	b.characters = append(b.characters, char)

	return b
}
func (gen *UniverseGenerator) characters(chars ...*Node) *characterBuilder {

	b := new(characterBuilder)
	b.db = gen.db

	b.characters = chars

	return b
}
func (b *characterBuilder) is(reltype string) *characterBuilder {
	b.relType = reltype
	return b
}
func (b *characterBuilder) are(reltype string) *characterBuilder {
	b.relType = reltype
	return b
}
func (b *characterBuilder) of(actors ...*Node) *characterBuilder {
	for _, actor := range actors {
		for _, character := range b.characters {
			character.RelateTo(actor, b.relType)
		}
	}
	return b
}
