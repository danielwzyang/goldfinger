package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"danielyang.cc/chess/internal/board"
	"danielyang.cc/chess/internal/engine"
	"danielyang.cc/chess/internal/gui"
)

func main() {
	// grab flags
	fen := flag.String("fen", board.DEFAULT_BOARD, "Board state in FEN format")
	depth := flag.Int("depth", 6, "Starting search depth (recommended 6-8, lower or increase according to required performance)")
	flag.Parse()

	// init
	board.ParseFEN(*fen)
	board.Init()

	engine.SetOptions(engine.Options{
		SearchDepth: *depth,
	})

	engineMoves := 0
	engineTime := 0
	maxTime := 0

	// Store thinking times for each ply
	var thinkingTimes []int

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

		move, ms := engine.FindMove()

		if move == 0 {
			fmt.Println("The engine resigns :(")
			break
		}

		board.MakeMove(move, board.ALL_MOVES)
		gui.UpdateBoard(move)

		thinkingTimes = append(thinkingTimes, ms)

		fmt.Println("The engine played:")
		board.PrintMove(move)
		fmt.Println()

		engineMoves++
		engineTime += ms
		maxTime = max(ms, maxTime)

		fmt.Printf("Thought for %d ms.\n(Avg: %dms | Max: %dms | Total: %dms)\n", ms, engineTime/engineMoves, maxTime, engineTime)
		fmt.Println()
	}

	fmt.Println()
	fmt.Printf("Finished in %d plies.\n", engineMoves)

	// csv used to analyze trends in thinking time
	file, err := os.Create("thinking_times.csv")
	if err != nil {
		fmt.Println("Error creating CSV file:", err)
	} else {
		defer file.Close()
		fmt.Fprintln(file, "ply,ms")
		for i, t := range thinkingTimes {
			fmt.Fprintf(file, "%d,%d\n", i+1, t)
		}
		fmt.Println("Thinking times exported to thinking_times.csv")
	}

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
		if board.MakeMove(move, board.ALL_MOVES) {
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
