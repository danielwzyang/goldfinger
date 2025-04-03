package board

import (
	"fmt"
	"math"
)

var (
	Board        [8][8]string
	DefaultBoard = [8][8]string{
		{"bR", "bN", "bB", "bQ", "bK", "bB", "bN", "bR"},
		{"bP", "bP", "bP", "bP", "bP", "bP", "bP", "bP"},
		{" ", " ", " ", " ", " ", " ", " ", " "},
		{" ", "wP", " ", "wP", " ", "wP", " ", " "},
		{" ", " ", " ", " ", " ", " ", " ", " "},
		{" ", " ", " ", " ", " ", " ", " ", " "},
		{"wP", "wP", "wP", "wP", "wP", "wP", "wP", "wP"},
		{"wR", "wN", "wB", "wQ", "wK", "wB", "wN", "wR"},
	}

	blackKing [2]int
	whiteKing [2]int
	ascii     = map[string]string{
		"wK": "♔",
		"wQ": "♕",
		"wR": "♖",
		"wB": "♗",
		"wN": "♘",
		"wP": "♙",

		"bK": "♚",
		"bQ": "♛",
		"bR": "♜",
		"bB": "♝",
		"bN": "♞",
		"bP": "♟",

		" ": " ",
	}

	bCastleKS = false // black can castle kingside until rook or king moves
	bCastleQS = false // black can castle queenside until rook or king moves
	wCastleKS = false // white can castle kingside until rook or king moves
	wCastleQS = false // white can castle queenside until rook or king moves

	enPassant = [2]int{-10, -10} // set to the position that a pawn can move to for en passant capturing
)

func Init(board [8][8]string) {
	Board = board

	for r, row := range Board {
		for c, piece := range row {
			if piece == "wK" {
				whiteKing = [2]int{r, c}
			}
			if piece == "bK" {
				blackKing = [2]int{r, c}
			}
		}
	}

	// castling
	// for white the king has to be at 7 4
	if whiteKing[0] == 7 && whiteKing[1] == 4 {
		// for kingside the white rook has to be at 7 7
		wCastleKS = Board[7][7] == "wR"

		// for queenside the white rook has to be at 7 0
		wCastleQS = Board[7][0] == "wR"
	}

	// for black the king has to be at 0 4
	if blackKing[0] == 0 && blackKing[1] == 4 {
		// for kingside the black rook has to be at 0 7
		bCastleKS = Board[0][7] == "bR"
		// for queenside the black rook has to be at 0 0
		bCastleQS = Board[0][0] == "bR"
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

		for _, char := range row {
			// pieces
			fmt.Printf("%s\ufe0e %s│%s", ascii[char], gray, reset)
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

func KingsideCastle(color byte) {
	row := 7
	if color == 'b' {
		row = 0
	}

	MakeMove(row, 4, row, 6) // move king to g
	MakeMove(row, 7, row, 5) // move rook to f
}

func QueensideCastle(color byte) {
	row := 7
	if color == 'b' {
		row = 0
	}

	MakeMove(row, 4, row, 2) // move king to c
	MakeMove(row, 0, row, 3) // move rook to d
}

func MakeMove(r1 int, c1 int, r2 int, c2 int) {
	// update new position
	piece := Board[r1][c1]
	Board[r2][c2] = piece
	Board[r1][c1] = " "

	// if moved king, update stored king position and invalidate castling
	if piece[1] == 'K' {
		if piece[0] == 'b' {
			blackKing = [2]int{r2, c2}
			bCastleKS = false
			bCastleQS = false
		} else {
			whiteKing = [2]int{r2, c2}
			wCastleKS = false
			wCastleQS = false
		}
	}

	// if rook, invalidate castling for that side
	if piece[1] == 'R' {
		if piece[0] == 'b' {
			if c1 == 7 {
				bCastleKS = false
			}
			if c1 == 0 {
				bCastleQS = false
			}
		} else {
			if c1 == 7 {
				wCastleKS = false
			}
			if c1 == 0 {
				wCastleQS = false
			}
		}
	}

	// if pawn handle en passant related things
	if piece[1] == 'P' {
		// pawn made double move
		if math.Abs(float64(r2-r1)) == 2 {
			// set en passant target to the space skipped over
			enPassant = [2]int{(r1 + r2) / 2, c2}
			return
		} else if math.Abs(float64(r2-r1)) == 1 && math.Abs(float64(c2-c1)) == 1 && r2 == enPassant[0] && c2 == enPassant[1] {
			// capturing en passant
			// moving diagonally so distance for rows and columns are both 1 and the space is the space for en passant

			var capturedPawnRow int
			if piece[0] == 'w' {
				// white pawn moves up so the captured pawn is one row below
				capturedPawnRow = r2 + 1
			} else {
				// black pawn moving down so the captured pawn is one row above
				capturedPawnRow = r2 - 1
			}

			// capture pawn
			Board[capturedPawnRow][c2] = " "

			// reset en passant target
			enPassant = [2]int{-10, -10}
			return
		}
	}

	// any move that isn't double pawn invalidates en passant
	enPassant = [2]int{-10, -10}
}
