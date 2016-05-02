// Copyright 2016 Steven Oud. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be found
// in the LICENSE file.

package evaler

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
	tokens []*token // tokenized lexemes
}

func isIdent(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_' || (c >= 0x80 && unicode.IsLetter(c))
}

func isNumber(c rune) bool {
	return (c >= '0' && c <= '9') || c == '.'
}

// Lex starts lexing an expression. We keep reading until EOL is found, which
// we add because we need a padding of 1 to always be able to peek().
//
// Returns the generated tokens and any error found.
func Lex(expr string) ([]*token, error) {
	l := &lexer{
		expr:  append([]rune(expr), eol), // add eol as padding
		pos:   0,
		start: 0,
	}

	return l.lex()
}

func (l *lexer) lex() ([]*token, error) {
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
					l.switchEq(LSH, LSH_EQ)
				} else {
					l.emit(ILLEGAL)
					l.eat()
					return nil, errors.New("expected <<, got <")
				}
			case '>':
				if l.peek() == '>' {
					l.eat()
					l.switchEq(RSH, RSH_EQ)
				} else {
					l.emit(ILLEGAL)
					l.eat()
					return nil, errors.New("expected >>, got >")
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
				return nil, errors.New("unexpected token " + string(l.ch))
			}
		}
	}

	return l.tokens, nil
}

func (l *lexer) peek() rune {
	return l.expr[l.pos]
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
	for l.peek() == '\t' || l.peek() == ' ' || l.peek() == '\r' {
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
	toktype := INT
	for isNumber(l.peek()) {
		if l.ch == '.' {
			toktype = FLOAT
		}
		l.eat()
	}

	l.emit(toktype)
}

func (l *lexer) switchEq(tokA, tokB tokenType) {
	if l.peek() == '=' {
		l.eat()
		l.emit(tokB)
	} else {
		l.emit(tokA)
	}
}
