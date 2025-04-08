package board

import (
	"fmt"
)

var (
	Board        [8][8]string
	DefaultBoard = [8][8]string{
		{"bR", "bN", "bB", "bQ", "bK", "bB", "bN", "bR"},
		{"bP", "bP", "bP", "bP", "bP", "bP", "bP", "bP"},
		{" ", " ", " ", " ", " ", " ", " ", " "},
		{" ", " ", " ", " ", " ", " ", " ", " "},
		{" ", " ", " ", " ", " ", " ", " ", " "},
		{" ", " ", " ", " ", " ", " ", " ", " "},
		{"wP", "wP", "wP", "wP", "wP", "wP", "wP", "wP"},
		{"wR", "wN", "wB", "wQ", "wK", "wB", "wN", "wR"},
	}

	BlackKing [2]int
	WhiteKing [2]int
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

	BCastleKS = false // black can castle kingside until rook or king moves
	BCastleQS = false // black can castle queenside until rook or king moves
	WCastleKS = false // white can castle kingside until rook or king moves
	WCastleQS = false // white can castle queenside until rook or king moves

	EnPassant = [2]int{-10, -10} // set to the position that a pawn can move to for en passant capturing
)

func Init(board [8][8]string) {
	Board = board

	for r, row := range Board {
		for c, piece := range row {
			if piece == "wK" {
				WhiteKing = [2]int{r, c}
			}
			if piece == "bK" {
				BlackKing = [2]int{r, c}
			}
		}
	}

	// castling
	// for white the king has to be at 7 4
	if WhiteKing[0] == 7 && WhiteKing[1] == 4 {
		// for kingside the white rook has to be at 7 7
		WCastleKS = Board[7][7] == "wR"

		// for queenside the white rook has to be at 7 0
		WCastleQS = Board[7][0] == "wR"
	}

	// for black the king has to be at 0 4
	if BlackKing[0] == 0 && BlackKing[1] == 4 {
		// for kingside the black rook has to be at 0 7
		BCastleKS = Board[0][7] == "bR"
		// for queenside the black rook has to be at 0 0
		BCastleQS = Board[0][0] == "bR"
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

func Draw(color byte) bool {
	_, n := GetAllValidMoves(color)
	return !InCheck(color) && n == 0 && !ValidKingSideCastle(color) && !ValidQueenSideCastle(color)
}

func Checkmate(color byte) bool {
	_, n := GetAllValidMoves(color)
	return InCheck(color) && n == 0 && !ValidKingSideCastle(color) && !ValidQueenSideCastle(color)
}
