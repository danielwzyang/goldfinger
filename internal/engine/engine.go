package engine

import (
	"fmt"
	"time"

	"danielyang.cc/chess/internal/board"
	"danielyang.cc/chess/internal/transposition"
)

var (
	type_ byte
	color int
	depth int
)

func Init(t byte, c int, d int) {
	type_ = t
	color = c
	depth = d
	transposition.Init()
}

func MakeMove() (string, string) {
	start := time.Now()

	switch type_ {
	case 'r':
		return makeRandomMove(), timeSince(start)
	case 'n':
		return alphaBeta(), timeSince(start)
	}

	return "", ""
}

func timeSince(start time.Time) string {
	return fmt.Sprintf("%d ms", time.Since(start).Milliseconds())
}

func numericToAlgebraic(position board.Position) string {
	piece := ""

	if board.Board[position.Rank][position.File].Type == board.PAWN {
		piece = ""
	} else {
		piece = board.PIECE_LETTERS[board.Board[position.Rank][position.File].Type]
	}

	letter := rune(position.File + 97) // 'a' is 97
	number := rune(56 - position.Rank) // '8' is 56

	return piece + string(letter) + string(number)
}
