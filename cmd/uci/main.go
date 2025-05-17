package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"danielyang.cc/chess/internal/board"
	"danielyang.cc/chess/internal/engine"
)

var (
	stopSearch = false
)

func main() {
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
			board.ParseFEN(board.DEFAULT_BOARD)
			stopSearch = false

		case "position":
			if len(tokens) < 2 {
				continue
			}

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
					board.MakeMove(board.StringToMove(moveStr), board.ALL_MOVES)
				}
			}

		case "go":
			stopSearch = false
			engine.SetStopFlag(false)

			var wtime, btime, winc, binc, depth int
			var infinite bool

			// parse time control params
			for i := 1; i < len(tokens); i++ {
				switch tokens[i] {
				case "infinite":
					infinite = true
				case "depth":
					if i+1 < len(tokens) {
						depth, _ = strconv.Atoi(tokens[i+1])
						i++
					}
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

			// determine search depth
			searchDepth := 6 // default depth
			if depth > 0 {
				searchDepth = depth
			} else if !infinite {
				// calculate search depth
				searchDepth = getSearchDepth(wtime, btime, winc, binc, board.Side)
			}

			// start search (go routine to allow for stop)
			go func() {
				// set depth
				engine.SetOptions(engine.Options{
					SearchDepth: searchDepth,
				})

				bestMove, _ := engine.FindMove()

				if !stopSearch {
					if bestMove != 0 {
						fmt.Printf("bestmove %s\n", board.MoveToString(bestMove))
					} else {
						fmt.Println("bestmove 0000") // no legal moves
					}
				}
			}()

		case "stop":
			stopSearch = true
			engine.SetStopFlag(true)

		case "quit":
			return
		}
	}
}

func getSearchDepth(wtime, btime, winc, binc int, side int) int {
	var timeLeft, increment int
	if side == board.WHITE {
		timeLeft = wtime
		increment = winc
	} else {
		timeLeft = btime
		increment = binc
	}

	if timeLeft == 0 {
		return 6 // default depth if no time control
	}

	// convert to seconds
	timeLeftSec := float64(timeLeft) / 1000.0
	incrementSec := float64(increment) / 1000.0

	// roughly estimating that ~40 moves left in average position
	remainingMoves := 40.0

	// calculate time left for move
	timeForMove := (timeLeftSec / remainingMoves) + incrementSec

	// values tuned based on performance
	switch {
	case timeForMove <= 1:
		return 7
	case timeForMove <= 5:
		return 8
	case timeForMove <= 10:
		return 9
	default:
		return 10
	}
}
