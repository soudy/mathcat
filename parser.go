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

// Run executes an expression on an existing parser instance. Useful for
// variable assignment.
//
// Example:
//     p.Run("a = 555")
//     p.Run("a += 45")
//     res, err := p.Run("a + a") // 1200
func (p *Parser) Run(expr string) (float64, error) {
	tokens, err := Lex(expr)

	if err != nil {
		return -1, err
	}

	p.reset()
	p.tokens = tokens

	return p.parse()
}

// Exec executes an expression and returns the result.
func Exec(expr string) (float64, error) {
	return 0, nil
}

// GetVar gets an existing variable.
func (p *Parser) GetVar(index string) (float64, error) {
	if val, ok := p.Variables[index]; ok {
		return val, nil
	}

	return -1, fmt.Errorf("Undefined variable '%s'", index)
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

	// Save the token in case of a assignment variable is used and we need to
	// save the result in a variable.
	lhsIdent := operands.Pop()
	if left, err = p.lookup(lhsIdent); err != nil {
		return -1, err
	}

	result, err = execute(operator, left, right)
	if err != nil {
		return -1, err
	}

	switch operator.Type {
	case ADD_EQ, SUB_EQ, DIV_EQ, MUL_EQ, POW_EQ, REM_EQ, AND_EQ, OR_EQ, XOR_EQ, LSH_EQ, RSH_EQ:
		// Save result in variable
		p.Variables[lhsIdent.(*token).Value] = result
	}

	return result, nil
}

func execute(operator *token, lhs, rhs float64) (float64, error) {
	var result float64

	if operator.IsBitwise() && (!isWholeNumber(lhs) || !isWholeNumber(rhs)) {
		// Both lhs and rhs have to be whole numbers for bitwise operations
		return -1, fmt.Errorf("Unsupported type (float) for '%s'", operator.Type)
	}

	switch operator.Type {
	case ADD, ADD_EQ:
		result = lhs + rhs
	case SUB, SUB_EQ:
		result = lhs - rhs
	case DIV, DIV_EQ:
		if rhs == 0 {
			return -1, divisionByZeroErr
		}
		result = lhs / rhs
	case MUL, MUL_EQ:
		result = lhs * rhs
	case POW, POW_EQ:
		result = math.Pow(lhs, rhs)
	case REM, REM_EQ:
		if rhs == 0 {
			return -1, divisionByZeroErr
		}
		result = math.Mod(lhs, rhs)
	case AND, AND_EQ:
		result = float64(int64(lhs) & int64(rhs))
	case OR, OR_EQ:
		result = float64(int64(lhs) | int64(rhs))
	case XOR, XOR_EQ:
		result = float64(int64(lhs) ^ int64(rhs))
	case LSH, LSH_EQ:
		result = float64(uint64(lhs) << uint64(rhs))
	case RSH, RSH_EQ:
		result = float64(uint64(lhs) >> uint64(rhs))
	default:
		return -1, fmt.Errorf("Invalid operator '%s'", operator)
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
	case NUMBER:
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

func (p *Parser) reset() {
	p.tokens = nil
	p.pos = 0
}

func (p *Parser) peek() *token {
	return p.tokens[p.pos]
}

func (p *Parser) eat() *token {
	p.tok = p.peek()
	p.pos++
	return p.tok
}

func isWholeNumber(n float64) bool {
	epsilon := 1e-9
	_, frac := math.Modf(math.Abs(n))

	return frac < epsilon || frac > 1.0-epsilon
}
