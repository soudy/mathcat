// Copyright 2016 Steven Oud. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be found
// in the LICENSE file.

package mathcat

import (
	"errors"
	"unicode"
)

const eol rune = -1

// lexer holds the lexer's state while scanning an expression. If any error
// occurs, the scanning stops immediatly and returns the error.
type lexer struct {
	expr   []rune   // the input expression
	ch     rune     // current character
	pos    int      // current character position
	start  int      // current read offset
	tokens []*Token // tokenized lexemes
}

// https://en.wikipedia.org/wiki/Whitespace_character
var whitespaceChars = []rune{
	'\u0009', '\u000A', '\u000B', '\u000C', '\u000D', '\u0020', '\u0085',
	'\u00A0', '\u1680', '\u2000', '\u2001', '\u2002', '\u2003', '\u2004',
	'\u2005', '\u2006', '\u2007', '\u2008', '\u2009', '\u200A', '\u2028',
	'\u2029', '\u202F', '\u205F', '\u3000', '\u180E', '\u200B', '\u200C',
	'\u200D', '\u2060', '\uFEFF',
}

func isIdent(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_' || (c >= 0x80 && unicode.IsLetter(c))
}

func isNumber(c rune) bool {
	return (c >= '0' && c <= '9') || c == '.'
}

func isHex(c rune) bool {
	return (c >= '0' && c <= '9') || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

func isBinary(c rune) bool {
	return c == '0' || c == '1'
}

func isWhitespace(c rune) bool {
	if c == '\t' || c == ' ' || c == '\r' {
		return true
	}

	for _, v := range whitespaceChars {
		if v == c {
			return true
		}
	}

	return false
}

// Lex starts lexing an expression. We keep reading until EOL is found, which
// we add because we need a padding of 1 to always be able to peek().
//
// Returns the generated tokens and any error found.
func Lex(expr string) ([]*Token, error) {
	l := &lexer{
		expr:  append([]rune(expr), eol), // add eol as padding
		pos:   0,
		start: 0,
	}

	return l.lex()
}

func (l *lexer) lex() ([]*Token, error) {
loop:
	for l.ch != eol {
		l.start = l.pos

		l.eat()

		switch {
		case isIdent(l.ch):
			l.readIdent()
		case isNumber(l.ch):
			l.readNumber()
		case isWhitespace(l.ch):
			l.skipWhitespace()
		default:
			switch l.ch {
			case '+':
				l.switchEq(ADD, ADD_EQ)
			case '-':
				if l.isNegation() {
					l.emit(UNARY_MIN)
					break
				}
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
					l.switchEq(LSH, LSH_EQ)
				} else {
					l.switchEq(LT, LT_EQ)
				}
			case '>':
				if l.peek() == '>' {
					l.eat()
					l.switchEq(RSH, RSH_EQ)
				} else {
					l.switchEq(GT, GT_EQ)
				}
			case '~':
				l.emit(NOT)
			case '=':
				l.switchEq(EQ, EQ_EQ)
			case '!':
				if l.peek() != '=' {
					return nil, errors.New("Invalid operation " + string(l.ch))
				}
				l.eat()
				l.emit(BANG_EQ)
			case '(':
				l.emit(LPAREN)
			case ')':
				l.emit(RPAREN)
			case ',':
				l.emit(COMMA)
			case '#':
				// Comment, stop scanning for tokens
				l.emit(EOL)
				break loop
			case eol:
				l.emit(EOL)
			default:
				l.emit(ILLEGAL)
				return nil, errors.New("Invalid token " + string(l.ch))
			}
		}
	}

	return l.tokens, nil
}

func (l *lexer) peek() rune {
	return l.expr[l.pos]
}

func (l *lexer) prev() *Token {
	return l.tokens[len(l.tokens)-1]
}

func (l *lexer) eat() rune {
	l.ch = l.peek()
	l.pos++
	return l.ch
}

func (l *lexer) emit(toktype tokenType) {
	l.tokens = append(l.tokens, newToken(toktype, string(l.expr[l.start:l.pos]), l.start))
}

func (l *lexer) skipWhitespace() {
	for isWhitespace(l.peek()) {
		l.eat()
	}
}

func (l *lexer) readIdent() {
	for isIdent(l.peek()) || isNumber(l.peek()) {
		l.eat()
	}

	l.emit(IDENT)
}

func (l *lexer) readNumber() {
	// Hex literals
	if l.ch == '0' && (l.peek() == 'x' || l.peek() == 'X') {
		l.eat()

		for isHex(l.peek()) {
			l.eat()
		}

		l.emit(HEX)
		return
	}

	// Binary literals
	if l.ch == '0' && (l.peek() == 'b' || l.peek() == 'B') {
		l.eat()

		for isBinary(l.peek()) {
			l.eat()
		}

		l.emit(BINARY)
		return
	}

	// Normal literals
	for isNumber(l.peek()) || l.peek() == 'e' || l.peek() == 'E' {
		l.eat()
		if (l.ch == 'e' || l.ch == 'E') && l.peek() == '-' {
			l.eat()
		}
	}

	l.emit(NUMBER)
}

func (l *lexer) isNegation() bool {
	return l.tokens == nil || l.prev().Type == LPAREN || l.prev().IsOperator()
}

func (l *lexer) switchEq(tokA, tokB tokenType) {
	if l.peek() == '=' {
		l.eat()
		l.emit(tokB)
	} else {
		l.emit(tokA)
	}
}
