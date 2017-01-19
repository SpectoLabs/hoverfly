# try [![GoDoc](https://godoc.org/github.com/matryer/try?status.svg)](https://godoc.org/github.com/matryer/try)

Idiomatic Go retry package. Thanks to [@rowland](https://github.com/rowland) for code review.

```
go get gopkg.in/matryer/try.v1
```

### Usage

Just call `try.Do` with the function you want to retry in the event of an error:

  * Call `try.Do` that returns a `bool` indicating whether to retry or not, and an `error` 
  * The `attempt` argument will start at 1 and count up
  * `try.Do` blocks until you return `false`, or a `nil` error
  * `try.Do` returns the last error or `nil` if it was successful

```
var value string
err := try.Do(func(attempt int) (bool, error) {
  var err error
  value, err = SomeFunction()
  return attempt < 5, err // try 5 times
})
if err != nil {
  log.Fatalln("error:", err)
}
```

In the above example the function will be called repeatedly until error is `nil`, while `attempt < 5` (i.e. try 5 times)

#### Retrying panics

Try supports retrying in the event of a panic.

  * Use named return parameters
  * Set `retry` first
  * Defer the recovery code, and set `err` manually in the case of a panic
  * Use empty `return` statement at the end

```
var value string
err := try.Do(func(attempt int) (retry bool, err error) {
  retry = attempt < 5 // try 5 times
  defer func() {
    if r := recover(); r != nil {
      err = errors.New(fmt.Sprintf("panic: %v", r))
    }
  }()
  value, err = SomeFunction()
  return
})
if err != nil {
  log.Fatalln("error:", err)
}
```

#### Delay between retries

To introduce a delay between retries, just make a `time.Sleep` call before you return from the function if you are returning an error.

```
var value string
err := try.Do(func(attempt int) (bool, error) {
  var err error
  value, err = SomeFunction()
  if err != nil {
    time.Sleep(1 * time.Minute) // wait a minute
  }
  return attempt < 5, err
})
if err != nil {
  log.Fatalln("error:", err)
}
```

#### Maximum retry limit

To avoid infinite loops, Try will ensure it only makes `try.MaxRetries` attempts. By default, this value is `10`, but you can change it:

```
try.MaxRetries = 20
```

To see if a `Do` operation failed due to reaching the limit, you can check the `error` with `try.IsMaxRetries(err)`.
