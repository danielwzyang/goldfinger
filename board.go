package main

import "fmt"

var board [8][8]string

func InitBoard() {
	board = [8][8]string{
		{"bR", "bN", "bB", "bQ", "bK", "bB", "bN", "bR"},
		{"bP", "bP", "bP", "bP", "bP", "bP", "bP", "bP"},
		{" ", " ", " ", " ", " ", " ", " ", " "},
		{" ", " ", " ", " ", " ", " ", " ", " "},
		{" ", " ", " ", " ", " ", " ", " ", " "},
		{" ", " ", " ", " ", " ", " ", " ", " "},
		{"wP", "wP", "wP", "wP", "wP", "wP", "wP", "wP"},
		{"wR", "wN", "wB", "wQ", "wK", "wB", "wN", "wR"},
	}
}

var ascii = map[string]string{
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
}

func PrintBoard() {
	alternate := true

	for i, row := range board {
		fmt.Printf("%-*d", 2, 8-i)
		for _, char := range row {
			fmt.Printf("%-*s", 2, ascii[char])

			alternate = !alternate
		}

		fmt.Println()

		alternate = !alternate
	}

	for _, char := range " abcdefgh" {
		fmt.Printf("%-*c", 2, char)
	}
}

func ClearTerminal() {
	fmt.Print("\033[H\033[2J")
}
