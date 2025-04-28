package engine

import (
	"danielyang.cc/chess/internal/board"
)

func quiesce(alpha int, beta int, currentColor int) int {
	standPat := Evaluate(currentColor)
	bestScore := standPat

	// fail hard
	if standPat >= beta {
		return standPat
	}

	captures := board.GetCaptureMoves(currentColor)

	moveScores := make([]int, len(captures))

	for i, move := range captures {
		attacker := board.Board[move.From.Rank][move.From.File]
		victim := board.Board[move.To.Rank][move.To.File]
		// basic static exchange eval
		moveScores[i] = pieceWeights[victim.Type]*13 - pieceWeights[attacker.Type]
	}

	if len(captures) > 2 {
		insertionSort(captures, moveScores)
	}

	for _, capture := range captures {
		// delta pruning skips captures that dont raise alpha + prune if see < 0
		attacker := board.Board[capture.From.Rank][capture.From.File]
		victim := board.Board[capture.To.Rank][capture.To.File]
		see := pieceWeights[victim.Type] - pieceWeights[attacker.Type]
		if see < 0 || standPat+see+200 <= alpha {
			continue
		}

		board.MakeMove(capture)

		// pawn promotion
		if board.Board[capture.To.Rank][capture.To.File].Type == board.PAWN && (capture.To.Rank == 0 || capture.To.Rank == 7) {
			// automatically promote to queen
			board.Board[capture.To.Rank][capture.To.File] = board.Piece{
				Type:  board.QUEEN,
				Color: currentColor,
				Key:   currentColor*6 + board.QUEEN + 1,
			}
		}

		score := -quiesce(-beta, -alpha, currentColor^1)

		board.UndoMove()

		if score >= beta {
			return score
		}
		if score > bestScore {
			bestScore = score
		}
		if score > alpha {
			alpha = score
		}
	}

	return bestScore
}
