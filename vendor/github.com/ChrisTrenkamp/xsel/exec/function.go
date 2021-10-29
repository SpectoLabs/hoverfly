package exec

import (
	"fmt"
	"math"
	"strings"

	"github.com/ChrisTrenkamp/xsel/node"
	"github.com/ChrisTrenkamp/xsel/store"
	"golang.org/x/text/language"
)

type Function func(context Context, args ...Result) (Result, error)

type overloadHelper map[int]Function

var errBadArgs = fmt.Errorf("incorrect number of arguments")

func (o overloadHelper) dispatch(context Context, args ...Result) (Result, error) {
	fn := o[len(args)]

	if fn == nil {
		return nil, errBadArgs
	}

	return fn(context, args...)
}

func (o overloadHelper) build() Function {
	return func(context Context, args ...Result) (Result, error) {
		return o.dispatch(context, args...)
	}
}

var localNameDispatch = overloadHelper{
	0: localName0,
	1: localName1,
}

var namespaceUriDispatch = overloadHelper{
	0: namespaceUri0,
	1: namespaceUri1,
}

var nameDispatch = overloadHelper{
	0: name0,
	1: name1,
}

var stringDispatch = overloadHelper{
	0: string0,
	1: string1,
}

var stringLengthDispatch = overloadHelper{
	0: stringLength0,
	1: stringLength1,
}

var normalizeSpaceDispatch = overloadHelper{
	0: normalizeSpace0,
	1: normalizeSpace1,
}

var numberDispatch = overloadHelper{
	0: number0,
	1: number1,
}

var builtinFunctions = map[XmlName]Function{
	{"", "last"}:             last,
	{"", "position"}:         position,
	{"", "count"}:            count,
	{"", "local-name"}:       localNameDispatch.build(),
	{"", "namespace-uri"}:    namespaceUriDispatch.build(),
	{"", "name"}:             nameDispatch.build(),
	{"", "string"}:           stringDispatch.build(),
	{"", "concat"}:           concat,
	{"", "starts-with"}:      startsWith,
	{"", "contains"}:         contains,
	{"", "substring-before"}: substringBefore,
	{"", "substring-after"}:  substringAfter,
	{"", "substring"}:        substring,
	{"", "string-length"}:    stringLengthDispatch.build(),
	{"", "normalize-space"}:  normalizeSpaceDispatch.build(),
	{"", "translate"}:        translate,
	{"", "not"}:              not,
	{"", "true"}:             true0,
	{"", "false"}:            false0,
	{"", "lang"}:             lang,
	{"", "number"}:           numberDispatch.build(),
	{"", "sum"}:              sum,
	{"", "floor"}:            floor,
	{"", "ceiling"}:          ceiling,
	{"", "round"}:            round,
}

func last(context Context, args ...Result) (Result, error) {
	nodeSet, ok := context.Result().(NodeSet)

	if !ok {
		return nil, errQueryNonNodeset
	}

	return Number(len(nodeSet)) + 1, nil
}

func position(context Context, args ...Result) (Result, error) {
	return Number(context.ContextPosition() + 1), nil
}

func count(context Context, args ...Result) (Result, error) {
	if len(args) != 1 {
		return nil, errBadArgs
	}

	nodeSet, ok := args[0].(NodeSet)

	if !ok {
		return nil, errQueryNonNodeset
	}

	return Number(len(nodeSet)), nil
}

type nameType int

const (
	localOnly nameType = iota
	namespaceOnly
	localAndNamespace
)

func localName0(context Context, args ...Result) (Result, error) {
	nodeSet, ok := context.Result().(NodeSet)

	return getName(nodeSet, ok, localOnly)
}

func localName1(context Context, args ...Result) (Result, error) {
	nodeSet, ok := args[0].(NodeSet)

	return getName(nodeSet, ok, localOnly)
}

func namespaceUri0(context Context, args ...Result) (Result, error) {
	nodeSet, ok := context.Result().(NodeSet)

	return getName(nodeSet, ok, namespaceOnly)
}

func namespaceUri1(context Context, args ...Result) (Result, error) {
	nodeSet, ok := args[0].(NodeSet)

	return getName(nodeSet, ok, namespaceOnly)
}

func name0(context Context, args ...Result) (Result, error) {
	nodeSet, ok := context.Result().(NodeSet)

	return getName(nodeSet, ok, localAndNamespace)
}

func name1(context Context, args ...Result) (Result, error) {
	nodeSet, ok := args[0].(NodeSet)

	return getName(nodeSet, ok, localAndNamespace)
}

func getName(nodeSet NodeSet, ok bool, nameType nameType) (Result, error) {
	if !ok {
		return nil, errQueryNonNodeset
	}

	if len(nodeSet) == 0 {
		return String(""), nil
	}

	firstNode := nodeSet[0]

	if n, ok := firstNode.Node().(node.NamedNode); ok {
		if nameType == localOnly || (nameType == localAndNamespace && n.Space() == "") {
			return String(n.Local()), nil
		}

		if nameType == namespaceOnly {
			return String(n.Space()), nil
		}

		return String(fmt.Sprintf("{%s}%s", n.Space(), n.Local())), nil
	}

	return String(""), nil
}

func string0(context Context, args ...Result) (Result, error) {
	return String(context.Result().String()), nil
}

func string1(context Context, args ...Result) (Result, error) {
	return String(args[0].String()), nil
}

func concat(context Context, args ...Result) (Result, error) {
	ret := strings.Builder{}

	for _, i := range args {
		ret.WriteString(i.String())
	}

	return String(ret.String()), nil
}

func startsWith(context Context, args ...Result) (Result, error) {
	if len(args) != 2 {
		return nil, errBadArgs
	}

	str := args[0].String()
	prefix := args[1].String()

	return Bool(strings.HasPrefix(str, prefix)), nil
}

func contains(context Context, args ...Result) (Result, error) {
	if len(args) != 2 {
		return nil, errBadArgs
	}

	str := args[0].String()
	substr := args[1].String()

	return Bool(strings.Contains(str, substr)), nil
}

func substringBefore(context Context, args ...Result) (Result, error) {
	if len(args) != 2 {
		return nil, errBadArgs
	}

	str := args[0].String()
	substr := args[1].String()

	ind := strings.Index(str, substr)

	if ind < 0 {
		return String(""), nil
	}

	return String(str[:ind]), nil
}

func substringAfter(context Context, args ...Result) (Result, error) {
	if len(args) != 2 {
		return nil, errBadArgs
	}

	str := args[0].String()
	substr := args[1].String()

	ind := strings.Index(str, substr)

	if ind < 0 {
		return String(""), nil
	}

	return String(str[ind+len(substr):]), nil
}

func substring(context Context, args ...Result) (Result, error) {
	if len(args) != 2 && len(args) != 3 {
		return nil, errBadArgs
	}

	str := args[0].String()
	begin := getRound(args[1].Number())

	if float64(begin-1) >= float64(len(str)) || math.IsNaN(float64(begin)) {
		return String(""), nil
	}

	if len(args) == 2 {
		if begin <= 1 {
			begin = 1
		}

		return String(str[int(begin)-1:]), nil
	}

	end := getRound(args[2].Number())

	if end <= 0 || math.IsNaN(float64(end)) || (math.IsInf(float64(begin), 0) && math.IsInf(float64(end), 0)) {
		return String(""), nil
	}

	if begin <= 1 {
		end = begin + end - 1
		begin = 1
	}

	if float64(begin+end-1) >= float64(len(str)) {
		end = float64(len(str)) - begin + 1
	}

	return String(str[int(begin)-1 : int(begin+end)-1]), nil
}

func stringLength0(context Context, args ...Result) (Result, error) {
	return Number(len(context.Result().String())), nil
}

func stringLength1(context Context, args ...Result) (Result, error) {
	return Number(len(args[0].String())), nil
}

func normalizeSpace0(context Context, args ...Result) (Result, error) {
	return String(strings.TrimSpace(context.Result().String())), nil
}

func normalizeSpace1(context Context, args ...Result) (Result, error) {
	return String(strings.TrimSpace(args[0].String())), nil
}

func translate(context Context, args ...Result) (Result, error) {
	if len(args) != 3 {
		return nil, errBadArgs
	}

	src := args[0].String()
	old := args[1].String()
	new := args[2].String()

	for i := range old {
		r := ""

		if i < len(new) {
			r = string(new[i])
		}

		src = strings.Replace(src, string(old[i]), r, -1)
	}

	return String(src), nil
}

func not(context Context, args ...Result) (Result, error) {
	if len(args) != 1 {
		return nil, errBadArgs
	}

	return Bool(!args[0].Bool()), nil
}

func true0(context Context, args ...Result) (Result, error) {
	if len(args) != 0 {
		return nil, errBadArgs
	}

	return Bool(true), nil
}

func false0(context Context, args ...Result) (Result, error) {
	if len(args) != 0 {
		return nil, errBadArgs
	}

	return Bool(false), nil
}

func lang(context Context, args ...Result) (Result, error) {
	if len(args) != 1 {
		return nil, errBadArgs
	}

	nodeSet, ok := context.Result().(NodeSet)

	if !ok {
		return nil, errQueryNonNodeset
	}

	lStr := args[0].String()

	var n store.Cursor

	for _, i := range nodeSet {
		if _, ok := i.Node().(node.Element); ok {
			n = i
		} else {
			n = i.Parent()
		}

		for n.Pos() != 0 {
			if attr, ok := store.GetAttribute(n, "http://www.w3.org/XML/1998/namespace", "lang"); ok {
				return checkLang(lStr, attr.AttributeValue()), nil
			}

			n = n.Parent()
		}
	}

	return Bool(false), nil
}

func checkLang(srcStr, targStr string) Bool {
	srcLang := language.Make(srcStr)
	srcRegion, srcRegionConf := srcLang.Region()

	targLang := language.Make(targStr)
	targRegion, targRegionConf := targLang.Region()

	if srcRegionConf == language.Exact && targRegionConf != language.Exact {
		return Bool(false)
	}

	if srcRegion != targRegion && srcRegionConf == language.Exact && targRegionConf == language.Exact {
		return Bool(false)
	}

	_, _, conf := language.NewMatcher([]language.Tag{srcLang}).Match(targLang)
	return Bool(conf >= language.High)
}

func number0(context Context, args ...Result) (Result, error) {
	return Number(context.Result().Number()), nil
}

func number1(context Context, args ...Result) (Result, error) {
	return Number(args[0].Number()), nil
}

func sum(context Context, args ...Result) (Result, error) {
	if len(args) != 1 {
		return nil, errBadArgs
	}

	nodeSet, ok := args[0].(NodeSet)

	if !ok {
		return nil, errQueryNonNodeset
	}

	sum := 0

	for _, i := range nodeSet {
		sum += int(NodeSet{i}.Number())
	}

	return Number(sum), nil
}

func floor(context Context, args ...Result) (Result, error) {
	if len(args) != 1 {
		return nil, errBadArgs
	}

	return Number(math.Floor(float64(args[0].Number()))), nil
}

func ceiling(context Context, args ...Result) (Result, error) {
	if len(args) != 1 {
		return nil, errQueryNonNodeset
	}

	return Number(math.Ceil(float64(args[0].Number()))), nil
}

func round(context Context, args ...Result) (Result, error) {
	if len(args) != 1 {
		return nil, errBadArgs
	}

	return Number(getRound(float64(args[0].Number()))), nil
}

func getRound(n float64) float64 {
	if math.IsNaN(float64(n)) || math.IsInf(float64(n), 0) {
		return n
	}

	if n < -0.5 {
		n = float64(int(n - 0.5))
	} else if n > 0.5 {
		n = float64(int(n + 0.5))
	} else {
		n = 0
	}

	return n
}
