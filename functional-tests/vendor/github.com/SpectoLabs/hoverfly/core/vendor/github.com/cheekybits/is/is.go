package is

import (
	"fmt"
	"reflect"
	"sync"
)

// I represents the is interface.
type I interface {
	// OK asserts that the specified objects are all OK.
	OK(o ...interface{})
	// Equal asserts that the two values are
	// considered equal. Non strict.
	Equal(a, b interface{})
	// NotEqual asserts that the two values are not
	// considered equal. Non strict.
	NotEqual(a, b interface{})
	// NoErr asserts that the value is not an
	// error.
	NoErr(err ...error)
	// Err asserts that the value is an error.
	Err(err ...error)
	// Nil asserts that the specified objects are
	// all nil.
	Nil(obj ...interface{})
	// NotNil asserts that the specified objects are
	// all nil.
	NotNil(obj ...interface{})
	// True asserts that the specified objects are
	// all true
	True(obj ...interface{})
	// False asserts that the specified objects are
	// all true
	False(obj ...interface{})
	// Panic asserts that the specified function
	// panics.
	Panic(fn func())
	// PanicWith asserts that the specified function
	// panics with the specific message.
	PanicWith(m string, fn func())
	// Fail indicates that the test has failed with the
	// specified formatted arguments.
	Fail(args ...interface{})
	// Failf indicates that the test has failed with the
	// formatted arguments.
	Failf(format string, args ...interface{})
}

// New creates a new I capable of making
// assertions.
func New(t T) I {
	return &i{t: t}
}

// Relaxed creates a new I capable of making
// assertions, but will not fail immediately
// allowing all assertions to run.
func Relaxed(t T) I {
	return &i{t: t, relaxed: true}
}

// T represents the an interface for reporting
// failures.
// testing.T satisfied this interface.
type T interface {
	FailNow()
}

// i represents an implementation of interface I.
type i struct {
	t       T
	fails   []string
	l       sync.Mutex
	relaxed bool
}

func (i *i) Log(args ...interface{}) {
	i.l.Lock()
	fail := fmt.Sprint(args...)
	i.fails = append(i.fails, fail)
	fmt.Print(decorate(fail))
	i.l.Unlock()
	if !i.relaxed {
		i.t.FailNow()
	}
}
func (i *i) Logf(format string, args ...interface{}) {
	i.l.Lock()
	fail := fmt.Sprintf(format, args...)
	i.fails = append(i.fails, fail)
	fmt.Print(decorate(fail))
	i.l.Unlock()
	if !i.relaxed {
		i.t.FailNow()
	}
}

// OK asserts that the specified objects are all OK.
func (i *i) OK(o ...interface{}) {
	for _, obj := range o {

		if isNil(obj) {
			i.Log("unexpected nil")
		}

		switch co := obj.(type) {
		case func():
			// shouldn't panic
			var r interface{}
			func() {
				defer func() {
					r = recover()
				}()
				co()
			}()
			if r != nil {
				i.Logf("unexpected panic: %v", r)
			}
			return
		case string:
			if len(co) == 0 {
				i.Log("unexpected \"\"")
			}
			return
		case bool:
			// false
			if co == false {
				i.Log("unexpected false")
				return
			}
		}
		if isNil(o) {
			if _, ok := obj.(error); ok {
				// nil errors are ok
				return
			}
			i.Log("unexpected nil")
			return
		}

		if obj == 0 {
			i.Log("unexpected zero")
		}
	}
}

func (i *i) NoErr(errs ...error) {
	for n, err := range errs {
		if !isNil(err) {
			p := "unexpected error"
			if len(errs) > 1 {
				p += fmt.Sprintf(" (%d)", n)
			}
			p += ": " + err.Error()
			i.Logf(p)
		}
	}
}

func (i *i) Err(errs ...error) {
	for n, err := range errs {
		if isNil(err) {
			p := "error expected"
			if len(errs) > 1 {
				p += fmt.Sprintf(" (%d)", n)
			}
			i.Logf(p)
		}
	}
}

func (i *i) Nil(o ...interface{}) {
	for n, obj := range o {
		if !isNil(obj) {
			p := "expected nil"
			if len(o) > 1 {
				p += fmt.Sprintf(" (%d)", n)
			}
			p += ": " + fmt.Sprintf("%#v", obj)
			i.Logf(p)
		}
	}
}

func (i *i) NotNil(o ...interface{}) {
	for n, obj := range o {
		if isNil(obj) {
			p := "unexpected nil"
			if len(o) > 1 {
				p += fmt.Sprintf(" (%d)", n)
			}
			p += ": " + fmt.Sprintf("%#v", obj)
			i.Logf(p)
		}
	}
}

func (i *i) True(o ...interface{}) {
	for n, obj := range o {
		if b, ok := obj.(bool); !ok || b != true {
			p := "expected true"
			if len(o) > 1 {
				p += fmt.Sprintf(" (%d)", n)
			}
			p += ": " + fmt.Sprintf("%#v", obj)
			i.Logf(p)
		}
	}
}

func (i *i) False(o ...interface{}) {
	for n, obj := range o {
		if b, ok := obj.(bool); !ok || b != false {
			p := "expected false"
			if len(o) > 1 {
				p += fmt.Sprintf(" (%d)", n)
			}
			p += ": " + fmt.Sprintf("%#v", obj)
			i.Logf(p)
		}
	}
}

// Equal asserts that the two values are
// considered equal. Non strict.
func (i *i) Equal(a, b interface{}) {
	if !areEqual(a, b) {
		i.Logf("%v != %v", a, b)
	}
}

// NotEqual asserts that the two values are not
// considered equal. Non strict.
func (i *i) NotEqual(a, b interface{}) {
	if areEqual(a, b) {
		i.Logf("%v == %v", a, b)
	}
}

func (i *i) Fail(args ...interface{}) {
	i.Log(args...)
}

func (i *i) Failf(format string, args ...interface{}) {
	i.Logf(format, args...)
}

// Panic asserts that the specified function
// panics.
func (i *i) Panic(fn func()) {
	var r interface{}
	func() {
		defer func() {
			r = recover()
		}()
		fn()
	}()
	if r == nil {
		i.Log("expected panic")
	}
}

// PanicWith asserts that the specified function
// panics with the specific message.
func (i *i) PanicWith(m string, fn func()) {
	var r interface{}
	func() {
		defer func() {
			r = recover()
		}()
		fn()
	}()
	if r != m {
		i.Logf("expected panic: \"%s\"", m)
	}
}

// isNil gets whether the object is nil or not.
func isNil(object interface{}) bool {
	if object == nil {
		return true
	}
	value := reflect.ValueOf(object)
	kind := value.Kind()
	if kind >= reflect.Chan && kind <= reflect.Slice && value.IsNil() {
		return true
	}
	return false
}

// areEqual gets whether a equals b or not.
func areEqual(a, b interface{}) bool {
	if isNil(a) || isNil(b) {
		if isNil(a) && !isNil(b) {
			return false
		}
		if !isNil(a) && isNil(b) {
			return false
		}
		return a == b
	}
	if reflect.DeepEqual(a, b) {
		return true
	}
	aValue := reflect.ValueOf(a)
	bValue := reflect.ValueOf(b)
	if aValue == bValue {
		return true
	}

	// Attempt comparison after type conversion
	if bValue.Type().ConvertibleTo(aValue.Type()) {
		return reflect.DeepEqual(a, bValue.Convert(aValue.Type()).Interface())
	}
	// Last ditch effort
	if fmt.Sprintf("%#v", a) == fmt.Sprintf("%#v", b) {
		return true
	}

	return false
}
