package board

import (
	"math"
)

var (
	mgTable = [2][6][64]int{} // [color][piece][square]
	egTable = [2][6][64]int{} // [color][piece][square]

	isolatedPawns = [64]uint64{}    // [square]
	doubledPawns  = [64]uint64{}    // [square]
	passedPawns   = [2][64]uint64{} // [color][square]
	outposts      = [2][64]uint64{} // [color][square]
)

func InitEvalTables() {
	for piece := 0; piece < 6; piece++ {
		for square := 0; square < 64; square++ {
			psqIndex := indexMap[square]

			mgTable[WHITE][piece][square] = taperedPieceWeights[0][piece] + psq[0][piece][psqIndex]
			egTable[WHITE][piece][square] = taperedPieceWeights[1][piece] + psq[1][piece][psqIndex]

			mgTable[BLACK][piece][square] = -(taperedPieceWeights[0][piece] + psq[0][piece][psqIndex^56])
			egTable[BLACK][piece][square] = -(taperedPieceWeights[1][piece] + psq[1][piece][psqIndex^56])
		}
	}

	for square := 0; square < 64; square++ {
		isolatedPawns[square] = isolatedMask(square)
		doubledPawns[square] = doubledMask(square)
		passedPawns[WHITE][square] = passedMask(WHITE, square)
		passedPawns[BLACK][square] = passedMask(BLACK, square)
		outposts[WHITE][square] = outpostMask(WHITE, square)
		outposts[BLACK][square] = outpostMask(BLACK, square)
	}
}

func Evaluate() int {
	gamePhase := 0

	mgScore := 0
	egScore := 0

	var bitboard uint64
	var square, p int

	color := WHITE

	// Calculate material and piece-square scores
	for piece := WHITE_PAWN; piece <= BLACK_KING; piece++ {
		if piece == BLACK_PAWN {
			color = BLACK
		}

		bitboard = Bitboards[piece]

		p = piece % 6

		for bitboard > 0 {
			square = LS1B(bitboard)

			mgScore += mgTable[color][p][square] + pieceBonus(piece, square)
			egScore += egTable[color][p][square] + pieceBonus(piece, square)

			gamePhase += GAME_PHASE_INCREMENT[p]

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

	kingSafetyScore := 0
	whiteKingSquare := LS1B(Bitboards[WHITE_KING])
	blackKingSquare := LS1B(Bitboards[BLACK_KING])

	kingSafetyWeight := (24 - gamePhase) * 2
	kingSafetyScore = (kingBonus(WHITE, whiteKingSquare) - kingBonus(BLACK, blackKingSquare)) * kingSafetyWeight / 24

	score += kingSafetyScore

	if Side == WHITE {
		return score
	}

	return -score
}

func pieceBonus(piece int, square int) int {
	switch piece {
	case WHITE_PAWN:
		return pawnBonus(WHITE, square)
	case BLACK_PAWN:
		return pawnBonus(BLACK, square)
	case WHITE_KNIGHT:
		return knightBonus(WHITE, square)
	case BLACK_KNIGHT:
		return knightBonus(BLACK, square)
	case WHITE_BISHOP:
		return bishopBonus(WHITE, square)
	case BLACK_BISHOP:
		return bishopBonus(BLACK, square)
	case WHITE_ROOK:
		return rookBonus(WHITE, square)
	case BLACK_ROOK:
		return rookBonus(BLACK, square)
	case WHITE_QUEEN:
		return queenBonus(WHITE, square)
	case BLACK_QUEEN:
		return queenBonus(BLACK, square)
	case WHITE_KING:
		return kingBonus(WHITE, square)
	case BLACK_KING:
		return kingBonus(BLACK, square)
	}

	return 0
}

func pawnBonus(color int, square int) int {
	value := 0
	if isProtected(color, square) {
		value += PAWN_PROTECTED_WEIGHT
	}
	if isDoubled(color, square) {
		value += PAWN_DOUBLED_WEIGHT
	}
	if isIsolated(color, square) {
		value += PAWN_ISOLATED_WEIGHT
	}
	if isPassed(color, square) {
		value += PAWN_PASSED_WEIGHT
	}

	if color == WHITE {
		return value
	}

	return -value
}

func isProtected(color int, square int) bool {
	pawn := WHITE_PAWN
	if color == BLACK {
		pawn = BLACK_PAWN
	}

	return PAWN_ATTACKS[color^1][square]&Bitboards[pawn] != 0
}

func isIsolated(color int, square int) bool {
	pawn := WHITE_PAWN
	if color == BLACK {
		pawn = BLACK_PAWN
	}

	return Bitboards[pawn]&isolatedPawns[square] != 0
}

func isDoubled(color int, square int) bool {
	pawn := WHITE_PAWN
	if color == BLACK {
		pawn = BLACK_PAWN
	}

	return Bitboards[pawn]&doubledPawns[square] != 0
}

func isPassed(color int, square int) bool {
	pawn := WHITE_PAWN
	if color == BLACK {
		pawn = BLACK_PAWN
	}

	return Bitboards[pawn]&passedPawns[color][square] != 0
}

func knightBonus(color int, square int) int {
	value := 0

	moves := KNIGHT_ATTACKS[square] & ^Occupancies[color]

	pawn := WHITE_PAWN
	enemyPawn := BLACK_PAWN
	enemyKing := BLACK_KING
	if color == BLACK {
		pawn = BLACK_PAWN
		enemyPawn = WHITE_PAWN
		enemyKing = WHITE_KING
	}

	// square not attacked by enemy pawn pawn and square is protected by friendly pawn
	if outposts[color][square]&Bitboards[enemyPawn] == 0 && PAWN_ATTACKS[color^1][square]&Bitboards[pawn] != 0 {
		if color == WHITE {
			value += KNIGHT_OUTPOST_WEIGHTS[square]
		} else {
			value += KNIGHT_OUTPOST_WEIGHTS[square^56]
		}
	}

	// add mobility weight, capture weight, and threat to king weight
	value += CountBits(moves)*KNIGHT_MOBILITY_WEIGHT + CountBits(moves&Occupancies[color^1])*CAPTURE_WEIGHT + CountBits(moves&Bitboards[enemyKing])*KNIGHT_THREAT_WEIGHT

	if color == WHITE {
		return value
	}

	return -value
}

func bishopBonus(color int, square int) int {
	value := 0
	moves := GetBishopAttacks(square, Occupancies[BOTH])

	pawn := WHITE_PAWN
	bishop := WHITE_BISHOP
	enemyPawn := BLACK_PAWN
	enemyKing := BLACK_KING
	if color == BLACK {
		pawn = BLACK_PAWN
		bishop = BLACK_BISHOP
		enemyPawn = WHITE_PAWN
		enemyKing = WHITE_KING
	}

	// square not attacked by enemy pawn pawn and square is protected by friendly pawn
	if outposts[color][square]&Bitboards[enemyPawn] == 0 && PAWN_ATTACKS[color^1][square]&Bitboards[pawn] != 0 {
		if color == WHITE {
			value += BISHOP_OUTPOST_WEIGHTS[square]
		} else {
			value += BISHOP_OUTPOST_WEIGHTS[square^56]
		}
	}

	// increase eval if more than 1 bishop because double bishop attack is strong
	if CountBits(Bitboards[bishop]) > 1 {
		value += 35
	}

	// add mobility weight, capture weight, and threat to king weight
	value += CountBits(moves)*BISHOP_MOBILITY_WEIGHT + CountBits(moves&Occupancies[color^1])*CAPTURE_WEIGHT + CountBits(moves&Bitboards[enemyKing])*BISHOP_THREAT_WEIGHT

	if color == WHITE {
		return value
	}

	return -value
}

func rookBonus(color int, square int) int {
	moves := GetRookAttacks(square, Occupancies[BOTH])

	enemyKing := BLACK_KING
	if color == BLACK {
		enemyKing = WHITE_KING
	}

	// add mobility weight, capture weight, and threat to king weight
	value := CountBits(moves)*ROOK_MOBILITY_WEIGHT + CountBits(moves&Occupancies[color^1])*CAPTURE_WEIGHT + CountBits(moves&Bitboards[enemyKing])*ROOK_THREAT_WEIGHT

	if color == WHITE {
		return value
	}

	return -value
}

func queenBonus(color int, square int) int {
	moves := GetQueenAttacks(square, Occupancies[BOTH])

	enemyKing := BLACK_KING
	if color == BLACK {
		enemyKing = WHITE_KING
	}

	// add mobility weight, capture weight, and threat to king weight
	value := CountBits(moves)*QUEEN_MOBILITY_WEIGHT + CountBits(moves&Occupancies[color^1])*CAPTURE_WEIGHT + CountBits(moves&Bitboards[enemyKing])*QUEEN_THREAT_WEIGHT

	if color == WHITE {
		return value
	}

	return -value
}

func kingBonus(color int, square int) int {
	gamePhase := calculateGamePhase()

	if gamePhase > 24 {
		gamePhase = 24
	}

	midgamePhase := gamePhase
	endgamePhase := 24 - gamePhase

	// king safety (middlegame)
	safetyScore := 0

	// distance from center penalty
	safetyScore -= KING_DISTANCE_CENTER_WEIGHT * DISTANCE_FROM_CENTER[square]

	// evaluate pawn shield
	pawnShield := evaluatePawnShield(color, square)
	safetyScore += pawnShield

	// castling bonus
	if color == WHITE {
		if square == G1 || square == C1 {
			safetyScore += KING_CASTLED_BONUS
		}
	} else {
		if square == G8 || square == C8 {
			safetyScore += KING_CASTLED_BONUS
		}
	}

	enemyQueen := BLACK_QUEEN
	enemyRook := BLACK_ROOK
	enemyBishop := BLACK_BISHOP
	enemyKnight := BLACK_KNIGHT
	if color == BLACK {
		enemyQueen = WHITE_QUEEN
		enemyRook = WHITE_ROOK
		enemyBishop = WHITE_BISHOP
		enemyKnight = WHITE_KNIGHT
	}

	// king tropism - consider all enemy pieces
	queenTropism := 0
	rookTropism := 0
	bishopTropism := 0
	knightTropism := 0

	// queen tropism
	queenBitboard := Bitboards[enemyQueen]
	for queenBitboard > 0 {
		queenSquare := LS1B(queenBitboard)
		dist := manhattanDistance(square, queenSquare)
		queenTropism += QUEEN_TROPISM_RANGE - dist
		PopBit(&queenBitboard, queenSquare)
	}

	// rook tropism
	rookBitboard := Bitboards[enemyRook]
	for rookBitboard > 0 {
		rookSquare := LS1B(rookBitboard)
		dist := manhattanDistance(square, rookSquare)
		rookTropism += ROOK_TROPISM_RANGE - dist
		PopBit(&rookBitboard, rookSquare)
	}

	// bishop tropism
	bishopBitboard := Bitboards[enemyBishop]
	for bishopBitboard > 0 {
		bishopSquare := LS1B(bishopBitboard)
		dist := manhattanDistance(square, bishopSquare)
		bishopTropism += BISHOP_TROPISM_RANGE - dist
		PopBit(&bishopBitboard, bishopSquare)
	}

	// knight tropism
	knightBitboard := Bitboards[enemyKnight]
	for knightBitboard > 0 {
		knightSquare := LS1B(knightBitboard)
		dist := manhattanDistance(square, knightSquare)
		knightTropism += KNIGHT_TROPISM_RANGE - dist
		PopBit(&knightBitboard, knightSquare)
	}

	// apply tropism weights
	safetyScore -= queenTropism * QUEEN_TROPISM_WEIGHT
	safetyScore -= rookTropism * ROOK_TROPISM_WEIGHT
	safetyScore -= bishopTropism * BISHOP_TROPISM_WEIGHT
	safetyScore -= knightTropism * KNIGHT_TROPISM_WEIGHT

	// king mobility (negative in middlegame for safety)
	moves := KING_ATTACKS[square] & ^Occupancies[color]
	safetyScore -= CountBits(moves) * KING_MOBILITY_WEIGHT

	// king activity (endgame)
	activityScore := 0
	enemyKingSquare := LS1B(Bitboards[color^1+WHITE_KING])

	// king centralization and distance to enemy king
	activityScore -= DISTANCE_FROM_CENTER[square] * KING_CENTRALIZATION_WEIGHT
	activityScore -= manhattanDistance(square, enemyKingSquare) * KING_DISTANCE_ENEMY_WEIGHT

	// king mobility (positive in endgame for activity)
	activityScore += CountBits(moves) * KING_MOBILITY_WEIGHT

	value := (safetyScore*midgamePhase + activityScore*endgamePhase) / 24

	if color == WHITE {
		return value
	}
	return -value
}

func evaluatePawnShield(color int, square int) int {
	value := 0
	pawn := WHITE_PAWN
	if color == BLACK {
		pawn = BLACK_PAWN
	}

	// create shield mask for the king's file and adjacent files
	file := square % 8
	fileMask := uint64(0)
	if file > 0 {
		fileMask |= FILE_A << (file - 1)
	}
	fileMask |= FILE_A << file
	if file < 7 {
		fileMask |= FILE_A << (file + 1)
	}

	// get forward mask for the king's rank
	shieldMask := forwardMask(color, square) & fileMask

	// count pawns in shield area
	shieldPawns := Bitboards[pawn] & shieldMask
	shieldCount := CountBits(shieldPawns)

	// evaluate central pawns
	centralPawns := shieldPawns & CENTER_FILES_MASK
	value += shieldCount * PAWN_SHIELD_BASE_VALUE
	value += CountBits(centralPawns) * PAWN_SHIELD_CENTRAL_BONUS

	// bonus for complete shield
	if shieldCount == 3 {
		value += PAWN_SHIELD_COMPLETE_BONUS
	}

	return value
}

func manhattanDistance(sq1, sq2 int) int {
	file1 := sq1 % 8
	rank1 := sq1 / 8
	file2 := sq2 % 8
	rank2 := sq2 / 8
	return int(math.Abs(float64(file1-file2))) + int(math.Abs(float64(rank1-rank2)))
}

func calculateGamePhase() int {
	phase := 0
	phase += CountBits(Bitboards[WHITE_KNIGHT]|Bitboards[BLACK_KNIGHT]) * GAME_PHASE_INCREMENT[WHITE_KNIGHT]
	phase += CountBits(Bitboards[WHITE_BISHOP]|Bitboards[BLACK_BISHOP]) * GAME_PHASE_INCREMENT[WHITE_BISHOP]
	phase += CountBits(Bitboards[WHITE_ROOK]|Bitboards[BLACK_ROOK]) * GAME_PHASE_INCREMENT[WHITE_ROOK]
	phase += CountBits(Bitboards[WHITE_QUEEN]|Bitboards[BLACK_QUEEN]) * GAME_PHASE_INCREMENT[WHITE_QUEEN]
	return phase
}
