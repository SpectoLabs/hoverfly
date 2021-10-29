package exec

import (
	"fmt"
	"strings"

	"github.com/ChrisTrenkamp/xsel/grammar"
	"github.com/ChrisTrenkamp/xsel/grammar/parser/bsr"
	"github.com/ChrisTrenkamp/xsel/grammar/parser/symbols"
	"github.com/ChrisTrenkamp/xsel/node"
)

var errQueryNonNodeset = fmt.Errorf("cannot query nodes on non-NodeSet's")

func init() {
	contextFunctions[symbols.NT_AbsoluteLocationPathOnly] = execAbsoluteLocationPathOnly
	contextFunctions[symbols.NT_RelativeLocationPathWithStep] = leftRightDependentResult
	contextFunctions[symbols.NT_Step] = execStep
	contextFunctions[symbols.NT_NodeTestAndPredicate] = leftRightDependentResult
	contextFunctions[symbols.NT_Predicate] = execPredicate
	contextFunctions[symbols.NT_NodeTestNodeTypeNoArgTest] = execNodeTestNodeTypeNoArgTest
	contextFunctions[symbols.NT_NodeTestProcInstTargetTest] = execNodeTestProcInstTargetTest
	contextFunctions[symbols.NT_NameTestAnyElement] = execNameTestAnyElement
	contextFunctions[symbols.NT_NameTestNamespaceAnyLocal] = execNameTestNamespaceAnyLocal
	contextFunctions[symbols.NT_NameTestNamespaceAnyLocalReservedNameConflict] = execNameTestNamespaceAnyLocalReservedNameConflict
	contextFunctions[symbols.NT_NameTestLocalAnyNamespace] = execNameTestLocalAnyNamespace
	contextFunctions[symbols.NT_NameTestLocalAnyNamespaceReservedNameConflict] = execNameTestLocalAnyNamespaceReservedNameConflict
	contextFunctions[symbols.NT_NameTestQNameNamespaceWithLocal] = execNameTestQNameNamespaceWithLocal
	contextFunctions[symbols.NT_NameTestQNameNamespaceWithLocalReservedNameConflictNamespace] = execNameTestQNameNamespaceWithLocalReservedNameConflictNamespace
	contextFunctions[symbols.NT_NameTestQNameNamespaceWithLocalReservedNameConflictLocal] = execNameTestQNameNamespaceWithLocalReservedNameConflictLocal
	contextFunctions[symbols.NT_NameTestQNameNamespaceWithLocalReservedNameConflictBoth] = execNameTestQNameNamespaceWithLocalReservedNameConflictBoth
	contextFunctions[symbols.NT_NameTestQNameLocalOnly] = execNameTestQNameLocalOnly
	contextFunctions[symbols.NT_NameTestQNameLocalOnlyReservedNameConflict] = execNameTestQNameLocalOnly
	contextFunctions[symbols.NT_StepWithAxisAndNodeTest] = leftRightDependentResult
	contextFunctions[symbols.NT_AxisName] = execAxisName
	contextFunctions[symbols.NT_AbbreviatedStepParent] = execAbbreviatedStepParent
	contextFunctions[symbols.NT_AbbreviatedAxisSpecifier] = execAbbreviatedAxisSpecifier
	contextFunctions[symbols.NT_AbbreviatedAbsoluteLocationPath] = execAbbreviatedAbsoluteLocationPath
	contextFunctions[symbols.NT_AbbreviatedRelativeLocationPath] = execAbbreviatedRelativeLocationPath
}

func execAbsoluteLocationPathOnly(context *exprContext, expr *grammar.Grammar) error {
	context.result = NodeSet{context.root}
	return nil
}

func execStep(context *exprContext, expr *grammar.Grammar) error {
	var nextBsr *bsr.BSR

	for _, cn := range expr.BSR.GetAllNTChildren() {
		for _, c := range cn {
			nextBsr = &c
			break
		}
	}

	switch nextBsr.Label.Slot().NT {
	case symbols.NT_NodeTest,
		symbols.NT_NodeTestAndPredicate,
		symbols.NT_NodeTestNodeTypeNoArgTest,
		symbols.NT_NodeTestProcInstTargetTest,
		symbols.NT_NameTestAnyElement,
		symbols.NT_NameTestNamespaceAnyLocal,
		symbols.NT_NameTestNamespaceAnyLocalReservedNameConflict,
		symbols.NT_NameTestLocalAnyNamespace,
		symbols.NT_NameTestLocalAnyNamespaceReservedNameConflict,
		symbols.NT_NameTestQNameNamespaceWithLocal,
		symbols.NT_NameTestQNameNamespaceWithLocalReservedNameConflictNamespace,
		symbols.NT_NameTestQNameNamespaceWithLocalReservedNameConflictLocal,
		symbols.NT_NameTestQNameNamespaceWithLocalReservedNameConflictBoth,
		symbols.NT_NameTestQNameLocalOnly,
		symbols.NT_NameTestQNameLocalOnlyReservedNameConflict:
		nodeSet, ok := context.result.(NodeSet)

		if !ok {
			return errQueryNonNodeset
		}

		context.result = selectChild(nodeSet)
	}

	return execContext(context, expr.Next(nextBsr))
}

func execPredicate(context *exprContext, expr *grammar.Grammar) error {
	nodeSet, ok := context.result.(NodeSet)

	if !ok {
		return errQueryNonNodeset
	}

	nextResult := make(NodeSet, 0)

	for i := range nodeSet {
		nextContext := context.copy()
		nextContext.result = NodeSet{nodeSet[i]}
		nextContext.contextPosition = i
		left, err := leftOnlyIndependentResult(&nextContext, expr)

		if err != nil {
			return err
		}

		if n, ok := left.(Number); ok {
			if (i + 1) == int(n) {
				nextResult = append(nextResult, nodeSet[i])
			}
		} else if b, ok := left.(Bool); ok {
			if bool(b) {
				nextResult = append(nextResult, nodeSet[i])
			}
		} else if left.Bool() {
			nextResult = append(nextResult, nodeSet[i])
		}
	}

	context.result = nextResult
	return nil
}

func execNodeTestNodeTypeNoArgTest(context *exprContext, expr *grammar.Grammar) error {
	nodeSet, ok := context.result.(NodeSet)

	if !ok {
		return nil
	}

	nodeType := expr.GetString()
	parenIndex := strings.LastIndex(nodeType, "(")
	nodeType = nodeType[:parenIndex]
	nodeType = strings.TrimSpace(nodeType)

	result := make(NodeSet, 0)

	switch nodeType {
	case "comment":
		for _, i := range nodeSet {
			if _, ok := i.Node().(node.Comment); ok {
				result = append(result, i)
			}
		}
	case "text":
		for _, i := range nodeSet {
			if _, ok := i.Node().(node.CharData); ok {
				result = append(result, i)
			}
		}
	case "processing-instruction":
		for _, i := range nodeSet {
			if _, ok := i.Node().(node.ProcInst); ok {
				result = append(result, i)
			}
		}
	case "node":
		return nil
	}

	context.result = result
	return nil
}

func execNodeTestProcInstTargetTest(context *exprContext, expr *grammar.Grammar) error {
	nodeSet, ok := context.result.(NodeSet)

	if !ok {
		return nil
	}

	literal, err := leftOnlyIndependentResult(context, expr)

	if err != nil {
		return err
	}

	literalString := literal.String()
	result := make(NodeSet, 0)

	for _, i := range nodeSet {
		if pi, ok := i.Node().(node.ProcInst); ok && pi.Target() == literalString {
			result = append(result, i)
		}
	}

	context.result = result
	return nil
}

func execNameTestAnyElement(context *exprContext, expr *grammar.Grammar) error {
	nodeSet, ok := context.result.(NodeSet)

	if !ok {
		return nil
	}

	result := make(NodeSet, 0)

	for _, i := range nodeSet {
		if _, ok := i.Node().(node.NamedNode); ok {
			result = append(result, i)
		}

		if _, ok := i.Node().(node.Namespace); ok {
			result = append(result, i)
		}
	}

	context.result = result
	return nil
}

func execNameTestNamespaceAnyLocal(context *exprContext, expr *grammar.Grammar) error {
	namespaceLookup := expr.BSR.GetTChildI(0).LiteralString()

	return nameTestNamespaceAnyLocal(namespaceLookup, context, expr)
}

func execNameTestNamespaceAnyLocalReservedNameConflict(context *exprContext, expr *grammar.Grammar) error {
	children := make([]*bsr.BSR, 0, 1)

	for _, cn := range expr.BSR.GetAllNTChildren() {
		for _, c := range cn {
			children = append(children, &c)
		}
	}

	namespaceLookup := expr.GetStringExtents(children[0].LeftExtent(), children[0].RightExtent())
	return nameTestNamespaceAnyLocal(namespaceLookup, context, expr)
}

func nameTestNamespaceAnyLocal(namespaceLookup string, context *exprContext, expr *grammar.Grammar) error {
	namespaceValue, ok := context.NamespaceDecls[namespaceLookup]

	if !ok {
		return fmt.Errorf("unknown namespace binding '%s'", namespaceLookup)
	}

	nodeSet, ok := context.result.(NodeSet)

	if !ok {
		return nil
	}

	result := make(NodeSet, 0)

	for _, i := range nodeSet {
		if node, ok := i.Node().(node.NamedNode); ok {
			if node.Space() == namespaceValue {
				result = append(result, i)
			}
		}
	}

	context.result = result
	return nil
}

func execNameTestLocalAnyNamespace(context *exprContext, expr *grammar.Grammar) error {
	localValue := expr.BSR.GetTChildI(2).LiteralString()

	return nameTestLocalAnyNamespace(localValue, context, expr)
}

func execNameTestLocalAnyNamespaceReservedNameConflict(context *exprContext, expr *grammar.Grammar) error {
	children := make([]*bsr.BSR, 0, 1)

	for _, cn := range expr.BSR.GetAllNTChildren() {
		for _, c := range cn {
			children = append(children, &c)
			break
		}
	}

	localValue := expr.GetStringExtents(children[0].LeftExtent(), children[0].RightExtent())
	return nameTestLocalAnyNamespace(localValue, context, expr)
}

func nameTestLocalAnyNamespace(localValue string, context *exprContext, expr *grammar.Grammar) error {
	nodeSet, ok := context.result.(NodeSet)

	if !ok {
		return nil
	}

	result := make(NodeSet, 0)

	for _, i := range nodeSet {
		if node, ok := i.Node().(node.NamedNode); ok {
			if node.Local() == localValue {
				result = append(result, i)
			}
		}
	}

	context.result = result
	return nil
}

func execNameTestQNameNamespaceWithLocal(context *exprContext, expr *grammar.Grammar) error {
	namespaceLookup := expr.BSR.GetTChildI(0).LiteralString()
	local := expr.BSR.GetTChildI(2).LiteralString()

	return nameTestQNameNamespaceWithLocal(namespaceLookup, local, context, expr)
}

func execNameTestQNameNamespaceWithLocalReservedNameConflictNamespace(context *exprContext, expr *grammar.Grammar) error {
	children := make([]*bsr.BSR, 0, 1)

	for _, cn := range expr.BSR.GetAllNTChildren() {
		for _, c := range cn {
			children = append(children, &c)
			break
		}
	}

	namespaceLookup := expr.GetStringExtents(children[0].LeftExtent(), children[0].RightExtent())
	local := expr.BSR.GetTChildI(2).LiteralString()
	return nameTestQNameNamespaceWithLocal(namespaceLookup, local, context, expr)
}

func execNameTestQNameNamespaceWithLocalReservedNameConflictLocal(context *exprContext, expr *grammar.Grammar) error {
	children := make([]*bsr.BSR, 0, 1)

	for _, cn := range expr.BSR.GetAllNTChildren() {
		for _, c := range cn {
			children = append(children, &c)
			break
		}
	}

	namespaceLookup := expr.BSR.GetTChildI(0).LiteralString()
	local := expr.GetStringExtents(children[0].LeftExtent(), children[0].RightExtent())
	return nameTestQNameNamespaceWithLocal(namespaceLookup, local, context, expr)
}

func execNameTestQNameNamespaceWithLocalReservedNameConflictBoth(context *exprContext, expr *grammar.Grammar) error {
	children := make([]*bsr.BSR, 0, 2)

	for _, cn := range expr.BSR.GetAllNTChildren() {
		for _, c := range cn {
			children = append(children, &c)
			break
		}
	}

	namespaceLookup := expr.GetStringExtents(children[0].LeftExtent(), children[0].RightExtent())
	local := expr.GetStringExtents(children[1].LeftExtent(), children[1].RightExtent())
	return nameTestQNameNamespaceWithLocal(namespaceLookup, local, context, expr)
}

func nameTestQNameNamespaceWithLocal(namespaceLookup, local string, context *exprContext, expr *grammar.Grammar) error {
	namespaceValue, ok := context.NamespaceDecls[namespaceLookup]

	if !ok {
		return fmt.Errorf("unknown namespace binding '%s'", namespaceLookup)
	}

	nodeSet, ok := context.result.(NodeSet)

	if !ok {
		return nil
	}

	result := make(NodeSet, 0)

	for _, i := range nodeSet {
		if node, ok := i.Node().(node.NamedNode); ok {
			if node.Local() == local && node.Space() == namespaceValue {
				result = append(result, i)
			}
		}
	}

	context.result = result
	return nil
}

func execNameTestQNameLocalOnly(context *exprContext, expr *grammar.Grammar) error {
	nodeSet, ok := context.result.(NodeSet)

	if !ok {
		return errQueryNonNodeset
	}

	nextResult := make(NodeSet, 0)
	queryName := expr.GetString()

	for _, child := range nodeSet {
		if elem, ok := child.Node().(node.NamedNode); ok {
			if elem.Space() == "" && elem.Local() == queryName {
				nextResult = append(nextResult, child)
			}
		}

		if ns, ok := child.Node().(node.Namespace); ok {
			namespaceValue := context.NamespaceDecls[queryName]

			if ns.NamespaceValue() == namespaceValue {
				nextResult = append(nextResult, child)
			}
		}
	}

	context.result = nextResult

	return execChildren(context, expr)
}

func execAxisName(context *exprContext, expr *grammar.Grammar) error {
	nodeSet, ok := context.result.(NodeSet)

	if !ok {
		return errQueryNonNodeset
	}

	axis := expr.GetString()
	var result Result

	switch axis {
	case "child":
		result = selectChild(nodeSet)
	case "attribute":
		result = selectAttributes(nodeSet)
	case "ancestor":
		result = selectAncestor(nodeSet)
	case "ancestor-or-self":
		result = selectAncestorOrSelf(nodeSet)
	case "descendant":
		result = selectDescendant(nodeSet)
	case "descendant-or-self":
		result = selectDescendantOrSelf(nodeSet)
	case "following":
		result = selectFollowing(nodeSet)
	case "following-sibling":
		result = selectFollowingSibling(nodeSet)
	case "namespace":
		result = selectNamespace(nodeSet)
	case "parent":
		result = selectParent(nodeSet)
	case "preceding":
		result = selectPreceding(nodeSet)
	case "preceding-sibling":
		result = selectPrecedingSibling(nodeSet)
	default: // self
		return nil
	}

	context.result = result

	return nil
}

func execAbbreviatedStepParent(context *exprContext, expr *grammar.Grammar) error {
	nodeSet, ok := context.result.(NodeSet)

	if !ok {
		return errQueryNonNodeset
	}

	context.result = selectParent(nodeSet)
	return nil
}

func execAbbreviatedAxisSpecifier(context *exprContext, expr *grammar.Grammar) error {
	nodeSet, ok := context.result.(NodeSet)

	if !ok {
		return errQueryNonNodeset
	}

	context.result = selectAttributes(nodeSet)
	return nil
}

func execAbbreviatedAbsoluteLocationPath(context *exprContext, expr *grammar.Grammar) error {
	context.result = selectDescendantOrSelf(NodeSet{context.root})

	for _, cn := range expr.BSR.GetAllNTChildren() {
		for _, c := range cn {
			return execContext(context, expr.Next(&c))
		}
	}

	return nil
}

func execAbbreviatedRelativeLocationPath(context *exprContext, expr *grammar.Grammar) error {
	children := make([]*bsr.BSR, 0, 2)

	for _, cn := range expr.BSR.GetAllNTChildren() {
		for _, c := range cn {
			children = append(children, &c)
		}
	}

	if err := execContext(context, expr.Next(children[0])); err != nil {
		return err
	}

	nodeSet, ok := context.result.(NodeSet)

	if !ok {
		return errQueryNonNodeset
	}

	context.result = selectDescendantOrSelf(nodeSet)
	return execContext(context, expr.Next(children[1]))
}
