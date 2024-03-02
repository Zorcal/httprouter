# httprouter

A thin wrapper around [http.ServeMux](https://pkg.go.dev/net/http#ServeMux)
with support for middleware and router groups. Defines a custom handler type
that returns an error, which allows for bubbling up all errors to middlewares.

```go
type Handler func(w http.ResponseWriter, r *http.Request) error
```

All of the documentation can be found on the [go.dev](https://pkg.go.dev/github.com/zorcal/httprouter?tab=doc) website.
