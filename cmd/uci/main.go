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
	plies = 0.0
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

			go func() {
				engine.TimeForMove = getTimeForMove(wtime, btime, winc, binc)

				bestMove, ms, depth, nodes := engine.FindMove()

				fmt.Printf("info depth %d time %d nodes %d nps %d\n", depth, ms, nodes, nodes*1000/ms)
				if bestMove != 0 {
					fmt.Printf("bestmove %s\n", board.MoveToString(bestMove))
				} else {
					fmt.Println("bestmove 0000") // no legal moves
				}
			}()
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

	// estimating around 40 moves per game with a minimum of 10 moves to end
	remainingMoves := max(20, 80.0-plies) / 2

	// reserve 2 seconds as overhead
	TimeForMove := (float64(timeLeft) / remainingMoves) + float64(increment)

	// don't allocate more than 50% of remaining time
	TimeForMove = min(TimeForMove, float64(timeLeft)*0.5)

	return int(TimeForMove)
}
