// Copyright 2016 Steven Oud. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be found
// in the LICENSE file.

package evaler

import (
	"fmt"
	"testing"
)

func TestParser(t *testing.T) {
	res, err := Eval("2 ** 10 / 5 * 2 - 6")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print(res)
}
