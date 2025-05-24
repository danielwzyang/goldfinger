package engine

import (
	"fmt"
	"time"

	"danielyang.cc/chess/internal/board"
)

const maxSearchDepth = 20

type Options struct {
	TimeForMove int
}

var (
	timeForMove int
	stopFlag    bool
	ply         int
	nodes       int
)

func SetOptions(options Options) {
	timeForMove = options.TimeForMove
}

func SetStopFlag(flag bool) {
	stopFlag = flag
}

func FindMove() (int, int, int) {
	start := time.Now()

	alpha := -board.LIMIT_SCORE
	beta := board.LIMIT_SCORE

	move := 0
	score := 0
	nodes = 0
	depthReached := maxSearchDepth

	lastDepthTime := float64(0)

	for depth := 1; depth <= maxSearchDepth; depth++ {
		if stopFlag {
			break
		}

		timeLeft := timeForMove - int(time.Since(start).Milliseconds())

		if depth > 1 && lastDepthTime > 0 {
			// this assumes the next depth takes 5x the time of the last depth
			// this is a rough estimate and isn't always accurate
			// but it helps to avoid searching too deep if time is running out
			estimatedNextDepthTime := lastDepthTime * 5.0
			if float64(timeLeft) < estimatedNextDepthTime {
				depthReached = depth - 1
				break
			}
		}

		depthStart := time.Now()
		move, score = alphaBeta(alpha, beta, depth)
		lastDepthTime = float64(time.Since(depthStart).Milliseconds())

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

	return move, timeSince(start), depthReached
}

func timeSince(start time.Time) int {
	return int(time.Since(start).Milliseconds())
}
