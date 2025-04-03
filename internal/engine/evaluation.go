package engine

import (
	"danielyang.cc/chess/internal/board"
)

var (
	pieceWeights = map[byte]int{
		'P': 1,
		'N': 3,
		'B': 3,
		'R': 5,
		'Q': 9,
	}
	mobilityWeight = 0.1
)

func Evaluate(color byte) int {
	multiplier := 1
	if color == 'b' {
		multiplier = -1
	}

	return (material() + mobility()) * multiplier
}

func material() int {
	score := 0

	for _, row := range board.Board {
		for _, piece := range row {
			if piece != " " && piece[1] != 'K' {
				if piece[0] == 'w' {
					score += pieceWeights[piece[1]]
				} else {
					score -= pieceWeights[piece[1]]
				}
			}
		}
	}

	return score
}

func mobility() int {
	_, wn := board.GetAllValidMoves('w')
	_, bn := board.GetAllValidMoves('w')
	return wn - bn
}
