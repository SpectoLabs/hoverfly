// Package bsr is generated by gogll. Do not edit.

/*
Package bsr implements a Binary Subtree Representation set as defined in

	Scott et al
	Derivation representation using binary subtree sets,
	Science of Computer Programming 175 (2019)
*/
package bsr

import (
	"bytes"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/ChrisTrenkamp/xsel/grammar/lexer"
	"github.com/ChrisTrenkamp/xsel/grammar/parser/slot"
	"github.com/ChrisTrenkamp/xsel/grammar/parser/symbols"
	"github.com/ChrisTrenkamp/xsel/grammar/token"
)

type bsr interface {
	LeftExtent() int
	RightExtent() int
	Pivot() int
}

/*
Set contains the set of Binary Subtree Representations (BSR).
*/
type Set struct {
	slotEntries   map[BSR]bool
	ntSlotEntries map[ntSlot][]BSR
	stringEntries map[stringBSR]bool
	rightExtent   int
	lex           *lexer.Lexer

	startSym symbols.NT
}

type ntSlot struct {
	nt          symbols.NT
	leftExtent  int
	rightExtent int
}

// BSR is the binary subtree representation of a parsed nonterminal
type BSR struct {
	Label       slot.Label
	leftExtent  int
	pivot       int
	rightExtent int
	set         *Set
}

type stringBSR struct {
	Label       slot.Label
	leftExtent  int
	pivot       int
	rightExtent int
	set         *Set
}

// New returns a new initialised BSR Set
func New(startSymbol symbols.NT, l *lexer.Lexer) *Set {
	return &Set{
		slotEntries:   make(map[BSR]bool),
		ntSlotEntries: make(map[ntSlot][]BSR),
		stringEntries: make(map[stringBSR]bool),
		rightExtent:   0,
		lex:           l,
		startSym:      startSymbol,
	}
}

/*
Add a bsr to the set. (i,j) is the extent. k is the pivot.
*/
func (s *Set) Add(l slot.Label, i, k, j int) {
	// fmt.Printf("bsr.Add(%s,%d,%d,%d)\n", l,i,k,j)
	if l.EoR() {
		s.insert(BSR{l, i, k, j, s})
	} else {
		if l.Pos() > 1 {
			s.insert(stringBSR{l, i, k, j, s})
		}
	}
}

// AddEmpty adds a grammar slot: X : ϵ•
func (s *Set) AddEmpty(l slot.Label, i int) {
	s.insert(BSR{l, i, i, i, s})
}

/*
Contain returns true iff the BSR Set contains the NT symbol with left and
right extent.
*/
func (s *Set) Contain(nt symbols.NT, left, right int) bool {
	// fmt.Printf("bsr.Contain(%s,%d,%d)\n",nt,left,right)
	for e := range s.slotEntries {
		// fmt.Printf("  (%s,%d,%d)\n",e.Label.Head(),e.leftExtent,e.rightExtent)
		if e.Label.Head() == nt && e.leftExtent == left && e.rightExtent == right {
			// fmt.Println("  true")
			return true
		}
	}
	// fmt.Println("  false")
	return false
}

/*
Dump prints the NT symbols of the parse forest.
*/
func (s *Set) Dump() {
	for _, root := range s.GetRoots() {
		s.dump(root, 0)
	}
}

func (s *Set) dump(b BSR, level int) {
	fmt.Print(indent(level, " "))
	fmt.Println(b)
	for _, cn := range b.GetAllNTChildren() {
		for _, c := range cn {
			s.dump(c, level+1)
		}
	}
}

func indent(n int, c string) string {
	buf := new(bytes.Buffer)
	for i := 0; i < 4*n; i++ {
		fmt.Fprint(buf, c)
	}
	return buf.String()
}

// GetAll returns all BSR grammar slot entries
func (s *Set) GetAll() (bsrs []BSR) {
	for b := range s.slotEntries {
		bsrs = append(bsrs, b)
	}
	return
}

// GetRightExtent returns the right extent of the BSR set
func (s *Set) GetRightExtent() int {
	return s.rightExtent
}

// GetRoot returns the root of the parse tree of an unambiguous parse.
// GetRoot fails if the parse was ambiguous. Use GetRoots() for ambiguous parses.
func (s *Set) GetRoot() BSR {
	rts := s.GetRoots()
	if len(rts) != 1 {
		failf("%d parse trees exist for start symbol %s", len(rts), s.startSym)
	}
	return rts[0]
}

// GetRoots returns all the roots of parse trees of the start symbol of the grammar.
func (s *Set) GetRoots() (roots []BSR) {
	for b := range s.slotEntries {
		if b.Label.Head() == s.startSym && b.leftExtent == 0 && s.rightExtent == b.rightExtent {
			roots = append(roots, b)
		}
	}
	return
}

func (s *Set) getString(l slot.Label, leftExtent, rightExtent int) stringBSR {
	for str := range s.stringEntries {
		if str.Label == l && str.leftExtent == leftExtent && str.rightExtent == rightExtent {
			return str
		}
	}
	fmt.Printf("Error: no string %s left extent=%d right extent=%d pos=%d\n",
		strings.Join(l.Symbols()[:l.Pos()].Strings(), " "), leftExtent, rightExtent, l.Pos())
	panic("must not happen")
}

func (s *Set) insert(bsr bsr) {
	if bsr.RightExtent() > s.rightExtent {
		s.rightExtent = bsr.RightExtent()
	}
	switch b := bsr.(type) {
	case BSR:
		s.slotEntries[b] = true
		nt := ntSlot{b.Label.Head(), b.leftExtent, b.rightExtent}
		s.ntSlotEntries[nt] = append(s.ntSlotEntries[nt], b)
	case stringBSR:
		s.stringEntries[b] = true
	default:
		panic(fmt.Sprintf("Invalid type %T", bsr))
	}
}

// Alternate returns the index of the grammar rule alternate.
func (b BSR) Alternate() int {
	return b.Label.Alternate()
}

// GetAllNTChildren returns all the NT Children of b. If an NT child of b has
// ambiguous parses then all parses of that child are returned.
func (b BSR) GetAllNTChildren() [][]BSR {
	children := [][]BSR{}
	for i, s := range b.Label.Symbols() {
		if s.IsNonTerminal() {
			sChildren := b.GetNTChildrenI(i)
			children = append(children, sChildren)
		}
	}
	return children
}

// GetNTChild returns the BSR of occurrence i of nt in s.
// GetNTChild fails if s has ambiguous subtrees of occurrence i of nt.
func (b BSR) GetNTChild(nt symbols.NT, i int) BSR {
	bsrs := b.GetNTChildren(nt, i)
	if len(bsrs) != 1 {
		ambiguousSlots := []string{}
		for _, c := range bsrs {
			ambiguousSlots = append(ambiguousSlots, c.String())
		}
		b.set.fail(b, "%s is ambiguous in %s\n  %s", nt, b, strings.Join(ambiguousSlots, "\n  "))
	}
	return bsrs[0]
}

// GetNTChildI returns the BSR of NT symbol[i] in the BSR set.
// GetNTChildI fails if the BSR set has ambiguous subtrees of NT i.
func (b BSR) GetNTChildI(i int) BSR {
	bsrs := b.GetNTChildrenI(i)
	if len(bsrs) != 1 {
		b.set.fail(b, "NT %d is ambiguous in %s", i, b)
	}
	return bsrs[0]
}

// GetNTChildren returns all the BSRs of occurrence i of nt in s
func (b BSR) GetNTChildren(nt symbols.NT, i int) []BSR {
	// fmt.Printf("GetNTChild(%s,%d) %s\n", nt, i, b)
	positions := []int{}
	for j, s := range b.Label.Symbols() {
		if s == nt {
			positions = append(positions, j)
		}
	}
	if len(positions) == 0 {
		b.set.fail(b, "Error: %s has no NT %s", b, nt)
	}
	return b.GetNTChildrenI(positions[i])
}

// GetNTChildrenI returns all the BSRs of NT symbol[i] in s
func (b BSR) GetNTChildrenI(i int) []BSR {
	// fmt.Printf("bsr.GetNTChildI(%d) %s\n", i, b)
	if i >= len(b.Label.Symbols()) {
		b.set.fail(b, "Error: cannot get NT child %d of %s", i, b)
	}
	if len(b.Label.Symbols()) == 1 {
		return b.set.getNTSlot(b.Label.Symbols()[i], b.pivot, b.rightExtent)
	}
	if len(b.Label.Symbols()) == 2 {
		if i == 0 {
			return b.set.getNTSlot(b.Label.Symbols()[i], b.leftExtent, b.pivot)
		}
		return b.set.getNTSlot(b.Label.Symbols()[i], b.pivot, b.rightExtent)
	}
	idx := b.Label.Index()
	str := stringBSR{b.Label, b.leftExtent, b.pivot, b.rightExtent, b.set}
	for idx.Pos > i+1 && idx.Pos > 2 {
		idx.Pos--
		str = b.set.getString(slot.GetLabel(idx.NT, idx.Alt, idx.Pos), str.leftExtent, str.pivot)
		// fmt.Printf("  %s\n", str)
	}
	if i == 0 {
		return b.set.getNTSlot(b.Label.Symbols()[i], str.leftExtent, str.pivot)
	}
	return b.set.getNTSlot(b.Label.Symbols()[i], str.pivot, str.rightExtent)
}

// func (b BSR) GetString() string {
// 	return set.lex.GetString(b.LeftExtent(),b.RightExtent())
// }

// GetTChildI returns the terminal symbol at position i in b.
// GetTChildI panics if symbol i is not a valid terminal
func (b BSR) GetTChildI(i int) *token.Token {
	symbols := b.Label.Symbols()

	if i >= len(symbols) {
		panic(fmt.Sprintf("%s has no T child %d", b, i))
	}
	if symbols[i].IsNonTerminal() {
		panic(fmt.Sprintf("symbol %d in %s is an NT", i, b))
	}

	lext := b.leftExtent
	for j := 0; j < i; j++ {
		if symbols[j].IsNonTerminal() {
			nt := b.GetNTChildI(j)
			lext += nt.rightExtent - nt.leftExtent
		} else {
			lext++
		}
	}
	return b.set.lex.Tokens[lext]
}

func deleteNTSlotEntry(b BSR) {
	// fmt.Printf("deletNTSlotEntry(%s)\n", b)
	nts := ntSlot{b.Label.Head(), b.leftExtent, b.rightExtent}
	slots := b.set.ntSlotEntries[nts]
	slots1 := make([]BSR, 0, len(slots))
	bi := -1
	for i, s := range slots {
		if s == b && bi != -1 {
			panic(fmt.Sprintf("Duplicate slot entries: %d and %d", bi, i))
		} else {
			slots1 = append(slots1, s)
		}
	}
	b.set.ntSlotEntries[nts] = slots1
}

// LeftExtent returns the left extent of the BSR
func (b BSR) LeftExtent() int {
	return b.leftExtent
}

// RightExtent returns the right extent of the BSR
func (b BSR) RightExtent() int {
	return b.rightExtent
}

// Pivot returns the pivot of the BSR
func (b BSR) Pivot() int {
	return b.pivot
}

func (b BSR) String() string {
	return fmt.Sprintf("%s,%d,%d,%d - %s", b.Label, b.leftExtent, b.pivot, b.rightExtent,
		b.set.lex.GetString(b.LeftExtent(), b.RightExtent()-1))
}

func (s stringBSR) LeftExtent() int {
	return s.leftExtent
}

func (s stringBSR) RightExtent() int {
	return s.rightExtent
}

func (s stringBSR) Pivot() int {
	return s.pivot
}

func (s stringBSR) Empty() bool {
	return s.leftExtent == s.pivot && s.pivot == s.rightExtent
}

// String returns a string representation of s
func (s stringBSR) String() string {
	// fmt.Printf("bsr.stringBSR.stringBSR(): %s, %d, %d, %d\n",
	// 	s.Label.Symbols(), s.leftExtent, s.pivot, s.rightExtent)
	ss := s.Label.Symbols()[:s.Label.Pos()].Strings()
	str := strings.Join(ss, " ")
	return fmt.Sprintf("%s,%d,%d,%d - %s", str, s.leftExtent, s.pivot,
		s.rightExtent, s.set.lex.GetString(s.LeftExtent(), s.RightExtent()))
}

func (s *Set) getNTSlot(sym symbols.Symbol, leftExtent, rightExtent int) (bsrs []BSR) {
	nt, ok := sym.(symbols.NT)
	if !ok {
		line, col := s.getLineColumn(leftExtent)
		failf("%s is not an NT at line %d col %d", sym, line, col)
	}
	return s.ntSlotEntries[ntSlot{nt, leftExtent, rightExtent}]
}

func (s *Set) fail(b BSR, format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	line, col := s.getLineColumn(b.LeftExtent())
	panic(fmt.Sprintf("Error in BSR: %s at line %d col %d\n", msg, line, col))
}

func failf(format string, args ...interface{}) {
	panic(fmt.Sprintf("Error in BSR: %s\n", fmt.Sprintf(format, args...)))
}

func decodeRune(str []byte) (string, rune, int) {
	if len(str) == 0 {
		return "$", '$', 0
	}
	r, sz := utf8.DecodeRune(str)
	if r == utf8.RuneError {
		panic(fmt.Sprintf("Rune error: %s", str))
	}
	switch r {
	case '\t', ' ':
		return "space", r, sz
	case '\n':
		return "\\n", r, sz
	}
	return string(str[:sz]), r, sz
}

func (s *Set) getLineColumn(cI int) (line, col int) {
	return s.lex.GetLineColumnOfToken(cI)
}

// ReportAmbiguous lists the ambiguous subtrees of the parse forest
func (s *Set) ReportAmbiguous() {
	fmt.Println("Ambiguous BSR Subtrees:")
	rts := s.GetRoots()
	if len(rts) != 1 {
		fmt.Printf("BSR has %d ambigous roots\n", len(rts))
	}
	for i, b := range s.GetRoots() {
		fmt.Println("In root", i)
		if !s.report(b) {
			fmt.Println("No ambiguous BSRs")
		}
	}
}

// report return true iff at least one ambigous BSR was found
func (s *Set) report(b BSR) bool {
	ambiguous := false
	for i, sym := range b.Label.Symbols() {
		ln, col := s.getLineColumn(b.LeftExtent())
		if sym.IsNonTerminal() {
			if len(b.GetNTChildrenI(i)) != 1 {
				ambiguous = true
				fmt.Printf("  Ambigous: in %s: NT %s (%d) at line %d col %d \n",
					b, sym, i, ln, col)
				fmt.Println("   Children:")
				for _, c := range b.GetNTChildrenI(i) {
					fmt.Printf("     %s\n", c)
				}
			}
			for _, b1 := range b.GetNTChildrenI(i) {
				s.report(b1)
			}
		}
	}
	return ambiguous
}

// IsAmbiguous returns true if the BSR set does not have exactly one root, or
// if any BSR in the set has an NT symbol, which does not have exactly one
// sub-tree.
func (s *Set) IsAmbiguous() bool {
	if len(s.GetRoots()) != 1 {
		return true
	}
	return isAmbiguous(s.GetRoot())
}

// isAmbiguous returns true if b or any of its NT children is ambiguous.
// A BSR is ambigous if any of its NT symbols does not have exactly one
// subtrees (children).
func isAmbiguous(b BSR) bool {
	for i, s := range b.Label.Symbols() {
		if s.IsNonTerminal() {
			if len(b.GetNTChildrenI(i)) != 1 {
				return true
			}
			for _, b1 := range b.GetNTChildrenI(i) {
				if isAmbiguous(b1) {
					return true
				}
			}
		}
	}
	return false
}
