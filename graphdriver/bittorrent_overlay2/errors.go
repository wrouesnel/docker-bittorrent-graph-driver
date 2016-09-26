package bittorrent_overlay2

// Error type returned for unimplemented structs (temporarily).
type ErrUnimplemented struct {
}

func (this *ErrUnimplemented) Error() string {
	return "Feature unimplemented"
}
