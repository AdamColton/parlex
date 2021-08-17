package pike

import (
	"fmt"
	"math"
	"strings"
)

func printCode(code []byte) string {
	w := wrapper{
		slice: code,
	}
	var out []string

	type line struct {
		number, idx int
	}
	lineMap := make(map[uint32]line)

	backfill := make(map[uint32][]int)
	pos := func() {
		idx := w.idxUint32()
		backfill[idx] = append(backfill[idx], len(out))
		out = append(out, "")
	}

	var n int
	for !w.done() {
		lineMap[w.idx] = line{
			number: n,
			idx:    len(out),
		}
		out = append(out, "")
		n++

		i := w.inst()
		switch i {
		case i_match: // rune
			args := fmt.Sprintf("'%s'\n", string(w.rune()))
			out = append(out, " match ", args)
		case i_match_range: // startRune,endRune
			args := fmt.Sprintf("%v, %v\n", w.rune(), w.rune())
			out = append(out, " match ", args)
		case i_branch: // pos
			out = append(out, "branch ")
			pos()
			out = append(out, "\n")
		case i_jump: // pos
			out = append(out, "  jump ")
			pos()
			out = append(out, "\n")
		case i_stop:
			out = append(out, "  stop\n")
		case i_accept:
			out = append(out, "accept\n")
		case i_startGroup:
			args := fmt.Sprintf("%d\n", w.idxUint32())
			out = append(out, "groupS ", args)
		case i_closeGroup:
			out = append(out, "groupE\n")
		case i_wait:
			out = append(out, "  wait\n")
		case i_startCounter:
			out = append(out, "countS\n")
		case i_incCounter:
			out = append(out, "countI\n")
		case i_closeCounter:
			out = append(out, "countE\n")
		case i_ck_lt_c:
			args := fmt.Sprintf("%d\n", w.idxUint32())
			out = append(out, " c_gte ", args)
		case i_ck_gte_c:
			args := fmt.Sprintf("%d\n", w.idxUint32())
			out = append(out, "  c_lt ", args)
		}
	}

	ln := int(math.Log10(float64(n))) + 1
	f := fmt.Sprintf("%%%dd ", ln)
	for idx, l := range lineMap {
		out[l.idx] = fmt.Sprintf(f, l.number)
		s := fmt.Sprint(l.number)
		for _, sIdx := range backfill[idx] {
			out[sIdx] = s
		}
	}

	return strings.Join(out, "")
}
