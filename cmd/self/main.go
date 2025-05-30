package main

import (
	"flag"
	"fmt"
	"time"

	"danielyang.cc/chess/internal/board"
	"danielyang.cc/chess/internal/engine"
	"danielyang.cc/chess/internal/gui"
)

func main() {
	// grab flags
	fen := flag.String("fen", board.DEFAULT_BOARD, "Board state in FEN format")
	timeForMove := flag.Int("time", 1000, "Time for move in milliseconds. Search may exceed this time but it will not iteratively search deeper once time is up.")
	flag.Parse()

	// init
	board.Init()
	board.ParseFEN(*fen)

	engine.Init()

	engineMoves := 0
	engineTime := 0
	maxTime := 0

	fmt.Println("──────────────────────────────────────────────────────")
	fmt.Println("Goldfinger | danielyang.cc")
	fmt.Println("──────────────────────────────────────────────────────")

	// start gui
	gui.NewGame()
	go gui.Run()

	// game loop
	for {
		if over() {
			fmt.Println("No more legal moves!")
			break
		}

		if board.InsufficientMaterial() {
			fmt.Println("Draw by insufficient material!")
			break
		}

		if repetition() {
			fmt.Println("Draw by repetition!")
			break
		}

		if board.Fifty >= 100 {
			fmt.Println("Draw by fifty move rule!")
			break
		}

		result := engine.FindMove(*timeForMove, false)
		move, ms, depth, nodes := result.BestMove, result.Time, result.Depth, result.Nodes

		if move == 0 {
			fmt.Println("The engine resigns :(")
			break
		}

		board.MakeMove(move)
		gui.UpdateBoard(move)

		fmt.Println("The engine played:")
		board.PrintMove(move)
		fmt.Println()

		engineMoves++
		engineTime += ms
		maxTime = max(ms, maxTime)

		fmt.Printf("Thought for %d ms.\n(Avg: %dms | Max: %dms | Total: %dms)\n", ms, engineTime/engineMoves, maxTime, engineTime)
		fmt.Printf("Search depth: %d\n", depth)
		fmt.Printf("Nodes searched: %d\n", nodes)
		fmt.Printf("Nodes per second: %.0f\n", float64(nodes*1000)/float64(ms))
		fmt.Println()
	}

	fmt.Println()
	fmt.Printf("Finished in %d plies.\n", engineMoves)

	fmt.Println()
	fmt.Println("Stopping in 10 seconds..")
	start := time.Now()
	for time.Since(start).Seconds() < 10 {
	}
}

func over() bool {
	moves := board.MoveList{}
	board.GenerateAllMoves(&moves)

	legal := 0

	for _, move := range moves.Moves {
		if board.MakeMove(move) {
			legal++
			board.RestoreState()
		}
	}

	return legal == 0
}

func repetition() bool {
	count := 0
	for i := board.RepetitionIndex; i >= board.RepetitionIndex-8; i -= 4 {
		if i < 0 {
			break
		}

		if board.RepetitionTable[i] == board.ZobristHash {
			count++
			if count == 3 {
				return true
			}
		}
	}
	return false
}
