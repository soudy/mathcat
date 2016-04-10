// Copyright 2016 Steven Oud. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be found
// in the LICENSE file.

package eparser

import (
	"errors"
	"fmt"
	"math"
	"strconv"
)

type Parser struct {
	tokens    []*token
	pos       int
	Variables map[string]float64
	tok       *token
}

type association int

const (
	ASSOC_NONE association = iota
	ASSOC_LEFT
	ASSOC_RIGHT
)

type operator struct {
	prec  int
	assoc association
	unary bool
}

var allOperators = map[tokenType]operator{
	// Assignment operators
	EQ:     {0, ASSOC_RIGHT, false}, // =
	ADD_EQ: {0, ASSOC_RIGHT, false}, // +=
	SUB_EQ: {0, ASSOC_RIGHT, false}, // -=
	DIV_EQ: {0, ASSOC_RIGHT, false}, // /=
	MUL_EQ: {0, ASSOC_RIGHT, false}, // *=
	POW_EQ: {0, ASSOC_RIGHT, false}, // **=
	REM_EQ: {0, ASSOC_RIGHT, false}, // %=
	AND_EQ: {0, ASSOC_RIGHT, false}, // &=
	OR_EQ:  {0, ASSOC_RIGHT, false}, // |=
	XOR_EQ: {0, ASSOC_RIGHT, false}, // ^=
	LSH_EQ: {0, ASSOC_RIGHT, false}, // <<=
	RSH_EQ: {0, ASSOC_RIGHT, false}, // >>=

	// Bitwise operators
	OR:  {1, ASSOC_RIGHT, false}, // |
	XOR: {2, ASSOC_RIGHT, false}, // ^
	AND: {3, ASSOC_RIGHT, false}, // &
	LSH: {4, ASSOC_RIGHT, false}, // <<
	RSH: {4, ASSOC_RIGHT, false}, // >>
	NOT: {8, ASSOC_LEFT, true},   // ~

	// Mathematical operators
	ADD: {5, ASSOC_LEFT, false}, // +
	SUB: {5, ASSOC_LEFT, false}, // -
	MUL: {6, ASSOC_LEFT, false}, // *
	DIV: {6, ASSOC_LEFT, false}, // /
	REM: {6, ASSOC_LEFT, false}, // %
	POW: {7, ASSOC_LEFT, false}, // **
}

// Determine if operator 1 has higher precendence than operator 2
func (o1 operator) hasHigherPrecThan(o2 operator) bool {
	return (o2.assoc == ASSOC_LEFT && o2.prec <= o1.prec) ||
		(o2.assoc == ASSOC_RIGHT && o2.prec < o1.prec)
}

// Some useful predefined variables that can be used in expressions. These
// can be overwritten.
var constants = map[string]float64{
	"pi":  math.Pi,
	"tau": math.Pi / 2,
	"phi": math.Phi,
	"e":   math.E,
}

// New initializes a new Parser instance, useful when you want to run multiple
// expression and/or use variables.
//
// For example, you could declare and use multiple variables like so:
//     p := eparser.New()
//     p.Run("a = 150")
//     p.Run("b = 715")
//     res, errs := p.Exec("a**b - (a/b)")
func New() *Parser {
	return &Parser{
		pos:       0,
		Variables: constants,
	}
}

// Parse evaluates an expression and returns its result and any errors found.
//
// Example:
//     res, errs := eparser.Parse("2 * 2 * 2") // 8
func Parse(expr string) (float64, error) {
	tokens, errs := Lex(expr)

	// If lexer errors occured don't parse
	if errs != nil {
		return -1, errs[0]
	}

	p := &Parser{
		tokens:    tokens,
		pos:       0,
		Variables: constants,
		tok:       tokens[0],
	}

	return p.parse()
}

// GetVar gets an existing variable.
func (p *Parser) GetVar(index string) (float64, error) {
	if val, ok := p.Variables[index]; ok {
		return val, nil
	}

	return -1, fmt.Errorf("Undefined variable '%s'", index)
}

// Run executes an expression but returns no result. Useful for variable
// assignment.
//
// Example:
//     p.Run("a = 555")
//     p.Run("a += 45")
//     p.Run("a + a") // does nothing
func (p *Parser) Run(expr string) []error {
	return nil
}

// Exec executes an expression and returns the result.
func (p *Parser) Exec(expr string) (float64, []error) {
	return 0, nil
}

func (p *Parser) parse() (float64, error) {
	var operands, operators stack
	var o1, o2 operator

	for p.eat().Type != EOL {
		switch {
		case p.tok.IsLiteral():
			operands.Push(p.tok)
		case p.tok.Type == LPAREN:
			operators.Push(p.tok)
		case p.tok.IsOperator():
			o1 = allOperators[p.tok.Type]

			// FIXME: I don't work correctly!
			if !operators.Empty() {
				var ok bool
				if o2, ok = allOperators[operators.Top().(*token).Type]; !ok {
					break
				}

				if o1.hasHigherPrecThan(o2) {
					operator := operators.Pop().(*token)
					left := operands.Pop().(*token)
					right := operands.Pop().(*token)

					val, err := p.evaluate(operator, left, right)
					if err != nil {
						return -1, err
					}

					operands.Push(val)

				}
			}
			operators.Push(p.tok)
		case p.tok.Type == RPAREN:
			for {
				if operators.Empty() {
					return -1, errors.New("unmatched parentheses")
				}

				tok := operators.Pop().(*token)
				if tok.Type == LPAREN {
					break
				}
				operands.Push(tok)
			}
		}
	}

	fmt.Print(operands)

	return 0, nil
}

func (p *Parser) evaluate(operator, left, right *token) (float64, error) {
	var nleft, nright float64
	var err error

	if nleft, err = p.lookup(left); err != nil {
		return -1, err
	}

	if nright, err = p.lookup(right); err != nil {
		return -1, err
	}

	var result float64

	switch operator.Type {
	case ADD:
		result = nleft + nright
	case SUB:
		result = nleft - nright
	case DIV:
		if nright == 0 {
			return -1, errors.New("divison by zero")
		}
		result = nleft / nright
	case MUL:
		result = nleft * nright
	case POW:
		result = math.Pow(nleft, nright)
	case REM:
		if nright == 0 {
			return -1, errors.New("divison by zero")
		}
		result = math.Mod(nleft, nright)
	}

	return result, nil
}

func (p *Parser) lookup(tok *token) (float64, error) {
	switch tok.Type {
	case INT, FLOAT:
		return strconv.ParseFloat(tok.Value, 64)
	case IDENT:
		val, err := p.GetVar(tok.Value)
		if err != nil {
			return -1, err
		}

		return val, nil
	}

	return -1, fmt.Errorf("Invalid lookup type: %s", tok.Type)
}

func (p *Parser) peek() *token {
	return p.tokens[p.pos]
}

func (p *Parser) eat() *token {
	p.tok = p.peek()
	p.pos++
	return p.tok
}
