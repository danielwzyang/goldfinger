package board

var (
	RepetitionIndex = 0
	RepetitionTable = [4096]uint64{}
)

func IsRepetition() bool {
	for i := RepetitionIndex - 1; i >= 0; i-- {
		if RepetitionTable[i] == ZobristHash {
			return true
		}
	}
	return false
}
