package engine

import (
	"danielyang.cc/chess/internal/board"
)

const DELTA_MARGIN = 975

func quiesce(alpha, beta int) int {
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

	for _, move := range moves.Moves {
		mvvlva := getMVVLVA(move)
		if mvvlva < 0 {
			break
		}

		if standpat+mvvlva+delta < alpha {
			continue
		}

		board.REPETITION_INDEX++
		board.REPETITION_TABLE[board.REPETITION_INDEX] = board.ZobristHash

		if !board.MakeMove(move, board.ONLY_CAPTURES) {
			board.REPETITION_INDEX--
			continue
		}

		score := -quiesce(-beta, -alpha)

		board.REPETITION_INDEX--

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
