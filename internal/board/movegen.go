package board

func GenerateAllMoves(moves *MoveList) {
	moves.Count = 0

	GeneratePawnMoves(moves)
	GenerateCastleMoves(moves)
	GenerateKnightMoves(moves)
	GenerateBishopMoves(moves)
	GenerateRookMoves(moves)
}

func GeneratePawnMoves(moves *MoveList) {
	var source, target int
	var bitboard, attacks uint64

	// white pawns
	if Side == WHITE {
		bitboard = Bitboards[WHITE_PAWN]

		for bitboard > 0 {
			source = LS1B(bitboard)
			target = source + 8

			// quiet moves
			if GetBit(Occupancies[BOTH], target) == 0 {
				// promotion
				if source >= A7 && source <= H7 {
					moves.AddMove(EncodeMove(source, target, WHITE_PAWN, WHITE_KNIGHT, 0, 0, 0, 0))
					moves.AddMove(EncodeMove(source, target, WHITE_PAWN, WHITE_BISHOP, 0, 0, 0, 0))
					moves.AddMove(EncodeMove(source, target, WHITE_PAWN, WHITE_ROOK, 0, 0, 0, 0))
					moves.AddMove(EncodeMove(source, target, WHITE_PAWN, WHITE_QUEEN, 0, 0, 0, 0))
				} else {
					// move forward one
					moves.AddMove(EncodeMove(source, target, WHITE_PAWN, 0, 0, 0, 0, 0))

					// double move
					if source >= A2 && source <= H2 && GetBit(Occupancies[BOTH], target+8) == 0 {
						moves.AddMove(EncodeMove(source, target+8, WHITE_PAWN, 0, 0, 1, 0, 0))
					}
				}
			}

			// captures
			attacks = PAWN_ATTACKS[Side][source] & Occupancies[BLACK]
			for attacks > 0 {
				target = LS1B(attacks)

				// promotion
				if source >= A7 && source <= H7 {
					moves.AddMove(EncodeMove(source, target, WHITE_PAWN, WHITE_KNIGHT, 1, 0, 0, 0))
					moves.AddMove(EncodeMove(source, target, WHITE_PAWN, WHITE_BISHOP, 1, 0, 0, 0))
					moves.AddMove(EncodeMove(source, target, WHITE_PAWN, WHITE_ROOK, 1, 0, 0, 0))
					moves.AddMove(EncodeMove(source, target, WHITE_PAWN, WHITE_QUEEN, 1, 0, 0, 0))
				} else {
					// normal capture no promotion
					moves.AddMove(EncodeMove(source, target, WHITE_PAWN, 0, 1, 0, 0, 0))
				}

				PopBit(&attacks, target)
			}

			// valid enpassant
			if EnPassant != INVALID_SQUARE && PAWN_ATTACKS[Side][source]&(1<<EnPassant) != 0 {
				moves.AddMove(EncodeMove(source, EnPassant, WHITE_PAWN, 0, 1, 0, 1, 0))
			}

			PopBit(&bitboard, source)
		}
		return
	}

	// black pawns
	bitboard = Bitboards[BLACK_PAWN]

	for bitboard > 0 {
		source = LS1B(bitboard)
		target = source - 8

		// quiet moves
		if GetBit(Occupancies[BOTH], target) == 0 {
			// promotion
			if source >= A2 && source <= H2 {
				moves.AddMove(EncodeMove(source, target, BLACK_PAWN, BLACK_KNIGHT, 0, 0, 0, 0))
				moves.AddMove(EncodeMove(source, target, BLACK_PAWN, BLACK_BISHOP, 0, 0, 0, 0))
				moves.AddMove(EncodeMove(source, target, BLACK_PAWN, BLACK_ROOK, 0, 0, 0, 0))
				moves.AddMove(EncodeMove(source, target, BLACK_PAWN, BLACK_QUEEN, 0, 0, 0, 0))
			} else {
				// move forward one
				moves.AddMove(EncodeMove(source, target, BLACK_PAWN, 0, 0, 0, 0, 0))

				// double move
				if source >= A7 && source <= H7 && GetBit(Occupancies[BOTH], target-8) == 0 {
					moves.AddMove(EncodeMove(source, target-8, BLACK_PAWN, 0, 0, 1, 0, 0))
				}
			}
		}

		// captures
		attacks = PAWN_ATTACKS[Side][source] & Occupancies[WHITE]
		for attacks > 0 {
			target = LS1B(attacks)

			// promotion
			if source >= A2 && source <= H2 {
				moves.AddMove(EncodeMove(source, target, BLACK_PAWN, BLACK_KNIGHT, 1, 0, 0, 0))
				moves.AddMove(EncodeMove(source, target, BLACK_PAWN, BLACK_BISHOP, 1, 0, 0, 0))
				moves.AddMove(EncodeMove(source, target, BLACK_PAWN, BLACK_ROOK, 1, 0, 0, 0))
				moves.AddMove(EncodeMove(source, target, BLACK_PAWN, BLACK_QUEEN, 1, 0, 0, 0))
			} else {
				// normal capture no promotion
				moves.AddMove(EncodeMove(source, target, BLACK_PAWN, 0, 1, 0, 0, 0))
			}

			PopBit(&attacks, target)
		}

		// valid enpassant
		if EnPassant != INVALID_SQUARE && PAWN_ATTACKS[Side][source]&(1<<EnPassant) != 0 {
			moves.AddMove(EncodeMove(source, EnPassant, BLACK_PAWN, 0, 1, 0, 1, 0))
		}

		PopBit(&bitboard, source)
	}
}

func GenerateCastleMoves(moves *MoveList) {
	// white castling
	if Side == WHITE {
		if Castle&WK != 0 && // can castle
			GetBit(Occupancies[BOTH], F1) == 0 && GetBit(Occupancies[BOTH], G1) == 0 && // squares are empty
			!IsSquareAttacked(E1, BLACK) && !IsSquareAttacked(F1, BLACK) { // king and f1 aren't in check
			moves.AddMove(EncodeMove(E1, G1, WHITE_KING, 0, 0, 0, 0, 1))
		}

		if Castle&WQ != 0 && // can castle
			GetBit(Occupancies[BOTH], D1) == 0 && GetBit(Occupancies[BOTH], C1) == 0 && GetBit(Occupancies[BOTH], B1) == 0 && // squares are empty
			!IsSquareAttacked(E1, BLACK) && !IsSquareAttacked(D1, BLACK) { // king and d1 aren't in check
			moves.AddMove(EncodeMove(E1, C1, WHITE_KING, 0, 0, 0, 0, 1))
		}
		return
	}

	// black castling
	if Castle&BK != 0 && // can castle
		GetBit(Occupancies[BOTH], F8) == 0 && GetBit(Occupancies[BOTH], G8) == 0 && // squares are empty
		!IsSquareAttacked(E8, BLACK) && !IsSquareAttacked(F8, BLACK) { // king and f8 aren't in check
		moves.AddMove(EncodeMove(E8, G8, WHITE_KING, 0, 0, 0, 0, 1))
	}

	if Castle&BQ != 0 && // can castle
		GetBit(Occupancies[BOTH], D8) == 0 && GetBit(Occupancies[BOTH], C8) == 0 && GetBit(Occupancies[BOTH], B8) == 0 && // squares are empty
		!IsSquareAttacked(E8, BLACK) && !IsSquareAttacked(D8, BLACK) { // king and d8 aren't in check
		moves.AddMove(EncodeMove(E8, C8, WHITE_KING, 0, 0, 0, 0, 1))
	}
}

func GenerateKnightMoves(moves *MoveList) {
	var source, target, piece int
	var bitboard, attacks uint64

	if Side == WHITE {
		piece = WHITE_KNIGHT
	} else {
		piece = BLACK_KNIGHT
	}

	bitboard = Bitboards[piece]

	for bitboard > 0 {
		source = LS1B(bitboard)

		attacks = KNIGHT_ATTACKS[source] & ^Occupancies[Side]

		for attacks > 0 {
			target = LS1B(attacks)

			// quiet move
			if GetBit(Occupancies[Side^1], target) == 0 {
				moves.AddMove(EncodeMove(source, target, piece, 0, 0, 0, 0, 0))
			} else {
				// capture
				moves.AddMove(EncodeMove(source, target, piece, 0, 1, 0, 0, 0))
			}

			PopBit(&attacks, target)
		}

		PopBit(&bitboard, source)
	}
}

func GenerateBishopMoves(moves *MoveList) {
	var source, target, piece int
	var bitboard, attacks uint64

	if Side == WHITE {
		piece = WHITE_BISHOP
	} else {
		piece = BLACK_BISHOP
	}

	bitboard = Bitboards[piece]

	for bitboard > 0 {
		source = LS1B(bitboard)

		attacks = GetBishopAttacks(source, Occupancies[BOTH]) & ^Occupancies[Side]

		for attacks > 0 {
			target = LS1B(attacks)

			// quiet move
			if GetBit(Occupancies[Side^1], target) == 0 {
				moves.AddMove(EncodeMove(source, target, piece, 0, 0, 0, 0, 0))
			} else {
				// capture
				moves.AddMove(EncodeMove(source, target, piece, 0, 1, 0, 0, 0))
			}

			PopBit(&attacks, target)
		}

		PopBit(&bitboard, source)
	}
}

func GenerateRookMoves(moves *MoveList) {
	var source, target, piece int
	var bitboard, attacks uint64

	if Side == WHITE {
		piece = WHITE_ROOK
	} else {
		piece = BLACK_ROOK
	}

	bitboard = Bitboards[piece]

	for bitboard > 0 {
		source = LS1B(bitboard)

		attacks = GetRookAttacks(source, Occupancies[BOTH]) & ^Occupancies[Side]

		for attacks > 0 {
			target = LS1B(attacks)

			// quiet move
			if GetBit(Occupancies[Side^1], target) == 0 {
				moves.AddMove(EncodeMove(source, target, piece, 0, 0, 0, 0, 0))
			} else {
				// capture
				moves.AddMove(EncodeMove(source, target, piece, 0, 1, 0, 0, 0))
			}

			PopBit(&attacks, target)
		}

		PopBit(&bitboard, source)
	}
}

func GenerateQueenMoves(moves *MoveList) {
	var source, target, piece int
	var bitboard, attacks uint64

	if Side == WHITE {
		piece = WHITE_QUEEN
	} else {
		piece = BLACK_QUEEN
	}

	bitboard = Bitboards[piece]

	for bitboard > 0 {
		source = LS1B(bitboard)

		attacks = GetQueenAttacks(source, Occupancies[BOTH]) & ^Occupancies[Side]

		for attacks > 0 {
			target = LS1B(attacks)

			// quiet move
			if GetBit(Occupancies[Side^1], target) == 0 {
				moves.AddMove(EncodeMove(source, target, piece, 0, 0, 0, 0, 0))
			} else {
				// capture
				moves.AddMove(EncodeMove(source, target, piece, 0, 1, 0, 0, 0))
			}

			PopBit(&attacks, target)
		}

		PopBit(&bitboard, source)
	}
}

func GenerateKingMoves(moves *MoveList) {
	var source, target, piece int
	var bitboard, attacks uint64

	if Side == WHITE {
		piece = WHITE_KING
	} else {
		piece = BLACK_KING
	}

	bitboard = Bitboards[piece]

	for bitboard > 0 {
		source = LS1B(bitboard)

		attacks = KING_ATTACKS[source] & ^Occupancies[Side]

		for attacks > 0 {
			target = LS1B(attacks)

			// quiet move
			if GetBit(Occupancies[Side^1], target) == 0 {
				moves.AddMove(EncodeMove(source, target, piece, 0, 0, 0, 0, 0))
			} else {
				// capture
				moves.AddMove(EncodeMove(source, target, piece, 0, 1, 0, 0, 0))
			}

			PopBit(&attacks, target)
		}

		PopBit(&bitboard, source)
	}
}
