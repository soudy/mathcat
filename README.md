# eparser
eparser is an expression parser library for Go. It supports basic arithmetic,
bitwise operations and variable assignment.

## Installation
```bash
go get github.com/soudy/eparser
```

And to update:
```bash
go get -u github.com/soudy/eparser
```

## Usage
There are two different ways to evaluate expressions, the first way is by
calling `Parse`, and the second way is by creating a new instance and using
`Run` and `Exec`.

### Parse
If you're not planning on declaring variables, you can use `Parse`. `Parse`
will evaluate an expression and return its result.
```go
res, errs := eparser.Parse("2*pi*5") // pi is a predefined variable
if errs != nil {
    // handle errors
}
fmt.Printf("Result: %f\n", res) // Result: 31.41592653589793
```

### Run and Exec
Otherwise you can use `Run` and `Exec` for a more featureful approach. With
these you can assign values to variables and use them throughout the same
instance.

Example:
```go
p := eparser.New()
p.Run("a = 1")
p.Run("b = 3")
res, errs := p.Exec("a + b * b") // 10
```

Note that `Run` doesn't return the result of the expression. It does however
return any error(s) found.
