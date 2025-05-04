package board

import "math/bits"

// for bitboard a1 (position 0) is the least significant bit

func CountBits(bitboard uint64) int {
	return bits.OnesCount64(bitboard)
}

func GetBit(bitboard uint64, index int) uint64 {
	return bitboard >> index & 1
}

func SetBit(bitboard *uint64, index int) {
	*bitboard |= (1 << index)
}

func PopBit(bitboard *uint64, index int) {
	*bitboard ^= (1 << index)
}

func SwapBit(bitboard *uint64, index1 int, index2 int) {
	*bitboard ^= (1 << index1) | (1 << index2)
}

func LS1B(bitboard uint64) int {
	for i := 0; i < 64; i++ {
		if GetBit(bitboard, i) == 1 {
			return i
		}
	}

	return -1
}

type HistoryState struct {
	LastMove Move

	BitboardA int // board of the piece that moved
	BitboardB int // board of the piece that was captured (-1 if no capture)

	WCastleKS bool
	WCastleQS bool
	BCastleKS bool
	BCastleQS bool

	EnPassant int
}

type Move struct {
	From int
	To   int
}
