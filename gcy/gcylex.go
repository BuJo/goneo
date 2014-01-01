package gcy

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

type itemType int
type item struct {
	typ itemType
	val string
}

const (
	itemError itemType = iota

	itemIdentifier
	itemField

	itemNumber

	// symbols
	itemLParen
	itemRParen
	itemLBrace
	itemRBrace
	itemLBracket
	itemRBracket
	itemComma
	itemColon
	itemEqual
	itemStar
	itemPipe

	// ..
	itemRange

	// Relationship directions
	itemLRelDir
	itemRRelDir

	// Arithmetic
	itemMinus
	itemPlus

	// main query block keywords
	itemKeyword
	itemStart
	itemCreate
	itemDelete
	itemMatch
	itemReturn
	itemWith
	itemAs

	itemEOF
)

var key = map[string]itemType{
	"start":  itemStart,
	"create": itemCreate,
	"delete": itemDelete,
	"match":  itemMatch,
	"return": itemReturn,
	"with":   itemWith,
	"as":     itemAs,
}

// (partial) Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

func (i item) String() string {
	switch {
	case i.typ == itemEOF:
		return "EOF"
	case i.typ == itemError:
		return i.val
	case i.typ == itemNumber:
		if i, e := strconv.ParseInt(i.val, 0, 64); e == nil {
			return fmt.Sprintf("%d", i)
		}
		if f, e := strconv.ParseFloat(i.val, 64); e == nil {
			return fmt.Sprintf("%f", f)
		}
	case i.typ > itemKeyword:
		return fmt.Sprintf("<%s>", i.val)
	case len(i.val) > 10:
		return fmt.Sprintf("%.10q...", i.val)
	}
	s := fmt.Sprintf("%q", i.val)
	if len(s) == 0 || s == "\"\"" {
		s = fmt.Sprintf("[%d]", i.typ)
	}
	return s
}

type lexer struct {
	name  string    // only for error reporting
	input string    // string being scanned
	start int       // start position of this item
	pos   int       // current position in input
	width int       // width of last item read
	items chan item // channel of scanned items

	parenDepth int
}

// represents the state of the scanner as a function returning the next state
type stateFn func(*lexer) stateFn

// run executes the lexer
func (l *lexer) run() {
	for state := lexGcy; state != nil; {
		////fmt.Printf("state change\n")
		state = state(l)
	}
	close(l.items) // no more tokens
}

// emit passes an item back to the client
func (l *lexer) emit(t itemType) {
	l.items <- item{t, l.input[l.start:l.pos]}

	//fmt.Printf("emitted item from %d to %d: %s\n", l.start, l.pos, item{t, l.input[l.start:l.pos]})

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
	//fmt.Printf("error at: %d(%s)\n", l.pos, l.input[l.pos:])
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
	l.ignore()
	if numSkipped > 0 {
		return true
	}
	return false
}
func (l *lexer) atBoundary() bool {
	r := l.peek()
	if isSpace(r) || isEndOfLine(r) {
		return true
	}
	switch r {
	case eof, '.', ',', '|', ':', ')', '(', '{', '}', '[', ']', '=', '*', 0xFFFD:
		return true
	}

	return false
}

// peek returns but does not consume the next rune in the input.
func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

const (
	markIdentifier = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_"
	markNumbers    = "1234567890"
	markSpace      = "\n 	"
	markEOL        = "\n"
	eof            = 0
)

func lexGcy(l *lexer) stateFn {
	if l.pos >= len(l.input) {
		//fmt.Println("end lexing Gcy")
		l.emit(itemEOF)
		return nil
	}

	//fmt.Printf("peek: %c\n", l.peek())

	switch r := l.next(); {
	case r == eof || isEndOfLine(r):
		return l.errorf("unclosed action")
	case isSpace(r):
		return lexSpace
	case r == '=':
		l.emit(itemEqual)
	case r == ',':
		l.emit(itemComma)
	case r == '*':
		l.emit(itemStar)
	case r == ':':
		l.emit(itemColon)
	case r == '|':
		l.emit(itemPipe)
	case r == '.':
		if l.peek() == '.' {
			l.next()
			l.emit(itemRange)
			return lexGcy
		}
		return lexField
	case r == '<':
		l.next()
		l.emit(itemLRelDir)
	case r == '-':
		p := l.peek()
		if r == '-' && (p == '[' || p == '(' || p == '-' || p == '>') {
			if p == '>' {
				l.next()
			}
			l.emit(itemRRelDir)
			return lexGcy
		}
		fallthrough
	case r == '+' || r == '-' || ('0' <= r && r <= '9'):
		l.backup()
		return lexNumber
	case isAlphaNumeric(r):
		l.backup()
		return lexIdentifier
	case r == '(':
		l.emit(itemLParen)
		l.parenDepth++
		return lexGcy
	case r == ')':
		l.emit(itemRParen)
		l.parenDepth--
		if l.parenDepth < 0 {
			return l.errorf("unexpected right paren %#U", r)
		}
		return lexGcy
	case r == '[':
		l.emit(itemLBracket)
		l.parenDepth++
		return lexGcy
	case r == ']':
		l.emit(itemRBracket)
		l.parenDepth--
		if l.parenDepth < 0 {
			return l.errorf("unexpected right paren %#U", r)
		}
		return lexGcy
	case r == '{':
		l.emit(itemLBrace)
		l.parenDepth++
		return lexGcy
	case r == '}':
		l.emit(itemRBrace)
		l.parenDepth--
		if l.parenDepth < 0 {
			return l.errorf("unexpected right paren %#U", r)
		}
		return lexGcy
	default:
		return l.errorf("unrecognized character in action: %#U", r)
	}
	return lexGcy
}

func lexIdentifier(l *lexer) stateFn {
	//fmt.Println("lexing ident: ", l.input[l.start:l.pos])

Loop:
	for {
		switch r := l.next(); {
		case isAlphaNumeric(r):
			// absorb.
		default:
			l.backup()
			word := l.input[l.start:l.pos]
			if !l.atBoundary() {
				return l.errorf("bad character %#U", r)
			}
			switch {
			case key[word] > itemKeyword:
				l.emit(key[word])
			case word[0] == '.':
				l.emit(itemField)
			default:
				l.emit(itemIdentifier)
			}
			break Loop
		}
	}

	return lexGcy
}

// lexSpace scans a run of space characters.
// One space has already been seen.
func lexSpace(l *lexer) stateFn {
	l.skipSpace()
	return lexGcy
}

// lexField scans a field: .Alphanumeric.
// The . has been scanned.
func lexField(l *lexer) stateFn {
	return lexFieldOrVariable(l, itemField)
}

// lexVariable scans a Variable: $Alphanumeric.
// The $ has been scanned.
func lexVariable(l *lexer) stateFn {
	if l.atBoundary() { // Nothing interesting follows -> "$".
		l.emit(itemIdentifier)
		return lexGcy
	}
	return lexFieldOrVariable(l, itemIdentifier)
}

// lexVariable scans a field or variable: [.$]Alphanumeric.
// The . or $ has been scanned.
func lexFieldOrVariable(l *lexer, typ itemType) stateFn {
	if l.atBoundary() { // Nothing interesting follows -> "." or "$".
		if typ == itemIdentifier {
			l.emit(itemIdentifier)
		} else {
			l.emit(itemIdentifier)
		}
		return lexGcy
	}
	var r rune
	for {
		r = l.next()
		if !isAlphaNumeric(r) {
			l.backup()
			break
		}
	}
	if !l.atBoundary() {
		return l.errorf("bad character %#U", r)
	}
	l.emit(typ)
	return lexGcy
}

// lexNumber scans a number: decimal, octal, hex, float, or imaginary. This
// isn't a perfect number scanner - for instance it accepts "." and "0x0.2"
// and "089" - but when it's wrong the input is invalid and the parser (via
// strconv) will notice.
func lexNumber(l *lexer) stateFn {
	if !l.scanNumber() {
		return l.errorf("bad number syntax: %q", l.input[l.start:l.pos])
	}
	l.emit(itemNumber)
	return lexGcy
}

func (l *lexer) scanNumber() bool {
	// Optional leading sign.
	l.accept("+-")
	// Is it hex?
	digits := "0123456789"
	if l.accept("0") && l.accept("xX") {
		digits = "0123456789abcdefABCDEF"
	}
	l.acceptRun(digits)
	if l.accept(".") {
		// special case for ranges
		if l.peek() == '.' {
			l.backup()
		} else {
			l.acceptRun(digits)
		}
	}
	if l.accept("eE") {
		l.accept("+-")
		l.acceptRun("0123456789")
	}
	// Is it imaginary?
	l.accept("i")
	// Next thing mustn't be alphanumeric.
	if isAlphaNumeric(l.peek()) {
		l.next()
		return false
	}
	return true
}

// isSpace reports whether r is a space character.
func isSpace(r rune) bool {
	return r == ' ' || r == '\t'
}

// isEndOfLine reports whether r is an end-of-line character.
func isEndOfLine(r rune) bool {
	return r == '\r' || r == '\n'
}

// isAlphaNumeric reports whether r is an alphabetic, digit, or underscore.
func isAlphaNumeric(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}
