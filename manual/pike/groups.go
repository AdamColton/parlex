package pike

type group struct {
	idx        uint32
	start, end int
	prev       *group
}

func (g *group) toMap() map[uint32][][2]int {
	if g == nil {
		return make(map[uint32][][2]int)
	}
	m := g.prev.toMap()
	m[g.idx] = append(m[g.idx], [2]int{g.start, g.end})
	return m
}

func (g *group) close(end int) *group {
	return &group{
		idx:   g.idx,
		start: g.start,
		end:   end,
	}
}

func (g *group) open(idx uint32, start int) *group {
	return &group{
		idx:   idx,
		start: start,
		prev:  g,
	}
}
