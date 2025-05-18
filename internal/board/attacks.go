package board

var (
	PAWN_ATTACKS   = [2][64]uint64{}
	KNIGHT_ATTACKS = [64]uint64{}
	KING_ATTACKS   = [64]uint64{}

	// sliding pieces
	BISHOP_MASKS   = [64]uint64{}
	ROOK_MASKS     = [64]uint64{}
	BISHOP_ATTACKS = [64][512]uint64{}
	ROOK_ATTACKS   = [64][4096]uint64{}

	// RELEVANT OCCUPANCY BITS
	// generated from CountBits(MaskBishopAttacks()) i => [0-64)
	BISHOP_BITS = [64]uint64{
		6, 5, 5, 5, 5, 5, 5, 6,
		5, 5, 5, 5, 5, 5, 5, 5,
		5, 5, 7, 7, 7, 7, 5, 5,
		5, 5, 7, 9, 9, 7, 5, 5,
		5, 5, 7, 9, 9, 7, 5, 5,
		5, 5, 7, 7, 7, 7, 5, 5,
		5, 5, 5, 5, 5, 5, 5, 5,
		6, 5, 5, 5, 5, 5, 5, 6,
	}
	ROOK_BITS = [64]uint64{
		12, 11, 11, 11, 11, 11, 11, 12,
		11, 10, 10, 10, 10, 10, 10, 11,
		11, 10, 10, 10, 10, 10, 10, 11,
		11, 10, 10, 10, 10, 10, 10, 11,
		11, 10, 10, 10, 10, 10, 10, 11,
		11, 10, 10, 10, 10, 10, 10, 11,
		11, 10, 10, 10, 10, 10, 10, 11,
		12, 11, 11, 11, 11, 11, 11, 12,
	}
)

func InitNonSlidingAttacks() {
	for rank := 0; rank < 8; rank++ {
		for file := 0; file < 8; file++ {
			square := rank*8 + file

			PAWN_ATTACKS[WHITE][square] = MaskPawnAttacks(WHITE, square)
			PAWN_ATTACKS[BLACK][square] = MaskPawnAttacks(BLACK, square)

			KNIGHT_ATTACKS[square] = MaskKnightAttacks(square)

			KING_ATTACKS[square] = MaskKingAttacks(square)
		}
	}
}

func InitSlidingAttacks(bishop bool) {
	for square := 0; square < 64; square++ {
		BISHOP_MASKS[square] = MaskBishopAttacks(square)
		ROOK_MASKS[square] = MaskRookAttacks(square)

		var attackmask uint64
		if bishop {
			attackmask = BISHOP_MASKS[square]
		} else {
			attackmask = ROOK_MASKS[square]
		}

		relevantBitsCount := CountBits(attackmask)
		occupancyIndices := (1 << relevantBitsCount)

		for index := 0; index < occupancyIndices; index++ {
			occupancy := GetOccupancyMask(index, relevantBitsCount, attackmask)
			if bishop {
				magicIndex := (occupancy * BISHOP_MAGIC_NUMBERS[square]) >> (64 - BISHOP_BITS[square])

				BISHOP_ATTACKS[square][magicIndex] = PrecomputeBishopAttacks(square, occupancy)
			} else {
				magicIndex := (occupancy * ROOK_MAGIC_NUMBERS[square]) >> (64 - ROOK_BITS[square])

				ROOK_ATTACKS[square][magicIndex] = PrecomputeRookAttacks(square, occupancy)
			}
		}
	}
}

func MaskPawnAttacks(color int, square int) uint64 {
	var attacks uint64

	var bitboard uint64
	SetBit(&bitboard, square)

	if color == WHITE {
		attacks |= (bitboard & ^FILE_A) << 7 // capture left
		attacks |= (bitboard & ^FILE_H) << 9 // capture right
	} else {
		attacks |= (bitboard & ^FILE_A) >> 9 // capture left (from white's pov)
		attacks |= (bitboard & ^FILE_H) >> 7 // capture right (from white's pov)
	}

	return attacks
}

func MaskKnightAttacks(square int) uint64 {
	var attacks uint64

	var bitboard uint64
	SetBit(&bitboard, square)

	attacks |= (bitboard & ^FILE_AB) << 6  // capture left 2 up 1
	attacks |= (bitboard & ^FILE_A) << 15  // capture left 1 up 2
	attacks |= (bitboard & ^FILE_H) << 17  // capture right 1 up 2
	attacks |= (bitboard & ^FILE_GH) << 10 // capture right 2 up 1

	attacks |= (bitboard & ^FILE_GH) >> 6  // capture right 2 down 1
	attacks |= (bitboard & ^FILE_H) >> 15  // capture right 1 down 2
	attacks |= (bitboard & ^FILE_A) >> 17  // capture left 1 down 2
	attacks |= (bitboard & ^FILE_AB) >> 10 // capture left 2 down 1

	return attacks
}

func MaskKingAttacks(square int) uint64 {
	var attacks uint64

	var bitboard uint64
	SetBit(&bitboard, square)

	attacks |= (bitboard & ^FILE_A) << 7 // attack left 1 up 1
	attacks |= (bitboard & ^FILE_A) >> 1 // attack left 1
	attacks |= (bitboard & ^FILE_A) >> 9 // attack left 1 down 1

	attacks |= (bitboard & ^FILE_H) << 9 // attack right 1 up 1
	attacks |= (bitboard & ^FILE_H) << 1 // attack right 1
	attacks |= (bitboard & ^FILE_H) >> 7 // attack right 1 down 1

	attacks |= bitboard << 8 // attack up 1
	attacks |= bitboard >> 8 // attack down 1

	return attacks
}

func MaskBishopAttacks(square int) uint64 {
	var attacks uint64

	targetRank := square / 8
	targetFile := square % 8

	for rank, file := targetRank+1, targetFile+1; rank <= 6 && file <= 6; rank, file = rank+1, file+1 {
		attacks |= 1 << (rank*8 + file)
	}

	for rank, file := targetRank-1, targetFile+1; rank >= 1 && file <= 6; rank, file = rank-1, file+1 {
		attacks |= 1 << (rank*8 + file)
	}

	for rank, file := targetRank+1, targetFile-1; rank <= 6 && file >= 1; rank, file = rank+1, file-1 {
		attacks |= 1 << (rank*8 + file)
	}

	for rank, file := targetRank-1, targetFile-1; rank >= 1 && file >= 1; rank, file = rank-1, file-1 {
		attacks |= 1 << (rank*8 + file)
	}

	return attacks
}

func MaskRookAttacks(square int) uint64 {
	var attacks uint64

	targetRank := square / 8
	targetFile := square % 8

	for rank := targetRank + 1; rank <= 6; rank++ {
		attacks |= 1 << (rank*8 + targetFile)
	}

	for rank := targetRank - 1; rank >= 1; rank-- {
		attacks |= 1 << (rank*8 + targetFile)
	}

	for file := targetFile + 1; file <= 6; file++ {
		attacks |= 1 << (targetRank*8 + file)
	}

	for file := targetFile - 1; file >= 1; file-- {
		attacks |= 1 << (targetRank*8 + file)
	}

	return attacks
}

// ATTACK GENERATION FOR SLIDING PIECES
func PrecomputeBishopAttacks(square int, blocking uint64) uint64 {
	var attacks uint64

	targetRank := square / 8
	targetFile := square % 8

	for rank, file := targetRank+1, targetFile+1; rank <= 7 && file <= 7; rank, file = rank+1, file+1 {
		attacks |= 1 << (rank*8 + file)
		if GetBit(blocking, rank*8+file) == 1 {
			break
		}
	}

	for rank, file := targetRank-1, targetFile+1; rank >= 0 && file <= 7; rank, file = rank-1, file+1 {
		attacks |= 1 << (rank*8 + file)
		if GetBit(blocking, rank*8+file) == 1 {
			break
		}
	}

	for rank, file := targetRank+1, targetFile-1; rank <= 7 && file >= 0; rank, file = rank+1, file-1 {
		attacks |= 1 << (rank*8 + file)
		if GetBit(blocking, rank*8+file) == 1 {
			break
		}
	}

	for rank, file := targetRank-1, targetFile-1; rank >= 0 && file >= 0; rank, file = rank-1, file-1 {
		attacks |= 1 << (rank*8 + file)
		if GetBit(blocking, rank*8+file) == 1 {
			break
		}
	}

	return attacks
}

func PrecomputeRookAttacks(square int, blocking uint64) uint64 {
	var attacks uint64

	targetRank := square / 8
	targetFile := square % 8

	for rank := targetRank + 1; rank <= 7; rank++ {
		attacks |= 1 << (rank*8 + targetFile)
		if GetBit(blocking, rank*8+targetFile) == 1 {
			break
		}
	}

	for rank := targetRank - 1; rank >= 0; rank-- {
		attacks |= 1 << (rank*8 + targetFile)
		if GetBit(blocking, rank*8+targetFile) == 1 {
			break
		}
	}

	for file := targetFile + 1; file <= 7; file++ {
		attacks |= 1 << (targetRank*8 + file)
		if GetBit(blocking, targetRank*8+file) == 1 {
			break
		}
	}

	for file := targetFile - 1; file >= 0; file-- {
		attacks |= 1 << (targetRank*8 + file)
		if GetBit(blocking, targetRank*8+file) == 1 {
			break
		}
	}

	return attacks
}

func GetOccupancyMask(index int, bitCount int, mask uint64) uint64 {
	occupancy := uint64(0)

	for i := 0; i < bitCount; i++ {
		square := LS1B(mask)
		PopBit(&mask, square)

		// occupancy is on board (if out of bounds bit will be shifted off)
		if index&(1<<i) != 0 {
			occupancy |= 1 << square
		}
	}

	return occupancy
}

func GetBishopAttacks(square int, occupancy uint64) uint64 {
	occupancy &= BISHOP_MASKS[square]
	occupancy *= BISHOP_MAGIC_NUMBERS[square]
	occupancy >>= 64 - BISHOP_BITS[square]

	return BISHOP_ATTACKS[square][occupancy]
}

func GetRookAttacks(square int, occupancy uint64) uint64 {
	occupancy &= ROOK_MASKS[square]
	occupancy *= ROOK_MAGIC_NUMBERS[square]
	occupancy >>= 64 - ROOK_BITS[square]

	return ROOK_ATTACKS[square][occupancy]
}

func GetQueenAttacks(square int, occupancy uint64) uint64 {
	return GetBishopAttacks(square, occupancy) | GetRookAttacks(square, occupancy)
}

func IsSquareAttacked(square int, attacker int) bool {
	// attacked by white pawn
	if attacker == WHITE && (PAWN_ATTACKS[BLACK][square]&Bitboards[WHITE_PAWN] != 0) {
		return true
	}

	// attacked by black pawn
	if attacker == BLACK && (PAWN_ATTACKS[WHITE][square]&Bitboards[BLACK_PAWN] != 0) {
		return true
	}

	// attacked by white knight
	if attacker == WHITE && (KNIGHT_ATTACKS[square]&Bitboards[WHITE_KNIGHT] != 0) {
		return true
	}

	// attacked by black knight
	if attacker == BLACK && (KNIGHT_ATTACKS[square]&Bitboards[BLACK_KNIGHT] != 0) {
		return true
	}

	// attacked by white king
	if attacker == WHITE && (KING_ATTACKS[square]&Bitboards[WHITE_KING] != 0) {
		return true
	}

	// attacked by black king
	if attacker == BLACK && (KING_ATTACKS[square]&Bitboards[BLACK_KING] != 0) {
		return true
	}

	// attacked by white bishop
	if attacker == WHITE && (GetBishopAttacks(square, Occupancies[BOTH])&Bitboards[WHITE_BISHOP] != 0) {
		return true
	}

	// attacked by black bishop
	if attacker == BLACK && (GetBishopAttacks(square, Occupancies[BOTH])&Bitboards[BLACK_BISHOP] != 0) {
		return true
	}

	// attacked by white rook
	if attacker == WHITE && (GetRookAttacks(square, Occupancies[BOTH])&Bitboards[WHITE_ROOK] != 0) {
		return true
	}

	// attacked by black rook
	if attacker == BLACK && (GetRookAttacks(square, Occupancies[BOTH])&Bitboards[BLACK_ROOK] != 0) {
		return true
	}

	// attacked by white queen
	if attacker == WHITE && (GetQueenAttacks(square, Occupancies[BOTH])&Bitboards[WHITE_QUEEN] != 0) {
		return true
	}

	// attacked by black queen
	if attacker == BLACK && (GetQueenAttacks(square, Occupancies[BOTH])&Bitboards[BLACK_QUEEN] != 0) {
		return true
	}

	return false
}
