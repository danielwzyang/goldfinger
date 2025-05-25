package engine

import (
	"time"

	"danielyang.cc/chess/internal/board"
)

const maxSearchDepth = 20

var (
	TimeForMove int
	ply   int
	nodes int
)
// move, time, depth, nodes
func FindMove() (int, int, int, int) {
	searchStart = time.Now()

	alpha := -board.LIMIT_SCORE
	beta := board.LIMIT_SCORE

	nodes = 0
	lastDepthTime := float64(0)
	bestMove := 0
	bestDepth := 0

	for depth := 1; depth <= maxSearchDepth; depth++ {
		timeLeft := TimeForMove - int(time.Since(searchStart).Milliseconds())

		if depth > 1 && lastDepthTime > 0 {
			// this assumes the next depth takes 5x the time of the last depth
			// this is a rough estimate and isn't always accurate
			// but it helps to avoid searching too deep if time is running out
			estimatedNextDepthTime := lastDepthTime * 5.0
			if float64(timeLeft) < estimatedNextDepthTime {
				break
			}
		}

		depthStart := time.Now()
		move, score := alphaBeta(alpha, beta, depth)
		lastDepthTime = float64(time.Since(depthStart).Milliseconds())

		// alpha beta will return 0 if time is up
		if move == 0 {
			bestDepth = depth - 1
			break
		}

		bestMove = move
		bestDepth = depth

		// out of window
		if score <= alpha || score >= beta {
			alpha = -board.LIMIT_SCORE
			beta = board.LIMIT_SCORE
			continue
		}

		// narrow down window by 50 centipawns
		alpha = score - 50
		beta = score + 50
	}

	return bestMove, timeSince(searchStart), bestDepth, nodes
}

func timeSince(start time.Time) int {
	return int(time.Since(start).Milliseconds())
}
