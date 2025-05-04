package board

import "fmt"

var (
	Bitboards   [12]uint64
	WhitePieces uint64
	BlackPieces uint64

	BCastleKS = false // black can castle kingside until rook or king moves
	BCastleQS = false // black can castle queenside until rook or king moves
	WCastleKS = false // white can castle kingside until rook or king moves
	WCastleQS = false // white can castle queenside until rook or king moves

	EnPassant = INVALID_SQUARE // set to the position that a pawn can move to for en passant capturing
)

func Init(board [12]uint64) {
	InitNonSlidingAttacks()

	InitSlidingAttacks(true)  // bishops
	InitSlidingAttacks(false) // rooks

	Bitboards = board

	WhitePieces = Bitboards[WHITE_PAWN] |
		Bitboards[WHITE_KNIGHT] |
		Bitboards[WHITE_BISHOP] |
		Bitboards[WHITE_ROOK] |
		Bitboards[WHITE_QUEEN] |
		Bitboards[WHITE_KING]

	BlackPieces = Bitboards[BLACK_PAWN] |
		Bitboards[BLACK_KNIGHT] |
		Bitboards[BLACK_BISHOP] |
		Bitboards[BLACK_ROOK] |
		Bitboards[BLACK_QUEEN] |
		Bitboards[BLACK_KING]

	// castling
	// for white the king has to be at e1
	if LS1B(Bitboards[WHITE_KING]) == e1 {
		// for kingside the white rook has to be at h1
		WCastleKS = GetBit(Bitboards[WHITE_ROOK], h1) == 1

		// for queenside the white rook has to be at a1
		WCastleQS = GetBit(Bitboards[WHITE_ROOK], a1) == 1
	}

	// for black the king has to be at e8
	if LS1B(Bitboards[BLACK_KING]) == e8 {
		// for kingside the black rook has to be at h8
		BCastleKS = GetBit(Bitboards[BLACK_ROOK], h8) == 1

		// for queenside the black rook has to be at a8
		BCastleQS = GetBit(Bitboards[BLACK_ROOK], a8) == 1
	}
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
