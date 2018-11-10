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
	arity int
	fn    func(args []float64) float64
}

type functions map[string]function

// FunctionNames holds all the function names that are available for use
var FunctionNames []string

var funcs = make(functions)

func (f functions) register(name string, function function) {
	FunctionNames = append(FunctionNames, name)
	f[name] = function
}

func init() {
	funcs.register("abs", function{
		arity: 1,
		fn: func(args []float64) float64 {
			return math.Abs(args[0])
		},
	})
	funcs.register("ceil", function{
		arity: 1,
		fn: func(args []float64) float64 {
			return math.Ceil(args[0])
		},
	})
	funcs.register("floor", function{
		arity: 1,
		fn: func(args []float64) float64 {
			return math.Floor(args[0])
		},
	})
	funcs.register("sin", function{
		arity: 1,
		fn: func(args []float64) float64 {
			return math.Sin(args[0])
		},
	})
	funcs.register("cos", function{
		arity: 1,
		fn: func(args []float64) float64 {
			return math.Cos(args[0])
		},
	})
	funcs.register("tan", function{
		arity: 1,
		fn: func(args []float64) float64 {
			return math.Tan(args[0])
		},
	})
	funcs.register("asin", function{
		arity: 1,
		fn: func(args []float64) float64 {
			return math.Asin(args[0])
		},
	})
	funcs.register("acos", function{
		arity: 1,
		fn: func(args []float64) float64 {
			return math.Acos(args[0])
		},
	})
	funcs.register("atan", function{
		arity: 1,
		fn: func(args []float64) float64 {
			return math.Atan(args[0])
		},
	})
	funcs.register("ln", function{
		arity: 1,
		fn: func(args []float64) float64 {
			return math.Log(args[0])
		},
	})
	funcs.register("log", function{
		arity: 1,
		fn: func(args []float64) float64 {
			return math.Log10(args[0])
		},
	})
	funcs.register("logn", function{
		arity: 2,
		fn: func(args []float64) float64 {
			base := args[0]
			arg := args[1]
			return math.Log10(arg) / math.Log10(base)
		},
	})
	funcs.register("max", function{
		arity: 2,
		fn: func(args []float64) float64 {
			return math.Max(args[0], args[1])
		},
	})
	funcs.register("min", function{
		arity: 2,
		fn: func(args []float64) float64 {
			return math.Min(args[0], args[1])
		},
	})
	funcs.register("sqrt", function{
		arity: 1,
		fn: func(args []float64) float64 {
			return math.Sqrt(args[0])
		},
	})
	funcs.register("rand", function{
		arity: 0,
		fn: func(_ []float64) float64 {
			return rand.Float64()
		},
	})
	funcs.register("fact", function{
		arity: 1,
		fn: func(args []float64) float64 {
			return float64(Factorial(int64(args[0])))
		},
	})
	funcs.register("gcd", function{
		arity: 2,
		fn: func(args []float64) float64 {
			return Gcd(args[0], args[1])
		},
	})
	funcs.register("list", function{
		arity: 0,
		fn: func(_ []float64) float64 {
			for _, name := range FunctionNames {
				fmt.Print(name + " ")
			}
			fmt.Println()
			return 0
		},
	})
}

// Factorial calculates the factorial of number n
func Factorial(n int64) int64 {
	if n <= 1 {
		return 1
	}

	return n * Factorial(n-1)
}

// Gcd calculates the greatest common divisor of the numbers x and y
func Gcd(x, y float64) float64 {
	for y != 0 {
		x, y = y, math.Mod(x, y)
	}

	return x
}
