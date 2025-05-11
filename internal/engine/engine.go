package engine

import "time"

type Options struct {
	SearchDepth int
	Type        byte
}

var (
	searchDepth int
	t           byte
)

func Init(options Options) {
	searchDepth = options.SearchDepth

	t = options.Type
	if t != 'r' && t != 'n' {
		panic("Input r or n as an engine type.")
	}
}

func FindMove() (int, int) {
	start := time.Now()

	switch t {
	case 'r':
		return randomMove(), timeSince(start)
	case 'n':
		return alphaBetaWrapper(), timeSince(start)
	}

	return 0, 0
}

func timeSince(start time.Time) int {
	return int(time.Since(start).Milliseconds())
}
