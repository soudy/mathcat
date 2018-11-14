// Copyright 2016 Steven Oud. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be found
// in the LICENSE file.

package main

type Mode int

const (
	DECIMAL Mode = iota
	HEX
	BINARY
	OCTAL
)

var modes = map[string]Mode{
	"decimal": DECIMAL,
	"hex":     HEX,
	"binary":  BINARY,
	"octal":   OCTAL,
}
