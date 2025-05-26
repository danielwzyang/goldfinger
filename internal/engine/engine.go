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

		// ignore search if time ran out midway
		select {
		case <-ctx.Done():
			return result
		default:
		}

		result = SearchResult{
			move,
			timeSince(start),
			depth,
			nodes,
			score,
		}

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

	return result
}

func timeSince(start time.Time) int {
	return int(time.Since(start).Milliseconds())
}
