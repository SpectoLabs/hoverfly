package jsonpath

type Pos int
type stateFn func(lexer, *intStack) stateFn

const (
	lexError = 0 // must match jsonError and pathError
	lexEOF   = 1
	eof      = -1
	noValue  = -2
)

type Item struct {
	typ int
	pos Pos // The starting position, in bytes, of this Item in the input string.
	val []byte
}

// Used by evaluator
type tokenReader interface {
	next() (*Item, bool)
}

// Used by state functions
type lexer interface {
	tokenReader
	take() int
	takeString() error
	peek() int
	emit(int)
	ignore()
	errorf(string, ...interface{}) stateFn
	reset()
}

type lex struct {
	initialState   stateFn
	currentStateFn stateFn
	item           Item
	hasItem        bool
	stack          intStack
}

func newLex(initial stateFn) lex {
	return lex{
		initialState:   initial,
		currentStateFn: initial,
		item:           Item{},
		stack:          *newIntStack(),
	}
}

func (i *Item) clone() *Item {
	ic := Item{
		typ: i.typ,
		pos: i.pos,
		val: make([]byte, len(i.val)),
	}
	copy(ic.val, i.val)
	return &ic
}

func itemsDescription(items []Item, nameMap map[int]string) []string {
	vals := make([]string, len(items))
	for i, item := range items {
		vals[i] = itemDescription(&item, nameMap)
	}
	return vals
}

func itemDescription(item *Item, nameMap map[int]string) string {
	var found bool
	val, found := nameMap[item.typ]
	if !found {
		return string(item.val)
	}
	return val
}

func typesDescription(types []int, nameMap map[int]string) []string {
	vals := make([]string, len(types))
	for i, val := range types {
		vals[i] = nameMap[val]
	}
	return vals
}
