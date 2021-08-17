package pike

// wrapper around a slice with an index assuming a read-forward operation.
// helper for reading
type wrapper struct {
	slice []byte
	idx   uint32
}

func (w *wrapper) idxByte() byte {
	b := w.slice[w.idx]
	w.idx++
	return b
}

func (w *wrapper) idxUint32() uint32 {
	u := w.readUint32(w.idx)
	w.idx += 4
	return u
}

func (w *wrapper) setIdxUint32(v uint32) {
	setUint32(w.slice, w.idx, v)
	w.idx += 4
}

func (w *wrapper) readUint32(idx uint32) uint32 {
	return readUint32(w.slice, idx)
}

func (w *wrapper) inst() inst {
	if int(w.idx) >= len(w.slice) {
		return i_stop
	}
	return inst(w.idxByte())
}

func (w *wrapper) jump() {
	w.idx = w.readUint32(w.idx)
}

func (w *wrapper) done() bool {
	return w.idx >= uint32(len(w.slice))
}

func (w *wrapper) rune() rune {
	return rune(w.idxUint32())
}

func readUint32(s []byte, idx uint32) uint32 {
	return uint32(s[idx]) +
		uint32(s[idx+1])<<8 +
		uint32(s[idx+2])<<(8*2) +
		uint32(s[idx+3])<<(8*3)
}

func setUint32(s []byte, idx, val uint32) {
	for i := uint32(0); i < 4; i++ {
		b := byte(val)
		val >>= 8
		s[idx+i] = b
	}
}
