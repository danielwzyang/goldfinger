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
	plies      = 0.0
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
					plies++
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

				bestMove, ms := engine.FindMove()

				if !stopSearch {
					fmt.Printf("info depth %d time %d\n", searchDepth, ms)
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

func getSearchDepth(wtime, btime, winc, binc, side int) int {
	var timeLeft, increment int
	if side == board.WHITE {
		timeLeft = wtime
		increment = winc
	} else {
		timeLeft = btime
		increment = binc
	}

	// estimating around 40 moves per game with a minimum of 10 moves to end
	remainingPlies := max(20, 80.0-plies)

	timeForMove := (float64(timeLeft) / (remainingPlies / 2)) + float64(increment)

	// values tuned based on performance
	switch {
	case timeForMove >= 15000:
		// only in opening/midgame with a lot of time
		return 10
	case timeForMove >= 3000:
		// use in early endgame with good time
		return 9
	case timeForMove >= 500:
		// use for most positions
		return 8
	default:
		// fallback for time pressure
		return 7
	}
}
