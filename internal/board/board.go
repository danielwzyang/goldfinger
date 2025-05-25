package board

var (
	Bitboards [12]uint64

	// occupancies
	Occupancies = [3]uint64{
		0, // WHITE
		0, // BLACK
		0, // BOTH
	}

	Castle = 0
	/*
		0001 = 1 = white kingside
		0010 = 2 = white queenside
		0100 = 4 = black kingside
		1000 = 8 = black queenside
	*/
	WK = 1
	WQ = 2
	BK = 4
	BQ = 8

	EnPassant = INVALID_SQUARE
	Side      = WHITE
	Fifty     = 0
)

func Init() {
	InitNonSlidingAttacks()

	InitSlidingAttacks(true)  // bishops
	InitSlidingAttacks(false) // rooks

	InitEvalTables()

	InitZobristTables()

	ResetRepetition()

	ResetStateHistory()
}
