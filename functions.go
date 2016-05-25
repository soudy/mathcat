// Copyright 2016 Steven Oud. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be found
// in the LICENSE file.

package mathcat

import (
	"math"
	"math/rand"
)

type function struct {
	nargs int
	fn    func(args []float64) float64
}

var functions = map[string]*function{
	"abs": &function{
		nargs: 1,
		fn: func(args []float64) float64 {
			return math.Abs(args[0])
		},
	},
	"ceil": &function{
		nargs: 1,
		fn: func(args []float64) float64 {
			return math.Ceil(args[0])
		},
	},
	"floor": &function{
		nargs: 1,
		fn: func(args []float64) float64 {
			return math.Floor(args[0])
		},
	},
	"sin": &function{
		nargs: 1,
		fn: func(args []float64) float64 {
			return math.Sin(args[0])
		},
	},
	"cos": &function{
		nargs: 1,
		fn: func(args []float64) float64 {
			return math.Cos(args[0])
		},
	},
	"tan": &function{
		nargs: 1,
		fn: func(args []float64) float64 {
			return math.Tan(args[0])
		},
	},
	"asin": &function{
		nargs: 1,
		fn: func(args []float64) float64 {
			return math.Asin(args[0])
		},
	},
	"acos": &function{
		nargs: 1,
		fn: func(args []float64) float64 {
			return math.Acos(args[0])
		},
	},
	"atan": &function{
		nargs: 1,
		fn: func(args []float64) float64 {
			return math.Atan(args[0])
		},
	},
	"log": &function{
		nargs: 1,
		fn: func(args []float64) float64 {
			return math.Log(args[0])
		},
	},
	"max": &function{
		nargs: 2,
		fn: func(args []float64) float64 {
			return math.Max(args[0], args[1])
		},
	},
	"min": &function{
		nargs: 2,
		fn: func(args []float64) float64 {
			return math.Min(args[0], args[1])
		},
	},
	"sqrt": &function{
		nargs: 1,
		fn: func(args []float64) float64 {
			return math.Sqrt(args[0])
		},
	},
	"rand": &function{
		nargs: 0,
		fn: func(_ []float64) float64 {
			return rand.Float64()
		},
	},
}
