package engine

import (
	"danielyang.cc/chess/internal/board"
	"danielyang.cc/chess/internal/transposition"
)

var sortedMovesCache = map[uint64][]board.Move{}

func quiesce(alpha float64, beta float64, color int) float64 {
	standPat := Evaluate(color)
	bestScore := standPat

	// fail hard
	if standPat >= beta {
		return standPat
	}

	// delta pruning; 200 centipawns
	delta := 2.0
	if standPat < alpha-delta {
		return alpha
	}

	if alpha < standPat {
		alpha = standPat
	}

	// generate captures or get from cache
	var captures []board.Move

	cache, ok := sortedMovesCache[transposition.HashBoard(color)]
	if ok {
		captures = cache
	} else {
		captures = board.GetCaptureMoves(color)

		moveScores := make([]int, len(captures))

		for i, move := range captures {
			score := 0

			attacker := board.Board[move.From.Rank][move.From.File]
			victim := board.Board[move.To.Rank][move.To.File]
			score += int(pieceWeights[victim.Type]*13 - pieceWeights[attacker.Type])

			moveScores[i] = score
		}

		insertionSort(captures, moveScores)
	}

	for _, capture := range captures {
		board.MakeMove(capture)

		// pawn promotion
		if board.Board[capture.To.Rank][capture.To.File].Type == board.PAWN && (capture.To.Rank == 0 || capture.To.Rank == 7) {
			// automatically promote to queen
			board.Board[capture.To.Rank][capture.To.File] = board.Piece{
				Type:  board.QUEEN,
				Color: color,
				Key:   board.GetKey(board.QUEEN, color),
			}
		}

		score := -quiesce(-beta, -alpha, color^1)

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
