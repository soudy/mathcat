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

type association int

type operator struct {
	prec  int
	assoc association
	unary bool
}

const (
	AssocLeft association = iota
	AssocRight
)

var ErrDivisionByZero = errors.New("Division by zero")

var operators = map[TokenType]operator{
	// Assignment operators
	Eq:    {0, AssocRight, false}, // =
	AddEq: {0, AssocRight, false}, // +=
	SubEq: {0, AssocRight, false}, // -=
	DivEq: {0, AssocRight, false}, // /=
	MulEq: {0, AssocRight, false}, // *=
	PowEq: {0, AssocRight, false}, // **=
	RemEq: {0, AssocRight, false}, // %=
	AndEq: {0, AssocRight, false}, // &=
	OrEq:  {0, AssocRight, false}, // |=
	XorEq: {0, AssocRight, false}, // ^=
	LshEq: {0, AssocRight, false}, // <<=
	RshEq: {0, AssocRight, false}, // >>=

	// Relational operators
	EqEq:  {1, AssocRight, false}, // ==
	NotEq: {1, AssocRight, false}, // !=
	Gt:    {1, AssocRight, false}, // >
	GtEq:  {1, AssocRight, false}, // >=
	Lt:    {1, AssocRight, false}, // <
	LtEq:  {1, AssocRight, false}, // <=

	// Bitwise operators
	Or:  {2, AssocRight, false}, // |
	Xor: {3, AssocRight, false}, // ^
	And: {4, AssocRight, false}, // &
	Lsh: {5, AssocRight, false}, // <<
	Rsh: {5, AssocRight, false}, // >>
	Not: {9, AssocLeft, true},   // ~

	// Mathematical operators
	Add:      {6, AssocLeft, false}, // +
	Sub:      {6, AssocLeft, false}, // -
	Mul:      {7, AssocLeft, false}, // *
	Div:      {7, AssocLeft, false}, // /
	Pow:      {8, AssocLeft, false}, // **
	Rem:      {7, AssocLeft, false}, // %
	UnaryMin: {10, AssocLeft, true}, // -
}

// Determine if operator 1 has higher precedence than operator 2
func (o1 operator) hasHigherPrecThan(o2 operator) bool {
	return (o2.assoc == AssocLeft && o2.prec <= o1.prec) ||
		(o2.assoc == AssocRight && o2.prec < o1.prec)
}

// Execute a binary or unary expression
func executeExpression(operator *Token, lhs, rhs *big.Rat) (*big.Rat, error) {
	result := new(big.Rat)

	// Both lhs and rhs have to be integers for bitwise operations
	if operator.IsBitwise() {
		if (lhs == nil && !rhs.IsInt()) || (lhs != nil && (!rhs.IsInt() || !lhs.IsInt())) {
			return nil, fmt.Errorf("Expecting integers for ‘%s’", operator)
		}
	}

	switch operator.Type {
	case Add, AddEq:
		result.Add(lhs, rhs)
	case Sub, SubEq:
		result.Sub(lhs, rhs)
	case UnaryMin:
		result.Neg(rhs)
	case Div, DivEq:
		if rhs.Sign() == 0 {
			return nil, ErrDivisionByZero
		}
		result.Quo(lhs, rhs)
	case Mul, MulEq:
		result.Mul(lhs, rhs)
	case Pow, PowEq:
		if lhs.IsInt() && rhs.IsInt() {
			intResult := new(big.Int)
			intResult.Set(lhs.Num())
			intResult.Exp(intResult, rhs.Num(), nil)
			result.SetInt(intResult)
		} else {
			lhsFloat, _ := lhs.Float64()
			rhsFloat, _ := rhs.Float64()
			result.SetFloat64(math.Pow(lhsFloat, rhsFloat))
		}
	case Rem, RemEq:
		if rhs.Sign() == 0 {
			return nil, ErrDivisionByZero
		}
		result.Set(Mod(lhs, rhs))
	case And, AndEq:
		result.SetInt(new(big.Int).And(lhs.Num(), rhs.Num()))
	case Or, OrEq:
		result.SetInt(new(big.Int).Or(lhs.Num(), rhs.Num()))
	case Xor, XorEq:
		result.SetInt(new(big.Int).Xor(lhs.Num(), rhs.Num()))
	case Lsh, LshEq:
		shift := uint(rhs.Num().Uint64())
		result.SetInt(new(big.Int).Lsh(lhs.Num(), shift))
	case Rsh, RshEq:
		shift := uint(rhs.Num().Uint64())
		result.SetInt(new(big.Int).Rsh(lhs.Num(), shift))
	case Not:
		result.SetInt(new(big.Int).Not(rhs.Num()))
	case Eq:
		result = rhs
	case EqEq:
		result = boolToRat(lhs.Cmp(rhs) == 0)
	case NotEq:
		result = boolToRat(lhs.Cmp(rhs) != 0)
	case Gt:
		result = boolToRat(lhs.Cmp(rhs) == 1)
	case GtEq:
		result = boolToRat(lhs.Cmp(rhs) == 1 || lhs.Cmp(rhs) == 0)
	case Lt:
		result = boolToRat(lhs.Cmp(rhs) == -1)
	case LtEq:
		result = boolToRat(lhs.Cmp(rhs) == -1 || lhs.Cmp(rhs) == 0)
	default:
		return nil, fmt.Errorf("Invalid operator ‘%s’", operator)
	}

	return result, nil
}

func boolToRat(b bool) *big.Rat {
	if b {
		return RatTrue
	}
	return RatFalse
}
