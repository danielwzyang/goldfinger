package board

const ALL_MOVES = 0
const ONLY_CAPTURES = 1

type MoveList struct {
	Moves [256]int
	Count int
}

func (moveList *MoveList) AddMove(move int) {
	moveList.Moves[moveList.Count] = move
	moveList.Count++
}

func EncodeMove(source int, target int, piece int, promotion int, capture int, double int, enpassant int, castling int) int {
	return (source) | // 6 bits 2^6 = 64
		(target << 6) | // 6 bits
		(piece << 12) | // 4 bits 2^4 = 16 (12 pieces)
		(promotion << 16) | // 4 bits (12 pieces)
		(capture << 20) | // 1 bit (bool)
		(double << 21) | // 1 bit
		(enpassant << 22) | // 1 bit
		(castling << 23) // 1 bit
}

func GetSource(move int) int {
	return move & 0x3f
}

func GetTarget(move int) int {
	return (move & 0xfc0) >> 6
}

func GetPiece(move int) int {
	return (move & 0xf000) >> 12
}

func GetPromotion(move int) int {
	return (move & 0xf0000) >> 16
}

func GetCapture(move int) int {
	return move & 0x100000
}

func GetDouble(move int) int {
	return move & 0x200000
}

func GetEnPassant(move int) int {
	return move & 0x400000
}

func GetCastling(move int) int {
	return move & 0x800000
}

// flag == 0 for all moves
// flag == 1 for captures only (e.g quiescence search)
// returns true if legal, false if not
func MakeMove(move int, flag int) bool {
	if flag == ALL_MOVES {
		source := GetSource(move)
		target := GetTarget(move)
		piece := GetPiece(move)
		promotion := GetPromotion(move)
		capture := GetCapture(move) > 0
		enpassant := GetEnPassant(move) > 0
		castling := GetCastling(move) > 0
		double := GetDouble(move) > 0

		// move bit
		SwapBit(&Bitboards[piece], source, target)

		Fifty++

		// reset 50 move draw
		if piece == WHITE_PAWN || piece == BLACK_PAWN {
			Fifty = 0
		}

		if capture {
			// reset 50 move draw
			Fifty = 0

			// set piece range to see which piece is being captured
			var start, end int
			if Side == WHITE {
				start = BLACK_PAWN
				end = BLACK_KING
			} else {
				start = WHITE_PAWN
				start = WHITE_KING
			}

			for i := start; i <= end; i++ {
				// capture found, pop bit
				if GetBit(Bitboards[i], target) != 0 {
					PopBit(&Bitboards[i], target)
					break
				}
			}
		}

		if promotion > 0 {
			// pop pawn
			PopBit(&Bitboards[piece], target)

			SetBit(&Bitboards[promotion], target)
		}

		if enpassant {
			// pop captured pawn
			if Side == WHITE {
				PopBit(&Bitboards[BLACK_PAWN], target-8)
			} else {
				PopBit(&Bitboards[WHITE_PAWN], target+8)
			}
		}

		EnPassant = INVALID_SQUARE

		if double {
			// set enpassant target on double pawn push
			if Side == WHITE {
				EnPassant = target - 8
			} else {
				EnPassant = target + 8
			}
		}

		if castling {
			switch target {
			// white kingside
			case G1:
				SwapBit(&Bitboards[WHITE_ROOK], H1, F1)
			// white queenside
			case C1:
				SwapBit(&Bitboards[WHITE_ROOK], A1, D1)
			// black kingside
			case G8:
				SwapBit(&Bitboards[BLACK_ROOK], H8, F8)
			// black queenside
			case C8:
				SwapBit(&Bitboards[BLACK_ROOK], A8, D8)
			}
		}

		Castle &= CASTLING_RIGHTS[source]
		Castle &= CASTLING_RIGHTS[target]

		// occupancies
		Occupancies = [3]uint64{0, 0, 0}

		for i := WHITE_PAWN; i <= WHITE_KING; i++ {
			Occupancies[WHITE] |= Bitboards[i]
		}

		for i := BLACK_PAWN; i <= BLACK_KING; i++ {
			Occupancies[BLACK] |= Bitboards[i]
		}

		Occupancies[BOTH] = Occupancies[WHITE] | Occupancies[BLACK]

		// flip side
		Side ^= 1

		// to prevent pseudo legal moves from being made (i.e. still in check)
		var king int
		if Side == WHITE {
			king = WHITE_KING
		} else {
			king = BLACK_KING
		}

		// illegal move, return false
		if IsSquareAttacked(LS1B(Bitboards[king]), Side) {
			RestoreState()
			return false
		}
	}

	// flag is capture moves only
	if GetCapture(move) > 0 {
		return MakeMove(move, ALL_MOVES)
	}

	// not a capture + flag is captures only
	return false
}
