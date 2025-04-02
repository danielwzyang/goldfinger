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
	board.Init()

	StartOpts(&playerColor, &engineType)

	var enemyColor byte = 'b'
	if playerColor == 'b' {
		enemyColor = 'w'
	}

	engine.Init(engineType, enemyColor)

	var engineLastMove string
	var engineTime string

	tick := 0

	playerTick := 1
	if playerColor == 'b' {
		playerTick = 0
	}

	for {
		tick++

		Clear()

		board.Print()

		stop, endingText := Stop()
		if stop {
			fmt.Println(endingText)
			return
		}

		if tick%2 != playerTick {
			engineLastMove, engineTime = engine.MakeMove()
			continue
		}

		if tick != 1 {
			fmt.Println("The engine played " + engineLastMove + ".")
			fmt.Println("It thought for " + engineTime + ".")
		}

		if board.InCheck('b') {
			fmt.Println("Black is in check.")
		}

		if board.InCheck('w') {
			fmt.Println("White is in check.")
		}

		fmt.Println()
		fmt.Println("Enter your move. Read /docs/input.md for help.")
		InputMove()
	}
}

func Stop() (bool, string) {
	if Draw('b') || Draw('w') {
		return true, "The game has ended in a draw!"
	}

	if Checkmate('b') {
		return true, "White has won by checkmate!"
	}

	if Checkmate('w') {
		return true, "Black has won by checkmate!"
	}

	return false, ""
}

func Draw(color byte) bool {
	_, n := board.GetAllValidMoves(color)
	return !board.InCheck(color) && n == 0 && !board.ValidKingSideCastle(color) && !board.ValidQueenSideCastle(color)
}

func Checkmate(color byte) bool {
	_, n := board.GetAllValidMoves(color)
	return board.InCheck(color) && n == 0 && !board.ValidKingSideCastle(color) && !board.ValidQueenSideCastle(color)
}
