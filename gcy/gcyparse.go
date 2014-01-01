package gcy

import (
	"errors"
	"fmt"
	"strconv"
)

type (
	query struct {
	}

	SearchQuery struct {
		query

		Roots   []*Root
		Match   *Match
		Returns []*Return
	}

	Root struct {
		Name   string
		Typ    string
		IdVars []int
	}

	Match struct {
	}

	Return struct {
		Name  string
		Alias string
	}

	GcyQuery interface{}
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

func (p *parser) parseSearchQuery() *SearchQuery {
	fmt.Println("parsing search query")

	query := &SearchQuery{}

	fmt.Println("first in q: ", p.tok)
	if p.tok.typ == itemStart {
		p.expect("start")

		query.Roots = p.parseRoots()
	}

	query.Returns = p.parseReturns()

	fmt.Println("returning from search query: ", query)

	return query
}

func (p *parser) parseRoots() []*Root {
	fmt.Println("parsing search query roots")

	roots := make([]*Root, 0)
	roots = append(roots, p.parseRoot())
	return roots
}

func (p *parser) parseRoot() (r *Root) {
	varname := p.tok.val
	p.expectType(itemIdentifier)
	p.expectType(itemEqual)

	r = &Root{Name: varname}

	switch p.tok.val {
	case "node", "relation":
		typ := p.tok.val
		p.expect(typ)
		p.expect("(")

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
				for i := start; i < end; i += 1 {
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
func (p *parser) parseReturns() []*Return {
	p.expect("return")

	rets := make([]*Return, 0)

	varname := p.tok.val
	alias := varname
	p.expectType(itemIdentifier)

	if p.tok.typ == itemAs {
		p.expectType(itemAs)

		alias = p.tok.val

		p.expectType(itemIdentifier)
	}

	rets = append(rets, &Return{varname, alias})

	return rets
}

func (p *parser) parseDeleteQuery() *SearchQuery {
	return nil
}

func (p *parser) parseCreateQuery() *SearchQuery {
	return nil
}

func (p *parser) parseQuery() GcyQuery {

	if p.tok.typ == itemStart || p.tok.typ == itemMatch {
		return p.parseSearchQuery()
	} else if p.tok.typ == itemDelete {
		return p.parseDeleteQuery()
	} else if p.tok.typ == itemCreate {
		return p.parseCreateQuery()
	}

	return nil
}

func (p *parser) parse(filename string, channel chan item) GcyQuery {
	p.scanner = channel

	p.next() // initializes first token

	query := p.parseQuery()
	return query
}

func Parse(filename string, src string) (GcyQuery, error) {
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
