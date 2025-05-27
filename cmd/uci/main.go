package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"danielyang.cc/chess/internal/board"
	"danielyang.cc/chess/internal/engine"
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

			go func() {
				engine.FindMove(timeForMove, true)
			}()
		case "stop":
			if engine.Stop != nil {
				engine.Stop()
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

	// >= 1ms
	return max(1, int(timeForMove))
}
