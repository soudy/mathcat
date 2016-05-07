// Copyright 2016 Steven Oud. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be found
// in the LICENSE file.

package evaler

import "testing"

func TestFloatBitwise(t *testing.T) {
	badExpressions := []string{
		"2.4 | 2", "5.5 & 32", "7.7 ^ 2.1", "9 << 20.1", "7 >> 21.2",
	}

	for _, expr := range badExpressions {
		_, err := Eval(expr)
		if err == nil {
			t.Error("expected error on using bitwise operator on float")
		}
	}

	okExpressions := []string{
		"2.0 | 2", "5 & 32.0", "7 ^ 2", "9 << 20", "7.0 >> 21.0",
	}

	for _, expr := range okExpressions {
		_, err := Eval(expr)
		if err != nil {
			t.Error("unexpected error on using bitwise operator on int")
		}
	}
}

func TestParser(t *testing.T) {
	_, err := Eval("2 ** 10 / 5 * 2 - 6")

	if err != nil {
		t.Error(err)
	}
}
