package pike

import "unsafe"

type cursor struct {
	ip            uint32
	partialGroups groupID
	groups        groupID
	counter       *counter
}

func (c *cursor) cmpr(c2 *cursor) int8 {
	if c.ip != c2.ip {
		if c.ip < c2.ip {
			return 0
		}
		return 1
	}
	if c.partialGroups != c2.partialGroups {
		if c.partialGroups < c2.partialGroups {
			return 0
		}
		return 1
	}
	if c.groups != c2.groups {
		if c.groups < c2.groups {
			return 0
		}
		return 1
	}
	cp := uintptr(unsafe.Pointer(c.counter))
	c2p := uintptr(unsafe.Pointer(c2.counter))
	if cp != c2p {
		if cp < c2p {
			return 0
		}
		return 1
	}
	return -1
}

type cursors struct {
	m     map[cursor]struct{}
	slice []cursor
	idx   int
}

func newCursors() *cursors {
	return &cursors{
		m: make(map[cursor]struct{}),
	}
}

func (c *cursors) add(cur cursor) {
	_, found := c.m[cur]
	if !found {
		c.m[cur] = struct{}{}
		c.slice = append(c.slice, cur)
	}
}

func (c *cursors) pop() (cur cursor, found bool) {
	if c.idx >= len(c.slice) {
		return
	}
	cur = c.slice[c.idx]
	found = true
	c.idx++
	return
}

func (c *cursors) reset() {
	c.idx = 0
	c.slice = c.slice[:0]
	for k := range c.m {
		delete(c.m, k)
	}
}

type cursorNode struct {
	self             *cursor
	children         [2]*cursorNode
	popVisit, popped bool
}

func (cn *cursorNode) add(cur *cursor) bool {
	if cn.self == nil {
		cn.self = cur
		cn.popVisit = true
		return true
	}
	c := cn.self.cmpr(cur)
	if c == -1 {
		return false
	}
	cn.popVisit = true
	if cn.children[c] == nil {
		cn.children[c] = &cursorNode{}
	}
	return cn.children[c].add(cur)
}

func (cn *cursorNode) pop() (cur *cursor, found bool) {
	cn.popVisit = false
	for _, c := range cn.children {
		if c != nil && c.popVisit {
			if !found {
				cur, found = c.pop()
			}
			cn.popVisit = cn.popVisit || c.popVisit
		}
	}
	if !found && cn.self != nil && !cn.popped {
		cp := *(cn.self)
		cur = &cp
		found = true
		cn.popped = true
	}
	cn.popVisit = cn.popVisit || !cn.popped
	return
}
