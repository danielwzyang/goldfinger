package board

var (
	taperedPieceWeights = [2][6]int{
		{82, 337, 365, 477, 1025, 12000}, // midgame
		{94, 281, 297, 512, 936, 12000},  // endgame
	} // [phase][piece]

	indexMap = [64]int{
		56, 57, 58, 59, 60, 61, 62, 63,
		48, 49, 50, 51, 52, 53, 54, 55,
		40, 41, 42, 43, 44, 45, 46, 47,
		32, 33, 34, 35, 36, 37, 38, 39,
		24, 25, 26, 27, 28, 29, 30, 31,
		16, 17, 18, 19, 20, 21, 22, 23,
		8, 9, 10, 11, 12, 13, 14, 15,
		0, 1, 2, 3, 4, 5, 6, 7,
	}

	psq = [2][6][64]int{
		{
			PAWN_MG,
			KNIGHT_MG,
			BISHOP_MG,
			ROOK_MG,
			QUEEN_MG,
			KING_MG,
		},
		{
			PAWN_EG,
			KNIGHT_EG,
			BISHOP_EG,
			ROOK_EG,
			QUEEN_EG,
			KING_EG,
		},
	}

	gamePhaseInc = [12]int{0, 1, 1, 2, 4, 0, 0, 1, 1, 2, 4, 0} // [piece]

	mgTable = [2][12][64]int{} // [color][piece][square]
	egTable = [2][12][64]int{} // [color][piece][square]

	MATE        = 10000
	LIMIT_SCORE = 20000
)

func InitEvalTables() {
	for piece := 0; piece < 6; piece++ {
		for square := 0; square < 64; square++ {
			psqIndex := indexMap[square]

			mgTable[WHITE][piece][square] = taperedPieceWeights[0][piece] + psq[0][piece][psqIndex]
			egTable[WHITE][piece][square] = taperedPieceWeights[1][piece] + psq[1][piece][psqIndex]

			mgTable[BLACK][piece][square] = -(taperedPieceWeights[0][piece] + psq[0][piece][psqIndex^56])
			egTable[BLACK][piece][square] = -(taperedPieceWeights[1][piece] + psq[1][piece][psqIndex^56])

			mgTable[WHITE][piece+6][square] = taperedPieceWeights[0][piece] + psq[0][piece][psqIndex]
			egTable[WHITE][piece+6][square] = taperedPieceWeights[1][piece] + psq[1][piece][psqIndex]

			mgTable[BLACK][piece+6][square] = -(taperedPieceWeights[0][piece] + psq[0][piece][psqIndex^56])
			egTable[BLACK][piece+6][square] = -(taperedPieceWeights[1][piece] + psq[1][piece][psqIndex^56])
		}
	}
}

func Evaluate() int {
	gamePhase := 0

	mgScore := 0
	egScore := 0

	color := WHITE

	for piece := WHITE_PAWN; piece <= BLACK_KING; piece++ {
		if piece == BLACK_PAWN {
			color = BLACK
		}

		bitboard := Bitboards[piece]

		for bitboard > 0 {
			square := LS1B(bitboard)

			mgScore += mgTable[color][piece][square]
			egScore += egTable[color][piece][square]

			gamePhase += gamePhaseInc[piece]

			PopBit(&bitboard, square)
		}
	}

	if gamePhase > 24 {
		gamePhase = 24
	}

	midgamePhase := gamePhase
	endgamePhase := 24 - gamePhase

	score := (mgScore*midgamePhase + egScore*endgamePhase) / 24

	score = (score * (100 - Fifty) / 100)

	if Side == WHITE {
		return score
	}

	return -score
}

func InsufficientMaterial() bool {
	// king vs king
	if Occupancies[BOTH] == (Bitboards[WHITE_KING] | Bitboards[BLACK_KING]) {
		return true
	}

	// king and knight vs king
	if Occupancies[BOTH] == (Bitboards[WHITE_KING]|Bitboards[BLACK_KING]|Bitboards[WHITE_KNIGHT]) ||
		Occupancies[BOTH] == (Bitboards[WHITE_KING]|Bitboards[BLACK_KING]|Bitboards[BLACK_KNIGHT]) {
		return true
	}

	// knight and bishop vs king
	if Occupancies[BOTH] == (Bitboards[WHITE_KING]|Bitboards[BLACK_KING]|Bitboards[WHITE_BISHOP]) ||
		Occupancies[BOTH] == (Bitboards[WHITE_KING]|Bitboards[BLACK_KING]|Bitboards[BLACK_BISHOP]) {
		return true
	}

	// king and bishop vs king and bishop
	if Occupancies[BOTH] == (Bitboards[WHITE_KING] | Bitboards[BLACK_KING] | Bitboards[WHITE_BISHOP] | Bitboards[BLACK_BISHOP]) {
		// bishops on the same color squares
		whiteBishop := LS1B(Bitboards[WHITE_BISHOP])
		blackBishop := LS1B(Bitboards[BLACK_BISHOP])
		if (whiteBishop+blackBishop)%2 == 0 {
			return true
		}
	}

	return false
}
