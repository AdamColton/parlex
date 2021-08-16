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
	p          prog
	h          hash.Hash64
	flow, wait *cursors
	r          *Reader
	best       int
	bestState  state
	groups     map[uint32][][2]int
}

func (p prog) run(input string) *runOp {
	op := &runOp{
		p:    p,
		h:    crc64.New(crc64.MakeTable(crc64.ECMA)),
		r:    NewStringReader(input),
		flow: newCursors(),
		wait: newCursors(),
		best: -1,
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
		op.waitOps()
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
				switch i {
				case i_wait:
					c.ip = w.idx
					op.wait.add(s.state(), c, op.h)
					break cursorLoop
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
				case i_match:
					expect := rune(w.idxUint32())
					if op.r.R != expect {
						break cursorLoop
					}
				case i_match_range:
					start := rune(w.idxUint32())
					end := rune(w.idxUint32())
					if op.r.R < start || op.r.R > end {
						break cursorLoop
					}
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

func (op *runOp) waitOps() {
	// TODO: just flip
	for sc := op.wait.pop(); sc != nil; sc = op.wait.pop() {
		for _, c := range sc.cursors {
			op.flow.add(sc.state, c, op.h)
		}
	}
}
