// Copyright 2016 Steven Oud. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be found
// in the LICENSE file.

package mathcat

import "testing"

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
