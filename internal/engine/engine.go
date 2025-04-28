package engine

import (
	"time"

	"danielyang.cc/chess/internal/board"
	"danielyang.cc/chess/internal/transposition"
)

var (
	type_       byte
	engineColor int
	searchDepth int
)

func Init(t byte, c int, d int) {
	type_ = t
	engineColor = c
	searchDepth = d
	transposition.Init()
	InitEvalTables()
}

func MakeMove() (string, int) {
	start := time.Now()

	switch type_ {
	case 'r':
		return makeRandomMove(), timeSince(start)
	case 'n':
		return iterativeDeepening(), timeSince(start)
	}

	return "", 0
}

func timeSince(start time.Time) int {
	return int(time.Since(start).Milliseconds())
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
