package engine

import (
	"fmt"
	"time"

	"danielyang.cc/chess/internal/board"
)

var (
	type_ byte
	color byte
)

func Init(t byte, c byte) {
	type_ = t
	color = c
}

func MakeMove() (string, string) {
	start := time.Now()

	switch type_ {
	case 'r':
		return makeRandomMove(), timeSince(start)
	case 'n':
		return negamax(), timeSince(start)
	}

	return "", ""
}

func timeSince(start time.Time) string {
	return fmt.Sprintf("%d ms", time.Since(start).Milliseconds())
}

func numericToAlgebraic(position [2]int) string {
	piece := ""

	if board.Board[position[0]][position[1]][1] == 'P' {
		piece = ""
	} else {
		piece = string(board.Board[position[0]][position[1]][1])
	}

	x := rune(position[1] + 97) // 'a' is 97
	y := rune(56 - position[0]) // '8' is 56

	return piece + string(x) + string(y)
}
