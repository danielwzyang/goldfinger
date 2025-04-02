package game

import (
	"fmt"

	"danielyang.cc/chess/internal/board"
)

func Start() {
	board.Init()

	var playerColor byte
	StartOpts(&playerColor)

	var enemyColor byte = 'b'
	if playerColor == 'b' {
		enemyColor = 'w'
	}

	for {
		Clear()

		board.Print()

		stop, endingText := Stop()
		if stop {
			fmt.Println(endingText)
			return
		}

		if board.InCheck('b') {
			fmt.Println("Black is in check.")
		}

		if board.InCheck('w') {
			fmt.Println("White is in check.")
		}

		fmt.Println("Enter your move. Read /docs/input.md for help.")

		InputMove(playerColor)
		InputMove(enemyColor)
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
