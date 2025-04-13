package board

import "fmt"

var (
	Board [8][8]Piece

	BlackKing Position
	WhiteKing Position

	BCastleKS = true // black can castle kingside until rook or king moves
	BCastleQS = true // black can castle queenside until rook or king moves
	WCastleKS = true // white can castle kingside until rook or king moves
	WCastleQS = true // white can castle queenside until rook or king moves

	EnPassant = Position{-10, -10} // set to the position that a pawn can move to for en passant capturing

	BoardHistory  = []BoardState{}
	HistoryLength = 0
)

func Init(board [8][8]Piece) {
	Board = board

	for r, row := range Board {
		for c, piece := range row {
			if piece.Color == KING {
				if piece.Color == WHITE {
					WhiteKing = Position{r, c}
				}
				if piece.Color == BLACK {
					BlackKing = Position{r, c}
				}
			}
		}
	}

	// castling
	// for white the king has to be at 7 4
	if WhiteKing.Rank == 7 && WhiteKing.File == 4 {
		// for kingside the white rook has to be at 7 7
		WCastleKS = Board[7][7].Color == WHITE && Board[7][7].Type == ROOK

		// for queenside the white rook has to be at 7 0
		WCastleQS = Board[7][0].Color == WHITE && Board[7][0].Type == ROOK
	}

	// for black the king has to be at 0 4
	if BlackKing.Rank == 0 && BlackKing.File == 4 {
		// for kingside the black rook has to be at 0 7
		BCastleKS = Board[0][7].Color == BLACK && Board[0][7].Type == ROOK

		// for queenside the black rook has to be at 0 0
		BCastleQS = Board[0][0].Color == BLACK && Board[0][0].Type == ROOK
	}
}

func Print() {
	gray := "\033[2;37m"
	reset := "\033[0m"

	// top border
	fmt.Printf("%s   ┌──┬──┬──┬──┬──┬──┬──┬──┐%s\n", gray, reset)

	for i, row := range Board {
		// numbers on side
		fmt.Printf(" %d %s│%s", 8-i, gray, reset)

		for _, piece := range row {
			// pieces
			fmt.Printf("%s\ufe0e %s│%s", ascii[piece.Key], gray, reset)
		}
		fmt.Println()

		// middle borders or bottom border
		if i < 7 {
			fmt.Printf("%s   ├──┼──┼──┼──┼──┼──┼──┼──┤%s\n", gray, reset)
		} else {
			fmt.Printf("%s   └──┴──┴──┴──┴──┴──┴──┴──┘%s\n", gray, reset)
		}
	}

	// letters on bottom
	fmt.Println("    a  b  c  d  e  f  g  h")
	fmt.Println()
}

func Draw(color int) bool {
	_, n := GetAllValidMoves(color)
	return !InCheck(color) && n == 0 && !ValidKingSideCastle(color) && !ValidQueenSideCastle(color)
}

func Checkmate(color int) bool {
	_, n := GetAllValidMoves(color)
	return InCheck(color) && n == 0 && !ValidKingSideCastle(color) && !ValidQueenSideCastle(color)
}
