package pike

import "unsafe"

type cursor struct {
	ip            uint32
	partialGroups groupID
	groups        groupID
	counter       *counter
	noMod         bool
}

func (c *cursor) setCounter(ctr *counter) *cursor {
	if c.noMod {
		return &cursor{
			partialGroups: c.partialGroups,
			groups:        c.groups,
			counter:       ctr,
		}
	}
	c.counter = ctr
	return c
}

func (c *cursor) setIP(ip uint32) *cursor {
	if c.noMod {
		return &cursor{
			ip:            ip,
			partialGroups: c.partialGroups,
			groups:        c.groups,
			counter:       c.counter,
		}
	}
	c.ip = ip
	return c
}

func (c *cursor) setGroups(p, g groupID) *cursor {
	if c.noMod {
		return &cursor{
			partialGroups: p,
			groups:        g,
			counter:       c.counter,
		}
	}
	c.partialGroups, c.groups = p, g
	return c
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

type cursorNode struct {
	self             *cursor
	children         [2]*cursorNode
	popVisit, popped bool
}

func (cn *cursorNode) add(cur *cursor) bool {
	if cn.self == nil {
		cn.self = cur
		cn.popVisit = true
		cur.noMod = true
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
		cp := *cn.self
		cur = &cp
		found = true
		cn.popped = true
	}
	cn.popVisit = cn.popVisit || !cn.popped
	return
}
