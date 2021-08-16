package pike

import "unicode/utf8"

func NewStringReader(input string) *Reader {
	return NewReader([]byte(input))
}

func NewReader(input []byte) *Reader {
	r := &Reader{
		input: input,
	}
	return r
}

type Reader struct {
	input    []byte
	R        rune
	Idx, Ln  int
	row, col int
}

func (r *Reader) Inc() {
	if r.Done() {
		return
	}

	if r.R == '\n' {
		r.row++
		r.col = 0
	} else {
		r.col++
	}

	r.Idx += r.Ln
	r.R, r.Ln = utf8.DecodeRune(r.input[r.Idx:])
}

func (r *Reader) Done() bool {
	return r.R == utf8.RuneError
}

func (r *Reader) Consume(c rune) {
	for !r.Done() && r.R == c {
		r.Inc()
	}
}

func (r *Reader) ConsumeRange(s, e rune) {
	for !r.Done() && r.R >= s && r.R <= e {
		r.Inc()
	}
}
