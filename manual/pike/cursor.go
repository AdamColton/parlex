package pike

type cursor struct {
	ip            uint32
	partialGroups groupID
	groups        groupID
	counter       *counter
}

type cursors struct {
	m map[cursor]struct{}
}

func newCursors() *cursors {
	return &cursors{
		m: make(map[cursor]struct{}),
	}
}

func (c *cursors) add(cur cursor) {
	c.m[cur] = struct{}{}
}

func (c *cursors) pop() (cur cursor, found bool) {
	for getCur := range c.m {
		cur = getCur
		delete(c.m, cur)
		found = true
		return
	}
	return
}
