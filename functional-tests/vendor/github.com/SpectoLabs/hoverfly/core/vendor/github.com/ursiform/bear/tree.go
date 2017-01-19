// Copyright 2015 Afshin Darian. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package bear

import (
	"fmt"
	"strings"
)

type tree struct {
	children map[string]*tree
	handlers []HandlerFunc
	name     string
	pattern  string
}

func parsePattern(s string) (pattern string, components []string, last int) {
	if slashr != s[0] {
		s = slash + s // start with slash
	}
	if slashr != s[len(s)-1] {
		s = s + slash // end with slash
	}
	pattern = dbl.ReplaceAllString(s, slash)
	components = strings.SplitAfter(pattern, slash)
	components = components[1 : len(components)-1]
	last = len(components) - 1
	return pattern, components, last
}

func (tr *tree) set(verb string, pattern string, handlers []HandlerFunc,
	wildcards *bool, err *error) {
	if pattern == slash || pattern == empty {
		if nil != tr.handlers {
			*err = fmt.Errorf("bear: %s %s exists, ignoring", verb, pattern)
			return
		} else {
			tr.pattern = slash
			tr.handlers = handlers
			return
		}
	}
	if nil == tr.children {
		tr.children = make(map[string]*tree)
	}
	current := &tr.children
	pattern, components, last := parsePattern(pattern)
	for index, component := range components {
		var (
			match []string = dyn.FindStringSubmatch(component)
			key   string   = component
			name  string
		)
		if 0 < len(match) {
			key, name = dynamic, match[1]
		} else if key == lasterisk {
			key, name = wildcard, asterisk
			*wildcards = true
		}
		if nil == (*current)[key] {
			(*current)[key] = &tree{
				children: make(map[string]*tree), name: name}
		}
		if index == last {
			if nil != (*current)[key].handlers {
				*err = fmt.Errorf("bear: %s %s exists, ignoring", verb, pattern)
				return
			}
			(*current)[key].pattern = pattern
			(*current)[key].handlers = handlers
			return
		} else if key == wildcard {
			*err = fmt.Errorf("bear: %s %s wildcard (%s) token must be last",
				verb, pattern, asterisk)
			return
		}
		current = &(*current)[key].children
	}
}
