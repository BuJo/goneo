package gcy

import (
	"errors"
	"fmt"
)

// ----------------------------------------------------------------------------
// Internal representation

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
		Name    string
		Typ     string
		IdRange string
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
		p.errorExpected("expected type, got type")
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
	p.expectType(itemVariable)
	p.expectType(itemAssign)

	r = &Root{Name: varname}

	if p.tok.val == "node" {
		p.expect("node")
		p.expect("(")
		r.Typ = "node"
		r.IdRange = p.tok.val
		p.expectType(itemRange)
	} else {
		p.expect("relation")
		p.expect("(")
		r.Typ = "relation"
		r.IdRange = p.tok.val
		p.expectType(itemRange)
	}
	p.expect(")")

	return r
}
func (p *parser) parseReturns() []*Return {
	p.expect("return")

	rets := make([]*Return, 0)

	varname := p.tok.val
	alias := varname
	p.expectType(itemVariable)

	if p.tok.typ == itemAs {
		p.expectType(itemAs)

		alias = p.tok.val

		p.expectType(itemVariable)
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

	return query, p.errors
}
