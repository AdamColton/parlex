package stacklexer

import (
	"github.com/adamcolton/parlex"
	"github.com/adamcolton/parlex/lexeme"
	"github.com/adamcolton/parlex/symbol/setsymbol"
	"strings"
)

type lexOp struct {
	*subLexer
	stack []*subLexer
	b     []byte
	lxs   []parlex.Lexeme
	next  [][]int // next match [kind.Idx] => [0]:start, [1]: end
	err   struct {
		flag  bool
		start int
		kind  *setsymbol.Symbol
	}
	cur   int
	lines int
}

func (l *StackLexer) Lex(str string) []parlex.Lexeme {
	op := &lexOp{
		subLexer: l.start,
		b:        []byte(str),
		lines:    1,
	}
	op.err.kind = l.set.Str(op.Error)
	op.populateNext()

	op.lex()

	return op.lxs
}

func (op *lexOp) lex() {
	for {
		lx, lxEnd, r := op.findNextMatch()
		if lxEnd == op.cur {
			if !op.err.flag {
				op.err.flag = true
				op.err.start = op.cur
			}
			op.cur++
		} else {
			op.checkError()
			if !r.discard {
				op.lxs = append(op.lxs, lx)
			}
			op.cur = lxEnd
		}
		if op.cur >= len(op.b) {
			break
		}
		if r != nil && r.pop {
			ln := len(op.stack)
			if ln == 0 {
				break
			}
			ln--
			op.subLexer = op.stack[ln]
			op.stack = op.stack[:ln]
			op.populateNext()
		} else if r != nil && r.push != "" {
			op.stack = append(op.stack, op.subLexer)
			op.subLexer = op.lexers[r.push]
			op.populateNext()
		} else {
			op.updateNext()
		}
	}
	op.checkError()
	op.consumeRemainingAsError()
}

func (op *lexOp) populateNext() {
	op.next = make([][]int, op.set.Size())
	for kind, r := range op.rules {
		if r != nil {
			loc := op.rules[kind].re.FindIndex(op.b[op.cur:])
			if loc != nil {
				loc[0] += op.cur
				loc[1] += op.cur
			}
			op.next[kind] = loc
		}
	}
}

func (op *lexOp) findNextMatch() (*lexeme.Lexeme, int, *rule) {
	var r *rule
	lx := &lexeme.Lexeme{}
	lxEnd := op.cur
	lxP := -1

	// look in next for matches and take the longest one
	for kind, loc := range op.next {
		if loc != nil && loc[0] == op.cur {
			tr := op.rules[kind]
			if tr != nil && op.compare(loc[1], tr.priority, lxEnd, lxP) {
				lx.K = op.set.ByIdx(kind)
				lx.V = string(op.b[loc[0]:loc[1]])
				lxEnd = loc[1]
				lxP = tr.priority
				r = tr
			}
		}
	}

	op.handleLineCol(lx)

	return lx, lxEnd, r
}

func (op *lexOp) handleLineCol(lx *lexeme.Lexeme) {
	lx.L = op.lines
	lx.C = strings.LastIndex(string(op.b[:op.cur]), "\n")
	lx.C = op.cur - lx.C
	op.lines += strings.Count(lx.V, "\n")
}

func (op *lexOp) checkError() {
	if !op.err.flag {
		return
	}
	op.err.flag = false
	val := string(op.b[op.err.start:op.cur])
	lx := lexeme.New(op.err.kind).Set(val)
	op.handleLineCol(lx)
	op.lxs = append(op.lxs, errLexeme{lx})
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

func (op *lexOp) consumeRemainingAsError() {
	if op.cur == len(op.b) {
		return
	}
	op.err.flag = true
	op.err.start = op.cur
	op.cur = len(op.b)
	op.checkError()
}
