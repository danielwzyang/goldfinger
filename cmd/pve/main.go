package main

import (
	"fmt"
	"regexp"

	"danielyang.cc/chess/internal/board"
	"danielyang.cc/chess/internal/engine"
)

func main() {
	board.ParseFEN(board.DEFAULT_BOARD)
	board.Init()
	board.Print()

	engine.Init(engine.Options{
		SearchDepth: 6,
		Type:        'n',
	})

	var input string
	for {
		possibleMoves := board.MoveList{}
		board.GenerateAllMoves(&possibleMoves)

		for {
			fmt.Scanln(&input)
			if !validInput(input) {
				println("Invalid input!")
				continue
			}

			move := board.StringToWhiteMove(input)

			if !possibleMoves.ContainsMove(move) {
				board.PrintMove(move)
				println("Move not possible!")
				continue
			}

			if !board.MakeMove(move, board.ALL_MOVES) {
				board.PrintMove(move)
				println("Still in check!")
				continue
			}

			board.PrintMove(move)
			break
		}

		println()
		move, ms := engine.FindMove()
		board.MakeMove(move, board.ALL_MOVES)
		board.PrintMove(move)
		fmt.Printf("thought for %d ms\n", ms)
		board.Print()
	}
}

var regex = regexp.MustCompile("^[a-h][1-8][a-h][1-8][qrnb]?$")

func validInput(input string) bool {
	return regex.MatchString(input)
}
