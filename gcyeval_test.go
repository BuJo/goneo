package goneo

import (
	"fmt"
	"testing"
)

func TestBasicStartQuery(t *testing.T) {
	db := NewUniverseGenerator().Generate()
	count := len(db.GetAllNodes())

	table, err := db.Evaluate("start n=node(*) return n as node")
	NewTableTester(t, table, err).HasColumns("node").HasLen(count)
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
	NewTableTester(t, table, err).HasLen(6).HasColumns("v")
}

func TestTwoReturns(t *testing.T) {
	db := NewUniverseGenerator().Generate()

	table, err := db.Evaluate("match (n:Tag)<-[r:IS_TAGGED]-(v) return v, n")
	NewTableTester(t, table, err).HasLen(6).HasColumns("v", "n")
}

func TestPropertyRetMatch(t *testing.T) {
	db := NewUniverseGenerator().Generate()

	table, err := db.Evaluate("match (n:Person {actor: \"Joss Whedon\"}) return n.actor as actor")
	NewTableTester(t, table, err).Has("actor", "Joss Whedon")
}

func TestStartMatch(t *testing.T) {
	db := NewUniverseGenerator().Generate()

	creator, _ := db.GetNode(0)
	fmt.Println(creator, creator.Relations(Both))

	table, err := db.Evaluate("start joss=node(0) match (joss)-->(o) return o.series")
	NewTableTester(t, table, err).Has("o.series", "Firefly")
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

	table, err := db.Evaluate("match (e1)-[:ARCS_TO]->(e2), (e1)<-[:APPEARED_IN]-(niska {character: \"Adelai Niska\"})-[:APPEARED_IN]->(e2) return e1.episode, e2.episode")
	NewTableTester(t, table, err).Has("e1.episode", "2")
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

func TestFunctionCount(t *testing.T) {
	db := NewUniverseGenerator().Generate()

	table, err := db.Evaluate("match (e:Episode)-[:ARCS_TO]->(e2) return count(e) as nrArcs")
	NewTableTester(t, table, err).Has("nrArcs", 1)
}

type TableTester struct {
	t          *testing.T
	table      *TabularData
	currentRow int
	err        error
}

func NewTableTester(t *testing.T, table *TabularData, err error) *TableTester {
	tester := new(TableTester)
	tester.t = t
	tester.table = table
	tester.currentRow = -1
	tester.err = err

	if err != nil {
		t.Error(err)
	}

	return tester
}
func (t *TableTester) Has(column string, expected interface{}) *TableTester {
	t.currentRow += 1

	if t.err != nil {
		return t
	}

	if t.table.Len() < t.currentRow+1 {
		t.t.Error("Table not big enough, want: ", t.currentRow+1)
		return t
	}

	if actual := t.table.Get(t.currentRow, column); actual != expected {
		t.t.Error("Bad cell, expected ", expected, " got ", actual)
	}

	return t
}
func (t *TableTester) HasLen(l int) *TableTester {
	if t.err != nil {
		return t
	}

	if t.table.Len() != l {
		t.t.Error("Bad table length, expected ", l, " got ", t.table.Len())
	}
	return t
}
func (t *TableTester) HasColumns(cols ...string) *TableTester {
	if t.err != nil {
		return t
	}

	if len(t.table.Columns()) != len(cols) {
		t.t.Error("Bad number of columns, expected ", len(cols), " got ", len(t.table.Columns()), ": ", t.table.Columns())
	}

	for _, expectedCol := range cols {
		found := false
		for _, actualCol := range t.table.Columns() {
			if expectedCol == actualCol {
				found = true
			}
		}
		if !found {
			t.t.Error("Bad column length, expected ", expectedCol, " in ", t.table.Columns())
		}
	}
	return t
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