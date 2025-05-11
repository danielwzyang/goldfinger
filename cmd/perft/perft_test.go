package perft

import (
	"fmt"
	"testing"

	"danielyang.cc/chess/internal/board"
)

func TestPerft(t *testing.T) {
	board.Init()

	tests := []struct {
		depth int
		want  int
	}{
		{0, 1},
		{1, 20},
		{2, 400},
		{3, 8902},
		{4, 197281},
		{5, 4865609},
	}

	for _, test := range tests {
		board.ParseFEN(board.DEFAULT_BOARD)

		t.Run(fmt.Sprintf("depth=%d", test.depth), func(t *testing.T) {
			got := Perft(test.depth)

			if got != test.want {
				t.Errorf("Perft(%d) = %d; want %d", test.depth, got, test.want)
			}
		})
	}
}

func Perft(depth int) int {
	if depth == 0 {
		return 1
	}

	nodes := 0

	moves := board.MoveList{}
	board.GenerateAllMoves(&moves)

	for i := 0; i < moves.Count; i++ {
		move := moves.Moves[i]
		if !board.MakeMove(move, board.ALL_MOVES) {
			continue
		}

		nodes += Perft(depth - 1)
		board.RestoreState()
	}

	return nodes
}
