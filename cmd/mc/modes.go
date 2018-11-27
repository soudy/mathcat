// Copyright 2016 Steven Oud. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be found
// in the LICENSE file.

package main

type Mode int

const (
	Decimal Mode = iota
	Hex
	Binary
	Octal
)

var modes = map[string]Mode{
	"decimal": Decimal,
	"hex":     Hex,
	"binary":  Binary,
	"octal":   Octal,
}
