// Copyright 2016 Steven Oud. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be found
// in the LICENSE file.

package mathcat

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

var operators = map[TokenType]*operator{
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

	// Relational operators
	EQ_EQ:   {1, ASSOC_RIGHT, false}, // ==
	BANG_EQ: {1, ASSOC_RIGHT, false}, // !=
	GT:      {1, ASSOC_RIGHT, false}, // >
	GT_EQ:   {1, ASSOC_RIGHT, false}, // >=
	LT:      {1, ASSOC_RIGHT, false}, // <
	LT_EQ:   {1, ASSOC_RIGHT, false}, // <=

	// Bitwise operators
	OR:  {2, ASSOC_RIGHT, false}, // |
	XOR: {3, ASSOC_RIGHT, false}, // ^
	AND: {4, ASSOC_RIGHT, false}, // &
	LSH: {5, ASSOC_RIGHT, false}, // <<
	RSH: {5, ASSOC_RIGHT, false}, // >>
	NOT: {9, ASSOC_LEFT, true},   // ~

	// Mathematical operators
	ADD:       {6, ASSOC_LEFT, false}, // +
	SUB:       {6, ASSOC_LEFT, false}, // -
	MUL:       {7, ASSOC_LEFT, false}, // *
	DIV:       {7, ASSOC_LEFT, false}, // /
	POW:       {8, ASSOC_LEFT, false}, // **
	REM:       {7, ASSOC_LEFT, false}, // %
	UNARY_MIN: {10, ASSOC_LEFT, true}, // -
}

// Determine if operator 1 has higher precedence than operator 2
func (o1 operator) hasHigherPrecThan(o2 *operator) bool {
	return (o2.assoc == ASSOC_LEFT && o2.prec <= o1.prec) ||
		(o2.assoc == ASSOC_RIGHT && o2.prec < o1.prec)
}
