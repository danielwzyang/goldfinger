package engine

import (
	"fmt"
	"time"

	"danielyang.cc/chess/internal/board"
)

type Options struct {
	SearchDepth int
}

var (
	searchDepth int
	stopFlag    bool
	ply         int
	nodes       int
)

func SetOptions(options Options) {
	searchDepth = options.SearchDepth
}

func SetStopFlag(flag bool) {
	stopFlag = flag
}

func FindMove() (int, int) {
	start := time.Now()

	alpha := -board.LIMIT_SCORE
	beta := board.LIMIT_SCORE

	move := 0
	score := 0
	nodes = 0

	for depth := 1; depth <= searchDepth; depth++ {
		if stopFlag {
			break
		}

		move, score = alphaBeta(alpha, beta, depth)

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

	// print nps for total iterative deepening search
	fmt.Printf("Nodes: %d | per second: %.0f\n", nodes, float64(nodes)/time.Since(start).Seconds())

	return move, timeSince(start)
}

func timeSince(start time.Time) int {
	return int(time.Since(start).Milliseconds())
}
