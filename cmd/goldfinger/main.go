package main

import (
	"flag"
	"fmt"
	"regexp"

	"danielyang.cc/chess/internal/board"
	"danielyang.cc/chess/internal/engine"
)

func main() {
	// grab flags
	fen := flag.String("fen", board.DEFAULT_BOARD, "Board state in FEN format")
	depth := flag.Int("depth", 8, "Starting search depth (recommended 7-9, lower or increase according to required performance)")
	playerSide := flag.Int("side", 8, "Side that the player plays (0 for white, 1 for black)")
	flag.Parse()

	if *playerSide != 0 && *playerSide != 1 {
		panic("Please input 0 or 1 for the side flag.")
	}

	// init
	board.ParseFEN(*fen)
	board.Init()

	engine.Init(engine.Options{
		SearchDepth: *depth,
	})

	engineMoves := 0
	engineTime := 0
	maxTime := 0

	// game loop
	fmt.Println("──────────────────────────────────────────────────────")
	fmt.Println("Goldfinger | danielyang.cc")
	fmt.Println("──────────────────────────────────────────────────────")

	board.Print(0)

	var input string
	for {
		if over() {
			fmt.Println("No more legal moves!")
			break
		}

		if board.Fifty >= 100 {
			fmt.Println("Draw by fifty move rule!")
			break
		}

		if board.Side == *playerSide {
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

				board.Print(move)
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

			board.Print(move)
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
