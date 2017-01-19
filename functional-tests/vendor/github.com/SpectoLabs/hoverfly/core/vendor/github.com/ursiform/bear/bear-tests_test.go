// Copyright 2015 Afshin Darian. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package bear

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type tester func(*testing.T)

// Generate tests for param requests using HandlerFunc.
func paramBearTest(
	label string, method string, path string,
	pattern string, want map[string]string) tester {
	return func(t *testing.T) {
		var (
			mux *Mux = New()
			req *http.Request
			res *httptest.ResponseRecorder
		)
		req, _ = http.NewRequest(method, path, nil)
		res = httptest.NewRecorder()
		handler := func(_ http.ResponseWriter, _ *http.Request, ctx *Context) {
			if !reflect.DeepEqual(want, ctx.Params) {
				t.Errorf(
					"%s %s (%s) %s got %v want %v",
					method, path, pattern, label, ctx.Params, want)
			}
		}
		mux.On(method, pattern, HandlerFunc(handler))
		mux.ServeHTTP(res, req)
		if res.Code != http.StatusOK {
			t.Errorf(
				"%s %s (%s) %s got %d want %d",
				method, path, pattern, label, res.Code, http.StatusOK)
		}
	}
}

// Generate tests for param requests using anonymous HandlerFunc
// compatible functions.
func paramBearAnonTest(
	label string, method string, path string,
	pattern string, want map[string]string) tester {
	return func(t *testing.T) {
		var (
			mux *Mux = New()
			req *http.Request
			res *httptest.ResponseRecorder
		)
		req, _ = http.NewRequest(method, path, nil)
		res = httptest.NewRecorder()
		handler := func(_ http.ResponseWriter, _ *http.Request, ctx *Context) {
			if !reflect.DeepEqual(want, ctx.Params) {
				t.Errorf(
					"%s %s (%s) %s got %v want %v",
					method, path, pattern, label, ctx.Params, want)
			}
		}
		mux.On(method, pattern, handler)
		mux.ServeHTTP(res, req)
		if res.Code != http.StatusOK {
			t.Errorf(
				"%s %s (%s) %s got %d want %d",
				method, path, pattern, label, res.Code, http.StatusOK)
		}
	}
}

// Generate tests for requests (i.e. no *Context) using http.HandlerFunc.
func simpleHttpTest(
	label string, method string, path string, pattern string, want int) tester {
	return func(t *testing.T) {
		var (
			mux *Mux = New()
			req *http.Request
			res *httptest.ResponseRecorder
		)
		req, _ = http.NewRequest(method, path, nil)
		res = httptest.NewRecorder()
		handler := func(http.ResponseWriter, *http.Request) {}
		mux.On(method, pattern, http.HandlerFunc(handler))
		mux.ServeHTTP(res, req)
		if res.Code != want {
			t.Errorf(
				"%s %s (%s) %s got %d want %d",
				method, path, pattern, label, res.Code, want)
		}
	}
}

// Generate tests for param AND no-param requests (i.e. no *Context) using
// anonymous http.HandlerFunc compatible func.
func simpleHttpAnonTest(
	label string, method string, path string, pattern string, want int) tester {
	return func(t *testing.T) {
		var (
			mux *Mux = New()
			req *http.Request
			res *httptest.ResponseRecorder
		)
		req, _ = http.NewRequest(method, path, nil)
		res = httptest.NewRecorder()
		mux.On(method, pattern, func(http.ResponseWriter, *http.Request) {})
		mux.ServeHTTP(res, req)
		if res.Code != want {
			t.Errorf(
				"%s %s (%s) %s got %d want %d",
				method, path, pattern, label, res.Code, want)
		}
	}
}

// Generate tests for simple no-param requests using HandlerFunc.
func simpleBearTest(
	label string, method string, path string, pattern string, want int) tester {
	return func(t *testing.T) {
		var (
			mux *Mux = New()
			req *http.Request
			res *httptest.ResponseRecorder
		)
		req, _ = http.NewRequest(method, path, nil)
		res = httptest.NewRecorder()
		handler := func(http.ResponseWriter, *http.Request, *Context) {}
		mux.On(method, pattern, HandlerFunc(handler))
		mux.ServeHTTP(res, req)
		if res.Code != want {
			t.Errorf(
				"%s %s (%s) %s got %d want %d",
				method, path, pattern, label, res.Code, want)
		}
	}
}

// Generate tests for simple no-param requests using anonymous functions that
// only accept a *Context param.
func simpleContextTest(
	label string, method string, path string, pattern string, want int) tester {
	return func(t *testing.T) {
		var (
			mux *Mux = New()
			req *http.Request
			res *httptest.ResponseRecorder
		)
		req, _ = http.NewRequest(method, path, nil)
		res = httptest.NewRecorder()
		handler := func(*Context) {}
		mux.On(method, pattern, handler)
		mux.ServeHTTP(res, req)
		if res.Code != want {
			t.Errorf(
				"%s %s (%s) %s got %d want %d",
				method, path, pattern, label, res.Code, want)
		}
	}
}

// Generate tests for simple no-param requests using anonymous HandlerFunc
// compatible functions.
func simpleBearAnonTest(label string, method string, path string,
	pattern string, want int) tester {
	return func(t *testing.T) {
		var (
			mux *Mux = New()
			req *http.Request
			res *httptest.ResponseRecorder
		)
		req, _ = http.NewRequest(method, path, nil)
		res = httptest.NewRecorder()
		handler := func(http.ResponseWriter, *http.Request, *Context) {}
		mux.On(method, pattern, handler)
		mux.ServeHTTP(res, req)
		if res.Code != want {
			t.Errorf(
				"%s %s (%s) %s got %d want %d",
				method, path, pattern, label, res.Code, want)
		}
	}
}

func TestBadHandler(t *testing.T) {
	mux := New()
	handlerBad := func(res http.ResponseWriter) {}
	if err := mux.On("*", "*", handlerBad); nil == err {
		t.Errorf("handlerBad should not be accepted")
	}
}
func TestBadVerbOn(t *testing.T) {
	mux := New()
	verb := "BLUB"
	handler := func(http.ResponseWriter, *http.Request) {}
	if err := mux.On(verb, "*", handler); nil == err {
		t.Errorf("%s should not be accepted", verb)
	}
}
func TestBadVerbServe(t *testing.T) {
	var (
		method  string = "BLUB"
		mux     *Mux   = New()
		pattern string = "/"
		path    string = "/"
		want    int    = http.StatusNotFound
		req     *http.Request
		res     *httptest.ResponseRecorder
	)
	handler := func(res http.ResponseWriter, req *http.Request) {}
	mux.On("*", pattern, handler)
	req, _ = http.NewRequest(method, path, nil)
	res = httptest.NewRecorder()
	mux.ServeHTTP(res, req)
	if status := res.Code; status != want {
		t.Errorf("%s %s got %d want %d", method, path, res.Code, want)
	}
}
func TestDuplicateFailure(t *testing.T) {
	var (
		handler HandlerFunc
		mux     *Mux   = New()
		pattern string = "/foo/{bar}"
	)
	handler = HandlerFunc(
		func(http.ResponseWriter, *http.Request, *Context) {})
	for _, verb := range verbs {
		if err := mux.On(verb, pattern, handler); err != nil {
			t.Error(err)
		} else if err := mux.On(verb, pattern, handler); err == nil {
			t.Errorf(
				"%s %s addition must fail because it is a duplicate",
				verb, pattern)
		}
	}
}
func TestDuplicateFailureRoot(t *testing.T) {
	var (
		handler HandlerFunc
		mux     *Mux   = New()
		pattern string = "/"
	)
	handler = HandlerFunc(
		func(http.ResponseWriter, *http.Request, *Context) {})
	for _, verb := range verbs {
		if err := mux.On(verb, pattern, handler); err != nil {
			t.Error(err)
		} else if err := mux.On(verb, pattern, handler); err == nil {
			t.Errorf(
				"%s %s addition must fail because it is a duplicate",
				verb, pattern)
		}
	}
}
func TestMiddleware(t *testing.T) {
	var (
		middlewares int  = 4
		mux         *Mux = New()
		params      map[string]string
		path        string = "/foo/BAR/baz/QUX"
		pattern     string = "/foo/{bar}/baz/{qux}"
		state       map[string]interface{}
		stateOne    = "one"
		stateTwo    = "two"
		stateThree  = "three"
		stateNil    = "nope"
	)
	params = map[string]string{"bar": "BAR", "qux": "QUX"}
	state = map[string]interface{}{"one": 1, "two": 2, "three": 3}
	run := func(method string) {
		var (
			req     *http.Request
			res     *httptest.ResponseRecorder
			visited int = 0
		)
		one := func(_ http.ResponseWriter, _ *http.Request, ctx *Context) {
			visited++
			if ctx.Get(stateNil) != nil {
				t.Errorf(
					"%s %s (%s) got %v want %v",
					method, path, pattern, ctx.Get(stateNil), nil)
			}
			ctx.Set("one", 1).Next()
		}
		two := func(_ http.ResponseWriter, _ *http.Request, ctx *Context) {
			visited++
			ctx.Set("two", 2).Next()
		}
		three := func(ctx *Context) {
			visited++
			ctx.Set("three", 3).Next()
		}
		last := func(_ http.ResponseWriter, _ *http.Request, ctx *Context) {
			visited++
			if !reflect.DeepEqual(params, ctx.Params) {
				t.Errorf(
					"%s %s (%s) got %v want %v",
					method, path, pattern, ctx.Params, params)
			}
			if ctx.Get(stateOne) != state[stateOne] {
				t.Errorf(
					"%s %s (%s) got %v want %v",
					method, path, pattern, ctx.Get(stateOne), state[stateOne])
			}
			if ctx.Get(stateTwo) != state[stateTwo] {
				t.Errorf(
					"%s %s (%s) got %v want %v",
					method, path, pattern, ctx.Get(stateTwo), state[stateTwo])
			}
			if ctx.Get(stateThree) != state[stateThree] {
				t.Errorf(
					"%s %s (%s) got %v want %v",
					method, path, pattern, ctx.Get(stateThree), state[stateThree])
			}
		}
		req, _ = http.NewRequest(method, path, nil)
		res = httptest.NewRecorder()
		mux.On(method, pattern, one, two, three, last)
		mux.ServeHTTP(res, req)
		if visited != middlewares {
			t.Errorf(
				"%s %s (%s) expected %d middlewares, visited %d",
				method, path, pattern, middlewares, visited)
		}
	}
	for _, verb := range verbs {
		run(verb)
	}
}
func TestMiddlewareRejectionA(t *testing.T) {
	mux := New()
	one := func(http.ResponseWriter, *http.Request, *Context) {}
	two := func(http.ResponseWriter, *http.Request) {}
	last := func(http.ResponseWriter, *http.Request, *Context) {}
	if err := mux.On("*", "*", one, two, last); err == nil {
		t.Errorf("middleware with wrong signature was accepted")
	}
}
func TestMiddlewareRejectionB(t *testing.T) {
	mux := New()
	one := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})
	two := HandlerFunc(
		func(http.ResponseWriter, *http.Request, *Context) {})
	last := HandlerFunc(
		func(http.ResponseWriter, *http.Request, *Context) {})
	if err := mux.On("*", "*", one, two, last); err == nil {
		t.Errorf("middleware with wrong signature was accepted")
	}
}
func TestMiddlewareRejectionC(t *testing.T) {
	mux := New()
	one := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})
	two := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})
	if err := mux.On("*", "*", one, two); err == nil {
		t.Errorf("middleware with wrong signature was accepted")
	}
}
func TestMiddlewareRejectionD(t *testing.T) {
	mux := New()
	one := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})
	two := func(http.ResponseWriter, *http.Request) {}
	if err := mux.On("*", "*", one, two); err == nil {
		t.Errorf("middleware with wrong signature was accepted")
	}
}
func TestMiddlewareRejectionE(t *testing.T) {
	mux := New()
	var (
		one   HandlerFunc
		two   http.HandlerFunc
		three func(http.ResponseWriter, *http.Request)
		four  func(http.ResponseWriter, *http.Request, *Context)
		five  func(*Context)
	)
	if err := mux.On("GET", "/", one); err == nil {
		t.Errorf("nil middleware was accepted")
	}
	if err := mux.On("GET", "/", two); err == nil {
		t.Errorf("nil middleware was accepted")
	}
	if err := mux.On("GET", "/", three); err == nil {
		t.Errorf("nil middleware was accepted")
	}
	if err := mux.On("GET", "/", four); err == nil {
		t.Errorf("nil middleware was accepted")
	}
	if err := mux.On("GET", "/", five); err == nil {
		t.Errorf("nil middleware was accepted")
	}
}
func TestNoHandlers(t *testing.T) {
	var (
		mux  *Mux   = New()
		path string = "/foo/bar"
		req  *http.Request
		res  *httptest.ResponseRecorder
		want int = http.StatusNotFound
	)
	for _, verb := range verbs {
		req, _ = http.NewRequest(verb, path, nil)
		res = httptest.NewRecorder()
		mux.ServeHTTP(res, req)
		if res.Code != want {
			t.Errorf("%s %s got %d want %d", verb, path, res.Code, want)
		}
	}
}
func TestNotFoundCustom(t *testing.T) {
	var (
		method       string = "GET"
		mux          *Mux   = New()
		pathFound    string = "/foo/bar"
		pathLost     string = "/foo/bar/baz"
		patternFound string = "/foo/bar"
		patternLost  string = "/*"
		req          *http.Request
		res          *httptest.ResponseRecorder
		wantFound    int = http.StatusOK
		wantLost     int = http.StatusTeapot
	)
	handlerFound := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		res.WriteHeader(http.StatusOK)
		res.Write([]byte("found!"))
	})
	handlerLost := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		res.WriteHeader(http.StatusTeapot)
		res.Write([]byte("not found!"))
	})
	// Test found to make sure wildcard doesn't overtake everything.
	req, _ = http.NewRequest(method, pathFound, nil)
	res = httptest.NewRecorder()
	mux.On(method, patternFound, handlerFound)
	mux.ServeHTTP(res, req)
	if res.Code != wantFound {
		t.Errorf(
			"%s %s (%s) got %d want %d",
			method, pathFound, patternFound, res.Code, wantFound)
	}
	// Test lost to make sure wildcard can gets non-pattern-matching paths.
	req, _ = http.NewRequest(method, pathLost, nil)
	res = httptest.NewRecorder()
	mux.On(method, patternLost, handlerLost)
	mux.ServeHTTP(res, req)
	if res.Code != wantLost {
		t.Errorf(
			"%s %s (%s) got %d want %d",
			method, pathLost, patternLost, res.Code, wantLost)
	}
}
func TestNotFoundNoParams(t *testing.T) {
	var (
		path    string = "/foo/bar"
		pattern string = "/foo"
		want    int    = http.StatusNotFound
	)
	for _, verb := range verbs {
		simpleHttpTest(
			"http.HandlerFunc",
			verb, path, pattern, want)(t)
		simpleHttpAnonTest(
			"anonymous http.HandlerFunc",
			verb, path, pattern, want)(t)
		simpleBearTest(
			"HandlerFunc",
			verb, path, pattern, want)(t)
		simpleBearAnonTest(
			"anonymous http.HandlerFunc",
			verb, path, pattern, want)(t)
		simpleContextTest(
			"anonymous *Context func",
			verb, path, pattern, want)(t)
	}
}
func TestNotFoundParams(t *testing.T) {
	var (
		path    string = "/foo/BAR/baz"
		pattern string = "/foo/{bar}/baz/{qux}"
		want    int    = http.StatusNotFound
	)
	for _, verb := range verbs {
		simpleHttpTest(
			"http.HandlerFunc",
			verb, path, pattern, want)(t)
		simpleHttpAnonTest(
			"anonymous http.HandlerFunc",
			verb, path, pattern, want)(t)
		simpleBearTest(
			"HandlerFunc",
			verb, path, pattern, want)(t)
		simpleBearAnonTest(
			"anonymous http.HandlerFunc",
			verb, path, pattern, want)(t)
		simpleContextTest(
			"anonymous *Context func",
			verb, path, pattern, want)(t)
	}
}
func TestNotFoundRoot(t *testing.T) {
	var (
		method string = "GET"
		mux    *Mux   = New()
		path   string = "/"
		req    *http.Request
		res    *httptest.ResponseRecorder
		want   int = http.StatusNotFound
	)
	req, _ = http.NewRequest(method, path, nil)
	res = httptest.NewRecorder()
	mux.ServeHTTP(res, req)
	if res.Code != want {
		t.Errorf("%s %s got %d want %d", method, path, res.Code, want)
	}
}
func TestNotFoundWildA(t *testing.T) {
	var (
		path    string = "/foo"
		pattern string = "/foo/*"
		want    int    = http.StatusNotFound
	)
	for _, verb := range verbs {
		simpleHttpTest(
			"http.HandlerFunc",
			verb, path, pattern, want)(t)
		simpleHttpAnonTest(
			"anonymous http.HandlerFunc",
			verb, path, pattern, want)(t)
		simpleBearTest(
			"HandlerFunc",
			verb, path, pattern, want)(t)
		simpleBearAnonTest(
			"anonymous http.HandlerFunc",
			verb, path, pattern, want)(t)
		simpleContextTest(
			"anonymous *Context func",
			verb, path, pattern, want)(t)
	}
}
func TestNotFoundWildB(t *testing.T) {
	var (
		method     string = "GET"
		mux        *Mux   = New()
		path       string = "/bar/baz"
		patternOne string = "/foo/*"
		patternTwo string = "/bar"
		req        *http.Request
		res        *httptest.ResponseRecorder
		want       int = http.StatusNotFound
	)
	handler := func(http.ResponseWriter, *http.Request) {}
	mux.On(method, patternOne, handler)
	mux.On(method, patternTwo, handler)
	req, _ = http.NewRequest(method, path, nil)
	res = httptest.NewRecorder()
	mux.ServeHTTP(res, req)
	if res.Code != want {
		t.Errorf("%s %s got %d want %d", method, path, res.Code, want)
	}
}
func TestOKMultiRuleParams(t *testing.T) {
	var (
		method     string = "GET"
		mux        *Mux   = New()
		path       string = "/bar/baz"
		patternOne string = "/foo/*"
		patternTwo string = "/bar/baz"
		req        *http.Request
		res        *httptest.ResponseRecorder
		want       int = http.StatusOK
	)
	handler := func(http.ResponseWriter, *http.Request) {}
	mux.On(method, patternOne, handler)
	mux.On(method, patternTwo, handler)
	req, _ = http.NewRequest(method, path, nil)
	res = httptest.NewRecorder()
	mux.ServeHTTP(res, req)
	if res.Code != want {
		t.Errorf("%s %s got %d want %d", method, path, res.Code, want)
	}
}
func TestOKNoParams(t *testing.T) {
	var (
		path    string = "/foo/bar"
		pattern string = "/foo/bar"
		want    int    = http.StatusOK
	)
	for _, verb := range verbs {
		simpleHttpTest(
			"http.HandlerFunc",
			verb, path, pattern, want)(t)
		simpleHttpAnonTest(
			"anonymous http.HandlerFunc",
			verb, path, pattern, want)(t)
		simpleBearTest(
			"HandlerFunc",
			verb, path, pattern, want)(t)
		simpleBearAnonTest(
			"anonymous http.HandlerFunc",
			verb, path, pattern, want)(t)
	}
}
func TestOKParams(t *testing.T) {
	var (
		path    string = "/foo/BAR/baz/QUX"
		pattern string = "/foo/{bar}/baz/{qux}"
		want    map[string]string
	)
	want = map[string]string{"bar": "BAR", "qux": "QUX"}
	for _, verb := range verbs {
		simpleHttpTest(
			"http.HandlerFunc",
			verb, path, pattern, http.StatusOK)(t)
		simpleHttpAnonTest(
			"anonymous http.HandlerFunc",
			verb, path, pattern, http.StatusOK)(t)
		paramBearTest(
			"HandlerFunc",
			verb, path, pattern, want)(t)
		paramBearAnonTest(
			"anonymous http.HandlerFunc",
			verb, path, pattern, want)(t)
	}
}
func TestOKParamsTrailingSlash(t *testing.T) {
	var (
		path    string = "/foo/BAR/baz/QUX/"
		pattern string = "/foo/{bar}/baz/{qux}"
		want    map[string]string
	)
	want = map[string]string{"bar": "BAR", "qux": "QUX"}
	for _, verb := range verbs {
		simpleHttpTest(
			"http.HandlerFunc",
			verb, path, pattern, http.StatusOK)(t)
		simpleHttpAnonTest(
			"anonymous http.HandlerFunc",
			verb, path, pattern, http.StatusOK)(t)
		paramBearTest(
			"HandlerFunc",
			verb, path, pattern, want)(t)
		paramBearAnonTest(
			"anonymous http.HandlerFunc",
			verb, path, pattern, want)(t)
	}
}
func TestOKRoot(t *testing.T) {
	var (
		path    string = "/"
		pattern string = "/"
		want    int    = http.StatusOK
	)
	for _, verb := range verbs {
		simpleHttpTest(
			"http.HandlerFunc",
			verb, path, pattern, want)(t)
		simpleHttpAnonTest(
			"anonymous http.HandlerFunc",
			verb, path, pattern, want)(t)
		simpleBearTest(
			"HandlerFunc",
			verb, path, pattern, want)(t)
		simpleBearAnonTest(
			"anonymous http.HandlerFunc",
			verb, path, pattern, want)(t)
	}
}
func TestOKWildRoot(t *testing.T) {
	var (
		path    string = "/"
		pattern string = "*"
		want    int    = http.StatusOK
	)
	for _, verb := range verbs {
		simpleHttpTest(
			"http.HandlerFunc",
			verb, path, pattern, want)(t)
		simpleHttpAnonTest(
			"anonymous http.HandlerFunc",
			verb, path, pattern, want)(t)
		simpleBearTest(
			"HandlerFunc",
			verb, path, pattern, want)(t)
		simpleBearAnonTest(
			"anonymous http.HandlerFunc",
			verb, path, pattern, want)(t)
	}
}
func TestSanitizePatternPrefixSuffix(t *testing.T) {
	var (
		method  string = "GET"
		mux     *Mux   = New()
		pattern string = "foo/{bar}/*"
		path    string = "/foo/ABC/baz"
		want    string = "/foo/{bar}/*/"
		req     *http.Request
		res     *httptest.ResponseRecorder
	)
	handler := func(res http.ResponseWriter, _ *http.Request, ctx *Context) {
		res.WriteHeader(http.StatusOK)
		res.Write([]byte(ctx.tree.pattern))
	}
	mux.On(method, pattern, handler)
	req, _ = http.NewRequest(method, path, nil)
	res = httptest.NewRecorder()
	mux.ServeHTTP(res, req)
	if body := res.Body.String(); body != want {
		t.Errorf("%s %s (%s) got %s want %s", method, path, pattern, body, want)
	}
}
func TestSanitizePatternDoubleSlash(t *testing.T) {
	var (
		method  string = "GET"
		mux     *Mux   = New()
		pattern string = "///foo///{bar}////*//"
		path    string = "/foo/ABC/baz"
		want    string = "/foo/{bar}/*/"
		req     *http.Request
		res     *httptest.ResponseRecorder
	)
	handler := func(res http.ResponseWriter, _ *http.Request, ctx *Context) {
		res.WriteHeader(http.StatusOK)
		res.Write([]byte(ctx.tree.pattern))
	}
	mux.On(method, pattern, handler)
	req, _ = http.NewRequest(method, path, nil)
	res = httptest.NewRecorder()
	mux.ServeHTTP(res, req)
	if body := res.Body.String(); body != want {
		t.Errorf("%s %s (%s) got %s want %s", method, path, pattern, body, want)
	}
}
func TestUnreachableA(t *testing.T) {
	mux := New()
	one := func(http.ResponseWriter, *http.Request, *Context) {}
	two := func(http.ResponseWriter, *http.Request) {}
	three := func(http.ResponseWriter, *http.Request, *Context) {}
	err := mux.On("*", "*", one, two, three)
	if nil == err {
		t.Errorf("unreachable A")
	}
}
func TestUnreachableB(t *testing.T) {
	mux := New()
	one := func(http.ResponseWriter, *http.Request, *Context) {}
	two := func(http.ResponseWriter, *http.Request) {}
	three := func(*Context) {}
	err := mux.On("*", "*", one, two, three)
	if nil == err {
		t.Errorf("unreachable B")
	}
}
func TestWildcardCompeting(t *testing.T) {
	var (
		method       string = "GET"
		mux          *Mux   = New()
		patternOne   string = "/*"
		pathOne      string = "/bar/baz"
		wantOne      string = "bar/baz"
		patternTwo   string = "/foo/*"
		pathTwo      string = "/foo/baz"
		wantTwo      string = "baz"
		patternThree string = "/foo/bar/*"
		pathThree    string = "/foo/bar/bar/baz"
		wantThree    string = "bar/baz"
		req          *http.Request
		res          *httptest.ResponseRecorder
	)
	handler := func(res http.ResponseWriter, _ *http.Request, ctx *Context) {
		res.WriteHeader(http.StatusOK)
		res.Write([]byte(ctx.Params["*"]))
	}
	mux.On(method, patternOne, handler)
	mux.On(method, patternTwo, handler)
	mux.On(method, patternThree, handler)
	req, _ = http.NewRequest(method, pathOne, nil)
	res = httptest.NewRecorder()
	mux.ServeHTTP(res, req)
	if body := res.Body.String(); body != wantOne {
		t.Errorf(
			"%s %s (%s) got %s want %s",
			method, pathOne, patternOne, body, wantOne)
	}
	req, _ = http.NewRequest(method, pathTwo, nil)
	res = httptest.NewRecorder()
	mux.ServeHTTP(res, req)
	if body := res.Body.String(); body != wantTwo {
		t.Errorf(
			"%s %s (%s) got %s want %s",
			method, pathTwo, patternTwo, body, wantTwo)
	}
	req, _ = http.NewRequest(method, pathThree, nil)
	res = httptest.NewRecorder()
	mux.ServeHTTP(res, req)
	if body := res.Body.String(); body != wantThree {
		t.Errorf(
			"%s %s (%s) got %s want %s",
			method, pathThree, patternThree, body, wantThree)
	}
}
func TestWildcardMethod(t *testing.T) {
	var (
		mux     *Mux   = New()
		path    string = "/foo/bar"
		pattern string = "/foo/bar"
		req     *http.Request
		res     *httptest.ResponseRecorder
		want    int = http.StatusOK
	)
	mux.On(
		"*",
		pattern,
		func(http.ResponseWriter, *http.Request, *Context) {})
	for _, verb := range verbs {
		req, _ = http.NewRequest(verb, path, nil)
		res = httptest.NewRecorder()
		mux.ServeHTTP(res, req)
		if res.Code != want {
			t.Errorf(
				"%s %s (%s) got %d want %d",
				verb, path, pattern, res.Code, want)
		}
	}
}
func TestWildcardMethodWarning(t *testing.T) {
	var (
		mux     *Mux   = New()
		pattern string = "/foo/bar/"
		verb    string = "GET"
	)
	handler := func(http.ResponseWriter, *http.Request, *Context) {}
	if err := mux.On(verb, pattern, handler); err != nil {
		t.Errorf("%s %s should have registered", verb, pattern)
	}
	if err := mux.On("*", pattern, handler); err == nil {
		t.Errorf("* %s should have registered, but with an error", pattern)
	}
}
func TestWildcardNotLast(t *testing.T) {
	var (
		mux     *Mux   = New()
		pattern string = "/foo/*/bar"
	)
	handler := func(res http.ResponseWriter, req *http.Request) {}
	err := mux.On("*", pattern, handler)
	if err == nil {
		t.Errorf(
			"wildcard pattern (%s) with non-final wildcard was accepted",
			pattern)
	}
}
func TestWildcardParams(t *testing.T) {
	var (
		method  string = "GET"
		mux     *Mux   = New()
		pattern string = "/foo/{bar}/*"
		path    string = "/foo/ABC/baz"
		want    string = "ABC"
		req     *http.Request
		res     *httptest.ResponseRecorder
	)
	handler := func(res http.ResponseWriter, _ *http.Request, ctx *Context) {
		res.WriteHeader(http.StatusOK)
		res.Write([]byte(ctx.Params["bar"]))
	}
	mux.On(method, pattern, handler)
	req, _ = http.NewRequest(method, path, nil)
	res = httptest.NewRecorder()
	mux.ServeHTTP(res, req)
	if body := res.Body.String(); body != want {
		t.Errorf("%s %s (%s) got %s want %s", method, path, pattern, body, want)
	}
}
func TestAlways(t *testing.T) {
	var (
		mux      *Mux = New()
		pattern       = "/foo"
		method        = "GET"
		path          = "/foo"
		req      *http.Request
		res      *httptest.ResponseRecorder
		keyOne   = "one"
		stateOne = 1
		keyTwo   = "two"
		stateTwo = 2
	)
	one := func(_ http.ResponseWriter, _ *http.Request, ctx *Context) {
		ctx.Set(keyOne, stateOne).Next()
	}
	two := func(_ http.ResponseWriter, _ *http.Request, ctx *Context) {
		ctx.Set(keyTwo, stateTwo).Next()
	}
	three := func(res http.ResponseWriter, _ *http.Request, ctx *Context) {
		first := reflect.DeepEqual(ctx.Get(keyOne), stateOne)
		second := reflect.DeepEqual(ctx.Get(keyTwo), stateTwo)
		if !first || !second {
			t.Errorf("Always middlewares did not execute before On middleware")
		}
	}
	mux.Always(one)
	mux.Always(HandlerFunc(two))
	mux.On(method, pattern, three)
	req, _ = http.NewRequest(method, path, nil)
	res = httptest.NewRecorder()
	mux.ServeHTTP(res, req)
}

func TestAlwaysBeforeNotFound(t *testing.T) {
	var (
		mux      *Mux = New()
		pattern       = "/foo"
		method        = "GET"
		path          = "/bar"
		req      *http.Request
		res      *httptest.ResponseRecorder
		keyOne   = "one"
		stateOne = 1
	)
	always := func(_ http.ResponseWriter, _ *http.Request, ctx *Context) {
		ctx.Set(keyOne, stateOne).Next()
	}
	handler := func(_ http.ResponseWriter, _ *http.Request, _ *Context) {
		t.Errorf("handler should not be fired because path != pattern")
	}
	notFound := func(_ http.ResponseWriter, _ *http.Request, ctx *Context) {
		first := reflect.DeepEqual(ctx.Get(keyOne), stateOne)
		if !first {
			t.Errorf("Always middleware did not execute before notFound")
		}
	}
	mux.Always(always)
	mux.On(method, pattern, handler)
	mux.On("*", "/*", notFound)
	req, _ = http.NewRequest(method, path, nil)
	res = httptest.NewRecorder()
	mux.ServeHTTP(res, req)
}
func TestAlwaysRejection(t *testing.T) {
	var (
		mux   *Mux = New()
		one   HandlerFunc
		two   http.HandlerFunc
		three func(http.ResponseWriter, *http.Request)
		four  func(http.ResponseWriter, *http.Request, *Context)
	)
	if err := mux.Always(one); err == nil {
		t.Errorf("Always requires non-nil HandlerFunc or its signature")
	}
	if err := mux.Always(two); err == nil {
		t.Errorf("Always requires non-nil HandlerFunc or its signature")
	}
	if err := mux.Always(three); err == nil {
		t.Errorf("Always requires non-nil HandlerFunc or its signature")
	}
	if err := mux.Always(four); err == nil {
		t.Errorf("Always requires non-nil HandlerFunc or its signature")
	}
}
