package data

import (
	"fmt"

	. "github.com/BuJo/goneo/db"
)

type universeGenerator struct {
	db DatabaseService

	crew     []Node
	episodes []Node
	enemies  []Node
	series   Node
}

// NewUniverseGenerator for an in-memory database containing information about
// a reasonably popular sci-fi TV series.
func NewUniverseGenerator(db DatabaseService) *universeGenerator {
	gen := new(universeGenerator)

	gen.db = db

	gen.addMeta()
	gen.addCharacters()
	gen.addEpisodes()

	return gen
}

func (gen *universeGenerator) Generate() DatabaseService {
	return gen.db
}

func (gen *universeGenerator) addMeta() {
	creator := gen.actor("Joss Whedon").actor
	creator.SetProperty("creator", "Joss Whedon")

	gen.series = gen.db.NewNode("Series")
	gen.series.SetProperty("series", "Firefly")

	movie := gen.db.NewNode("Movie")
	movie.SetProperty("movie", "Serenity")

	creator.RelateTo(gen.series, "CREATED")

	for _, tag := range []string{"Adventure", "Drama", "Sci-Fi"} {
		t := gen.db.NewNode("Tag")
		t.SetProperty("tag", tag)
		gen.series.RelateTo(t, "IS_TAGGED")
		movie.RelateTo(t, "IS_TAGGED")
	}
}
func (gen *universeGenerator) addCharacters() {
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
	blue1 := gen.actor("Jeff Ricketts").played("Blue Glove Man #1")
	blue2 := gen.actor("Dennis Cockrum").played("Blue Glove Man #2")
	niska := gen.actor("Michael Fairman").played("Adelai Niska")
	operative := gen.actor("Chiwetel Ejiofor").played("The Operative")

	gen.character(river).is("SISTER").of(simon)
	gen.character(simon).is("BROTHER").of(river)

	gen.crew = []Node{mal, zoe, wash, inara, jayne, kaylee, simon, river, sheperd}
	gen.enemies = []Node{blue1, blue2, niska, operative}

	gen.characters(gen.enemies...).are("ENEMY").of(ship)

	gen.characters(gen.crew...).are("CREW").of(ship)
	gen.character(mal).is("CAPTAIN").of(ship)
}
func (gen *universeGenerator) addEpisodes() {

	season := gen.db.NewNode("Season")
	gen.series.RelateTo(season, "HAS_SEASON")

	ep01 := gen.createEpisode(1, "Serenity")
	ep02 := gen.createEpisode(2, "Train Job")
	ep03 := gen.createEpisode(3, "Bushwhacked")
	ep04 := gen.createEpisode(4, "Shindig")
	ep05 := gen.createEpisode(5, "Safe")
	ep06 := gen.createEpisode(6, "Our Mrs. Reynolds")
	ep07 := gen.createEpisode(7, "Jaynestown")
	ep08 := gen.createEpisode(8, "Out of Gas")
	ep09 := gen.createEpisode(9, "Ariel")
	ep10 := gen.createEpisode(10, "War Stories")
	ep11 := gen.createEpisode(11, "Trash")
	ep12 := gen.createEpisode(12, "The Message")
	ep13 := gen.createEpisode(13, "Heart of Gold")
	ep14 := gen.createEpisode(14, "Objects in Space")

	niska := gen.enemies[2]

	niska.RelateTo(ep02, "APPEARED_IN")
	ep02.RelateTo(ep10, "ARCS_TO")
	niska.RelateTo(ep10, "APPEARED_IN")

	funman := gen.actor("Joss Whedon").played("Man at Funeral")
	funman.RelateTo(ep12, "APPEARED_IN")

	gen.episodes = []Node{ep01, ep02, ep03, ep04, ep05, ep06, ep07, ep08, ep09, ep10, ep11, ep12, ep13, ep14}

	e0 := gen.episodes[0]
	for _, e1 := range gen.episodes[1:] {
		e0.RelateTo(e1, "LEADS_TO")
		e0 = e1
	}

	season.RelateTo(ep01, "BEGINS_WITH")
}

func (gen *universeGenerator) createEpisode(nr int, title string) Node {
	ep := gen.db.NewNode("Episode")
	ep.SetProperty("episode", fmt.Sprintf("%d", nr))
	ep.SetProperty("title", title)
	return ep
}

type actorBuilder struct {
	actor Node

	db DatabaseService
}

func (gen *universeGenerator) actor(name string) *actorBuilder {
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

func (b *actorBuilder) played(names ...string) Node {

	var character Node

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
	characters []Node
	relType    string

	db DatabaseService
}

func (gen *universeGenerator) character(char Node) *characterBuilder {
	b := new(characterBuilder)
	b.db = gen.db

	b.characters = append(b.characters, char)

	return b
}
func (gen *universeGenerator) characters(chars ...Node) *characterBuilder {

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
func (b *characterBuilder) of(actors ...Node) *characterBuilder {
	for _, actor := range actors {
		for _, character := range b.characters {
			character.RelateTo(actor, b.relType)
		}
	}
	return b
}
