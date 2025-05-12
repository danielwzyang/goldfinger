package board

var (
	RepetitionIndex = 0
	RepetitionTable = [4096]uint64{}
)

func IsRepetition() bool {
	for i := 0; i < RepetitionIndex; i++ {
		if RepetitionTable[RepetitionIndex] == ZobristHash {
			return true
		}
	}

	return false
}
