package board

var (
	RepetitionIndex = 0
	RepetitionTable = [4096]uint64{}
)

func IsRepetition() bool {
	count := 1
	historyDepth := max(0, RepetitionIndex-int(Fifty))
	for i := RepetitionIndex - 2; i >= historyDepth; i -= 2 {
		if RepetitionTable[i] == ZobristHash {
			count++
			if count == 3 {
				return true
			}
		}
	}
	return false
}
