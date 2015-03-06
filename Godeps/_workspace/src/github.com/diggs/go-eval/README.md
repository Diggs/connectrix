## What is go-eval?

go-eval evaluates boolean or basic arithmetic expressions from strings. It's analogous to a basic form of the eval function often found in dynamic languages.

Examples:
```go
res, err := EvalBool("1 > 2")
log.Print(res)
> false
```

```go
res, err := EvalBool(`(1 + 3) >= 4 && ("FOO" == "BAR" || "FOO" == "FOO")`)
log.Print(res)
> true
```

```go
res, err := EvalArithmetic("1 + 2")
log.Print(res)
> 3
```

```go
res, err := EvalArithmetic("2 - -1")
log.Print(res)
> 3
```

## How does it work?

go-eval leverges the awesome AST and Parser packages that ship in the Go standard library. It constructs an abstract syntax tree from the specified expression and walks the tree to determine the result.

## Tests

```go
go get github.com/tools/godep
godep restore
go test

PASS
ok  	github.com/diggs/go-eval	0.005s
```

## Limitations

Floats are not supported yet.
Only boolean and arithmetic expressions are supported, this doesn't support conditionals etc.