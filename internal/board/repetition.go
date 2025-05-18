package board

var (
	RepetitionIndex = 0
	RepetitionTable = [4096]uint64{}
)

func IsRepetition() bool {
	count := 0
	for i := RepetitionIndex; i >= 0; i-- {
		if RepetitionTable[i] == ZobristHash {
			count++
			if count == 3 {
				return true
			}
		}
	}
	return false
}
