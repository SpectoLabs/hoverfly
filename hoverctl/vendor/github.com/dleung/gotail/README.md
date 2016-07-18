# Gotail

[![Build Status](https://travis-ci.org/dleung/gotail.svg?branch=master)](https://travis-ci.org/dleung/gotail)

Blog: (http://capykoa.com/articles/11)

### Background
I needed a simple utility to follow rapidly growing adserver logrotating logs that I can process dynamically and in real time.

Features
- Default behavior supports logrotating files.
- Timeout option to wait on new files that hasn't been created yet.
- Supports high reads.  On my 2013 Macbook Pro, I was able benchmark at more than 150,000 lines/second (`go test`).  However, expected performance on production should be much higher.
- Lightweight and low memory footprint
- Partly inspired by [ActiveState/tail](https://github.com/ActiveState/tail)

### Usage
```go
package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/dleung/gotail"
)

var fname string

func main() {
	flag.StringVar(&fname, "file", "", "File to tail")
	flag.Parse()

	// Set timeout to 0 to automatically fail if file isn't found
	tail, err := gotail.NewTail(fname, gotail.Config{Timeout: 10})
	if err != nil {
		log.Fatalln(err)
	}

	// lines on the tail.Lines channel for new lines.
	for line := range tail.Lines {
		fmt.Println(line)
	}
}
```
