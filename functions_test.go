// Copyright 2016 Steven Oud. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be found
// in the LICENSE file.

package mathcat

import (
	"math/big"
	"testing"
)

func TestFunctions(t *testing.T) {
	badCalls := []string{
		"a()", "a(1, 2, 3)", "2 + 6 * (a(1, 2))", "abs(1, 2)", "abs()",
		"max(1)", "min(2)",
	}

	for _, expr := range badCalls {
		_, err := Eval(expr)
		if err == nil {
			t.Error("expected error on bad function call")
		}
	}

	okCalls := []string{
		"abs(-300)", "max(8, 8)", "8 * cos(pi) - 6", "tan(8 * 8 * (7**7))",
		"tan(cos(8) / sin(3))",
	}

	for _, expr := range okCalls {
		_, err := Eval(expr)
		if err != nil {
			t.Errorf("unexpected error on ok function call: %s", err)
		}
	}
}

func TestFunctionsResult(t *testing.T) {
	calls := map[string]*big.Rat{
		"abs(-700)":                         big.NewRat(700, 1),
		"ceil(813.23)":                      big.NewRat(814, 1),
		"ceil(ceil(10 ** 16 + 0.1))":        big.NewRat(10000000000000001, 1),
		"floor(813.23)":                     big.NewRat(813, 1),
		"floor(-50.23)":                     big.NewRat(-51, 1),
		"floor(-50)":                        big.NewRat(-50, 1),
		"sin(74)":                           big.NewRat(-8873408663100473, 9007199254740992),
		"cos(74)":                           big.NewRat(6186769253457135, 36028797018963968),
		"tan(74)":                           big.NewRat(-6459313142528259, 1125899906842624),
		"asin(-1)":                          big.NewRat(-884279719003555, 562949953421312),
		"acos(-1)":                          big.NewRat(884279719003555, 281474976710656),
		"atan(-1)":                          big.NewRat(-884279719003555, 1125899906842624),
		"ln(3*100)":                         big.NewRat(802736019608251, 140737488355328),
		"log(50)":                           big.NewRat(59777192800323, 35184372088832),
		"logn(2, 50)":                       big.NewRat(6354417158300529, 1125899906842624),
		"max(5, 8)":                         big.NewRat(8, 1),
		"min(5, 8)":                         big.NewRat(5, 1),
		"sqrt(144)":                         big.NewRat(12, 1),
		"tan(144) + tan(-3) + sin(5)":       big.NewRat(-49720712606960177, 36028797018963968),
		"fact(6) * fact(7) == fact(10)":     big.NewRat(1, 1),
		"fact(6.5) * fact(7.3) == fact(10)": big.NewRat(1, 1),
	}

	for expr, expected := range calls {
		res, err := Eval(expr)
		if err != nil {
			t.Errorf("unexpected error on ok function call: %s", err)
		}

		if res.Cmp(expected) != 0 {
			t.Errorf("wrong result in function call '%s' (expected %s, got %s)",
				expr, expected, res)
		}
	}
}
