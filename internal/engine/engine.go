package engine

import (
	"context"
	"time"

	"danielyang.cc/chess/internal/board"
)

const maxSearchDepth = 64

type SearchResult struct {
	BestMove int
	Time     int
	Depth    int
	Nodes    int
	Score    int
}

var (
	ply   int
	nodes int
)

// move, time, depth, nodes
func FindMove(ctx context.Context) SearchResult {
	start := time.Now()

	alpha := -board.LIMIT_SCORE
	beta := board.LIMIT_SCORE

	nodes = 0

	result := SearchResult{}

	for depth := 1; depth <= maxSearchDepth; depth++ {
		move, score := alphaBeta(ctx, alpha, beta, depth)

		// out of window
		if score <= alpha || score >= beta {
			alpha = -board.LIMIT_SCORE
			beta = board.LIMIT_SCORE
			move, score = alphaBeta(ctx, alpha, beta, depth)
		}

		// narrow window
		alpha = score - 50
		beta = score + 100

		// ignore search if time ran out midway
		select {
		case <-ctx.Done():
			// search was cut off
			return SearchResult{
				result.BestMove, // ignore the move from the cut off search
				timeSince(start),
				depth,
				nodes,
				result.Score, // ignore the score from the cut off search
			}
		default:
			// search completed successfully
			result = SearchResult{
				move,
				timeSince(start),
				depth,
				nodes,
				score,
			}
		}

	}

	return result
}

func timeSince(start time.Time) int {
	return int(time.Since(start).Milliseconds())
}
