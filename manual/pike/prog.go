package pike

import (
	"hash"
	"hash/crc64"
)

type prog struct {
	stateLen uint32
	code     []byte
}

type runOp struct {
	p           prog
	h           hash.Hash64
	flow, match *cursors
	r           *Reader
	best        int
	bestState   state
	groups      map[uint32][][2]int
}

func (p prog) run(input string) *runOp {
	op := &runOp{
		p:     p,
		h:     crc64.New(crc64.MakeTable(crc64.ECMA)),
		r:     NewStringReader(input),
		flow:  newCursors(),
		match: newCursors(),
		best:  -1,
	}
	op.flow.add(newState(p.stateLen), &cursor{}, op.h)

	op.run()
	return op
}

func (op *runOp) run() {
	for len(op.flow.keys) > 0 {
		// run all until match
		if op.flowOps() {
			op.best = op.r.Idx + op.r.Ln
		}
		op.r.Inc()
		op.matchOps()
	}

}

func (op *runOp) flowOps() bool {
	accept := false
	for sc := op.flow.pop(); sc != nil; sc = op.flow.pop() {
		for _, c := range sc.cursors {
			s := sc.state.workingState()
			w := op.wrapper(c)
		cursorLoop:
			for {
				i := w.inst()
				if i < startFlowOps {
					c.ip = w.idx - 1
					op.match.add(s.state(), c, op.h)
					break
				}
				switch i {
				case i_branch:
					cp := c.copy()
					cp.ip = w.idx + 4
					op.flow.add(s.state(), cp, op.h)
					w.jump()
				case i_jump:
					w.jump()
				case i_stop:
					break cursorLoop
				case i_accept:
					accept = true
					op.bestState = s.state()
					op.groups = c.groups.toMap()
				case i_inc:
					s.inc(w.idxUint32())
				case i_set_rv:
					s.set(w.idxUint32(), w.idxUint32())
				case i_set_rr:
					s.set(w.idxUint32(), s.readUint32(w.idxUint32()))
				case i_ck_lt_rv:
					r := s.readUint32(w.idxUint32())
					v := w.idxUint32()
					if !(r < v) {
						break cursorLoop
					}
				case i_ck_gte_rv:
					r := s.readUint32(w.idxUint32())
					v := w.idxUint32()
					if !(r >= v) {
						break cursorLoop
					}
				case i_startGroup:
					idx := w.idxUint32()
					start := op.r.Idx + op.r.Ln
					c.partialGroups = c.partialGroups.open(idx, start)
				case i_closeGroup:
					g := c.partialGroups.close(op.r.Idx + op.r.Ln)
					c.partialGroups = c.partialGroups.prev
					g.prev = c.groups
					c.groups = g
				}
			}
		}
	}
	return accept
}

func (op *runOp) wrapper(c *cursor) *wrapper {
	return &wrapper{
		slice: op.p.code,
		idx:   c.ip,
	}
}

func (op *runOp) matchOps() {
	r := op.r.R
	// run all matches
	for sc := op.match.pop(); sc != nil; sc = op.match.pop() {
		for _, c := range sc.cursors {
			w := op.wrapper(c)
			switch w.inst() {
			case i_match:
				expect := rune(w.idxUint32())
				if r == expect {
					c.ip = w.idx
					op.flow.add(sc.state, c, op.h)
				}
			case i_match_range:
				start := rune(w.idxUint32())
				end := rune(w.idxUint32())
				if r >= start && r <= end {
					c.ip = w.idx
					op.flow.add(sc.state, c, op.h)
				}
			}
		}
	}
}
