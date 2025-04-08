package game

import (
	"fmt"

	"danielyang.cc/chess/internal/board"
	"danielyang.cc/chess/internal/engine"
)

var (
	playerColor byte
	engineType  byte
)

func Start() {
	board.Init(board.DefaultBoard)

	StartOpts(&playerColor, &engineType)

	var engineColor byte = 'b'
	if playerColor == 'b' {
		engineColor = 'w'
	}

	engine.Init(engineType, engineColor, 10)

	var engineLastMove string
	var engineTime string

	tick := 0

	playerTick := 1
	if playerColor == 'b' {
		playerTick = 0
	}

	for {
		tick++

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

		fmt.Printf("Eval: %.2f\n", engine.Evaluate('w'))

		if board.InCheck('b') {
			fmt.Println("Black is in check.")
		}

		if board.InCheck('w') {
			fmt.Println("White is in check.")
		}

		fmt.Println()

		if tick%2 != playerTick {
			engineLastMove, engineTime = engine.MakeMove()
			continue
		}

		if tick != 1 {
			fmt.Println("The engine played " + engineLastMove + ".")
			fmt.Println("It thought for " + engineTime + ".")
		}

		fmt.Println()
		fmt.Println("Enter your move. Read /docs/input.md for help.")
		InputMove()
	}
}

func Stop() (bool, string) {
	if board.Draw('b') || board.Draw('w') {
		return true, "The game has ended in a draw!"
	}

	if board.Checkmate('b') {
		return true, "White has won by checkmate!"
	}

	if board.Checkmate('w') {
		return true, "Black has won by checkmate!"
	}

	return false, " "
}
