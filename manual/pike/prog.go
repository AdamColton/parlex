package pike

import (
	"hash"
	"hash/crc64"
)

type prog struct {
	code []byte
}

func (p *prog) optimize() {
	w := &wrapper{
		slice: p.code,
	}
	for !w.done() {
		i := w.inst()
		switch i {
		case i_branch, i_jump:
			reset := w.idx - 1
			to := w.idxUint32()
			nextI := inst(w.slice[to])
			if nextI == i_jump {
				nextJmp := readUint32(w.slice, to+1)
				setUint32(w.slice, reset+1, nextJmp)
				w.idx = reset // go back on next pass and see if we're pointing to another jump
			}
		case i_match, i_ck_lt_c, i_ck_gte_c:
			w.idx += 4
		case i_match_range:
			w.idx += 8
		}
	}
}

type runOp struct {
	p           prog
	h           hash.Hash64
	flow, wait  *cursors
	r           *Reader
	best        int
	bestGroups  groupID
	counterRoot *counter
	groupMap    *groupMap
}

func (p prog) run(input string) *runOp {
	op := &runOp{
		p:           p,
		h:           crc64.New(crc64.MakeTable(crc64.ECMA)),
		r:           NewStringReader(input),
		flow:        newCursors(),
		wait:        newCursors(),
		best:        -1,
		counterRoot: &counter{},
		groupMap:    newGroupMap(),
	}
	op.flow.add(cursor{})

	op.run()
	return op
}

func (op *runOp) run() {
	for len(op.flow.m) > 0 {
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
	for c, found := op.flow.pop(); found; c, found = op.flow.pop() {
		w := op.wrapper(&c)
	cursorLoop:
		for {
			i := w.inst()
			switch i {
			case i_wait:
				c.ip = w.idx
				op.wait.add(c)
				break cursorLoop
			case i_branch:
				cp := c
				cp.ip = w.idx + 4
				op.flow.add(cp)
				w.jump()
			case i_jump:
				w.jump()
			case i_stop:
				break cursorLoop
			case i_accept:
				accept = true
				op.bestGroups = c.groups
			case i_startGroup:
				idx := w.idxUint32()
				start := uint32(op.r.Idx + op.r.Ln)
				c.partialGroups = op.groupMap.open(c.partialGroups, idx, start)
			case i_closeGroup:
				end := uint32(op.r.Idx + op.r.Ln)
				c.partialGroups, c.groups = op.groupMap.close(c.partialGroups, c.groups, end)
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
			case i_startCounter:
				if c.counter == nil {
					c.counter = op.counterRoot
				} else {
					c.counter = c.counter.newCounter()
				}
			case i_incCounter:
				c.counter = c.counter.inc()
			case i_closeCounter:
				c.counter = c.counter.pop()
			case i_ck_lt_c:
				v := w.idxUint32()
				if !(c.counter.val < v) {
					break cursorLoop
				}
			case i_ck_gte_c:
				v := w.idxUint32()
				if !(c.counter.val >= v) {
					break cursorLoop
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
	op.flow, op.wait = op.wait, op.flow
}
