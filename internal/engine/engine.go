package engine

import (
	"time"

	"danielyang.cc/chess/internal/board"
)

type Options struct {
	SearchDepth int
}

var (
	searchDepth int
)

func Init(options Options) {
	searchDepth = options.SearchDepth
}

func FindMove() (int, int) {
	start := time.Now()

	move, _ := alphaBeta(-board.LIMIT_SCORE, board.LIMIT_SCORE, searchDepth)

	return move, timeSince(start)
}

func timeSince(start time.Time) int {
	return int(time.Since(start).Milliseconds())
}
