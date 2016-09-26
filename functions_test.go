// Copyright 2016 Steven Oud. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be found
// in the LICENSE file.

package mathcat

import (
	"math"
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
	calls := map[string]float64{
		"abs(-700)":                         math.Abs(-700),
		"ceil(813.23)":                      math.Ceil(813.23),
		"floor(813.23)":                     math.Floor(813.23),
		"sin(74)":                           math.Sin(74),
		"cos(74)":                           math.Cos(74),
		"tan(74)":                           math.Tan(74),
		"asin(-1)":                          math.Asin(-1),
		"acos(-1)":                          math.Acos(-1),
		"atan(-1)":                          math.Atan(-1),
		"log(3*100)":                        math.Log(3 * 100),
		"max(5, 8)":                         math.Max(5, 8),
		"min(5, 8)":                         math.Min(5, 8),
		"sqrt(144)":                         math.Sqrt(144),
		"tan(144) + tan(-3) + sin(5)":       -1.380026998425437,
		"fact(6) * fact(7) == fact(10)":     1,
		"fact(6.5) * fact(7.3) == fact(10)": 1,
	}

	for expr, expected := range calls {
		res, err := Eval(expr)
		if err != nil {
			t.Errorf("unexpected error on ok function call: %s", err)
		}

		if res != expected {
			t.Errorf("wrong result in function call '%s' (expected %f, got %f)",
				expr, expected, res)
		}
	}
}
