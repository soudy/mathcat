package eparser

import (
	"fmt"
	"math"
)

type Parser struct {
	tokens       []*Token
	pos          int
	Variables    map[string]float64
	currentToken *Token
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
func Parse(expr string) (float64, []error) {
	tokens, errs := Lex(expr)

	// If lexer errors occured don't parse
	if errs != nil {
		return nil, errs
	}

	p := &Parser{
		tokens:       tokens,
		pos:          0,
		Variables:    constants,
		currentToken: tokens[0],
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

func (p *Parser) parse() (float64, []error) {
	return 0, nil
}

func (p *Parser) peek() *Token {
	return nil
}

func (p *Parser) eat() *Token {
	return nil
}

func (p *Parser) expect(token tokenType) {
}
