mathcat [![Build Status](https://travis-ci.org/soudy/mathcat.svg?branch=master)](https://travis-ci.org/soudy/mathcat) [![GoDoc](https://godoc.org/github.com/soudy/mathcat?status.svg)](https://godoc.org/github.com/soudy/mathcat)
===============
mathcat is an expression parser library for Go. It supports basic arithmetic,
bitwise operations and variable assignment.

## Features
mathcat doesn't just evaluate basic expressions, it has some tricks up its
sleeve. Here's a list with some of its features:

- Hex literals (0xDEADBEEF)
- Binary literals (0b1101001)
- Scientific notation (24e3)
- Variables
- Bitwise operators
- Relational operators
- Some handy [predefined variables](#predefined-variables)
- Its own [REPL](#repl)

## Installation
### Library
```bash
go get github.com/soudy/mathcat
```

### REPL
```bash
go get github.com/soudy/mathcat/cmd/mathcat
```

## REPL usage
The REPL can be used by simply launching `mathcat`, or can read from stdin like
so:

```bash
echo "3**pi * (6 - -7)" | mathcat
```

## Library usage
There are three different ways to evaluate expressions, the first way is by
calling `Eval`, the second way is by creating a new instance and using `Run`,
and the final way is to use `Exec` in which you can pass a map with variables to
use in the expression in which you can pass a map with variables to use in the
expression.

### Eval
If you're not planning on declaring variables, you can use `Eval`. `Eval`
will evaluate an expression and return its result.

```go
res, err := mathcat.Eval("2 * pi * 5") // pi is a predefined variable
if err != nil {
    // handle errors
}
fmt.Printf("Result: %f\n", res) // Result: 31.41592653589793
```

### Run
You can use `Run` and for a more featureful approach. With this method you can
assign and use variables across the `Parser` instance.

```go
p := mathcat.New()
p.Run("a = 1")
p.Run("b = 3")
res, err := p.Run("a + b * b") // 10
```

### Exec
To pass external variables to an expression without using `Run`, you can use
`Exec` to pass a map of variables.

```go
res, err := mathcat.Exec("a + b * b", map[string]float64{
    "a": 1,
    "b": 3,
}) // 10
```

Besides evaluating expressions, mathcat also offers some other handy functions.
### GetVar
You can get a defined variable at any time with `GetVar`.
```go
p := mathcat.New()
p.Run("a = 1")
fmt.Printf("%f\n", p.GetVar("a")) // 1
```

### IsWholeNumber
Check if a `float64` is a whole number.
```go
if mathcat.IsWholeNumber(res) {
    fmt.Printf("%d\n", int64(res))
} else {
    fmt.Printf("%f\n", res)
}
```

### Supported operators

| Operator   | Description           |
|:----------:|:---------------------:|
| =          | assignment            |
| +          | addition              |
| -          | subtraction           |
| /          | division              |
| *          | multiply              |
| **         | power                 |
| %          | remainder             |
| &          | bitwise and           |
| \|         | bitwise or            |
| ^          | bitwise xor           |
| <<         | bitwise left shift    |
| >>         | bitwise right shift   |
| ~          | bitwise not           |
| ==         | equal                 |
| >          | greater than          |
| >=         | greater than or equal |
| <          | less than             |
| <=         | less than or equal    |

All of these except `~` and relational operators also have an assignment
variant (`+=`, `-=`, `**=` etc.) that can be used to assign values to variables.

### Predefined variables
There are some handy predefined variables you can use (and change) throughout
your expressions:

- pi
- tau
- phi
- e

## Documentation
For a more technical description of mathcat, it's available [here](https://godoc.org/github.com/soudy/mathcat).

## License
This project is licensed under the MIT License. See the [LICENSE](https://github.com/soudy/mathcat/blob/master/LICENSE) file for the full license.
