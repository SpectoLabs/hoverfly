package jsonpath

import (
	"bytes"
	"fmt"
)

type queryStateFn func(*query, *Eval, *Item) queryStateFn

type query struct {
	Path
	state       queryStateFn
	start       int
	pos         int
	firstType   int // first json token type in buffer
	buffer      bytes.Buffer
	resultQueue *Results
	valLoc      stack // capture the current location stack at capture
	errors      []error
	buckets     stack // stack of exprBucket
}

type exprBucket struct {
	operatorLoc int
	expression  []Item
	queries     []*query
	results     *Results
}

type evalStateFn func(*Eval, *Item) evalStateFn

type Eval struct {
	tr         tokenReader
	levelStack intStack
	location   stack
	queries    map[string]*query
	state      evalStateFn
	prevIndex  int
	nextKey    []byte
	copyValues bool

	resultQueue *Results
	Error       error
}

func newEvaluation(tr tokenReader, paths ...*Path) *Eval {
	e := &Eval{
		tr:          tr,
		location:    *newStack(),
		levelStack:  *newIntStack(),
		state:       evalRoot,
		queries:     make(map[string]*query, 0),
		prevIndex:   -1,
		nextKey:     nil,
		copyValues:  true, // depends on which lexer is used
		resultQueue: newResults(),
	}

	for _, p := range paths {
		e.queries[p.stringValue] = newQuery(p)
	}
	// Determine whether to copy emitted item values ([]byte) from lexer
	switch tr.(type) {
	case *readerLexer:
		e.copyValues = true
	default:
		e.copyValues = false
	}

	return e
}

func newQuery(p *Path) *query {
	return &query{
		Path:        *p,
		state:       pathMatchOp,
		start:       -1,
		pos:         -1,
		buffer:      *bytes.NewBuffer(make([]byte, 0, 50)),
		valLoc:      *newStack(),
		errors:      make([]error, 0),
		resultQueue: newResults(),
		buckets:     *newStack(),
	}
}

func (e *Eval) Iterate() (*Results, bool) {
	e.resultQueue.clear()

	t, ok := e.tr.next()
	if !ok || e.state == nil {
		return nil, false
	}

	// run evaluator function
	e.state = e.state(e, t)

	anyRunning := false
	// run path function for each path
	for str, query := range e.queries {
		anyRunning = true
		query.state = query.state(query, e, t)
		if query.state == nil {
			delete(e.queries, str)
		}

		if query.resultQueue.len() > 0 {
			e.resultQueue.push(query.resultQueue.Pop())
		}

		for _, b := range query.buckets.values {
			bucket := b.(exprBucket)
			for _, dq := range bucket.queries {
				dq.state = dq.state(dq, e, t)

				if query.resultQueue.len() > 0 {
					e.resultQueue.push(query.resultQueue.Pop())
				}
			}
		}
	}

	if !anyRunning {
		return nil, false
	}

	if e.Error != nil {
		return nil, false
	}

	return e.resultQueue, true
}

func (e *Eval) Next() (*Result, bool) {
	if e.resultQueue.len() > 0 {
		return e.resultQueue.Pop(), true
	}

	for {
		if _, ok := e.Iterate(); ok {
			if e.resultQueue.len() > 0 {
				return e.resultQueue.Pop(), true
			}
		} else {
			break
		}

	}
	return nil, false
}

func (q *query) loc() int {
	return abs(q.pos-q.start) + q.start
}

func (q *query) trySpillOver() {
	if b, ok := q.buckets.peek(); ok {
		bucket := b.(exprBucket)
		if q.loc() < bucket.operatorLoc {
			q.buckets.pop()

			exprRes, err := bucket.evaluate()
			if err != nil {
				q.errors = append(q.errors, err)
			}
			if exprRes {
				next, ok := q.buckets.peek()
				var spillover *Results
				if !ok {
					// fmt.Println("Spilling over into end queue")
					spillover = q.resultQueue
				} else {
					// fmt.Println("Spilling over into lower bucket")
					nextBucket := next.(exprBucket)
					spillover = nextBucket.results
				}
				for {
					v := bucket.results.Pop()
					if v != nil {
						spillover.push(v)
					} else {
						break
					}
				}
			}
		}
	}
}

func pathMatchOp(q *query, e *Eval, i *Item) queryStateFn {
	curLocation := e.location.len() - 1

	if q.loc() > curLocation {
		q.pos -= 1
		q.trySpillOver()
	} else if q.loc() <= curLocation {
		if q.loc() == curLocation-1 {
			if len(q.operators)+q.start >= curLocation {
				current, _ := e.location.peek()
				nextOp := q.operators[abs(q.loc()-q.start)]
				if itemMatchOperator(current, i, nextOp) {
					q.pos += 1

					if nextOp.whereClauseBytes != nil && len(nextOp.whereClause) > 0 {
						bucket := exprBucket{
							operatorLoc: q.loc(),
							expression:  nextOp.whereClause,
							queries:     make([]*query, len(nextOp.dependentPaths)),
							results:     newResults(),
						}

						for i, p := range nextOp.dependentPaths {
							bucket.queries[i] = newQuery(p)
							bucket.queries[i].pos = q.loc()
							bucket.queries[i].start = q.loc()
							bucket.queries[i].captureEndValue = true
						}
						q.buckets.push(bucket)
					}
				}

			}
		}
	}

	if q.loc() == len(q.operators)+q.start && q.loc() <= curLocation {
		if q.captureEndValue {
			q.firstType = i.typ
			q.buffer.Write(i.val)
		}
		q.valLoc = *e.location.clone()
		return pathEndValue
	}

	if q.loc() < -1 {
		return nil
	} else {
		return pathMatchOp
	}
}

func pathEndValue(q *query, e *Eval, i *Item) queryStateFn {
	if e.location.len()-1 >= q.loc() {
		if q.captureEndValue {
			q.buffer.Write(i.val)
		}
	} else {
		r := &Result{Keys: q.valLoc.toArray()}
		if q.buffer.Len() > 0 {
			val := make([]byte, q.buffer.Len())
			copy(val, q.buffer.Bytes())
			r.Value = val

			switch q.firstType {
			case jsonBraceLeft:
				r.Type = JsonObject
			case jsonString:
				r.Type = JsonString
			case jsonBracketLeft:
				r.Type = JsonArray
			case jsonNull:
				r.Type = JsonNull
			case jsonBool:
				r.Type = JsonBool
			case jsonNumber:
				r.Type = JsonNumber
			default:
				r.Type = -1
			}
		}

		if q.buckets.len() == 0 {
			q.resultQueue.push(r)
		} else {
			b, _ := q.buckets.peek()
			b.(exprBucket).results.push(r)
		}

		q.valLoc = *newStack()
		q.buffer.Truncate(0)
		q.pos -= 1
		return pathMatchOp
	}
	return pathEndValue
}

func (b *exprBucket) evaluate() (bool, error) {
	values := make(map[string]Item)
	for _, q := range b.queries {
		result := q.resultQueue.Pop()
		if result != nil {
			t, err := getJsonTokenType(result.Value)
			if err != nil {
				return false, err
			}
			i := Item{
				typ: t,
				val: result.Value,
			}
			values[q.Path.stringValue] = i
		}
	}

	res, err := evaluatePostFix(b.expression, values)
	if err != nil {
		return false, err
	}
	res_bool, ok := res.(bool)
	if !ok {
		return false, fmt.Errorf(exprErrorFinalValueNotBool, res)
	}
	return res_bool, nil
}

func itemMatchOperator(loc interface{}, i *Item, op *operator) bool {
	topBytes, isKey := loc.([]byte)
	topInt, isIndex := loc.(int)
	if isKey {
		switch op.typ {
		case opTypeNameWild:
			return true
		case opTypeName, opTypeNameList:
			_, found := op.keyStrings[string(topBytes)]
			return found
		}
	} else if isIndex {
		switch op.typ {
		case opTypeIndexWild:
			return true
		case opTypeIndex, opTypeIndexRange:
			return topInt >= op.indexStart && (!op.hasIndexEnd || topInt <= op.indexEnd)
		}
	}
	return false
}
