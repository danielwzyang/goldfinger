package engine

import (
	"context"
	"fmt"
	"time"

	"danielyang.cc/chess/internal/board"
	"danielyang.cc/chess/internal/polyglot"
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
	Stop  context.CancelFunc
)

func Init() {
	polyglot.LoadBook("books/Perfect2023.bin")
}

// move, time, depth, nodes
func FindMove(timeForMove int, print bool) SearchResult {
	start := time.Now()

	if polyglot.HasBookMove() {
		return SearchResult{
			polyglot.GetBestMove(),
			timeSince(start),
			0,
			0,
			0,
		}
	}

	alpha := -board.LIMIT_SCORE
	beta := board.LIMIT_SCORE

	nodes = 0

	result := SearchResult{}

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(timeForMove))
	Stop = cancel
	defer cancel()
	defer func() {
		Stop = nil
		if print {
			if result.BestMove != 0 {
				fmt.Printf("bestmove %s\n", board.MoveToString(result.BestMove))
			} else {
				fmt.Println("bestmove 0000")
			}
		}
	}()

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
			return result
		default:
			// search completed successfully
			result = SearchResult{
				move,
				timeSince(start),
				depth,
				nodes,
				score,
			}

			if print {
				printInfo(result)
			}
		}

	}

	return result
}

func timeSince(start time.Time) int {
	return int(time.Since(start).Milliseconds())
}

func printInfo(result SearchResult) {
	fmt.Printf("info depth %d time %d nodes %d nps %d score cp %d\n",
		result.Depth, result.Time, result.Nodes, result.Nodes*1000/max(1, result.Time), result.Score)
}
