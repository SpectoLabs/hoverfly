is
==

[![GoDoc](https://godoc.org/github.com/cheekybits/is?status.png)](http://godoc.org/github.com/cheekybits/is)

A mini testing helper for Go.

  * Simple interface (`is.OK` and `is.Equal`)
  * Plugs into existing Go toolchain (uses `testing.T`)
  * Obvious for newcomers and newbs
  * Also gives you `is.Panic` and `is.PanicWith` helpers - because testing panics is ugly

### Usage

  1. Write test functions as usual
  1. Add `is := is.New(t)` at top of your test functions
  1. Call target code
  1. Make assertions using new `is` object

```
func TestSomething(t *testing.T) {
  is := is.New(t)

  // ensure not nil
  obj := SomeFunc()
  is.OK(obj)

  // ensure no error
  obj, err := SomeFunc()
  is.NoErr(err)

  // ensure not false
  b := SomeBool()
  is.OK(b)

  // ensure not ""
  s := SomeString()
  is.OK(s)

  // ensure not zero
  is.OK(len(something))

  // ensure doesn't panic
  is.OK(func(){
    MethodShouldNotPanic()
  })

  // ensure many things in one go
  is.OK(b, err, obj, "something")

  // ensure something does panic
  is.Panic(func(){
    MethodShouldPanic(1)
  })
  is.PanicWith("package: arg must be >0", func(){
    MethodShouldPanicWithSpecificMessage(0)
  })

  // make sure two values are equal
  is.Equal(1, 2)
  is.Equal(err, ErrSomething)
  is.Equal(a, b)

}
```

### Get started

Get it:

```
go get github.com/cheekybits/is
```

Then import it:

```
import (
  "testing"
  "github.com/cheekybits/is"
)
```
