// Package try provides retry functionality.
//     var value string
//     err := try.Do(func(attempt int) (bool, error) {
//       var err error
//       value, err = SomeFunction()
//       return attempt < 5, err // try 5 times
//     })
//     if err != nil {
//       log.Fatalln("error:", err)
//     }
package try
