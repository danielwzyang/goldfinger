package engine

import (
	"context"
	"time"

	"danielyang.cc/chess/internal/board"
)

const maxSearchDepth = 20

type SearchResult struct {
	BestMove int
	Time     int
	Depth    int
	Nodes    int
	Score    int
}

var (
	ply    int
	nodes  int
	Result SearchResult
)

// move, time, depth, nodes
func FindMove(ctx context.Context) {
	start := time.Now()

	alpha := -board.LIMIT_SCORE
	beta := board.LIMIT_SCORE

	nodes = 0

	for depth := 1; depth <= maxSearchDepth; depth++ {
		select {
		case <-ctx.Done():
			println("test")
			return
		default:
		}

		move, score := alphaBeta(ctx, alpha, beta, depth)

		// alpha beta will return 0 if time is up
		if move == 0 {
			break
		}

		Result = SearchResult{
			move,
			max(1, timeSince(start)),
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
}

func timeSince(start time.Time) int {
	return int(time.Since(start).Milliseconds())
}
