package transposition

import (
	"math/rand"

	"danielyang.cc/chess/internal/board"
)

var (
	pieces    [64][12]uint64 // 64 positions, 12 types of pieces
	castling  [4]uint64      // 4 castling booleans
	enPassant [8]uint64      // 8 files
	whiteTurn uint64         // only xor if white is moving

	pieceMap = map[string]int{
		"wP": 0,
		"wN": 1,
		"wB": 2,
		"wR": 3,
		"wQ": 4,
		"wK": 5,
		"bP": 6,
		"bN": 7,
		"bB": 8,
		"bR": 9,
		"bQ": 10,
		"bK": 11,
	}
)

func initZobrist() {
	// zobrist hashing uses psuedorandom numbers

	for i := 0; i < 64; i++ {
		for j := 0; j < 12; j++ {
			pieces[i][j] = rand.Uint64()
		}
	}

	for i := 0; i < 4; i++ {
		castling[i] = rand.Uint64()
	}

	for i := 0; i < 8; i++ {
		enPassant[i] = rand.Uint64()
	}

	whiteTurn = rand.Uint64()
}

func HashBoard(color byte) uint64 {
	hash := uint64(0)

	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if board.Board[i][j] != " " {
				hash ^= pieces[i*8+j][pieceMap[board.Board[i][j]]]
			}
		}
	}

	if color == 'w' {
		hash ^= whiteTurn
	}

	if board.WCastleKS {
		hash ^= castling[0]
	}

	if board.WCastleQS {
		hash ^= castling[1]
	}

	if board.BCastleKS {
		hash ^= castling[2]
	}

	if board.BCastleQS {
		hash ^= castling[3]
	}

	if board.EnPassant[0] != -10 {
		// based on column
		hash ^= enPassant[board.EnPassant[1]]
	}

	return hash
}
