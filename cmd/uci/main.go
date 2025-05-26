package main

import (
	"bufio"
	"context"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"danielyang.cc/chess/internal/board"
	"danielyang.cc/chess/internal/engine"
)

var plies = 0.0

const delay = 100

func main() {
	var cancelFunc context.CancelFunc

	board.Init()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()
		tokens := strings.Split(input, " ")

		if len(tokens) == 0 {
			continue
		}

		switch tokens[0] {
		case "uci":
			fmt.Println("id name Goldfinger")
			fmt.Println("id author Daniel Yang")
			fmt.Println("uciok")

		case "isready":
			fmt.Println("readyok")

		case "ucinewgame":
			board.Init()
			board.ParseFEN(board.DEFAULT_BOARD)
			engine.ResetHeuristics()

		case "position":
			if len(tokens) < 2 {
				continue
			}

			board.Init()
			engine.ResetHeuristics()

			moveIndex := 1
			if tokens[1] == "startpos" {
				board.ParseFEN(board.DEFAULT_BOARD)
				// moves token is right after startpos
				moveIndex = 2
			} else if tokens[1] == "fen" {
				// find moves token
				movesIndex := -1
				for i, token := range tokens {
					if token == "moves" {
						movesIndex = i
						break
					}
				}

				// construct fen string
				var fenString string
				if movesIndex == -1 {
					fenString = strings.Join(tokens[2:], " ")
					moveIndex = len(tokens)
				} else {
					fenString = strings.Join(tokens[2:movesIndex], " ")
					moveIndex = movesIndex
				}

				// parse fen
				board.ParseFEN(fenString)
			}

			// play moves
			if moveIndex < len(tokens) && tokens[moveIndex] == "moves" {
				for _, moveStr := range tokens[moveIndex+1:] {
					board.MakeMove(board.StringToMove(moveStr))
					plies++
				}
			}

		case "go":
			var wtime, btime, winc, binc int

			// parse time control params
			for i := 1; i < len(tokens); i++ {
				switch tokens[i] {
				case "wtime":
					if i+1 < len(tokens) {
						wtime, _ = strconv.Atoi(tokens[i+1])
						i++
					}
				case "btime":
					if i+1 < len(tokens) {
						btime, _ = strconv.Atoi(tokens[i+1])
						i++
					}
				case "winc":
					if i+1 < len(tokens) {
						winc, _ = strconv.Atoi(tokens[i+1])
						i++
					}
				case "binc":
					if i+1 < len(tokens) {
						binc, _ = strconv.Atoi(tokens[i+1])
						i++
					}
				}
			}

			timeForMove := getTimeForMove(wtime, btime, winc, binc)

			// 100 ms delay for the engine to prepare
			ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(timeForMove+delay))
			cancelFunc = cancel

			go func() {
				result := engine.FindMove(ctx)
				<-ctx.Done()
				printResult(result)
			}()
		case "stop":
			if cancelFunc != nil {
				cancelFunc()
			}

		case "quit":
			return
		}
	}
}

func getTimeForMove(wtime, btime, winc, binc int) int {
	var timeLeft, increment int
	if board.Side == board.WHITE {
		timeLeft = wtime
		increment = winc
	} else {
		timeLeft = btime
		increment = binc
	}

	// 0 (endgame) to 24 (opening)
	phase := float64(board.CalculateGamePhase())

	// midgame = ~20, endgame = ~10
	remainingMoves := max(10.0, 24.0-phase+10.0)

	// safety margin of 0.9 and reserve for 2 seconds
	reserve := 2000
	if timeLeft < int(reserve) {
		reserve = 0
	}
	baseTime := 0.9*(float64(timeLeft-reserve)/remainingMoves) + float64(increment)

	// tapered time scaling from 0.5x (opening/endgame) to 1.0x (midgame)
	// quadratic peak at midgame (phase = 12)
	scaling := 1.0 - math.Pow((phase-12.0)/12.0, 2)
	scaling = max(.5, min(1.0, scaling))

	timeForMove := baseTime * scaling

	// cap to 50% of time left
	timeForMove = min(timeForMove, float64(timeLeft)*0.5)

	return int(timeForMove)
}

func printResult(result engine.SearchResult) {
	fmt.Printf("info depth %d time %d nodes %d nps %d score cp %d\n",
		result.Depth, result.Time, result.Nodes, result.Nodes*1000/max(1, result.Time), result.Score)
	if result.BestMove != 0 {
		fmt.Printf("bestmove %s\n", board.MoveToString(result.BestMove))
	} else {
		fmt.Println("bestmove 0000")
	}
}
