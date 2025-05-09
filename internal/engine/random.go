package engine

import (
	"math/rand"

	"danielyang.cc/chess/internal/board"
)

func randomMove() int {
	moves := board.MoveList{}
	board.GenerateAllMoves(&moves)

	move := moves.Moves[rand.Intn(moves.Count)]
	board.MakeMove(move, board.ALL_MOVES)

	return move
}
