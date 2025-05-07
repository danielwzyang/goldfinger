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

	for i := 0; i < moves.Count; i++ {
		fmt.Println(board.POSITIONTOSTRING[board.GetSource(moves.Moves[i])], board.POSITIONTOSTRING[board.GetTarget(moves.Moves[i])])
	}
}
