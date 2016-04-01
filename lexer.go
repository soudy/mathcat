package eparser

import (
	"fmt"
	"unicode"
)

const eol rune = -1

// Lexer holds the lexer's state while scanning an expression. If an error is
// encountered it's appended to errors and the error count gets increased. If
// there are 5 errors already it will stop emitting them to prevent spam.
type Lexer struct {
	expr   []rune   // the input expression
	ch     rune     // current character
	pos    int      // current character offset
	start  int      // current read offset
	tokens []*Token // tokenized lexemes

	errors     []error // errors
	ErrorCount int     // error count
}

func isIdent(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_' || (c >= 0x80 && unicode.IsLetter(c))
}

func isNumber(c rune) bool {
	return c >= '0' && c <= '9'
}

func newLexer() *Lexer {
	return &Lexer{
		pos:   0,
		start: 0,

		errors:     nil,
		ErrorCount: 0,
	}
}

func (l *Lexer) error(msg string) {
	if l.ErrorCount > 4 {
		// At this point we're just spamming output
		return
	}
	l.ErrorCount++
	l.errors = append(l.errors, fmt.Errorf("Syntax Error: %s at position %d", msg, l.start+1))
}

// Lex starts lexing an expression. We keep reading until EOL is found, which
// we add because we need a padding of 1 to always be able to peek().
//
// Returns the generated tokens and any error found.
func (l *Lexer) Lex(expr string) ([]*Token, []error) {
	l.expr = append([]rune(expr), eol) // add eol as padding
	l.ch = []rune(expr)[0]
	for l.ch != eol {
		l.start = l.pos

		l.eat()

		switch {
		case isIdent(l.ch):
			l.readIdent()
		case isNumber(l.ch):
			l.readNumber()
		default:
			switch l.ch {
			case '+':
				l.switchEq(ADD, ADD_EQ)
			case '-':
				l.switchEq(SUB, SUB_EQ)
			case '/':
				l.switchEq(DIV, DIV_EQ)
			case '*':
				if l.peek() == '*' {
					l.eat()
					l.switchEq(POW, POW_EQ)
				} else {
					l.switchEq(MUL, MUL_EQ)
				}
			case '%':
				l.switchEq(REM, REM_EQ)
			case '&':
				l.switchEq(AND, AND_EQ)
			case '|':
				l.switchEq(OR, OR_EQ)
			case '^':
				l.switchEq(XOR, XOR_EQ)
			case '<':
				if l.peek() == '<' {
					l.eat()
					l.emit(LSH)
				} else {
					l.emit(ILLEGAL)
					l.eat()
					l.error("expected <<, got <")
				}
			case '>':
				if l.peek() == '>' {
					l.eat()
					l.emit(RSH)
				} else {
					l.emit(ILLEGAL)
					l.eat()
					l.error("expected >>, got >")
				}
			case '~':
				l.emit(NOT)
			case '=':
				l.emit(EQ)
			case '(':
				l.emit(LPAREN)
			case ')':
				l.emit(RPAREN)
			case '\r', ' ', '\t':
				l.skipWhitespace()
			case eol:
				l.emit(EOL)
			default:
				l.emit(ILLEGAL)
				l.error("unexpected token " + string(l.ch))
			}
		}
	}

	return l.tokens, l.errors
}

func (l *Lexer) peek() rune {
	return l.expr[l.pos]
}

func (l *Lexer) eat() rune {
	l.ch = l.peek()
	l.pos++
	return l.ch
}

func (l *Lexer) emit(toktype tokenType) {
	l.tokens = append(l.tokens, newToken(toktype, string(l.expr[l.start:l.pos]), l.start))
}

func (l *Lexer) skipWhitespace() {
	for l.peek() == '\t' || l.peek() == ' ' || l.peek() == '\r' {
		l.eat()
	}
}

func (l *Lexer) readIdent() {
	for isIdent(l.peek()) || isNumber(l.peek()) {
		l.eat()
	}

	l.emit(IDENT)
}

func (l *Lexer) readNumber() {
	toktype := INT
	for isNumber(l.peek()) || l.peek() == '.' {
		if l.ch == '.' {
			toktype = FLOAT
		}
		l.eat()
	}

	l.emit(toktype)
}

func (l *Lexer) switchEq(tokA, tokB tokenType) {
	if l.peek() == '=' {
		l.eat()
		l.emit(tokB)
	} else {
		l.emit(tokA)
	}
}
