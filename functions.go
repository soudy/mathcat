// Copyright 2016 Steven Oud. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be found
// in the LICENSE file.

package mathcat

import (
	"fmt"
	"math"
	"math/rand"
)

type function struct {
	nargs int
	fn    func(args []float64) float64
}

type functions map[string]*function

var (
	funcs        functions = make(map[string]*function)
	allFunctions []string
)

func (f functions) register(name string, function *function) {
	allFunctions = append(allFunctions, name)
	f[name] = function
}

func init() {
	funcs.register("abs", &function{
		nargs: 1,
		fn: func(args []float64) float64 {
			return math.Abs(args[0])
		},
	})
	funcs.register("ceil", &function{
		nargs: 1,
		fn: func(args []float64) float64 {
			return math.Ceil(args[0])
		},
	})
	funcs.register("floor", &function{
		nargs: 1,
		fn: func(args []float64) float64 {
			return math.Floor(args[0])
		},
	})
	funcs.register("sin", &function{
		nargs: 1,
		fn: func(args []float64) float64 {
			return math.Sin(args[0])
		},
	})
	funcs.register("cos", &function{
		nargs: 1,
		fn: func(args []float64) float64 {
			return math.Cos(args[0])
		},
	})
	funcs.register("tan", &function{
		nargs: 1,
		fn: func(args []float64) float64 {
			return math.Tan(args[0])
		},
	})
	funcs.register("tan", &function{
		nargs: 1,
		fn: func(args []float64) float64 {
			return math.Asin(args[0])
		},
	})
	funcs.register("acos", &function{
		nargs: 1,
		fn: func(args []float64) float64 {
			return math.Acos(args[0])
		},
	})
	funcs.register("atan", &function{
		nargs: 1,
		fn: func(args []float64) float64 {
			return math.Atan(args[0])
		},
	})
	funcs.register("log", &function{
		nargs: 1,
		fn: func(args []float64) float64 {
			return math.Log(args[0])
		},
	})
	funcs.register("max", &function{
		nargs: 2,
		fn: func(args []float64) float64 {
			return math.Max(args[0], args[1])
		},
	})
	funcs.register("max", &function{
		nargs: 2,
		fn: func(args []float64) float64 {
			return math.Min(args[0], args[1])
		},
	})
	funcs.register("sqrt", &function{
		nargs: 1,
		fn: func(args []float64) float64 {
			return math.Sqrt(args[0])
		},
	})
	funcs.register("rand", &function{
		nargs: 0,
		fn: func(_ []float64) float64 {
			return rand.Float64()
		},
	})
	funcs.register("fact", &function{
		nargs: 1,
		fn: func(args []float64) float64 {
			return factorial(args[0])
		},
	})
	funcs.register("list", &function{
		nargs: 0,
		fn: func(_ []float64) float64 {
			for _, name := range allFunctions {
				fmt.Printf(name + " ")
			}
			fmt.Println()
			return 0
		},
	})
}

func factorial(n float64) float64 {
	if n == 0 {
		return 1
	}

	return n * factorial(n-1)
}
