// Copyright 2015 Afshin Darian. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package bear

import (
	"fmt"
	"net/http"
)

// HandlerFunc is similar to http.HandlerFunc, except it requires
// an extra argument for the *Context of a request.
type HandlerFunc func(http.ResponseWriter, *http.Request, *Context)

// handlerize takes one of handler formats that Mux accepts.
// It returns a HandlerFunc, a flag indicating whether the HandlerFunc
// can be follwed by other handlers, and any error that may have arisen in
// conversion.
func handlerize(function interface{}) (HandlerFunc, bool, error) {
	followable := true
	unfollowable := false
	switch function.(type) {
	case HandlerFunc:
		handler := function.(HandlerFunc)
		if handler == nil {
			return nil, unfollowable, fmt.Errorf("nil middleware")
		} else {
			return HandlerFunc(handler), followable, nil
		}
	case func(*Context):
		handler := function.(func(*Context))
		if handler == nil {
			return nil, unfollowable, fmt.Errorf("nil middleware")
		} else {
			handler := HandlerFunc(
				func(_ http.ResponseWriter, _ *http.Request, ctx *Context) {
					handler(ctx)
				})
			return handler, followable, nil
		}
	case func(http.ResponseWriter, *http.Request, *Context):
		handler := function.(func(http.ResponseWriter, *http.Request, *Context))
		if handler == nil {
			return nil, unfollowable, fmt.Errorf("nil middleware")
		} else {
			return HandlerFunc(handler), followable, nil
		}
	case http.HandlerFunc:
		handler := function.(http.HandlerFunc)
		if handler == nil {
			return nil, unfollowable, fmt.Errorf("nil middleware")
		} else {
			handler := HandlerFunc(
				func(res http.ResponseWriter, req *http.Request, _ *Context) {
					handler(res, req)
				})
			return handler, unfollowable, nil
		}
	case func(http.ResponseWriter, *http.Request):
		handler := function.(func(http.ResponseWriter, *http.Request))
		if handler == nil {
			return nil, unfollowable, fmt.Errorf("nil middleware")
		} else {
			handler := HandlerFunc(
				func(res http.ResponseWriter, req *http.Request, _ *Context) {
					handler(res, req)
				})
			return handler, unfollowable, nil
		}
	default:
		err := fmt.Errorf(
			"handler must match: %s, %s, or %s",
			"http.HandlerFunc", "bear.HandlerFunc", "func(*Context)")
		return nil, unfollowable, err
	}
}

func handlerizeLax(
	verb string, pattern string, functions []interface{}) ([]HandlerFunc, error) {
	var handlers []HandlerFunc
	unreachable := false
	for _, function := range functions {
		if unreachable {
			err := fmt.Errorf(
				"bear: %s %s has unreachable middleware",
				verb, pattern)
			return nil, err
		}
		if handler, followable, err := handlerize(function); err != nil {
			return nil, fmt.Errorf("bear: %s %s: %s", verb, pattern, err)
		} else {
			if !followable {
				unreachable = true
			}
			handlers = append(handlers, handler)
		}
	}
	return handlers, nil
}

func handlerizeStrict(functions []interface{}) (handlers []HandlerFunc, err error) {
	for _, function := range functions {
		switch function.(type) {
		case HandlerFunc:
			handler := function.(HandlerFunc)
			if handler == nil {
				return nil, fmt.Errorf("bear: nil middleware")
			} else {
				handlers = append(handlers, HandlerFunc(handler))
			}
		case func(http.ResponseWriter, *http.Request, *Context):
			handler := function.(func(http.ResponseWriter, *http.Request, *Context))
			if handler == nil {
				return nil, fmt.Errorf("bear: nil middleware")
			} else {
				handlers = append(handlers, HandlerFunc(handler))
			}
		default:
			return nil, fmt.Errorf(
				"bear: handler must be a bear.HandlerFunc or match its signature")
		}
	}
	return handlers, nil
}
