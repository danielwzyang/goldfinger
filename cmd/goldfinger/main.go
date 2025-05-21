package main

import (
	"flag"
	"fmt"
	"regexp"
	"time"

	"danielyang.cc/chess/internal/board"
	"danielyang.cc/chess/internal/engine"
	"danielyang.cc/chess/internal/gui"
)

func main() {
	// grab flags
	fen := flag.String("fen", board.DEFAULT_BOARD, "Board state in FEN format")
	depth := flag.Int("depth", 6, "Starting search depth (recommended 6-8, lower or increase according to required performance)")
	playBlack := flag.Bool("black", false, "Player plays black")
	flag.Parse()

	playerSide := board.WHITE
	if *playBlack {
		playerSide = board.BLACK
	}

	// init
	board.ParseFEN(*fen)
	board.Init()

	engine.SetOptions(engine.Options{
		SearchDepth: *depth,
	})

	engineMoves := 0
	engineTime := 0
	maxTime := 0

	// start gui
	gui.NewGame()
	go gui.Run()

	// game loop
	fmt.Println("──────────────────────────────────────────────────────")
	fmt.Println("Goldfinger | danielyang.cc")
	fmt.Println("──────────────────────────────────────────────────────")

	var input string
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

		if board.Side == playerSide {
			// player's turn

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

				move := board.StringToMove(input)

				if !possibleMoves.ContainsMove(move) {
					fmt.Println("This move is not possible!")
					continue
				}

				if !board.MakeMove(move, board.ALL_MOVES) {
					fmt.Println("You are still in check!")
					continue
				}

				gui.UpdateBoard(move)

				fmt.Println("You played:")
				board.PrintMove(move)
				break
			}
		} else {
			// engine's turn

			move, ms := engine.FindMove()

			if move == 0 {
				fmt.Println("The engine resigns :(")
				break
			}

			board.MakeMove(move, board.ALL_MOVES)
			gui.UpdateBoard(move)

			fmt.Println("The engine played:")
			board.PrintMove(move)
			fmt.Println()

			engineMoves++
			engineTime += ms
			maxTime = max(ms, maxTime)

			fmt.Printf("Thought for %d ms.\n(Avg: %dms | Max: %dms | Total: %dms)\n", ms, engineTime/engineMoves, maxTime, engineTime)
		}

		fmt.Println()
	}

	fmt.Println()
	fmt.Println("Stopping in 10 seconds..")
	start := time.Now()
	for time.Since(start).Seconds() < 10 {
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
