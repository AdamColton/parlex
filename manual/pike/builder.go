package pike

type builder struct {
	prog *prog
}

func newBuilder() *builder {
	return &builder{
		prog: &prog{},
	}
}

func (b *builder) close() *prog {
	return b.prog
}

func (b *builder) wait() {
	b.inst(i_wait)
}

func (b *builder) match(r rune) {
	b.inst(i_match)
	b.u32(uint32(r))
}

func (b *builder) match_range(start, end rune) {
	b.inst(i_match_range)
	b.u32(uint32(start))
	b.u32(uint32(end))
}

func (b *builder) branch(pos uint32) {
	b.inst(i_branch)
	b.u32(pos)
}

func (b *builder) jump(pos uint32) {
	b.inst(i_jump)
	b.u32(pos)
}

func (b *builder) stop() {
	b.inst(i_stop)
}

func (b *builder) accept() {
	b.inst(i_accept)
}

func (b *builder) startGroup(idx uint32) {
	b.inst(i_startGroup)
	b.u32(idx)
}

func (b *builder) closeGroup() {
	b.inst(i_closeGroup)
}

func (b *builder) startCounter() {
	b.inst(i_startCounter)
}
func (b *builder) incCounter() {
	b.inst(i_incCounter)
}

func (b *builder) closeCounter() {
	b.inst(i_closeCounter)
}

func (b *builder) ck_lt_c(val uint32) {
	b.inst(i_ck_lt_c)
	b.u32(val)
}

func (b *builder) ck_gte_c(val uint32) {
	b.inst(i_ck_gte_c)
	b.u32(val)
}

func (b *builder) u32(v uint32) {
	for i := uint32(0); i < 4; i++ {
		b.prog.code = append(b.prog.code, byte(v))
		v >>= 8
	}
}

func (b *builder) inst(i inst) {
	b.prog.code = append(b.prog.code, byte(i))
}

func (b *builder) deferU32() func(uint32) {
	idx := uint32(len(b.prog.code))
	b.u32(0)
	return func(v uint32) {
		for i := uint32(0); i < 4; i++ {
			b.prog.code[idx+i] = byte(v)
			v >>= 8
		}
	}
}

func (b *builder) loc() uint32 {
	return uint32(len(b.prog.code))
}

func (b *builder) deferLoc() func() {
	fn := b.deferU32()
	return func() {
		fn(b.loc())
	}
}

func (b *builder) defer_jump() func() {
	b.inst(i_jump)
	return b.deferLoc()
}

func (b *builder) defer_branch() func() {
	b.inst(i_branch)
	return b.deferLoc()
}
