// Copyright 2016 Steven Oud. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be found
// in the LICENSE file.

package mathcat

import (
	"fmt"
	"strings"
	"unicode"
)

// eol indicates the end of an expression
const eol rune = -1

// lexer holds the lexer's state while scanning an expression. If any error
// occurs, the scanning stops immediately and returns the error.
type lexer struct {
	expr   []rune // the input expression
	ch     rune   // current character
	pos    int    // current character position
	start  int    // current read offset
	tokens Tokens // tokenized lexemes
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

func isOctal(c rune) bool {
	return c >= '0' && c <= '7'
}

func isWhitespace(c rune) bool {
	return c == '\t' || c == ' ' || c == '\r' || c == '\n'
}

// IsValidIdent checks if a string qualifies as a valid identifier.
func IsValidIdent(s string) bool {
	checkIdent := func(c rune) bool { return isIdent(c) || isNumber(c) }

	return isIdent(rune(s[0])) && strings.IndexFunc(s, checkIdent) != -1
}

// Lex starts lexing an expression, converting an input string into a stream
// of tokens later passed on to the parser.
//
// Returns the generated tokens and any error found.
func Lex(expr string) (Tokens, error) {
	l := &lexer{
		expr:  append([]rune(expr), eol), // add eol as padding
		pos:   0,
		start: 0,
	}

	return l.lex()
}

func (l *lexer) lex() (Tokens, error) {
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
				l.switchEq(Add, AddEq)
			case '-':
				// Check for unary minus. We decide unaryness at lexer level to
				// make it easier for the parser to know the difference.
				if l.isNegation() && l.peek() != '=' {
					l.emit(UnaryMin)
					break
				}
				l.switchEq(Sub, SubEq)
			case '/':
				l.switchEq(Div, DivEq)
			case '*':
				if l.peek() == '*' {
					l.eat()
					l.switchEq(Pow, PowEq)
				} else {
					l.switchEq(Mul, MulEq)
				}
			case '%':
				l.switchEq(Rem, RemEq)
			case '&':
				l.switchEq(And, AndEq)
			case '|':
				l.switchEq(Or, OrEq)
			case '^':
				l.switchEq(Xor, XorEq)
			case '<':
				if l.peek() == '<' {
					l.eat()
					l.switchEq(Lsh, LshEq)
				} else {
					l.switchEq(Lt, LtEq)
				}
			case '>':
				if l.peek() == '>' {
					l.eat()
					l.switchEq(Rsh, RshEq)
				} else {
					l.switchEq(Gt, GtEq)
				}
			case '~':
				l.emit(Not)
			case '=':
				l.switchEq(Eq, EqEq)
			case '!':
				if l.peek() != '=' {
					return nil, fmt.Errorf("Invalid operation ‘%s’", string(l.ch))
				}
				l.eat()
				l.emit(NotEq)
			case '(':
				l.emit(Lparen)
			case ')':
				l.emit(Rparen)
			case ',':
				l.emit(Comma)
			case '#', eol:
				// Comment or EOL, stop scanning for tokens
				l.emit(Eol)
				break loop
			default:
				return nil, fmt.Errorf("Invalid token ‘%s’", string(l.ch))
			}
		}
	}

	return l.tokens, nil
}

func (l lexer) peek() rune {
	return l.expr[l.pos]
}

func (l lexer) prev() *Token {
	return l.tokens[len(l.tokens)-1]
}

func (l *lexer) eat() rune {
	l.ch = l.peek()
	l.pos++
	return l.ch
}

func (l *lexer) emit(toktype TokenType) {
	l.tokens = append(l.tokens, &Token{
		Type:  toktype,
		Value: string(l.expr[l.start:l.pos]),
		Pos:   l.start,
	})
}

func (l lexer) skipWhitespace() {
	for isWhitespace(l.peek()) {
		l.eat()
	}
}

func (l *lexer) readIdent() {
	for isIdent(l.peek()) || isNumber(l.peek()) {
		l.eat()
	}

	l.emit(Ident)
}

func (l *lexer) readNumber() {
	if l.ch == '0' {
		// Hex literals
		if l.peek() == 'x' || l.peek() == 'X' {
			l.eat()

			for isHex(l.peek()) {
				l.eat()
			}

			l.emit(Hex)
			return
		}

		// Binary literals
		if l.peek() == 'b' || l.peek() == 'B' {
			l.eat()

			for isBinary(l.peek()) {
				l.eat()
			}

			l.emit(Binary)
			return
		}

		// Octal literals
		if l.peek() == 'o' || l.peek() == 'O' {
			l.eat()

			for isOctal(l.peek()) {
				l.eat()
			}

			l.emit(Octal)
			return
		}
	}

	// Decimal literals
	for isNumber(l.peek()) || l.peek() == 'e' || l.peek() == 'E' {
		l.eat()
		if (l.ch == 'e' || l.ch == 'E') && l.peek() == '-' {
			l.eat()
		}
	}

	l.emit(Decimal)
}

func (l lexer) isNegation() bool {
	return l.tokens == nil || l.prev().Is(Lparen) || l.prev().IsOperator()
}

func (l *lexer) switchEq(tokA, tokB TokenType) {
	if l.peek() == '=' {
		l.eat()
		l.emit(tokB)
	} else {
		l.emit(tokA)
	}
}
