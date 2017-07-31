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
	next  [][]int // next match [kind.Idx]
	err   struct {
		flag  bool
		start int
		kind  *setsymbol.Symbol
	}
	cur   int
	lines int
}

// Lex fulfills parlex.Lexer. It uses the StackLexer to lex a string
func (l *StackLexer) Lex(str string) []parlex.Lexeme {
	op := &lexOp{
		subLexer: l.start,
		b:        []byte(str),
		lines:    1,
	}
	op.err.kind = l.set.Str(op.Error)
	if op.insert.startKind != "" {
		op.lxs = append(op.lxs, lexeme.String(op.insert.startKind).Set(op.insert.startVal))
	}
	op.populateNext()

	op.lex()

	if op.insert.endKind != "" {
		op.lxs = append(op.lxs, lexeme.String(op.insert.endKind).Set(op.insert.endVal))
	}

	return op.lxs
}

func (op *lexOp) lex() {
	for {
		lx, lxEnd, r := op.findNextMatch()
		if lxEnd == op.cur {
			if len(op.stack) == 0 {
				op.setError()
				op.cur++
				if op.cur >= len(op.b) {
					break
				}
			} else {
				op.pop(1)
			}
			continue
		}
		op.checkError()
		if !r.discard {
			op.lxs = append(op.lxs, lx)
		}
		op.cur = lxEnd
		if op.cur >= len(op.b) {
			break
		}
		if r != nil && r.pop > 0 {
			op.pop(r.pop)
		} else if r != nil && r.push != "" {
			op.stack = append(op.stack, op.subLexer)
			op.subLexer = op.lexers[r.push]
			op.populateNext()
		} else {
			op.updateNext()
		}
	}
	op.consumeRemainingAsError()
	op.checkError()
}

func (op *lexOp) pop(i int) {
	ln := len(op.stack) - i
	if ln < 0 {
		ln = 0
	}
	op.subLexer = op.stack[ln]
	op.stack = op.stack[:ln]
	op.populateNext()
}

func (op *lexOp) populateNext() {
	op.next = make([][]int, op.set.Size())
	for kind, r := range op.rules {
		if r != nil {
			loc := op.rules[kind].re.FindSubmatchIndex(op.b[op.cur:])
			if loc != nil {
				loc = append(loc, op.cur)
			}
			op.next[kind] = loc
		}
	}
}

func (op *lexOp) findNextMatch() (*lexeme.Lexeme, int, *rule) {
	var r *rule
	var idx []int

	// look in next for matches and take the longest one
	for kind, loc := range op.next {
		if loc != nil && loc[0]+loc[len(loc)-1] == op.cur {
			tr := op.rules[kind]
			if tr != nil && (r == nil || op.compare(loc[1], tr.priority, idx[1], r.priority)) {
				idx = loc
				r = tr
			}
		}
	}
	if r == nil {
		return nil, op.cur, nil
	}

	lx := &lexeme.Lexeme{
		K: op.set.ByIdx(r.kind),
	}
	offset := idx[len(idx)-1]
	if r.submatches != nil {
		for _, sub := range r.submatches {
			if sub.section == -1 {
				lx.V += sub.str
			} else {
				lx.V += string(op.b[offset+idx[sub.section*2] : offset+idx[sub.section*2+1]])
			}
		}
		op.handleLineCol(lx, string(op.b[offset+idx[0]:offset+idx[1]]))
	} else {
		lx.V = string(op.b[offset+idx[0] : offset+idx[1]])
		op.handleLineCol(lx, lx.V)
	}

	return lx, offset + idx[1], r
}

func (op *lexOp) handleLineCol(lx *lexeme.Lexeme, str string) {
	lx.L = op.lines
	lx.C = strings.LastIndex(string(op.b[:op.cur]), "\n")
	lx.C = op.cur - lx.C
	op.lines += strings.Count(str, "\n")
}

func (op *lexOp) checkError() {
	if !op.err.flag {
		return
	}
	op.err.flag = false
	val := string(op.b[op.err.start:op.cur])
	lx := lexeme.New(op.err.kind).Set(val)
	op.handleLineCol(lx, lx.V)
	op.lxs = append(op.lxs, &errLexeme{lx})
}

func (op *lexOp) updateNext() {
	for kind, loc := range op.next {
		if loc != nil && loc[0]+loc[len(loc)-1] <= op.cur {
			loc := op.rules[kind].re.FindSubmatchIndex(op.b[op.cur:])
			if loc != nil {
				loc = append(loc, op.cur)
			}
			op.next[kind] = loc
		}
	}
}

func (op *lexOp) consumeRemainingAsError() {
	if op.cur == len(op.b) {
		return
	}
	op.setError()
	op.cur = len(op.b)
}

func (op *lexOp) setError() {
	if op.err.flag {
		return
	}
	op.err.flag = true
	op.err.start = op.cur
}
