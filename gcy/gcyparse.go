package gcy

import (
	"errors"
	"fmt"
	"log"
	"strconv"
)

type (
	Query struct {
		Roots   []*Root
		Match   *Match
		Returns []*Returnable
		Deletes []*Returnable
		Creates []*Returnable
	}

	Root struct {
		Name   string
		Typ    string
		IdVars []int
	}

	Match struct {
		Paths []*Path
	}

	Path struct {
		Name  string
		Start *Node
	}

	Node struct {
		Name   string
		Labels []string
		Props  map[string]string

		LeftRel, RightRel *Relation
	}

	Relation struct {
		Name        string
		Direction   string
		Types       []string
		Cardinality string

		LeftNode, RightNode *Node
	}

	Returnable struct {
		Type string // int, string, function, variable

		Name          string
		Alias         string
		Object, Field string
		Vars          []*Returnable
		Value         interface{}
	}
)

type errorList []error

func (list errorList) Error() string {
	switch len(list) {
	case 0:
		return "no errors"
	case 1:
		return "error:" + list[0].Error()
	}
	return fmt.Sprintf("%s (and %d more errors)", list[0], len(list)-1)
}

func newError(msg string) error {
	return errors.New(fmt.Sprintf("%s", msg))
}

type parser struct {
	errors  errorList
	scanner chan item
	tok     item // one token look-ahead
}

func (p *parser) next() {
	p.tok = <-p.scanner

	if p.tok.typ == itemError {
		p.error(p.tok.val)
	}
}

func (p *parser) error(msg string) {
	p.errors = append(p.errors, newError(msg))
}

func (p *parser) errorExpected(msg string) {
	msg = `expected "` + msg + `"`
	p.error(msg)
}

func (p *parser) expect(tok string) {
	if p.tok.val != tok {
		p.errorExpected(tok + "(got:" + p.tok.String() + ")")
	}
	p.next() // make progress in any case
}
func (p *parser) expectType(tok itemType) {
	if p.tok.typ != tok {
		p.errorExpected(fmt.Sprintf("%s, got %s", item{typ: tok}, p.tok))
	}
	p.next() // make progress in any case
}

func (p *parser) parseStart() []*Root {
	log.Print("parsing search query")

	var roots []*Root

	roots = p.parseRoots()

	log.Print("returning from search query: ", roots)

	return roots
}

func (p *parser) parseRoots() []*Root {
	log.Print("parsing search query roots")

	roots := make([]*Root, 0)
	roots = append(roots, p.parseRoot())
	return roots
}

func (p *parser) parseRoot() (r *Root) {
	varname := p.tok.val
	p.expectType(itemIdentifier)
	p.expectType(itemEqual)

	r = &Root{Name: varname, IdVars: make([]int, 0)}

	switch p.tok.val {
	case "node", "relation":
		typ := p.tok.val
		p.expect(typ)
		p.expectType(itemLParen)

		r.Typ = typ

	Loop:
		for {
			var err error

			switch p.tok.typ {
			case itemRange:
				start, end := 0, 0
				p.expectType(itemRange)
				end, err = strconv.Atoi(p.tok.val)
				if err != nil {
					p.error(err.Error())
				}
				p.expectType(itemNumber)

				for i := start; i < end; i += 1 {
					r.IdVars = append(r.IdVars, i)
				}
			case itemNumber:
				start, end := 0, 0
				start, err = strconv.Atoi(p.tok.val)
				if err != nil {
					p.error(err.Error())
				}
				p.expectType(itemNumber)

				end = start

				if p.tok.typ == itemRange {
					p.expectType(itemRange)
					end, err = strconv.Atoi(p.tok.val)
					if err != nil {
						p.error(err.Error())
					}
					p.expectType(itemNumber)

				}
				log.Print("parsing item number in root: ", start, end)
				for i := start; i <= end; i += 1 {
					r.IdVars = append(r.IdVars, i)
				}
			case itemStar:
				p.expectType(itemStar)
				r.IdVars = append(r.IdVars, -1)

			default:
				break Loop
			}
			if p.tok.typ == itemComma {
				p.expectType(itemComma)
				continue
			}
			break
		}

	default:
		p.expect("node")
	}

	p.expect(")")

	return r
}
func (p *parser) parseReturns() []*Returnable {
	rets := make([]*Returnable, 0)

	for p.tok.typ == itemIdentifier {
		var ret *Returnable
		name := p.tok.val
		p.expectType(itemIdentifier)

		if p.tok.typ == itemLParen {
			ret = p.parseFunction(name)
		} else {
			ret = p.parseVariable(name)
		}

		ret.Alias = ret.Name

		if p.tok.typ == itemAs {
			p.expectType(itemAs)

			ret.Alias = p.tok.val

			p.expectType(itemIdentifier)
		}
		log.Print("adding var", ret)
		rets = append(rets, ret)

		if p.tok.typ == itemComma {
			p.expectType(itemComma)
			continue
		}
		break
	}

	return rets
}
func (p *parser) parseVariable(name string) *Returnable {
	variable := &Returnable{Type: "variable", Name: name, Object: name}

	if p.tok.typ == itemDot {
		p.expectType(itemDot)

		variable.Name += "."
		variable.Field = p.tok.val
		variable.Name += variable.Field

		p.expectType(itemField)
	}

	return variable
}
func (p *parser) parseFunction(name string) *Returnable {
	function := &Returnable{Type: "function", Name: name, Object: name}
	p.expectType(itemLParen)
	function.Name += "("

	function.Vars = make([]*Returnable, 0)

	for p.tok.typ == itemIdentifier {
		var ret *Returnable
		name := p.tok.val
		p.expectType(itemIdentifier)

		if p.tok.typ == itemLParen {
			ret = p.parseFunction(name)
		} else {
			ret = p.parseVariable(name)
		}

		function.Name += ret.Name

		function.Vars = append(function.Vars, ret)

		if p.tok.typ == itemComma {
			p.expectType(itemComma)
			function.Name += ", "
			continue
		}
		break
	}

	function.Name += ")"
	p.expectType(itemRParen)

	return function
}

func (p *parser) parseDelete() []*Returnable {
	return nil
}

func (p *parser) parseCreate() []*Returnable {
	return nil
}

func (p *parser) parseMatch() *Match {
	match := new(Match)

	for {
		path := new(Path)

		if p.tok.typ == itemLParen {
			path.Start = p.parsePath()
		} else {
			path.Name = p.tok.val
			p.expectType(itemIdentifier)

			p.expectType(itemEqual)

			path.Start = p.parsePath()
		}

		match.Paths = append(match.Paths, path)

		log.Print("added path to paths: ", match.Paths)

		if p.tok.typ == itemComma {
			p.expectType(itemComma)
			continue
		}
		break
	}
	return match
}

func (p *parser) parsePath() *Node {
	node := p.parseNode()

	currentNode := node

	for p.tok.typ == itemRelDir {
		rel := p.parseRelation()
		rel.LeftNode = currentNode
		rel.RightNode = p.parseNode()

		currentNode.RightRel = rel

		currentNode = currentNode.RightRel.RightNode

		currentNode.LeftRel = rel
	}

	return node
}

func (p *parser) parseNode() *Node {
	node := new(Node)
	p.expectType(itemLParen)

	node.Name = p.tok.val
	p.expectType(itemIdentifier)

	for p.tok.typ == itemColon {
		p.expectType(itemColon)
		node.Labels = append(node.Labels, p.tok.val)
		p.expectType(itemIdentifier)
	}

	if p.tok.typ == itemLBrace {
		p.expectType(itemLBrace)

		node.Props = make(map[string]string)

		for p.tok.typ == itemIdentifier {
			key := p.tok.val
			p.expectType(itemIdentifier)
			p.expectType(itemColon)
			val := p.tok.val[1 : len(p.tok.val)-1]
			p.expectType(itemString)

			node.Props[key] = val

			if p.tok.typ != itemComma {
				break
			}
		}

		p.expectType(itemRBrace)
	}

	p.expectType(itemRParen)

	return node
}

func (p *parser) parseRelation() *Relation {
	rel := new(Relation)

	rel.Direction = p.tok.val
	p.expectType(itemRelDir)

	if p.tok.typ == itemLBracket {
		p.expectType(itemLBracket)

		if p.tok.typ == itemIdentifier {
			rel.Name = p.tok.val
			p.expectType(itemIdentifier)
		}

		if p.tok.typ == itemColon {
			p.expectType(itemColon)
			rel.Types = make([]string, 0, 1)

			for p.tok.typ == itemIdentifier {
				rel.Types = append(rel.Types, p.tok.val)
				p.expectType(itemIdentifier)
				if p.tok.typ != itemPipe {
					break
				}
				p.expectType(itemPipe)
			}
		}

		if p.tok.typ == itemStar {
			p.expectType(itemStar)
			// TODO: implement properly
		}

		p.expectType(itemRBracket)
	}

	if rel.Direction == "<-" && p.tok.val == "->" {
		p.error("Relation has to point only in one direction or be undirected")
	} else if rel.Direction == "-" {
		rel.Direction = p.tok.val
	}
	p.expectType(itemRelDir)

	return rel
}

func (p *parser) parseQuery() *Query {

	query := new(Query)

	for p.tok.typ != itemEOF {
		switch p.tok.typ {
		case itemStart:
			p.expect("start")
			query.Roots = p.parseStart()
		case itemMatch:
			p.expect("match")
			query.Match = p.parseMatch()
		case itemDelete:
			p.expect("delete")
			query.Deletes = p.parseDelete()
		case itemCreate:
			p.expect("create")
			query.Creates = p.parseCreate()
		case itemReturn:
			p.expect("return")
			query.Returns = p.parseReturns()
		default:
			p.error("unknown top level type: " + p.tok.String())
			return nil
		}
	}

	return query
}

func (p *parser) parse(filename string, channel chan item) *Query {
	p.scanner = channel

	p.next() // initializes first token

	query := p.parseQuery()

	return query
}

func Parse(filename string, src string) (*Query, error) {
	var p parser

	_, channel := lex(filename, src)
	query := p.parse(filename, channel)

	if query == nil {
		p.error("Invalid query, Lexing might have failed")
	}

	if len(p.errors) == 0 {
		return query, nil
	}

	return query, p.errors
}
