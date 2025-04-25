package engine

import (
	"danielyang.cc/chess/internal/board"
	"danielyang.cc/chess/internal/transposition"
)

var sortedMovesCache = map[uint64][]board.Move{}

func quiesce(alpha int, beta int, currentColor int) int {
	standPat := Evaluate(currentColor)
	bestScore := standPat

	// fail hard
	if standPat >= beta {
		return standPat
	}

	// generate captures or get from cache
	var captures []board.Move

	hash := transposition.HashBoard(currentColor)
	cache, ok := sortedMovesCache[hash]
	if ok {
		captures = cache
	} else {
		captures = board.GetCaptureMoves(currentColor)

		moveScores := make([]int, len(captures))

		for i, move := range captures {
			score := 0

			attacker := board.Board[move.From.Rank][move.From.File]
			victim := board.Board[move.To.Rank][move.To.File]
			score += pieceWeights[victim.Type]*13 - pieceWeights[attacker.Type]

			moveScores[i] = score
		}

		insertionSort(captures, moveScores)

		sortedMovesCache[hash] = captures
	}

	for _, capture := range captures {
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
