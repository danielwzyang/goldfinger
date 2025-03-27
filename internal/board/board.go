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
