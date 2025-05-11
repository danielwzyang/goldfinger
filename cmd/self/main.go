package main

import (
	"fmt"

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
		SearchDepth: 6,
		Type:        'n',
	})

	engineMoves := 0
	engineTime := 0
	maxTime := 0

	for {
		if board.Fifty >= 100 {
			fmt.Println("Draw by fifty move rule!")
			break
		}

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
