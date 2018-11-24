// Copyright 2016 Steven Oud. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be found
// in the LICENSE file.

package mathcat

import (
	"math/big"
	"testing"
)

var RatZero = new(big.Rat)

func TestFloatBitwise(t *testing.T) {
	badExpressions := []string{
		"2.4 | -2", "5.5 & 32", "7.7 ^ 2.1", "9 << 20.1", "7 >> 21.2", "~5.3",
		"2 + 2 = 4",
	}

	for _, expr := range badExpressions {
		_, err := Eval(expr)
		if err == nil {
			t.Error("expected error on using bitwise operator on float")
		}
	}

	okExpressions := []string{
		"2.0 | 2", "5 & 32.0", "7 ^ 2", "9 << 20", "7.0 >> -21.0", "~255",
	}

	for _, expr := range okExpressions {
		_, err := Eval(expr)
		if err != nil {
			t.Errorf("unexpected error on using bitwise operator on int: %s", err)
		}
	}
}

func TestNumberLiterals(t *testing.T) {
	invalidNumbers := []string{
		"0x", "0X", "0x12p345", "0b", "0B", "0b2", "0x22.3", "0b10.1", "0b1e2",
		"0o8", "0oea",
	}

	for _, n := range invalidNumbers {
		_, err := Eval(n)
		if err == nil {
			t.Error("no error on invalid number literal")
		}
	}

	validNumbers := map[string]*big.Rat{
		"0xa":      big.NewRat(10, 1),
		"0Xaaabe":  big.NewRat(699070, 1),
		"0x12345":  big.NewRat(74565, 1),
		"0xe":      big.NewRat(14, 1),
		".200001":  big.NewRat(200001, 1000000),
		".2e5":     big.NewRat(20000, 1),
		"81e2":     big.NewRat(8100, 1),
		"32":       big.NewRat(32, 1),
		"100.0":    big.NewRat(100, 1),
		"0x0":      RatZero,
		"0":        RatZero,
		"0b110011": big.NewRat(51, 1),
		"0b1":      big.NewRat(1, 1),
		"0o666":    big.NewRat(438, 1),
		"0O6120":   big.NewRat(3152, 1),
		"0o0":      RatZero,
	}

	for n, expected := range validNumbers {
		res, err := Eval(n)
		if err != nil {
			t.Errorf("error on valid number literal: %s", err)
		}

		if res.Cmp(expected) != 0 {
			t.Errorf("invalid literal evaluation on %s, expected '%s', got '%s'", n, expected, res)
		}
	}

}

func TestEval(t *testing.T) {
	okExpressions := map[string]*big.Rat{
		"()":                                            RatZero,
		"-1":                                            big.NewRat(-1, 1),
		"(1)":                                           big.NewRat(1, 1),
		"12**12":                                        big.NewRat(8916100448256, 1),
		"~(~(1))":                                       big.NewRat(1, 1),
		"1000 > 10":                                     big.NewRat(1, 1),
		"1000 < 10":                                     RatZero,
		"55.0 == 55":                                    big.NewRat(1, 1),
		"55 <= 55":                                      big.NewRat(1, 1),
		"55 >= 55":                                      big.NewRat(1, 1),
		"-0 > 0":                                        RatZero,
		"0 < -0":                                        RatZero,
		"2 != 2":                                        RatZero,
		"5 % 25":                                        big.NewRat(5, 1),
		"((((((((((((1))))))))))))":                     big.NewRat(1, 1),
		"(1 + (2 + (3 + (4 + (5 + (6 + (7)))))))":       big.NewRat(28, 1),
		"(((((((1) + 2) + 3) + 4) + 5) + 6) + 7)":       big.NewRat(28, 1),
		"((2 + 2 - 3) / (5 + 5 * 8 / 9)) - (9 + 2)":     big.NewRat(-926, 85),
		"((2 * 4 - 6 / 3) * (3 * 5 + 8 / 4)) - (2 + 3)": big.NewRat(97, 1),
		"0xdeadbeef & 0xff000000":                       big.NewRat(3724541952, 1),
		"325-2*5+2":                                     big.NewRat(317, 1),
		"3**pi * (6 - -7)":                              big.NewRat(57713016890376237, 140737488355328),
	}

	for expr, expected := range okExpressions {
		res, err := Eval(expr)
		if err != nil {
			t.Errorf("parser error occured on correct expression '%s': %s", expr, err)
		}

		if expected.Cmp(res) != 0 {
			t.Errorf("wrong result in expression '%s' (expected %s, got %s)",
				expr, expected, res)
		}
	}

	badExpressions := []string{
		"2 / 0", "2 % 0", "+", "2 + 2 +", ")", "(2 + 2 * 8", "@#%@#*%&@#",
		"a + a", "~~2", "2 == ()", "5 < -", "2 * (9 ** 2))", "5 ~ 3",
	}

	for _, expr := range badExpressions {
		_, err := Eval(expr)
		if err == nil {
			t.Error("no parser error occured on bad expression")
		}
	}
}

func TestExec(t *testing.T) {
	type execTest struct {
		expr     string
		vars     map[string]*big.Rat
		expected *big.Rat
	}

	okExpressions := []execTest{
		{"酷 + b * b", map[string]*big.Rat{
			"酷": big.NewRat(1, 1),
			"b": big.NewRat(3, 1)},
			big.NewRat(10, 1)},
		{"a + b_ * pi", map[string]*big.Rat{
			"a":  big.NewRat(1, 1),
			"b_": big.NewRat(3, 1),
			"pi": big.NewRat(3, 1)},
			big.NewRat(10, 1)},
		{"a2 + b5 * pi3", map[string]*big.Rat{
			"a2":  big.NewRat(1, 1),
			"b5":  big.NewRat(3, 1),
			"pi3": big.NewRat(3, 1)},
			big.NewRat(10, 1)},
		{"Å ** Å", map[string]*big.Rat{"Å": big.NewRat(1, 1)}, big.NewRat(1, 1)},
	}

	for _, test := range okExpressions {
		res, err := Exec(test.expr, test.vars)
		if err != nil {
			t.Errorf("error on correct Exec: %s", err)
		}

		if res.Cmp(test.expected) != 0 {
			t.Error("wrong result in Exec")
		}
	}

	badExpressions := []execTest{
		{"", map[string]*big.Rat{"-1": nil}, nil},
		{"", map[string]*big.Rat{"55": nil}, nil},
		{"", map[string]*big.Rat{"55a": nil}, nil},
		{"", map[string]*big.Rat{".": nil}, nil},
		{"", map[string]*big.Rat{")": nil}, nil},
		{"", map[string]*big.Rat{"(": nil}, nil},
		{"", map[string]*big.Rat{"@": nil}, nil},
	}

	for _, test := range badExpressions {
		_, err := Exec(test.expr, test.vars)
		if err == nil {
			t.Error("no error on bad Exec")
		}
	}
}

func TestGetVar(t *testing.T) {
	p := New()
	p.Run("酷 = -33")

	if _, err := p.GetVar("酷"); err != nil {
		t.Error("GetVar failed: " + err.Error())
	}
}
