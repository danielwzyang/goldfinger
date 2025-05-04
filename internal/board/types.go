package board

import (
	"fmt"
	"math/bits"
)

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

func PrintBitboard(bitboard uint64) {
	for rank := 7; rank >= 0; rank-- {
		for file := 0; file < 8; file++ {
			square := rank*8 + file
			if GetBit(bitboard, square) == 1 {
				fmt.Print("1 ")
			} else {
				fmt.Print(". ")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

type HistoryState struct {
	LastMove Move

	BitboardA int // board of the piece that moved
	BitboardB int // board of the piece that was captured (-1 if no capture)

	Castle    int
	EnPassant int
}

type Move struct {
	From int
	To   int
}
