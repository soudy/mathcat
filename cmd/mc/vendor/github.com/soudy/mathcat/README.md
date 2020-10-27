mathcat [![Build Status](https://travis-ci.org/soudy/mathcat.svg?branch=master)](https://travis-ci.org/soudy/mathcat) [![GoDoc](https://godoc.org/github.com/soudy/mathcat?status.svg)](https://godoc.org/github.com/soudy/mathcat)
===============
mathcat is an expression evaluating library and REPL in Go with basic arithmetic,
functions, variables and more.

## Features
mathcat doesn't just evaluate basic expressions, it has some tricks up its
sleeve. Here's a list with some of its features:

- Hex literals (0xDEADBEEF)
- Binary literals (0b1101001)
- Octal literals (0o126632)
- Scientific notation (24e3)
- Variables (with UTF-8 support)
- Functions ([list](#functions))
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
go get github.com/soudy/mathcat/cmd/mc
```

## REPL usage
The REPL can be used by simply launching `mc`:
```bash
mc> 8**8
16777216
mc> (8**8) - e # Look, a comment!
16777213.281718171541
```

Or it can read from `stdin` like so:

```bash
$ echo "3**pi * (6 - -7)" | mc
```

### Arguments

| Name      | Description                                                          | Default |
|-----------|----------------------------------------------------------------------|---------|
| precision | bits of decimal precision used in decimal float results              | 64      |
| mode      | type of literal used as result. can be decimal, hex, binary or octal | decimal |

## Library usage
There are three different ways to evaluate expressions, the first way is by
calling `Eval`, the second way is by creating a new instance and using `Run`,
and the final way is to use `Exec` in which you can pass a map with variables to
use in the expression.

### Eval
If you're not planning on declaring variables, you can use `Eval`. `Eval`
will evaluate an expression and return its result.

```go
res, err := mathcat.Eval("2 * pi * 5") // pi is a predefined variable
if err != nil {
    // handle errors
}
fmt.Printf("Result: %s\n", res.FloatString(6)) // Result: 31.415927
```

### Run
You can use `Run` for a more featureful approach. With this method you can
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
res, err := mathcat.Exec("a + b * b", map[string]*big.Rat{
    "a": big.NewRat(1, 1),
    "b": big.NewRat(3, 1),
}) // 10
```

Besides evaluating expressions, mathcat also offers some other handy functions.
### GetVar
You can get a defined variable at any time with `GetVar`.
```go
p := mathcat.New()
p.Run("酷 = -33")
if val, err := p.GetVar("酷"); !err {
    fmt.Printf("%f\n", val) // -33
}
```

### IsValidIdent
Check if a string qualifies as a valid identifier
```go
mathcat.IsValidIdent("a2") // true
mathcat.IsValidIdent("6a") // false
```

### RationalToInteger
Convert a `big.Rat` to a `big.Int`. Useful for printing in other bases or when
you're only working with integers.
```go
integer := mathcat.RationalToInteger(big.NewRat(42, 1))
fmt.Printf("%#x\n", integer) // prints 0x2a
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
| !=         | not equal             |
| >          | greater than          |
| >=         | greater than or equal |
| <          | less than             |
| <=         | less than or equal    |

All of these except `~` and relational operators also have an assignment
variant (`+=`, `-=`, `**=` etc.) that can be used to assign values to variables.

### Functions
mathcat has a big list of functions you can use. A function call is invoked like
in most programming languages, with an identifier followed by a left parentheses
like this: `max(5, 10)`.

| Function        |     Arguments | Description                                                                      |
| :-------------: | :-----------: | -------------------------------------------------------------------------------- |
| abs(n)          |             1 | returns the absolute value of given number                                       |
| sin(n)          |             1 | returns the sine of given number                                                 |
| cos(n)          |             1 | returns the cosine of given number                                               |
| tan(n)          |             1 | returns the tangent of given number                                              |
| asin(n)         |             1 | returns the arcsine of given number                                              |
| acos(n)         |             1 | returns the acosine of given number                                              |
| atan(n)         |             1 | returns the arctangent of given number                                           |
| ceil(n)         |             1 | returns the smallest integer greater than or equal to a given number             |
| floor(n)        |             1 | returns the largest integer less than or equal to a given number                 |
| ln(n)           |             1 | returns the natural logarithm of given number                                    |
| log(n)          |             1 | returns the the decimal logarithm of given number                                |
| logn(k, n)      |             2 | returns the the k logarithm of n                                                 |
| max(a, b)       |             2 | returns the larger of the two given numbers                                      |
| min(a, b)       |             2 | returns the smaller of the two given numbers                                     |
| sqrt(n)         |             1 | returns the square root of given number                                          |
| rand()          |             0 | returns a random float between 0.0 and 1.0                                       |
| fact(n)         |             1 | returns the factorial of  given number                                           |
| list()          |             0 | list all functions                                                               |

### Predefined variables
There are some handy predefined variables you can use (and change) throughout
your expressions:

- pi
- tau
- phi
- e
- true (set to 1)
- false (set to 0)

## Documentation
For a more technical description of mathcat, see [here](https://godoc.org/github.com/soudy/mathcat).

## License
This project is licensed under the MIT License. See the [LICENSE](https://github.com/soudy/mathcat/blob/master/LICENSE) file for the full license.
