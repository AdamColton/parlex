package simplelexer

import (
	"fmt"
	"github.com/adamcolton/parlex"
	"github.com/adamcolton/parlex/lexeme"
	"github.com/adamcolton/parlex/symbol/stringsymbol"
	"regexp"
	"strings"
)

// Lexer implements parlex.Lexer. It can take a string and produce a slice of
// lexemes.
type Lexer struct {
	order           []stringsymbol.Symbol
	rules           map[stringsymbol.Symbol]*rule
	compare         func(e1, p1, e2, p2 int) bool
	priorityCounter int
	Error           stringsymbol.Symbol
}

// ByLength sets the lexer to choose the longest match and use priority to
// decide a tie. This is the default.
func (l *Lexer) ByLength() { l.compare = lengthThenPriority }

// ByPriority sets the lexer to choose the highest priority match and use the
// length to decide a tie.
func (l *Lexer) ByPriority() { l.compare = priorityThenLength }

func priorityThenLength(e1, p1, e2, p2 int) bool {
	return p1 < p2 || (p1 == p2 && e1 > e2)
}
func lengthThenPriority(e1, p1, e2, p2 int) bool {
	return e1 > e2 || (e1 == e2 && p1 < p2)
}

type rule struct {
	kind     stringsymbol.Symbol
	re       *regexp.Regexp
	discard  bool
	priority int
}

type errLexeme struct {
	*lexeme.Lexeme
}

func (e *errLexeme) Error() string {
	return fmt.Sprintf("Lex Error %d:%d) %s", e.L, e.C, e.Value())
}

type lexOp struct {
	*Lexer
	b        []byte
	lxs      []parlex.Lexeme
	next     map[stringsymbol.Symbol][]int
	errFlag  bool
	errStart int
	cur      int
	lines    int
}

// Lex takes a string and produces a slice of lexemes that can be consumed by a
// parser.
func (l *Lexer) Lex(str string) []parlex.Lexeme {
	op := &lexOp{
		Lexer: l,
		b:     []byte(str),
		next:  make(map[stringsymbol.Symbol][]int),
	}
	op.populateNext()

	for {
		lx, lxEnd := op.findNextMatch()
		if lxEnd == op.cur {
			if !op.errFlag {
				op.errFlag = true
				op.errStart = op.cur
			}
			op.cur++
		} else {
			op.checkError()
			if !op.rules[lx.K.(stringsymbol.Symbol)].discard {
				op.lxs = append(op.lxs, lx)
			}
			op.cur = lxEnd
		}
		if op.cur >= len(op.b) {
			break
		}
		op.updateNext()
	}
	op.checkError()

	return op.lxs
}

func (op *lexOp) checkError() {
	if !op.errFlag {
		return
	}
	op.errFlag = false
	val := string(op.b[op.errStart:op.cur])
	op.lxs = append(op.lxs, errLexeme{lexeme.New(op.Error).Set(val)})
}

func (op *lexOp) populateNext() {
	for kind, r := range op.rules {
		op.next[kind] = r.re.FindIndex(op.b)
	}
}

func (op *lexOp) updateNext() {
	for kind, loc := range op.next {
		if loc != nil && loc[0] <= op.cur {
			loc := op.rules[kind].re.FindIndex(op.b[op.cur:])
			if loc != nil {
				loc[0] += op.cur
				loc[1] += op.cur
			}
			op.next[kind] = loc
		}
	}
}

func (op *lexOp) findNextMatch() (*lexeme.Lexeme, int) {
	lx := &lexeme.Lexeme{}
	lxEnd := op.cur
	lxP := -1

	// look in next for matches and take the longest one
	for kind, loc := range op.next {
		if loc != nil && loc[0] == op.cur {
			p := op.rules[kind].priority
			if op.compare(loc[1], p, lxEnd, lxP) {
				lx.K = kind
				lx.V = string(op.b[loc[0]:loc[1]])
				lxEnd = loc[1]
				lxP = p
			}
		}
	}

	lx.L = op.lines
	lx.C = strings.LastIndex(string(op.b[:op.cur]), "\n")
	lx.C = op.cur - lx.C
	op.lines += strings.Count(lx.V, "\n")

	return lx, lxEnd
}
