package game

import (
	"fmt"

	"danielyang.cc/chess/internal/board"
	"danielyang.cc/chess/internal/engine"
)

var (
	currentColor = board.WHITE
	playerColor  int
	engineType   byte
)

func Start(searchDepth int) {
	board.Init(board.DefaultBoard)

	StartOpts(&playerColor, &engineType)

	var engineColor int = board.BLACK
	if playerColor == board.BLACK {
		engineColor = board.WHITE
	}

	engine.Init(engineType, engineColor, searchDepth)

	var (
		engineLastMove   string
		engineTime       int
		engineTotalTime  int
		engineTotalMoves int
	)
	for {
		Header()

		board.Print()

		stop, endingText := Stop()
		if stop {
			fmt.Println(endingText)
			fmt.Println("Type q to exit.")
			input := Input()
			for input != "q" {
				input = Input()
			}
			return
		}

		fmt.Printf("Eval: %d\n", engine.Evaluate(board.WHITE))

		if board.InCheck(board.BLACK) {
			fmt.Println("Black is in check.")
		}

		if board.InCheck(board.WHITE) {
			fmt.Println("White is in check.")
		}

		fmt.Println()

		if engineColor == currentColor {
			engineLastMove, engineTime = engine.MakeMove()
			engineTotalTime += engineTime
			engineTotalMoves++
			currentColor ^= 1
			continue
		}

		if engineLastMove != "" {
			fmt.Println("The engine played " + engineLastMove + ".")
			fmt.Printf("It thought for %d ms.\n", engineTime)
			fmt.Printf("Average thinking time: %d ms.\n", engineTotalTime/engineTotalMoves)
		}

		fmt.Println()
		fmt.Println("Enter your move. Read /docs/input.md for help.")
		InputMove()
		currentColor ^= 1
	}
}

func Stop() (bool, string) {
	if board.Draw(currentColor) {
		return true, "The game has ended in a draw!"
	}

	if board.Checkmate(board.BLACK) {
		return true, "White has won by checkmate!"
	}

	if board.Checkmate(board.WHITE) {
		return true, "Black has won by checkmate!"
	}

	return false, " "
}
