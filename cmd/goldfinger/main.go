package main

import (
	"fmt"
	"regexp"

	"danielyang.cc/chess/internal/board"
	"danielyang.cc/chess/internal/engine"
)

func main() {
	fmt.Println("──────────────────────────────────────────────────────")
	fmt.Println("Goldfinger | danielyang.cc")
	fmt.Println("──────────────────────────────────────────────────────")

	board.ParseFEN(board.DEFAULT_BOARD)
	board.Init()
	board.Print(0)

	engine.Init(engine.Options{
		SearchDepth: 9,
		Type:        'n',
	})

	engineMoves := 0
	engineTime := 0
	maxTime := 0

	var input string
	for {
		if over() {
			fmt.Println("No more legal moves!")
			break
		}

		possibleMoves := board.MoveList{}
		board.GenerateAllMoves(&possibleMoves)

		first := true
		for {
			if first {
				fmt.Println("Enter a move! (ex: e2e4, h7h8q)")
				first = false
			}

			fmt.Printf("> ")
			fmt.Scanln(&input)
			if !validInput(input) {
				fmt.Println("Invalid input!")
				continue
			}

			move := board.StringToWhiteMove(input)

			if !possibleMoves.ContainsMove(move) {
				fmt.Println("This move is not possible!")
				continue
			}

			if !board.MakeMove(move, board.ALL_MOVES) {
				fmt.Println("You are still in check!")
				continue
			}

			board.Print(move)
			fmt.Println("You played:")
			board.PrintMove(move)
			break
		}

		fmt.Println()

		if over() {
			fmt.Println("No more legal moves!")
			break
		}

		move, ms := engine.FindMove()

		if move == 0 {
			fmt.Println("The engine resigns :(")
			break
		}

		board.MakeMove(move, board.ALL_MOVES)

		board.Print(move)
		fmt.Println("The engine played:")
		board.PrintMove(move)
		fmt.Println()

		engineMoves++
		engineTime += ms
		maxTime = max(ms, maxTime)

		fmt.Printf("Thought for %d ms.\n(Avg: %dms | Max: %dms)\n", ms, engineTime/engineMoves, maxTime)
		fmt.Println()
	}
}

var regex = regexp.MustCompile("^[a-h][1-8][a-h][1-8][qrnb]?$")

func validInput(input string) bool {
	return regex.MatchString(input)
}

func over() bool {
	moves := board.MoveList{}
	board.GenerateAllMoves(&moves)

	legal := 0

	for _, move := range moves.Moves {
		if board.MakeMove(move, board.ALL_MOVES) {
			legal++
			board.RestoreState()
		}
	}

	return legal == 0
}
