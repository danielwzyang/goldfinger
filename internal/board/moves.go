package board

import "fmt"

const ALL_MOVES = 0
const ONLY_CAPTURES = 1

var LastCapture = -1

type MoveList struct {
	Moves [256]int
	Count int
}

func (moveList *MoveList) ContainsMove(move int) bool {
	for i := 0; i < moveList.Count; i++ {
		if moveList.Moves[i] == move {
			return true
		}
	}

	return false
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

func GetPieceOnSquare(square int) int {
	for i := WHITE_PAWN; i <= BLACK_KING; i++ {
		if GetBit(Bitboards[i], square) != 0 {
			return i
		}
	}
	return -1
}

// returns true if legal, false if not
func MakeMove(move int) bool {
	SaveState()

	// store current position in repetition table
	RepetitionIndex++
	RepetitionTable[RepetitionIndex] = ZobristHash

	source := GetSource(move)
	target := GetTarget(move)
	piece := GetPiece(move)
	promotion := GetPromotion(move)
	capture := GetCapture(move) > 0
	enpassant := GetEnPassant(move) > 0
	castling := GetCastling(move) > 0
	double := GetDouble(move) > 0

	// move piece
	PopBit(&Bitboards[piece], source)
	SetBit(&Bitboards[piece], target)

	// hash piece
	ZobristHash ^= PIECE_HASH[piece][source] // remove piece from source square
	ZobristHash ^= PIECE_HASH[piece][target] // set piece to target square

	Fifty++

	if piece == WHITE_PAWN || piece == BLACK_PAWN {
		Fifty = 0
	}

	if capture {
		Fifty = 0

		// set piece range to see which piece is being captured
		var start, end int
		if Side == WHITE {
			start = BLACK_PAWN
			end = BLACK_KING
		} else {
			start = WHITE_PAWN
			end = WHITE_KING
		}

		for i := start; i <= end; i++ {
			// capture found, pop bit
			if GetBit(Bitboards[i], target) != 0 {
				PopBit(&Bitboards[i], target)
				ZobristHash ^= PIECE_HASH[i][target]
				LastCapture = i
				StateStack[StateSize-1].LastCapture = int8(i)
				break
			}
		}
	} else {
		LastCapture = -1
	}

	if promotion > 0 {
		// pop pawn
		PopBit(&Bitboards[piece], target)
		ZobristHash ^= PIECE_HASH[piece][target]

		// add promoted piece
		SetBit(&Bitboards[promotion], target)
		ZobristHash ^= PIECE_HASH[promotion][target]
	}

	if enpassant {
		// pop captured pawn
		if Side == WHITE {
			PopBit(&Bitboards[BLACK_PAWN], target-8)
			ZobristHash ^= PIECE_HASH[BLACK_PAWN][target-8]
		} else {
			PopBit(&Bitboards[WHITE_PAWN], target+8)
			ZobristHash ^= PIECE_HASH[WHITE_PAWN][target+8]
		}
	}

	if EnPassant != INVALID_SQUARE {
		ZobristHash ^= ENPASSANT_HASH[EnPassant%8]
	}

	EnPassant = INVALID_SQUARE

	if double {
		// set enpassant target on double pawn push
		if Side == WHITE {
			EnPassant = target - 8
			ZobristHash ^= ENPASSANT_HASH[(target-8)%8]
		} else {
			EnPassant = target + 8
			ZobristHash ^= ENPASSANT_HASH[(target+8)%8]
		}
	}

	if castling {
		switch target {
		// white kingside
		case G1:
			PopBit(&Bitboards[WHITE_ROOK], H1)
			SetBit(&Bitboards[WHITE_ROOK], F1)
			ZobristHash ^= PIECE_HASH[WHITE_ROOK][H1]
			ZobristHash ^= PIECE_HASH[WHITE_ROOK][F1]
		// white queenside
		case C1:
			PopBit(&Bitboards[WHITE_ROOK], A1)
			SetBit(&Bitboards[WHITE_ROOK], D1)
			ZobristHash ^= PIECE_HASH[WHITE_ROOK][A1]
			ZobristHash ^= PIECE_HASH[WHITE_ROOK][D1]
		// black kingside
		case G8:
			PopBit(&Bitboards[BLACK_ROOK], H8)
			SetBit(&Bitboards[BLACK_ROOK], F8)
			ZobristHash ^= PIECE_HASH[BLACK_ROOK][H8]
			ZobristHash ^= PIECE_HASH[BLACK_ROOK][F8]
		// black queenside
		case C8:
			PopBit(&Bitboards[BLACK_ROOK], A8)
			SetBit(&Bitboards[BLACK_ROOK], D8)
			ZobristHash ^= PIECE_HASH[BLACK_ROOK][A8]
			ZobristHash ^= PIECE_HASH[BLACK_ROOK][D8]
		}
	}

	ZobristHash ^= CASTLE_HASH[Castle]

	// update castling rights
	Castle &= CASTLING_RIGHTS[source]
	Castle &= CASTLING_RIGHTS[target]

	ZobristHash ^= CASTLE_HASH[Castle]

	// update occupancies in O(1)
	sourceBit := uint64(1) << source
	targetBit := uint64(1) << target
	Occupancies[Side] = (Occupancies[Side] & ^sourceBit) | targetBit
	Occupancies[BOTH] = (Occupancies[BOTH] & ^sourceBit) | targetBit

	if capture && !enpassant {
		Occupancies[Side^1] &= ^targetBit
	}

	if enpassant {
		var enPassantSquare int
		if Side == WHITE {
			enPassantSquare = target - 8
		} else {
			enPassantSquare = target + 8
		}
		enPassantBit := uint64(1) << enPassantSquare
		Occupancies[Side^1] &= ^enPassantBit
		Occupancies[BOTH] &= ^enPassantBit
	}

	if castling {
		var rookSource, rookTarget int
		switch target {
		case G1:
			rookSource, rookTarget = H1, F1
		case C1:
			rookSource, rookTarget = A1, D1
		case G8:
			rookSource, rookTarget = H8, F8
		case C8:
			rookSource, rookTarget = A8, D8
		}
		rookSourceBit := uint64(1) << rookSource
		rookTargetBit := uint64(1) << rookTarget
		Occupancies[Side] = (Occupancies[Side] & ^rookSourceBit) | rookTargetBit
		Occupancies[BOTH] = (Occupancies[BOTH] & ^rookSourceBit) | rookTargetBit
	}

	// change side
	Side ^= 1

	// hash side
	ZobristHash ^= SIDE_HASH

	// to prevent pseudo legal moves from being made (i.e. still in check)
	var king int
	if Side == WHITE {
		king = LS1B(Bitboards[BLACK_KING])
	} else {
		king = LS1B(Bitboards[WHITE_KING])
	}

	// illegal move, return false
	if IsSquareAttacked(king, Side) {
		UndoMove(move)
		return false
	}

	return true
}

func MakeNullMove() {
	SaveState()

	// store current position in repetition table
	RepetitionIndex++
	RepetitionTable[RepetitionIndex] = ZobristHash

	if EnPassant != INVALID_SQUARE {
		ZobristHash ^= ENPASSANT_HASH[EnPassant%8]
	}
	EnPassant = INVALID_SQUARE

	Side ^= 1
	ZobristHash ^= SIDE_HASH
}

func UndoMove(move int) {
	StateSize--
	state := StateStack[StateSize]
	Side ^= 1

	source := GetSource(move)
	target := GetTarget(move)
	piece := GetPiece(move)
	promotion := GetPromotion(move)
	capture := GetCapture(move) > 0
	enpassant := GetEnPassant(move) > 0
	castling := GetCastling(move) > 0

	// undo promotion
	if promotion > 0 {
		PopBit(&Bitboards[promotion], target)
		SetBit(&Bitboards[piece], target)
	}

	// move piece back
	PopBit(&Bitboards[piece], target)
	SetBit(&Bitboards[piece], source)

	// restore captured piece
	if capture && !enpassant {
		SetBit(&Bitboards[state.LastCapture], target)
	}

	// restore enpassant piece
	if enpassant {
		if Side == WHITE {
			SetBit(&Bitboards[BLACK_PAWN], target-8)
		} else {
			SetBit(&Bitboards[WHITE_PAWN], target+8)
		}
	}

	// reverse rook castle
	if castling {
		switch target {
		case G1:
			PopBit(&Bitboards[WHITE_ROOK], F1)
			SetBit(&Bitboards[WHITE_ROOK], H1)
		case C1:
			PopBit(&Bitboards[WHITE_ROOK], D1)
			SetBit(&Bitboards[WHITE_ROOK], A1)
		case G8:
			PopBit(&Bitboards[BLACK_ROOK], F8)
			SetBit(&Bitboards[BLACK_ROOK], H8)
		case C8:
			PopBit(&Bitboards[BLACK_ROOK], D8)
			SetBit(&Bitboards[BLACK_ROOK], A8)
		}
	}

	EnPassant = int(state.EnPassant)
	Castle = int(state.Castle)
	ZobristHash = state.ZobristHash
	Fifty = int(state.Fifty)
	RepetitionIndex = int(state.RepetitionIndex)

	// update occupancies in O(1)
	sourceBit := uint64(1) << source
	targetBit := uint64(1) << target
	Occupancies[Side] = (Occupancies[Side] & ^targetBit) | sourceBit
	Occupancies[BOTH] = (Occupancies[BOTH] & ^targetBit) | sourceBit

	if capture && !enpassant {
		Occupancies[Side^1] |= targetBit
		Occupancies[BOTH] |= targetBit
	}

	if enpassant {
		var enPassantSquare int
		if Side == WHITE {
			enPassantSquare = target - 8
		} else {
			enPassantSquare = target + 8
		}
		enPassantBit := uint64(1) << enPassantSquare
		Occupancies[Side^1] |= enPassantBit
		Occupancies[BOTH] |= enPassantBit
	}

	if castling {
		var rookSource, rookTarget int
		switch target {
		case G1:
			rookSource, rookTarget = H1, F1
		case C1:
			rookSource, rookTarget = A1, D1
		case G8:
			rookSource, rookTarget = H8, F8
		case C8:
			rookSource, rookTarget = A8, D8
		}
		rookSourceBit := uint64(1) << rookSource
		rookTargetBit := uint64(1) << rookTarget
		Occupancies[Side] = (Occupancies[Side] & ^rookTargetBit) | rookSourceBit
		Occupancies[BOTH] = (Occupancies[BOTH] & ^rookTargetBit) | rookSourceBit
	}
}

func UndoNullMove() {
	StateSize--
	state := StateStack[StateSize]

	EnPassant = int(state.EnPassant)
	Side ^= 1
	ZobristHash = state.ZobristHash
	RepetitionIndex = int(state.RepetitionIndex)
}

func StringToPos(input string) int {
	file := int(input[0] - 'a')
	rank := int(input[1] - '1')
	return rank*8 + file
}

func MoveToString(move int) string {
	s := POSITIONTOSTRING[GetSource(move)] + POSITIONTOSTRING[GetTarget(move)]

	// handle promotions
	promotion := GetPromotion(move)
	if promotion != 0 {
		var promo byte
		switch promotion % 6 {
		case 1:
			promo = 'n'
		case 2:
			promo = 'b'
		case 3:
			promo = 'r'
		case 4:
			promo = 'q'
		}

		return s + string(promo)
	}

	return s
}

func StringToMove(input string) int {
	source := StringToPos(input[:2])
	target := StringToPos(input[2:4])

	piece := 0
	start := WHITE_PAWN
	end := WHITE_KING
	if Side == BLACK {
		start = BLACK_PAWN
		end = BLACK_KING
	}

	// find piece that's moving
	for i := start; i <= end; i++ {
		if GetBit(Bitboards[i], source) > 0 {
			piece = i
			break
		}
	}

	// promotion
	promotion := 0
	if len(input) == 5 {
		switch input[4] {
		case 'n':
			promotion = WHITE_KNIGHT
		case 'b':
			promotion = WHITE_BISHOP
		case 'r':
			promotion = WHITE_ROOK
		case 'q':
			promotion = WHITE_QUEEN
		}
		if Side == BLACK {
			promotion += 6
		}
	}

	// double pawn push
	double := 0
	if Side == WHITE && piece == WHITE_PAWN && target-source == 16 {
		double = 1
	} else if Side == BLACK && piece == BLACK_PAWN && source-target == 16 {
		double = 1
	}

	// capture + en passant
	capture := 0
	enpass := 0
	if target == EnPassant && (piece == WHITE_PAWN || piece == BLACK_PAWN) {
		capture = 1
		enpass = 1
	} else if GetBit(Occupancies[Side^1], target) > 0 {
		capture = 1
	}

	// castling
	castle := 0
	if Side == WHITE && piece == WHITE_KING {
		if source == E1 && (target == G1 || target == C1) {
			castle = 1
		}
	}
	if Side == BLACK && piece == BLACK_KING {
		if source == E8 && (target == G8 || target == C8) {
			castle = 1
		}
	}

	return EncodeMove(source, target, piece, promotion, capture, double, enpass, castle)
}

func PrintMove(move int) {
	source := POSITIONTOSTRING[GetSource(move)]
	target := POSITIONTOSTRING[GetTarget(move)]
	piece := ascii[GetPiece(move)+1]
	promotion := min(1, GetPromotion(move))
	capture := min(1, GetCapture(move))
	double := min(1, GetDouble(move))
	enpassant := min(1, GetEnPassant(move))
	castling := min(1, GetCastling(move))

	fmt.Printf("%s%s, %s, promotion: %d, capture: %d, double: %d, enpassant: %d, castling: %d\n", source, target, piece, promotion, capture, double, enpassant, castling)
}
