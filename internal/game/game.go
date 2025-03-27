package game

import (
	"fmt"
	"runtime"

	"danielyang.cc/chess/internal/board"
	"danielyang.cc/chess/internal/input"
)

func Start() {
	// lock to main thread
	runtime.LockOSThread()

	board.Init()

	var playerColor byte
	input.StartOpts(&playerColor)

	for {
		input.Clear()

		board.Print()

		fmt.Println("Enter your move in the format: xn-xn.")

		move := input.Input()
		valid, err := board.ValidInput(move, playerColor)
		for !valid {
			fmt.Println(err)
			move = input.Input()
			valid, err = board.ValidInput(move, playerColor)
		}

		board.ParseMove(move, playerColor)
	}
}
