// Copyright 2016 Steven Oud. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be found
// in the LICENSE file.

package mathcat

import (
	"math"
)

type function struct {
	name      string
	nargs     int
	operation func(args []float64) float64
}

var functions = map[string]*function{
	"abs": &function{
		name:  "abs",
		nargs: 1,
		operation: func(args []float64) float64 {
			return math.Abs(args[0])
		},
	},
	"ceil": &function{
		name:  "ceil",
		nargs: 1,
		operation: func(args []float64) float64 {
			return math.Ceil(args[0])
		},
	},
	"floor": &function{
		name:  "floor",
		nargs: 1,
		operation: func(args []float64) float64 {
			return math.Floor(args[0])
		},
	},
	"sin": &function{
		name:  "sin",
		nargs: 1,
		operation: func(args []float64) float64 {
			return math.Sin(args[0])
		},
	},
	"cos": &function{
		name:  "cos",
		nargs: 1,
		operation: func(args []float64) float64 {
			return math.Cos(args[0])
		},
	},
	"tan": &function{
		name:  "tan",
		nargs: 1,
		operation: func(args []float64) float64 {
			return math.Tan(args[0])
		},
	},
	"asin": &function{
		name:  "asin",
		nargs: 1,
		operation: func(args []float64) float64 {
			return math.Asin(args[0])
		},
	},
	"acos": &function{
		name:  "acos",
		nargs: 1,
		operation: func(args []float64) float64 {
			return math.Acos(args[0])
		},
	},
	"atan": &function{
		name:  "atan",
		nargs: 1,
		operation: func(args []float64) float64 {
			return math.Atan(args[0])
		},
	},
	"log": &function{
		name:  "log",
		nargs: 1,
		operation: func(args []float64) float64 {
			return math.Log(args[0])
		},
	},
	"max": &function{
		name:  "max",
		nargs: 2,
		operation: func(args []float64) float64 {
			return math.Max(args[0], args[1])
		},
	},
	"min": &function{
		name:  "min",
		nargs: 2,
		operation: func(args []float64) float64 {
			return math.Min(args[0], args[1])
		},
	},
	"sqrt": &function{
		name:  "sqrt",
		nargs: 1,
		operation: func(args []float64) float64 {
			return math.Sqrt(args[0])
		},
	},
}
