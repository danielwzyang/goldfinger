package engine

import "time"

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

	return iterativeDeepening(), timeSince(start)
}

func timeSince(start time.Time) int {
	return int(time.Since(start).Milliseconds())
}
