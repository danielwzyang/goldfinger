package board

import "math/rand/v2"

var (
	PIECE_HASH     [12][64]uint64 // [piece][position]
	CASTLE_HASH    [16]uint64
	ENPASSANT_HASH [8]uint64 // [file]
	SIDE_HASH      uint64    // only xor if black is moving

	ZobristHash uint64 // current state
)

func InitZobristTables() {
	// zobrist hashing uses psuedorandom numbers

	for i := 0; i < 12; i++ {
		for j := 0; j < 64; j++ {
			PIECE_HASH[i][j] = rand.Uint64()
		}
	}

	for i := 0; i < 16; i++ {
		CASTLE_HASH[i] = rand.Uint64()
	}

	for i := 0; i < 8; i++ {
		ENPASSANT_HASH[i] = rand.Uint64()
	}

	SIDE_HASH = rand.Uint64()
}
