package pike

import (
	"hash"
	"hash/crc64"
)

type groupID uint64

type group struct {
	idx        uint32
	start, end uint32
	prev       groupID
}

type groupMap struct {
	id2g map[groupID]group
	g2id map[group]groupID
	h    hash.Hash64
}

func newGroupMap() *groupMap {
	return &groupMap{
		id2g: make(map[groupID]group),
		g2id: make(map[group]groupID),
		h:    crc64.New(crc64.MakeTable(crc64.ECMA)),
	}
}

func (g *group) hash(h hash.Hash64) groupID {
	w := &wrapper{
		slice: make([]byte, 20),
	}
	w.setIdxUint32(g.idx)
	w.setIdxUint32(g.start)
	w.setIdxUint32(g.end)
	w.setIdxUint32(uint32(g.prev))
	w.setIdxUint32(uint32(g.prev >> 32))
	h.Reset()
	h.Write(w.slice)
	return groupID(h.Sum64())
}

func (gm groupMap) open(prev groupID, idx uint32, start uint32) groupID {
	g := group{
		idx:   idx,
		start: start,
		prev:  prev,
	}
	id, found := gm.g2id[g]
	if !found {
		id = g.hash(gm.h)
		gm.g2id[g] = id
		gm.id2g[id] = g
	}

	return id
}

func (gm groupMap) close(id, prev groupID, end uint32) (partial, complete groupID) {
	g := gm.id2g[id]
	partial = g.prev
	g.end = end
	g.prev = prev

	var found bool
	complete, found = gm.g2id[g]
	if !found {
		complete = g.hash(gm.h)
		gm.g2id[g] = complete
		gm.id2g[complete] = g
	}
	return
}

func (gm groupMap) toMap(id groupID) map[uint32][][2]int {
	if id == 0 {
		return make(map[uint32][][2]int)
	}
	g := gm.id2g[id]
	m := gm.toMap(g.prev)
	m[g.idx] = append(m[g.idx], [2]int{int(g.start), int(g.end)})
	return m
}
