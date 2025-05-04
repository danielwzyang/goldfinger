package main

import (
	"danielyang.cc/chess/internal/board"
)

func main() {
	board.Init(board.DEFAULTBOARD)
	board.ParseFEN("1B3K2/1P1pQ3/1nq3p1/PpP5/6p1/2b5/6b1/1k5B w - - 0 1")
	board.Print()
}
