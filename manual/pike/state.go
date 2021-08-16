package pike

import "hash"

type state []byte

func newState(ln uint32) state {
	return make(state, ln*4)
}

func (s state) hash(h hash.Hash64) uint64 {
	h.Reset()
	h.Write(s)
	return h.Sum64()
}

func (s state) workingState() *workingState {
	return &workingState{
		base: s,
	}
}

type workingState struct {
	base state
	mut  state
}

func (s *workingState) state() state {
	if s.mut != nil {
		s.base = s.mut
		s.mut = nil
	}
	return s.base
}

func (s *workingState) cp() {
	if s.mut != nil {
		return
	}
	s.mut = make([]byte, len(s.base))
	copy(s.mut, s.base)
}

func (s *workingState) set(r, v uint32) {
	s.cp()
	setUint32(s.mut, r*4, v)
}

func (s *workingState) inc(r uint32) {
	s.cp()
	idx := r * 4
	for i := uint32(0); i < 4; i++ {
		s.mut[idx]++
		if s.mut[idx] > 0 {
			break
		}
		idx++
	}
}

func (s *workingState) readUint32(r uint32) uint32 {
	if s.mut == nil {
		return readUint32(s.base, r*4)
	}
	return readUint32(s.mut, r*4)
}
