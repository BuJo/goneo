package gcy

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type itemType int
type item struct {
	typ itemType
	val string
}

const (
	itemError itemType = iota // error, val contains description
	itemStart
	itemCreate
	itemDelete
	itemMatch
	itemAs
	itemReturn
	itemNode
	itemNodeStart
	itemNodeEnd
	itemNodeId
	itemRelStart
	itemRelEnd
	itemRelDirection
	itemAssign
	itemVariable
	itemRange
	itemEOF
)

func (i item) String() string {
	switch i.typ {
	case itemEOF:
		return "EOF"
	case itemError:
		return i.val
	}
	if len(i.val) > 10 {
		return fmt.Sprintf("%.10q...", i.val)
	}
	return fmt.Sprintf("%q", i.val)
}

type lexer struct {
	name  string    // only for error reporting
	input string    // string being scanned
	start int       // start position of this item
	pos   int       // current position in input
	width int       // width of last item read
	items chan item // channel of scanned items
}

// represents the state of the scanner as a function returning the next state
type stateFn func(*lexer) stateFn

// run executes the lexer
func (l *lexer) run() {
	for state := lexGcy; state != nil; {
		fmt.Printf("state change\n")
		state = state(l)
	}
	close(l.items) // no more tokens
}

// emit passes an item back to the client
func (l *lexer) emit(t itemType) {
	l.items <- item{t, l.input[l.start:l.pos]}

	fmt.Printf("emitted item from %d to %d: %s\n", l.start, l.pos, item{t, l.input[l.start:l.pos]})

	l.start = l.pos
}

func lex(name, input string) (*lexer, chan item) {
	l := &lexer{
		name:  name,
		input: input,
		items: make(chan item),
	}

	go l.run()

	return l, l.items
}

func (l *lexer) next() (r rune) {
	if l.pos > len(l.input) {
		l.width = 0
		return eof
	}

	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width

	return r
}
func (l *lexer) ignore() {
	l.start = l.pos
}
func (l *lexer) backup() {
	l.pos -= l.width
}

func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	fmt.Printf("error at: %d(%s)\n", l.pos, l.input[l.pos:])
	l.items <- item{itemError, fmt.Sprintf(format, args)}
	return nil
}

func (l *lexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}
func (l *lexer) acceptRun(valid string) (count int) {
	for strings.IndexRune(valid, l.next()) >= 0 {
		count += 1
	}
	l.backup()
	return
}
func (l *lexer) acceptUntil(valid string) bool {
	var v rune
	for v = l.next(); strings.IndexRune(valid, v) < 0; {
	}

	return v == 0x20
}
func (l *lexer) skipSpace() bool {
	numSkipped := l.acceptRun(markSpace)
	if numSkipped > 0 {
		l.ignore()
		return true
	}
	return false
}

func printEmitted(items chan item) {
	for i := range items {
		fmt.Printf("got item: %s\n", i)
	}
}

const (
	markIdentifier = "abcdefghijklmnopqrstuvwxyz_"
	markNumbers    = "1234567890"
	markRange      = "1234567890*."
	markNodeStart  = "("
	markNodeEnd    = ")"
	markStart      = "start"
	markNode       = "node"
	markSpace      = "\n 	"
	markEOL        = "\n"
	eof            = 0
)

func lexGcy(l *lexer) stateFn {
	fmt.Println("lexing Gcy")

	for {
		if strings.HasPrefix(l.input[l.pos:], "start") {
			return lexQuery
		} else if strings.HasPrefix(l.input[l.pos:], "return") {
			return lexReturn
		} else if strings.HasPrefix(l.input[l.pos:], "match") {
			return lexMatch
		}

		if l.next() == eof || l.pos >= len(l.input) {
			break
		}
		l.backup()
		fmt.Printf("itm:%q\n", l.next())
	}
	l.emit(itemEOF)
	return nil
}

func lexQuery(l *lexer) stateFn {
	fmt.Println("lexing start")

	l.acceptRun(markStart)
	l.emit(itemStart)

	l.skipSpace()
	if l.acceptRun(markIdentifier) < 1 {
		return l.errorf("expected variable name after start")
	}
	l.emit(itemVariable)

	l.skipSpace()
	if !l.accept("=") {
		return l.errorf("expected assignment after start var")
	}
	l.emit(itemAssign)

	l.skipSpace()
	if l.acceptRun(markNode) < 1 {
		return l.errorf("expected node after start var=")
	}
	l.emit(itemNode)

	l.skipSpace()
	if !l.accept(markNodeStart) {
		return l.errorf("expected ( after start var=node")
	}
	l.emit(itemNodeStart)

	if l.acceptRun(markRange) < 1 {
		return l.errorf("expected node id/range after start var=node(")
	}
	l.emit(itemRange)

	if !l.accept(markNodeEnd) {
		return l.errorf("expected ) after start var=node(id")
	}
	l.emit(itemNodeEnd)

	l.skipSpace()

	return lexGcy
}

func lexMatch(l *lexer) stateFn {
	fmt.Println("lexing match")

	l.acceptRun("match")
	l.emit(itemMatch)

	l.skipSpace()

	if l.acceptRun(markIdentifier) > 0 {
		l.emit(itemVariable)
		l.skipSpace()
		l.accept("=")
		l.skipSpace()
	}

	l.skipSpace()

	return lexGcy
}

func lexReturn(l *lexer) stateFn {
	fmt.Println("lexing Return")

	l.acceptRun(markIdentifier)
	l.emit(itemReturn)

	l.skipSpace()
	if l.acceptRun(markIdentifier) < 1 {
		return l.errorf("expected variable name after start")
	}
	l.emit(itemVariable)

	l.skipSpace()

	if l.acceptRun("as") > 0 {
		l.emit(itemAs)
		l.acceptRun(markSpace)
		l.ignore()
		if l.acceptRun(markIdentifier) < 1 {
			return l.errorf("expected variable name after start")
		}
		l.emit(itemVariable)
	}

	l.emit(eof)

	return nil
}
