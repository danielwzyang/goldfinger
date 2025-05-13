package main

import (
	"flag"
	"fmt"

	"danielyang.cc/chess/internal/board"
	"danielyang.cc/chess/internal/engine"
)

func main() {
	// grab flags
	fen := flag.String("fen", board.DEFAULT_BOARD, "Board state in FEN format")
	depth := flag.Int("depth", 8, "Starting search depth (recommended 6-8, lower or increase according to required performance)")
	flag.Parse()

	// init
	board.ParseFEN(*fen)
	board.Init()

	engine.Init(engine.Options{
		SearchDepth: *depth,
	})

	engineMoves := 0
	engineTime := 0
	maxTime := 0

	fmt.Println("──────────────────────────────────────────────────────")
	fmt.Println("Goldfinger | danielyang.cc")
	fmt.Println("──────────────────────────────────────────────────────")

	board.Print(0)

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

		fmt.Printf("Thought for %d ms.\n(Avg: %dms | Max: %dms | Total: %dms)\n", ms, engineTime/engineMoves, maxTime, engineTime)
		fmt.Println()
	}

	fmt.Printf("Finished in %d plies.", engineMoves)
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
