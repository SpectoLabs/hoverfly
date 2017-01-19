package bear_test

import (
	"fmt"
	"net/http"

	"github.com/ursiform/bear"
)

func ExampleMux_On_error() {
	mux := bear.New()
	handlerOne := func(http.ResponseWriter, *http.Request) {}
	handlerTwo := func(http.ResponseWriter, *http.Request) {}
	if err := mux.On("GET", "/foo/", handlerOne); err != nil {
		fmt.Println(err)
	} else if err := mux.On("*", "/foo/", handlerTwo); err != nil {
		fmt.Println(err)
	}
	// Output: bear: GET /foo/ exists, ignoring
}
