package board

// can be optimized by keeping track of moves but im a bit lazy for now

var BoardStates = [4096]State{}
var StateSize = 0

type State struct {
	Bitboards       [12]uint64
	Occupancies     [3]uint64
	Side            int
	EnPassant       int
	Castle          int
	ZobristHash     uint64
	Fifty           int
	RepetitionIndex int
}

func ResetStateHistory() {
	StateSize = 0
}

func SaveState() {
	BoardStates[StateSize] = State{
		Bitboards,
		Occupancies,
		Side,
		EnPassant,
		Castle,
		ZobristHash,
		Fifty,
		RepetitionIndex,
	}

	StateSize++
}

func RestoreState() {
	state := BoardStates[StateSize-1]

	Bitboards = state.Bitboards
	Occupancies = state.Occupancies
	Side = state.Side
	EnPassant = state.EnPassant
	Castle = state.Castle
	ZobristHash = state.ZobristHash
	Fifty = state.Fifty
	RepetitionIndex = state.RepetitionIndex

	StateSize--
}
