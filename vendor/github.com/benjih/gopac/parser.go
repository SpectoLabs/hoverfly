// Copyright 2014 Jack Wakefield
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gopac

import (
	"errors"
	"io/ioutil"
	"net/http"
)

// Parser provides an interface to parse and utilise proxy auto-config (PAC)
// files.
type Parser struct {
	rt          *runtime
	initialised bool
}

func (parser *Parser) init() (err error) {
	var rt *runtime

	if rt, err = newRuntime(); err != nil {
		return
	}

	parser.rt = rt
	parser.initialised = true
	return
}

// Parse loads a proxy auto-config (PAC) file using the given path returning an
// error if the file fails to load.
func (parser *Parser) Parse(path string) error {
	if !parser.initialised {
		if err := parser.init(); err != nil {
			return err
		}
	}

	contents, err := ioutil.ReadFile(path)

	if err != nil {
		return err
	}

	return parser.rt.run(string(contents))
}

// ParseBytes loads a proxy auto-config (PAC) file using the given byte array
func (parser *Parser) ParseBytes(contents []byte) error {
	if !parser.initialised {
		if err := parser.init(); err != nil {
			return err
		}
	}

	return parser.rt.run(string(contents))
}

// ParseUrl downloads and parses a proxy auto-config (PAC) file using the given
// URL returning an error if the file fails to load.
func (parser *Parser) ParseUrl(url string) error {
	if !parser.initialised {
		if err := parser.init(); err != nil {
			return err
		}
	}

	response, err := http.Get(url)

	if err != nil {
		return err
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return err
	}

	return parser.rt.run(string(contents))
}

// FindProxy returns a proxy entry for the given URL or host returning an error
// if the attempt fails.
func (parser *Parser) FindProxy(url, host string) (string, error) {
	if parser.rt == nil {
		return "", errors.New("A proxy auto-config file has not been loaded")
	}

	return parser.rt.findProxyForURL(url, host)
}
