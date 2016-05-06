// Copyright 2016 Steven Oud. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be found
// in the LICENSE file.

package evaler

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

var (
	divisionByZeroErr       = errors.New("Divison by zero")
	unmatchedParenthesesErr = errors.New("Unmatched parentheses")
)

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
//     p := evaler.New()
//     p.Run("a = 150")
//     p.Run("b = 715")
//     res, err := p.Exec("a**b - (a/b)")
func New() *Parser {
	return &Parser{
		pos:       0,
		Variables: constants,
	}
}

// Eval evaluates an expression and returns its result and any errors found.
//
// Example:
//     res, err := evaler.Eval("2 * 2 * 2") // 8
func Eval(expr string) (float64, error) {
	tokens, err := Lex(expr)

	// If a lexer error occured don't parse
	if err != nil {
		return -1, err
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
func (p *Parser) Run(expr string) error {
	return nil
}

// Exec executes an expression and returns the result.
func (p *Parser) Exec(expr string) (float64, error) {
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

			if !operators.Empty() {
				var ok bool
				if o2, ok = allOperators[operators.Top().(*token).Type]; !ok {
					operators.Push(p.tok)
					break
				}

				if o2.hasHigherPrecThan(o1) {
					operator := operators.Pop().(*token)
					val, err := p.evaluate(operator, &operands)
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
					return -1, unmatchedParenthesesErr
				}

				tok := operators.Pop().(*token)
				if tok.Type == LPAREN {
					break
				}
				operands.Push(tok)
			}
		}

	}

	// Evaluate remaing operators
	for !operators.Empty() {
		operator := operators.Pop().(*token)

		if operator.Type == LPAREN {
			return -1, unmatchedParenthesesErr
		}

		val, err := p.evaluate(operator, &operands)
		if err != nil {
			return -1, err
		}
		operands.Push(val)
	}

	return operands[0].(float64), nil
}

func (p *Parser) evaluate(operator *token, operands *stack) (float64, error) {
	var result float64
	var left, right float64
	var err error

	if right, err = p.lookup(operands.Pop()); err != nil {
		return -1, err
	}

	// Save the token in case of a assignment variable is used so we need to
	// save the result in a variable.
	lhsIdent := operands.Pop()
	if left, err = p.lookup(lhsIdent); err != nil {
		return -1, err
	}

	switch operator.Type {
	case ADD, SUB, DIV, MUL, POW, REM, AND, OR, XOR, LSH, RSH, NOT:
		result, err = execute(operator.Type, left, right)
		if err != nil {
			return -1, err
		}
	case ADD_EQ, SUB_EQ, DIV_EQ, MUL_EQ, POW_EQ, REM_EQ, AND_EQ, OR_EQ, XOR_EQ, LSH_EQ, RSH_EQ:
		result, err = execute(operator.Type, left, right)
		if err != nil {
			return -1, err
		}

		// Save result in variable
		p.Variables[lhsIdent.(*token).Value] = result
	}

	return result, nil
}

func execute(operator tokenType, lhs, rhs float64) (float64, error) {
	var result float64
	switch operator {
	case ADD:
		result = lhs + rhs
	case SUB:
		result = lhs - rhs
	case DIV:
		if rhs == 0 {
			return -1, divisionByZeroErr
		}
		result = lhs / rhs
	case MUL:
		result = lhs * rhs
	case POW:
		result = math.Pow(lhs, rhs)
	case REM:
		if rhs == 0 {
			return -1, divisionByZeroErr
		}
		result = math.Mod(lhs, rhs)
	case AND:
		// TODO: check for int with bitwise operators
		// result = lhs & rhs
	case OR:
		// result = lhs | rhs
	case XOR:
		// result = lhs ^ rhs
	case LSH:
		// result = lhs << rhs
	case RSH:
		// result = lhs >> rhs
	}

	return result, nil
}

// Look up a literal. If it's an identifier, check the parser's variables map,
// otherwise convert the tokenized string to a float64.
func (p *Parser) lookup(val interface{}) (float64, error) {
	// val can be a token or a float64, if it's a float64 it has been already
	// evaluated and we don't need to do anything
	if v, ok := val.(float64); ok {
		return v, nil
	}

	tok := val.(*token)
	switch tok.Type {
	case INT, FLOAT:
		return strconv.ParseFloat(tok.Value, 64)
	case IDENT:
		res, err := p.GetVar(tok.Value)
		if err != nil {
			return -1, err
		}

		return res, nil
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
