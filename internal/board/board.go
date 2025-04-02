package board

import "fmt"

var (
	Board     [8][8]string
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
)

func Init() {
	Board = [8][8]string{
		{"bR", "bN", "bB", "bQ", "bK", "bB", "bN", "bR"},
		{"bP", "bP", "bP", "bP", "bP", "bP", "bP", "bP"},
		{" ", " ", " ", " ", " ", " ", " ", " "},
		{" ", " ", " ", " ", " ", " ", " ", " "},
		{" ", " ", " ", " ", " ", " ", " ", " "},
		{" ", " ", " ", " ", " ", " ", " ", " "},
		{"wP", "wP", "wP", "wP", "wP", "wP", "wP", "wP"},
		{"wR", "wN", "wB", "wQ", "wK", "wB", "wN", "wR"},
	}

	BlackKing = [2]int{0, 4}
	WhiteKing = [2]int{7, 4}
}

func Print() {
	gray := "\033[2;37m"
	reset := "\033[0m"
	fmt.Printf("%s   ┌──┬──┬──┬──┬──┬──┬──┬──┐%s\n", gray, reset)

	for i, row := range Board {
		fmt.Printf(" %d %s│%s", 8-i, gray, reset)

		for _, char := range row {
			fmt.Printf("%s\ufe0e %s│%s", ascii[char], gray, reset)
		}
		fmt.Println()

		if i < 7 {
			fmt.Printf("%s   ├──┼──┼──┼──┼──┼──┼──┼──┤%s\n", gray, reset)
		} else {
			fmt.Printf("%s   └──┴──┴──┴──┴──┴──┴──┴──┘%s\n", gray, reset)
		}
	}

	fmt.Println("    a  b  c  d  e  f  g  h")
	fmt.Println()
}

var (
	BCastleKS = true // black can castle kingside until rook or king moves
	BCastleQS = true // black can castle queenside until rook or king moves
	WCastleKS = true // white can castle kingside until rook or king moves
	WCastleQS = true // white can castle queenside until rook or king moves
)

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
	Board[r2][c2] = Board[r1][c1]

	// update original position
	Board[r1][c1] = " "

	// if moved king, update stored king position and invalidate castling
	if Board[r2][c2][1] == 'K' {
		if Board[r2][c2][0] == 'b' {
			BlackKing = [2]int{r2, c2}

			BCastleKS = false
			BCastleQS = false
		} else {
			WhiteKing = [2]int{r2, c2}

			WCastleKS = false
			WCastleQS = false
		}
	}

	// if rook, invalidate castling
	if Board[r2][c2][1] == 'R' {
		if Board[r2][c2][0] == 'b' {
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
}
