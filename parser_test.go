// Copyright 2016 Steven Oud. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be found
// in the LICENSE file.

package evaler

import "testing"

func TestFloatBitwise(t *testing.T) {
	badExpressions := []string{
		"2.4 | 2", "5.5 & 32", "7.7 ^ 2.1", "9 << 20.1", "7 >> 21.2", "~5.3",
	}

	for _, expr := range badExpressions {
		_, err := Eval(expr)
		if err == nil {
			t.Error("expected error on using bitwise operator on float")
		}
	}

	okExpressions := []string{
		"2.0 | 2", "5 & 32.0", "7 ^ 2", "9 << 20", "7.0 >> 21.0", "~255",
	}

	for _, expr := range okExpressions {
		_, err := Eval(expr)
		if err != nil {
			t.Error("unexpected error on using bitwise operator on int")
		}
	}
}

func TestParse(t *testing.T) {
	okExpressions := map[string]float64{
		"()":                                            0,
		"1":                                             1,
		"(1)":                                           1,
		"12**12":                                        8916100448256,
		"~(~(1))":                                       1,
		"((((((((((((1))))))))))))":                     1,
		"(1 + (2 + (3 + (4 + (5 + (6 + (7)))))))":       28,
		"(((((((1) + 2) + 3) + 4) + 5) + 6) + 7)":       28,
		"((2 + 2 - 3) / (5 + 5 * 8 / 9)) - (9 + 2)":     -10.894117647058824,
		"((2 * 4 - 6 / 3) * (3 * 5 + 8 / 4)) - (2 + 3)": 97,
	}

	for expr, expected := range okExpressions {
		res, err := Eval(expr)
		if err != nil {
			t.Error("parser error occured on correct expression")
		}

		if expected != res {
			t.Errorf("wrong result in expression '%s' (expected %.1f, got %.1f)",
				expr, expected, res)
		}
	}

	badExpressions := []string{
		"2 / 0", "2 % 0", "+", "2 + 2 +", ")", "(2 + 2 * 8", "@#%@#*%&@#",
		"a + a", "~~2",
	}

	for _, expr := range badExpressions {
		_, err := Eval(expr)
		if err == nil {
			t.Error("no parser error occured on bad expression")
		}
	}
}

func TestWholeNumber(t *testing.T) {
	wholeNumbers := []float64{2.0, 2, -2.0, -2, 100.0}
	for _, num := range wholeNumbers {
		if !IsWholeNumber(num) {
			t.Error("whole number not recognized")
		}
	}

	floats := []float64{2.00001, 2.1, -2.09999, -2.00009, 100.9}
	for _, num := range floats {
		if IsWholeNumber(num) {
			t.Error("float recognized as whole number")
		}
	}
}
