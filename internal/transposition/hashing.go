package transposition

import (
	"fmt"
	"math/rand"

	"danielyang.cc/chess/internal/board"
)

var (
	pieces    [64][12]uint64 // 64 positions, 12 types of pieces
	castling  [4]uint64      // 4 castling booleans
	enPassant [8]uint64      // 8 files
	whiteTurn uint64         // only xor if white is moving
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

func HashBoard(color int) uint64 {
	hash := uint64(0)

	for i, row := range board.Board {
		for j, piece := range row {
			if piece.Type != board.EMPTY {
				if piece.Key == 0 {
					fmt.Println(piece)
				}
				hash ^= pieces[i*8+j][piece.Key-1]
			}
		}
	}

	if color == board.WHITE {
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

	if board.EnPassant.File != -10 {
		// based on column
		hash ^= enPassant[board.EnPassant.File]
	}

	return hash
}
