/********************************
*** Multiplexer for Go        ***
*** Bone is under MIT license ***
*** Code by CodingFerret      ***
*** github.com/go-zoo         ***
*********************************/

package bone

import (
	"net/http"
	"regexp"
	"strings"
)

const (
	PARAM = 2
	SUB   = 4
	WC    = 8
	REGEX = 16
)

// Route content the required information for a valid route
// Path: is the Route URL
// Size: is the length of the path
// Token: is the value of each part of the path, split by /
// Pattern: is content information about the route, if it's have a route variable
// handler: is the handler who handle this route
// Method: define HTTP method on the route
type Route struct {
	Path    string
	Method  string
	Size    int
	Atts    int
	wildPos int
	Token   Token
	Pattern map[int]string
	Compile map[int]*regexp.Regexp
	Tag     map[int]string
	Handler http.Handler
}

// Token content all value of a spliting route path
// Tokens: string value of each token
// size: number of token
type Token struct {
	raw    []int
	Tokens []string
	Size   int
}

// NewRoute return a pointer to a Route instance and call save() on it
func NewRoute(url string, h http.Handler) *Route {
	r := &Route{Path: url, Handler: h}
	r.save()
	return r
}

// Save, set automatically the the Route.Size and Route.Pattern value
func (r *Route) save() {
	r.Size = len(r.Path)
	r.Token.Tokens = strings.Split(r.Path, "/")
	for i, s := range r.Token.Tokens {
		if len(s) >= 1 {
			switch s[:1] {
			case ":":
				if r.Pattern == nil {
					r.Pattern = make(map[int]string)
				}
				r.Pattern[i] = s[1:]
				r.Atts += PARAM
			case "#":
				if r.Compile == nil {
					r.Compile = make(map[int]*regexp.Regexp)
					r.Tag = make(map[int]string)
				}
				tmp := strings.Split(s, "^")
				r.Tag[i] = tmp[0][1:]
				r.Compile[i] = regexp.MustCompile("^" + tmp[1][:len(tmp[1])-1])
				r.Atts += REGEX
			case "*":
				r.wildPos = i
				r.Atts += WC
			default:
				r.Token.raw = append(r.Token.raw, i)
			}
		}
		r.Token.Size++
	}
}

// Match check if the request match the route Pattern
func (r *Route) Match(req *http.Request) bool {
	ss := strings.Split(req.URL.Path, "/")

	if r.matchRawTokens(&ss) {
		if len(ss) == r.Token.Size || r.Atts&WC != 0 {
			vars.Lock()
			if vars.v[req] == nil {
				vars.v[req] = make(map[string]string)
			}
			for k, v := range r.Pattern {
				vars.v[req][v] = ss[k]
			}
			vars.Unlock()
			if r.Atts&REGEX != 0 {
				for k, v := range r.Compile {
					if !v.MatchString(ss[k]) {
						return false
					}
					vars.Lock()
					vars.v[req][r.Tag[k]] = ss[k]
					vars.Unlock()
				}
			}
			return true
		}
	}
	return false
}

func (r *Route) parse(rw http.ResponseWriter, req *http.Request) bool {
	if r.Atts != 0 {
		if r.Atts&SUB != 0 {
			if len(req.URL.Path) >= r.Size {
				if req.URL.Path[:r.Size] == r.Path {
					req.URL.Path = req.URL.Path[r.Size:]
					r.Handler.ServeHTTP(rw, req)
					return true
				}
			}
		}
		if r.Match(req) {
			r.Handler.ServeHTTP(rw, req)
			vars.Lock()
			delete(vars.v, req)
			vars.Unlock()
			return true
		}
	}
	if req.URL.Path == r.Path {
		r.Handler.ServeHTTP(rw, req)
		return true
	}
	return false
}

func (r *Route) matchRawTokens(ss *[]string) bool {
	if len(*ss) >= r.Token.Size {
		for i, v := range r.Token.raw {
			if (*ss)[v] != r.Token.Tokens[v] {
				if r.Atts&WC != 0 && r.wildPos == i {
					return true
				}
				return false
			}
		}
		return true
	}
	return false
}

// Get set the route method to Get
func (r *Route) Get() *Route {
	r.Method = "GET"
	return r
}

// Post set the route method to Post
func (r *Route) Post() *Route {
	r.Method = "POST"
	return r
}

// Put set the route method to Put
func (r *Route) Put() *Route {
	r.Method = "PUT"
	return r
}

// Delete set the route method to Delete
func (r *Route) Delete() *Route {
	r.Method = "DELETE"
	return r
}

// Head set the route method to Head
func (r *Route) Head() *Route {
	r.Method = "HEAD"
	return r
}

// Patch set the route method to Patch
func (r *Route) Patch() *Route {
	r.Method = "PATCH"
	return r
}

// Options set the route method to Options
func (r *Route) Options() *Route {
	r.Method = "OPTIONS"
	return r
}

func (r *Route) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if r.Method != "" {
		if req.Method == r.Method {
			r.Handler.ServeHTTP(rw, req)
			return
		}
		http.NotFound(rw, req)
		return
	}
	r.Handler.ServeHTTP(rw, req)
}
