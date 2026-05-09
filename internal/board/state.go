package board

type State struct {
	EnPassant       int8
	Castle          uint8
	ZobristHash     uint64
	Fifty           uint8
	RepetitionIndex int16
	LastCapture		int8
}

var StateStack [4096]State
var StateSize = 0

func ResetStateHistory() {
	StateSize = 0
}

func SaveState() {
	StateStack[StateSize] = State{
		int8(EnPassant),
		uint8(Castle),
		ZobristHash,
		uint8(Fifty),
		int16(RepetitionIndex),
		int8(LastCapture),
	}

	StateSize++
}

