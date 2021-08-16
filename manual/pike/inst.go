package pike

type inst byte

const (
	undefined     inst = iota
	i_match            // rune
	i_match_range      // startRune,endRune
	i_branch           // pos
	i_jump             // pos
	i_stop
	i_accept
	i_inc       // reg
	i_set_rv    // val, reg
	i_set_rr    // regTo regFrom
	i_ck_lt_rv  // reg val
	i_ck_gte_rv // reg val
	i_startGroup
	i_closeGroup
)

const (
	startFlowOps = i_branch
)