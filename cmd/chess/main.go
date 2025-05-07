package main

import (
	"fmt"

	"danielyang.cc/chess/internal/board"
)

func main() {
	board.ParseFEN(board.SEARCH_TESTER)
	board.Init()
	board.Print()

	moves := board.MoveList{}
	board.GenerateAllMoves(&moves)

	move := moves.Moves[0]
	fmt.Println(board.GetSource(move), board.GetTarget(move))
	board.MakeMove(move, board.ALL_MOVES)

	board.Print()
}
