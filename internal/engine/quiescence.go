package engine

import (
	"danielyang.cc/chess/internal/board"
)

const DELTA_MARGIN = 975

var SIMPLE_PIECE_WEIGHTS = [6]int{100, 325, 325, 500, 975, 10000}

func quiesce(alpha, beta int) int {
	if stopFlag {
		return 0
	}

	standpat := board.Evaluate()

	if standpat >= beta {
		return beta
	}

	if standpat < alpha-DELTA_MARGIN {
		return alpha
	}

	if standpat > alpha {
		alpha = standpat
	}

	moves := board.MoveList{}
	board.GenerateAllMoves(&moves)

	scores := make([]int, moves.Count)
	for i := 0; i < moves.Count; i++ {
		scores[i] = getMVVLVA(moves.Moves[i])
	}

	sortMoves(&moves, scores)

	for moveCount := 0; moveCount < moves.Count; moveCount++ {
		if stopFlag {
			return alpha
		}

		move := moves.Moves[moveCount]

		if !board.MakeMove(move, board.ONLY_CAPTURES) {
			continue
		}

		score := -quiesce(-beta, -alpha)

		board.RestoreState()

		if score >= beta {
			return beta
		}
		if score > alpha {
			alpha = score
		}
	}

	return alpha
}
