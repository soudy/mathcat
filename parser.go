// Copyright 2016 Steven Oud. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be found
// in the LICENSE file.

package mathcat

import (
	"errors"
	"fmt"
	"math"
	"math/big"
)

// Parser holds the lexed tokens, token position, declared variables and stacks
// used throughout the parsing of an expression.
//
// By default, variables always contains the constants defined below. These can
// however be overwritten.
type Parser struct {
	Tokens    Tokens
	Variables map[string]*big.Rat

	pos int
	tok *Token

	operands, operators, arity stack
}

var (
	// RatTrue represents true in boolean operations
	RatTrue = big.NewRat(1, 1)
	// RatFalse represents false in boolean operations
	RatFalse = big.NewRat(0, 1)

	ErrUnmatchedParentheses = errors.New("Unmatched parentheses")
	ErrMisplacedComma       = errors.New("Misplaced ‘,’")
	ErrAssignToLiteral      = errors.New("Can't assign to literal")

	defaultVariables = map[string]*big.Rat{
		"pi":    new(big.Rat).SetFloat64(math.Pi),
		"tau":   new(big.Rat).Mul(new(big.Rat).SetFloat64(math.Pi), big.NewRat(2, 1)),
		"phi":   new(big.Rat).SetFloat64(math.Phi),
		"e":     new(big.Rat).SetFloat64(math.E),
		"true":  RatTrue,
		"false": RatFalse,
	}
)

// New initializes a new Parser instance, useful when you want to run multiple
// expression and/or use variables.
func New() *Parser {
	parser := &Parser{}

	parser.Variables = make(map[string]*big.Rat)

	for k, v := range defaultVariables {
		parser.Variables[k] = v
	}

	return parser
}

// Eval evaluates an expression and returns its result and any errors found.
//
// Example:
//     res, err := mathcat.Eval("2 * 2 * 2") // 8
func Eval(expr string) (*big.Rat, error) {
	tokens, err := Lex(expr)

	// If a lexer error occurred don't parse
	if err != nil {
		return nil, err
	}

	p := New()
	p.Tokens = tokens

	return p.parse()
}

// Run executes an expression on an existing parser instance. Useful for
// variable assignment.
//
// Example:
//     p.Run("a = 555")
//     p.Run("a += 45")
//     res, err := p.Run("a + a") // 1200
func (p *Parser) Run(expr string) (*big.Rat, error) {
	tokens, err := Lex(expr)

	if err != nil {
		return nil, err
	}

	p.reset()
	p.Tokens = tokens

	return p.parse()
}

// Exec executes an expression with a given map of variables.
//
// Example:
//     res, err := mathcat.Exec("a + b * b", map[string]*big.Rat{
//         "a": big.NewRat(1, 1),
//         "b": big.NewRat(3, 1),
//     }) // 10
func Exec(expr string, vars map[string]*big.Rat) (*big.Rat, error) {
	tokens, err := Lex(expr)

	if err != nil {
		return nil, err
	}

	p := New()
	p.Tokens = tokens

	for name, val := range vars {
		if !IsValidIdent(name) {
			return nil, fmt.Errorf("Invalid variable name: ‘%s’", name)
		}
		p.Variables[name] = val
	}

	return p.parse()
}

// GetVar gets an existing variable.
//
// Example:
//     p.Run("酷 = -33")
//     if val, err := p.GetVar("酷"); !err {
//         fmt.Printf("%f\n", val) // -33
//     }
func (p Parser) GetVar(index string) (*big.Rat, error) {
	if val, ok := p.Variables[index]; ok {
		return val, nil
	}

	return nil, fmt.Errorf("Undefined variable ‘%s’", index)
}

func (p *Parser) parse() (*big.Rat, error) {
	// Initializing current token value
	p.tok = p.Tokens[0]

	for !p.eat().Is(Eol) {
		switch {
		case p.tok.IsLiteral():
			if p.peek().Is(Lparen) {
				// It's a function call, push to operators stack instead
				p.operators.Push(p.tok)

				// Check ahead if the function call has any argument at all, so
				// we can do accurate tracking of arity
				if p.peekN(2).Is(Rparen) {
					p.arity.Push(0)
				} else {
					p.arity.Push(1)
				}
				break
			}
			p.operands.Push(p.tok)
		case p.tok.Is(Lparen):
			p.operators.Push(p.tok)
		case p.tok.Is(Comma):
			for {
				if p.operators.Empty() {
					return nil, ErrMisplacedComma
				}

				if p.operators.Top().(*Token).Is(Lparen) {
					break
				}

				val, err := p.evaluate(p.operators.Pop().(*Token))
				if err != nil {
					return nil, err
				}

				p.operands.Push(val)
			}
			p.arity.Push(p.arity.Pop().(int) + 1)
		case p.tok.IsOperator():
			if err := p.handleOperator(); err != nil {
				return nil, err
			}
		case p.tok.Is(Rparen):
			for {
				if p.operators.Empty() {
					return nil, ErrUnmatchedParentheses
				}

				top := p.operators.Pop().(*Token)
				if top.Is(Lparen) {
					break
				}

				val, err := p.evaluate(top)
				if err != nil {
					return nil, err
				}

				p.operands.Push(val)
			}
		}
	}

	// Evaluate remaining operators
	for !p.operators.Empty() {
		top := p.operators.Pop().(*Token)

		if top.Is(Lparen) {
			return nil, ErrUnmatchedParentheses
		}

		val, err := p.evaluate(top)
		if err != nil {
			return nil, err
		}

		p.operands.Push(val)
	}

	// If there are no operands, the expression is useless and doesn't do
	// anything, for example `()`
	if p.operands.Empty() {
		return new(big.Rat), nil
	}

	// Single operand left means the expression was evaluated successful
	if len(p.operands) == 1 {
		return p.lookup(p.operands[0])
	}

	// Leftover token on operand stack indicates invalid syntax
	return nil, fmt.Errorf("Unexpected ‘%s’", p.operands.Top())
}

func (p *Parser) handleOperator() error {
	var o1, o2 operator

	o1 = operators[p.tok.Type]

	// No operators yet, just push to operators stack
	if p.operators.Empty() {
		p.operators.Push(p.tok)
		return nil
	}

	// While there's a function at the top of the operator stack, or an operator
	// with higher precedence than o1, pop operators to operands
	for p.operators.Top().(*Token).Is(Ident) || p.operators.Top().(*Token).IsOperator() {
		// Function call, always take precedence over operator
		if p.operators.Top().(*Token).Is(Ident) {
			function := p.operators.Pop().(*Token)
			val, err := p.evaluateFunc(function)
			if err != nil {
				return err
			}

			p.operands.Push(val)
		} else {
			o2 = operators[p.operators.Top().(*Token).Type]

			// Another operator at top, check precedence
			if o2.hasHigherPrecThan(o1) {
				operator := p.operators.Pop().(*Token)
				val, err := p.evaluateOp(operator)
				if err != nil {
					return err
				}
				p.operands.Push(val)
			} else {
				break
			}
		}

		if p.operators.Empty() {
			break
		}
	}

	p.operators.Push(p.tok)

	return nil
}

// evaluate gets called when an operator or function call has to be evaluated
// for a result. In case of a function, evaluateFunc is called and in case of
// an operator evaluateOp is called.
func (p *Parser) evaluate(tok *Token) (*big.Rat, error) {
	if tok.IsOperator() {
		return p.evaluateOp(tok)
	}

	return p.evaluateFunc(tok)
}

func (p *Parser) evaluateFunc(tok *Token) (*big.Rat, error) {
	var (
		function function
		ok       bool
		i        int
	)

	if function, ok = funcs[tok.Value]; !ok {
		return nil, fmt.Errorf("Undefined function ‘%s’", tok)
	}

	if arity := p.arity.Pop().(int); arity != function.arity {
		return nil, fmt.Errorf("Invalid argument count for ‘%s’ (expected %d, got %d)", tok, function.arity, arity)
	}

	// Start popping off arguments for the function call
	args := make([]*big.Rat, function.arity)
	for i = function.arity - 1; i >= 0; i-- {
		if p.operands.Empty() {
			return nil, ErrMisplacedComma
		}

		arg, err := p.lookup(p.operands.Pop())
		if err != nil {
			return nil, err
		}

		args[i] = arg
	}

	return function.fn(args), nil
}

func (p *Parser) evaluateOp(operator *Token) (*big.Rat, error) {
	var (
		result   *big.Rat
		lhs, rhs *big.Rat
		err      error
		lhsToken interface{}
	)

	if p.operands.Empty() {
		return nil, fmt.Errorf("Unexpected ‘%s’", operator)
	}

	if rhs, err = p.lookup(p.operands.Pop()); err != nil {
		return nil, err
	}

	// Unary operators have no left hand side
	if op := operators[operator.Type]; !op.unary {
		if p.operands.Empty() {
			return nil, fmt.Errorf("Unexpected ‘%s’", operator)
		}
		// Save the token in case of a assignment variable is used and we need
		// to save the result in a variable
		lhsToken = p.operands.Pop()

		// Don't lookup the left hand side if = is used so we can do initial
		// assignment
		if !operator.Is(Eq) {
			lhs, err = p.lookup(lhsToken)
			if err != nil {
				return nil, err
			}
		}
	}

	result, err = executeExpression(operator, lhs, rhs)
	if err != nil {
		return nil, err
	}

	if operator.IsAssignment() {
		// Save result in variable
		if val, ok := lhsToken.(*Token); !(ok && val.Is(Ident)) {
			return nil, ErrAssignToLiteral
		}
		p.Variables[lhsToken.(*Token).Value] = result
	}

	return result, nil
}

// Look up a literal. If it's an identifier, check the parser's variables map,
// otherwise convert the tokenized string to a rational number.
func (p *Parser) lookup(val interface{}) (*big.Rat, error) {
	// val can be a token or a rational, if it's a rational it has been already
	// evaluated and we don't need to do anything
	if v, ok := val.(*big.Rat); ok {
		return v, nil
	}

	var (
		ok  bool
		res = new(big.Rat)
	)

	bases := [...]int{
		Hex:    16,
		Octal:  8,
		Binary: 2,
	}

	tok := val.(*Token)
	switch tok.Type {
	case Decimal:
		res, ok = res.SetString(tok.Value)

		if !ok {
			return nil, fmt.Errorf("Error parsing ‘%s’: invalid %s", tok.Value, tok.Type)
		}
	case Hex, Binary, Octal:
		tmpInt := new(big.Int)
		// Remove prefix of literal and convert to int first
		tmpInt, ok = tmpInt.SetString(tok.Value[2:], bases[tok.Type])

		if !ok {
			return nil, fmt.Errorf("Error parsing ‘%s’: invalid %s", tok.Value, tok.Type)
		}

		res.SetInt(tmpInt)
	case Ident:
		res, err := p.GetVar(tok.Value)
		if err != nil {
			return nil, err
		}

		return res, nil
	default:
		return nil, fmt.Errorf("Invalid lookup type ‘%s’", tok)
	}

	return res, nil
}

func (p *Parser) reset() {
	p.Tokens = nil
	p.pos = 0

	p.operators = nil
	p.operands = nil
	p.arity = nil
}

func (p *Parser) peek() *Token {
	return p.Tokens[p.pos]
}

func (p *Parser) peekN(n int) *Token {
	return p.Tokens[p.pos-1+n]
}

func (p *Parser) eat() *Token {
	p.tok = p.peek()
	p.pos++
	return p.tok
}
