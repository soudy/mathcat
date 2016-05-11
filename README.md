# evaler
[![Build Status](https://travis-ci.org/soudy/evaler.svg?branch=master)](https://travis-ci.org/soudy/evaler)
evaler is an expression parser library for Go. It supports basic arithmetic,
bitwise operations and variable assignment.

## Installation
```bash
go get github.com/soudy/evaler
```

And to update:
```bash
go get -u github.com/soudy/evaler
```

## Usage
There are three different ways to evaluate expressions, the first way is by
calling `Eval`, the second way is by creating a new instance and using `Run`,
and the final way is to use `Exec` in which you can pass a map with variables to
use in the expressio in which you can pass a map with variables to use in the
expression.

### Eval
If you're not planning on declaring variables, you can use `Eval`. `Eval`
will evaluate an expression and return its result.

```go
res, err := evaler.Eval("2 * pi * 5") // pi is a predefined variable
if err != nil {
    // handle errors
}
fmt.Printf("Result: %f\n", res) // Result: 31.41592653589793
```

### Run
You can use `Run` and for a more featureful approach. With this method you can
assign and use variables accross the `Parser` instance.

```go
p := evaler.New()
p.Run("a = 1")
p.Run("b = 3")
res, err := p.Run("a + b * b") // 10
```

### Exec
To pass external variables to an expression without using `Run`, you can use
`Exec` to pass a map of variables.

```go
res, err := evaler.Exec("a + b * b", map[string]float64{
    "a": 1,
    "b": 3,
}) // 10
```
