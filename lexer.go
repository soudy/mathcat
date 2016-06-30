// Copyright 2016 Steven Oud. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be found
// in the LICENSE file.

package mathcat

import (
	"fmt"
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
	if !isIdent(rune(s[0])) {
		return false
	}

	for _, c := range s {
		if !isIdent(c) && !isNumber(c) {
			return false
		}
	}

	return true
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
				l.switchEq(ADD, ADD_EQ)
			case '-':
				// Check for unary minus. We decide unaryness at lexer level to
				// make it easier for the parser to know the difference.
				if l.isNegation() && l.peek() != '=' {
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
					return nil, fmt.Errorf("Invalid operation ‘%s’", string(l.ch))
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
				return nil, fmt.Errorf("Invalid token ‘%s’", string(l.ch))
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

func (l *lexer) emit(toktype TokenType) {
	l.tokens = append(l.tokens, &Token{
		Type:  toktype,
		Value: string(l.expr[l.start:l.pos]),
		Pos:   l.start,
	})
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
	if l.ch == '0' {
		// Hex literals
		if l.peek() == 'x' || l.peek() == 'X' {
			l.eat()

			for isHex(l.peek()) {
				l.eat()
			}

			l.emit(HEX)
			return
		}

		// Binary literals
		if l.peek() == 'b' || l.peek() == 'B' {
			l.eat()

			for isBinary(l.peek()) {
				l.eat()
			}

			l.emit(BINARY)
			return
		}

		// Octal literals
		if l.peek() == 'o' || l.peek() == 'O' {
			l.eat()

			for isOctal(l.peek()) {
				l.eat()
			}

			l.emit(OCTAL)
			return
		}
	}

	// Numeral literals
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

func (l *lexer) switchEq(tokA, tokB TokenType) {
	if l.peek() == '=' {
		l.eat()
		l.emit(tokB)
	} else {
		l.emit(tokA)
	}
}
