package board

import "fmt"

var (
	Bitboards [12]uint64

	// occupancies
	Occupancies = [3]uint64{
		0, // WHITE
		0, // BLACK
		0, // BOTH
	}

	Castle = 0
	/*
		0001 = 1 = white kingside
		0010 = 2 = white queenside
		0100 = 4 = black kingside
		1000 = 8 = black queenside
	*/
	WK = 1
	WQ = 2
	BK = 4
	BQ = 8

	EnPassant = INVALID_SQUARE
	Side      = WHITE
	Fifty     = 0
)

func Init() {
	InitNonSlidingAttacks()

	InitSlidingAttacks(true)  // bishops
	InitSlidingAttacks(false) // rooks

	InitEvalTables()

	InitZobristTables()
}

func Print() {
	gray := "\033[2;37m"
	reset := "\033[0m"

	// top border
	fmt.Printf("%s   ┌──┬──┬──┬──┬──┬──┬──┬──┐%s\n", gray, reset)

	for rank := 7; rank >= 0; rank-- {
		// numbers on side
		fmt.Printf(" %d %s│%s", rank+1, gray, reset)

		for file := 0; file < 8; file++ {
			square := rank*8 + file
			occupied := false

			for piece, bitboard := range Bitboards {
				if GetBit(bitboard, square) == 1 {
					fmt.Printf("%s\ufe0e %s│%s", ascii[piece+1], gray, reset)
					occupied = true
					break
				}
			}

			if !occupied {
				fmt.Printf("%s\ufe0e %s│%s", ascii[0], gray, reset)
			}
		}

		fmt.Println()

		// middle borders or bottom border
		if rank > 0 {
			fmt.Printf("%s   ├──┼──┼──┼──┼──┼──┼──┼──┤%s\n", gray, reset)
		} else {
			fmt.Printf("%s   └──┴──┴──┴──┴──┴──┴──┴──┘%s\n", gray, reset)
		}
	}

	// letters on bottom
	fmt.Println("    a  b  c  d  e  f  g  h")
	fmt.Println()
}
