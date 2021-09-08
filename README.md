Package gerrors

`go get github.com/archomai/gerrors`

```go
if err != nil {
	return gerrors.New("New Error")
}
```

```go
if err != nil {
	gerrors.Trace(err)
}
```
