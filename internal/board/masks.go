package board

var (
	FILE_A = uint64(0x0101010101010101)
	FILE_B = uint64(0x0202020202020202)
	FILE_C = uint64(0x0404040404040404)
	FILE_D = uint64(0x0808080808080808)
	FILE_E = uint64(0x1010101010101010)
	FILE_F = uint64(0x2020202020202020)
	FILE_G = uint64(0x4040404040404040)
	FILE_H = uint64(0x8080808080808080)

	FILE_AB = FILE_A | FILE_B
	FILE_GH = FILE_G | FILE_H

	RANK_1 = uint64(0x00000000000000FF)
	RANK_2 = uint64(0x000000000000FF00)
	RANK_3 = uint64(0x0000000000FF0000)
	RANK_4 = uint64(0x00000000FF000000)
	RANK_5 = uint64(0x000000FF00000000)
	RANK_6 = uint64(0x0000FF0000000000)
	RANK_7 = uint64(0x00FF000000000000)
	RANK_8 = uint64(0xFF00000000000000)

	CENTER_FILES_MASK = FILE_C | FILE_D | FILE_E | FILE_F
)

func isolatedMask(square int) uint64 {
	file := square % 8

	// return adjacent files
	switch file {
	case 0:
		return FILE_B
	case 1:
		return FILE_A | FILE_C
	case 2:
		return FILE_B | FILE_D
	case 3:
		return FILE_C | FILE_E
	case 4:
		return FILE_D | FILE_F
	case 5:
		return FILE_E | FILE_G
	case 6:
		return FILE_F | FILE_H
	default:
		return FILE_G
	}
}

func doubledMask(square int) uint64 {
	file := square % 8

	// returns current file
	return FILE_A << file
}

func forwardMask(color, square int) uint64 {
	rank := square / 8

	var mask uint64

	if color == WHITE {
		for r := rank + 1; r < 8; r++ {
			mask |= RANK_1 << (r * 8)
		}
		return mask
	}

	for r := 0; r < rank; r++ {
		mask |= RANK_1 << (r * 8)
	}

	return mask
}

func passedMask(color, square int) uint64 {
	// squares that are in front of the pawn on the same and adjacent files
	return (isolatedPawns[square] | doubledPawns[square]) & forwardMask(color, square)
}

func outpostMask(color, square int) uint64 {
	return isolatedPawns[square] & forwardMask(color, square)
}
