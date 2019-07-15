// Package parser implements parser for anko.
package parser

import (
	"fmt"
	"unicode"

	"github.com/tesujiro/taplGo/ch03/ast"
)

const (
	// EOF is short for End of file.
	EOF = -1
	// EOL is short for End of line.
	EOL = '\n'
)

// Error provides a convenient interface for handling runtime error.
// It can be Error interface with type cast which can call Pos().
type Error struct {
	Message  string
	Pos      ast.Position
	Filename string
	Fatal    bool
}

var EOF_FLAG bool
var traceLexer bool

// Error returns the error message.
func (e *Error) Error() string {
	return e.Message
}

// Scanner stores informations for lexer.
type Scanner struct {
	src      []rune
	offset   int
	lineHead int
	line     int
}

// opName is correction of operation names.
var opName = map[string]int{
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"then":   THEN,
	"else":   ELSE,
	"succ":   SUCC,
	"pred":   PRED,
	"iszero": ISZERO,
}

// Scan analyses token, and decide identify or literals.
func (s *Scanner) Scan() (tok int, lit string, pos ast.Position, err error) {
	//retry:
	s.skipBlank()
	pos = s.pos()
	s.peek()
	switch ch := s.peek(); {
	case isLetter(ch):
		lit, err = s.scanIdentifier()
		if err != nil {
			return
		}
		if name, ok := opName[lit]; ok {
			tok = name
		} else {
			//tok = IDENT
			err = fmt.Errorf("syntax error on '%v' at %v:%v", string(ch), pos.Line, pos.Column)
			tok = int(ch)
			lit = string(ch)
			return
		}
	case ch == '0':
		tok = ZERO
		lit = "0"
		s.next()
	case ch == EOF:
		if !EOF_FLAG {
			tok = int(';')
			lit = string(';')
			EOF_FLAG = true
		} else {
			tok = EOF
			EOF_FLAG = true
		}
	default:
		err = fmt.Errorf("syntax error on '%v' at %v:%v", string(ch), pos.Line, pos.Column)
		tok = int(ch)
		lit = string(ch)
		return
	}
	return
}

// isLetter returns true if the rune is a letter for identity.
func isLetter(ch rune) bool {
	return unicode.IsLetter(ch) || ch == '_'
}

// isDigit returns true if the rune is a number.
func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

// isBlank returns true if the rune is empty character..
func isBlank(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\r'
}

// peek returns current rune in the code.
func (s *Scanner) peek() rune {
	if s.reachEOF() {
		return EOF
	}
	return s.src[s.offset]
}

// next moves offset to next.
func (s *Scanner) next() {
	if !s.reachEOF() {
		if s.peek() == '\n' {
			s.lineHead = s.offset + 1
			s.line++
		}
		s.offset++
	}
}

// current returns the current offset.
func (s *Scanner) current() int {
	return s.offset
}

// offset sets the offset value.
func (s *Scanner) set(o int) {
	s.offset = o
}

// back moves back offset once to top.
func (s *Scanner) back() {
	s.offset--
}

// reachEOF returns true if offset is at end-of-file.
func (s *Scanner) reachEOF() bool {
	return len(s.src) <= s.offset
}

// pos returns the position of current.
func (s *Scanner) pos() ast.Position {
	return ast.Position{Line: s.line + 1, Column: s.offset - s.lineHead + 1}
}

func (s *Scanner) skipBlank() string {
	str := ""
	for ch := s.peek(); isBlank(ch); ch = s.peek() {
		str = fmt.Sprintf("%s%c", str, ch)
		s.next()
	}
	return str
}

// scanIdentifier returns identifier beginning at current position.
func (s *Scanner) scanIdentifier() (string, error) {
	var ret []rune
	for {
		if !isLetter(s.peek()) && !isDigit(s.peek()) {
			break
		}
		ret = append(ret, s.peek())
		s.next()
	}
	return string(ret), nil
}

// Lexer provides interface to parse codes.
type Lexer struct {
	s      *Scanner
	lit    string
	pos    ast.Position
	e      error
	result ast.Term
}

// Lex scans the token and literals.
func (l *Lexer) Lex(lval *yySymType) int {
	tok, lit, pos, err := l.s.Scan()
	if traceLexer {
		fmt.Printf("tok:%v\tlit:%v\tpos:%v\terr:%v\n", tok, lit, pos, err)
	}
	if err != nil {
		l.e = &Error{Message: err.Error(), Pos: pos, Fatal: true}
	}
	lval.token = ast.Token{Token: tok, Literal: lit}
	lval.token.SetPosition(pos)
	l.lit = lit
	l.pos = pos
	return tok
}

// Error sets parse error.
func (l *Lexer) Error(msg string) {
	l.e = &Error{Message: msg, Pos: l.pos, Fatal: false}
}

// Parse provides way to parse the code using Scanner.
func Parse(s *Scanner) (ast.Term, error) {
	l := Lexer{s: s}
	if yyParse(&l) != 0 {
		return nil, l.e
	}
	return l.result, l.e
}

/*
// EnableErrorVerbose enabled verbose errors from the parser
func EnableErrorVerbose() {
	yyErrorVerbose = true
}
*/

func TraceLexer() {
	traceLexer = true
}

func TraceOffLexer() {
	traceLexer = false
}

func initialize() {
	EOF_FLAG = false
}

// ParseSrc provides way to parse the code from source.
func ParseSrc(src string) (ast.Term, error) {
	initialize()
	scanner := &Scanner{
		src: []rune(src),
	}
	return Parse(scanner)
}
