# Zeus

Zeus is a super-duper, simple and fast HTTP router for Go, nothing more, nothing less.

#### Install

    go get github.com/daryl/zeus

#### Usage

```go
package main

import (
    "fmt"
    "github.com/daryl/zeus"
    "net/http"
)

func main() {
    mux := zeus.New()
    // Supports named parameters.
    mux.GET("/users/:id", showUser)
    // Supports wildcards anywhere.
    mux.GET("/foo/*", catchFoo)
    // Custom 404 handler.
    mux.NotFound = notFound
    // Listen and serve.
    mux.Listen(":4545")
}

func showUser(w http.ResponseWriter, r *http.Request) {
    var id string

    // Get a map of all
    // route variables.
    vm := zues.Vars(r)

    id = vm["id"]

    // Or just one.
    id = zeus.Var(r, "id")

    fmt.Fprintf(w, "User ID: %s", id)
}

func catchFoo(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Gotta catch 'em all"))
}

func notFound(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Nothing to see here"))
}
```

#### Documentation

For further documentation, check out [GoDoc](http://godoc.org/github.com/daryl/zeus).
