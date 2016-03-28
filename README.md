# eparser
eparser is an expression parser library for Go. It supports basic arithmetic,
bitwise operations and variable assignment.


## Usage
There are two functions you'll be mostly using: `Run` and `Exec`. `Run` will
evaluate and process an expression, but won't return a result (besides errors).
`Exec` however, _will_ return a result, but it can't mutate variables.

For example:
```go
p := eparser.New()
p.Run("a = 5") // a is 5
res, errs := p.Exec("a = 20") // res is 0, a is still 5
```
