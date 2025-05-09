package board

var (
	REPETITION_TABLE = [1000]uint64{}
	REPETITION_INDEX = 0
)

func IsRepetition() bool {
	for i := 0; i < REPETITION_INDEX; i++ {
		if REPETITION_TABLE[REPETITION_INDEX] == ZobristHash {
			return true
		}
	}

	return false
}
