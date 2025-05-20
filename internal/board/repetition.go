package board

var (
	RepetitionIndex = 0
	RepetitionTable = [4096]uint64{}
)

func IsRepetition() bool {
	count := 0
	// current index = i (one repetition), then i - 4 = two repetitions, i - 8 = three repetitions.
	for i := RepetitionIndex; i >= RepetitionIndex-8; i -= 4 {
		if i < 0 {
			break
		}

		if RepetitionTable[i] == ZobristHash {
			count++
			if count == 3 {
				return true
			}
		}
	}
	return false
}
