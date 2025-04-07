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
			BlackKing = [2]int{r2, c2}
			BCastleKS = false
			BCastleQS = false
		} else {
			WhiteKing = [2]int{r2, c2}
			WCastleKS = false
			WCastleQS = false
		}
	}

	// if rook, invalidate castling for that side
	if piece[1] == 'R' {
		if piece[0] == 'b' {
			if c1 == 7 {
				BCastleKS = false
			}
			if c1 == 0 {
				BCastleQS = false
			}
		} else {
			if c1 == 7 {
				WCastleKS = false
			}
			if c1 == 0 {
				WCastleQS = false
			}
		}
	}

	// if pawn handle en passant related things
	if piece[1] == 'P' {
		// pawn made double move
		if math.Abs(float64(r2-r1)) == 2 {
			// set en passant target to the space skipped over
			EnPassant = [2]int{(r1 + r2) / 2, c2}
			return
		} else if math.Abs(float64(r2-r1)) == 1 && math.Abs(float64(c2-c1)) == 1 && r2 == EnPassant[0] && c2 == EnPassant[1] {
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
			EnPassant = [2]int{-10, -10}
			return
		}
	}

	// any move that isn't double pawn invalidates en passant
	EnPassant = [2]int{-10, -10}
}

func GetCaptureMoves(color byte) [][2][2]int {
	captures := [][2][2]int{}

	for r, row := range Board {
		for c, piece := range row {
			if piece[0] == color {
				for _, move := range getValidMoves(r, c) {
					if Board[move[0]][move[1]][0] != ' ' && Board[move[0]][move[1]][0] != color {
						captures = append(captures, [2][2]int{{r, c}, move})
					}
				}
			}
		}
	}

	return captures
}

func Draw(color byte) bool {
	_, n := GetAllValidMoves(color)
	return !InCheck(color) && n == 0 && !ValidKingSideCastle(color) && !ValidQueenSideCastle(color)
}

func Checkmate(color byte) bool {
	_, n := GetAllValidMoves(color)
	return InCheck(color) && n == 0 && !ValidKingSideCastle(color) && !ValidQueenSideCastle(color)
}
