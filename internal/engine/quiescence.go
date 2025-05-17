package engine

import (
	"danielyang.cc/chess/internal/board"
)

const DELTA_MARGIN = 975

func quiesce(alpha, beta int) int {
	if stopFlag {
		return 0
	}

	standpat := board.Evaluate()

	if standpat >= beta {
		return beta
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

	delta := DELTA_MARGIN

	for moveCount := 0; moveCount < moves.Count; moveCount++ {
		if stopFlag {
			return alpha
		}

		move := moves.Moves[moveCount]

		mvvlva := getMVVLVA(move)
		if mvvlva < 0 {
			break
		}

		if standpat+mvvlva+delta < alpha {
			continue
		}

		if !board.MakeMove(move, board.ONLY_CAPTURES) {
			continue
		}

		score := -quiesce(-beta, -alpha)

		board.RestoreState()

		if score > alpha {
			alpha = score
			if score >= beta {
				return beta
			}
		}
	}

	return alpha
}
