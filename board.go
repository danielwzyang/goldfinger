package main

import (
	"fmt"
)

var Board [8][8]string

func InitBoard() {
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
	for i, row := range Board {
		// 8 to 1 on left side
		fmt.Printf("%-*d", 2, 8-i)

		for _, char := range row {
			fmt.Printf("%-*s", 2, ascii[char])
		}

		fmt.Println()
	}

	// letters on bottom side

	for _, char := range " abcdefgh" {
		fmt.Printf("%-*c", 2, char)
	}

	fmt.Println()
}
