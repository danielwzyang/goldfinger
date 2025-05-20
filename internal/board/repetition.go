package board

var (
	RepetitionIndex = 0
	RepetitionTable = [4096]uint64{}
)

func IsRepetition() bool {
	count := 0
	for i := 0; i < RepetitionIndex; i++ {
		if RepetitionTable[i] == ZobristHash {
			count++
			if count == 3 {
				return true
			}
		}
	}
	return false
}
