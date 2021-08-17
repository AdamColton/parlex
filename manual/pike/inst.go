package pike

type inst byte

const (
	undefined inst = iota
	i_wait
	i_match       // rune
	i_match_range // startRune,endRune
	i_branch      // pos
	i_jump        // pos
	i_stop
	i_accept
	i_startGroup
	i_closeGroup
	i_startCounter
	i_incCounter
	i_closeCounter
	i_ck_lt_c  // val
	i_ck_gte_c // val
)
