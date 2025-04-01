package game

import (
	"fmt"

	"danielyang.cc/chess/internal/board"
)

func Start() {
	board.Init()

	var playerColor byte
	StartOpts(&playerColor)

	for {
		Clear()

		board.Print()

		fmt.Println("Enter your move. Read /docs/input.md for help.")

		InputMove(playerColor)
	}
}
