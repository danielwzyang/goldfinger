package board

import (
	"regexp"
	"strings"
)

var CharPieces = map[byte]int{
	'P': WHITE_PAWN,
	'N': WHITE_KNIGHT,
	'B': WHITE_BISHOP,
	'R': WHITE_ROOK,
	'Q': WHITE_QUEEN,
	'K': WHITE_KING,
	'p': BLACK_PAWN,
	'n': BLACK_KNIGHT,
	'b': BLACK_BISHOP,
	'r': BLACK_ROOK,
	'q': BLACK_QUEEN,
	'k': BLACK_KING,
}

var regexFEN = regexp.MustCompile(`^([rnbqkpRNBQKP1-8]{1,8}/){7}[rnbqkpRNBQKP1-8]{1,8}\s[bw]\s(-|[KQkq]{1,4})\s(-|[a-h][36])`)

var (
	DEFAULT_BOARD = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	SEARCH_TESTER = "rnbqkb1r/pp1p1pPp/8/2p1pP2/1P1P4/3P3P/P1P1P3/RNBQKBNR w KQkq e6 0 1"
)

func validFEN(fen string) bool {
	return regexFEN.MatchString(fen)
}

func ParseFEN(fen string) {
	// validate fen
	if !validFEN(fen) {
		panic("Invalid FEN format.")
	}

	// clear states
	for i := 0; i < 12; i++ {
		Bitboards[i] = 0
	}
	for i := 0; i < 3; i++ {
		Occupancies[i] = 0
	}

	Castle = 0
	EnPassant = INVALID_SQUARE
	Fifty = 0

	// break up into 4 parts
	parts := strings.Split(fen, " ")

	// piece placement
	rank := 7
	file := 0
	for i := 0; i < len(parts[0]); i++ {
		char := parts[0][i]

		if char == '/' {
			rank--
			file = 0
			continue
		}

		if char >= '1' && char <= '8' {
			file += int(char - '0')
			continue
		}

		square := rank*8 + file
		Bitboards[CharPieces[char]] |= 1 << square
		file++
	}

	// side to move
	if parts[1] == "w" {
		Side = WHITE
	} else {
		Side = BLACK
	}

	// castling rights
	castling := parts[2]
	if strings.Contains(castling, "K") {
		Castle |= WK
	}
	if strings.Contains(castling, "Q") {
		Castle |= WQ
	}
	if strings.Contains(castling, "k") {
		Castle |= BK
	}
	if strings.Contains(castling, "q") {
		Castle |= BQ
	}

	// en passant
	if parts[3] != "-" {
		EnPassant = StringToPos(parts[3])
	} else {
		EnPassant = INVALID_SQUARE
	}

	// occupancies
	Occupancies[WHITE] = Bitboards[WHITE_PAWN] |
		Bitboards[WHITE_KNIGHT] |
		Bitboards[WHITE_BISHOP] |
		Bitboards[WHITE_ROOK] |
		Bitboards[WHITE_QUEEN] |
		Bitboards[WHITE_KING]

	Occupancies[BLACK] = Bitboards[BLACK_PAWN] |
		Bitboards[BLACK_KNIGHT] |
		Bitboards[BLACK_BISHOP] |
		Bitboards[BLACK_ROOK] |
		Bitboards[BLACK_QUEEN] |
		Bitboards[BLACK_KING]

	Occupancies[BOTH] = Occupancies[WHITE] | Occupancies[BLACK]

	Fifty = 0
	RepetitionIndex = 0

	RepetitionTable[0] = ZobristHash

	// rehash
	ZobristHash = 0

	for piece := 0; piece < 12; piece++ {
		bitboard := Bitboards[piece]
		for square := 0; square < 64; square++ {
			if (bitboard & (1 << square)) != 0 {
				ZobristHash ^= PIECE_HASH[piece][square]
			}
		}
	}

	ZobristHash ^= CASTLE_HASH[Castle]

	if EnPassant != INVALID_SQUARE {
		ZobristHash ^= ENPASSANT_HASH[EnPassant%8]
	}

	if Side == BLACK {
		ZobristHash ^= SIDE_HASH
	}
}
