package pike

import "hash"

type cursor struct {
	ip            uint32
	partialGroups groupID
	groups        groupID
	counter       *counter
}

func (c *cursor) copy() *cursor {
	cp := *c
	return &cp
}

type stateCursors struct {
	state   state
	cursors []*cursor
}

type cursors struct {
	m    map[uint64]*stateCursors
	keys []uint64
}

func newCursors() *cursors {
	return &cursors{
		m: make(map[uint64]*stateCursors),
	}
}

func (c *cursors) add(s state, cur *cursor, h hash.Hash64) {
	key := s.hash(h)
	sc, found := c.m[key]
	if found {
		sc.cursors = append(sc.cursors, cur)
	} else {
		c.m[key] = &stateCursors{
			state:   s,
			cursors: []*cursor{cur},
		}
		c.keys = append(c.keys, key)
	}
}

func (c *cursors) pop() *stateCursors {
	ln := len(c.keys) - 1
	if ln < 0 {
		return nil
	}
	k := c.keys[ln]
	c.keys = c.keys[:ln]
	sc := c.m[k]
	delete(c.m, k)
	return sc
}
