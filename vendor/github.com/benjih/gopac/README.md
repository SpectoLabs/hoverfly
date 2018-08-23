# GoPac [![](http://img.shields.io/travis/jackwakefield/gopac.svg?style=flat-square)](http://travis-ci.org/jackwakefield/gopac) 

GoPac is a package for parsing and using proxy auto-config (PAC) files.

## Installation

```
go get github.com/jackwakefield/gopac
```

## Usage

[See GoDoc](http://godoc.org/github.com/jackwakefield/gopac) for documentation.

## Example

```go
package main

import (
	"log"

	"github.com/jackwakefield/gopac"
)

func main() {
	parser := new(gopac.Parser)

	// use parser.Parse(path) to parse a local file
	// or parser.ParseUrl(url) to parse a remote file
	if err := parser.ParseUrl("http://immun.es/pac"); err != nil {
		log.Fatalf("Failed to parse PAC (%s)", err)
	}

	// find the proxy entry for host check.immun.es
	entry, err := parser.FindProxy("", "check.immun.es")

	if err != nil {
		log.Fatalf("Failed to find proxy entry (%s)", err)
	}

	log.Println(entry)
}
```

## License

> Copyright 2014 Jack Wakefield
>
> Licensed under the Apache License, Version 2.0 (the "License");
> you may not use this file except in compliance with the License.
> You may obtain a copy of the License at
>
>     http://www.apache.org/licenses/LICENSE-2.0
>
> Unless required by applicable law or agreed to in writing, software
> distributed under the License is distributed on an "AS IS" BASIS,
> WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
> See the License for the specific language governing permissions and
> limitations under the License.