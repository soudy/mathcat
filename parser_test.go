// Copyright 2016 Steven Oud. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be found
// in the LICENSE file.

package evaler

import (
	"testing"
)

func TestParser(t *testing.T) {
	Eval("2 + 2 * 5")
}
