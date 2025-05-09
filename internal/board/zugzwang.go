package board

func HasNonPawnMaterial() bool {
	var start, end int
	if Side == WHITE {
		start = WHITE_PAWN + 1
		end = WHITE_KING - 1
	} else {
		start = BLACK_PAWN + 1
		end = BLACK_KING - 1
	}

	for i := start; i <= end; i++ {
		if Bitboards[i] != 0 {
			return true
		}
	}

	return false
}
