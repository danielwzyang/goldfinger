package engine

import (
	"math"

	"danielyang.cc/chess/internal/board"
)

var (
	pieceWeights = map[byte]float64{
		'P': 1.0,
		'N': 3.0,
		'B': 3.0,
		'R': 5.0,
		'Q': 9.0,
	}
	pawnStructureWeight = 0.5
	mobilityWeight      = 0.1
)

func Evaluate(color byte) float64 {
	multiplier := 1.0
	if color == 'b' {
		multiplier = -1.0
	}

	if board.Checkmate('w') || board.Checkmate('b') {
		return math.Inf(int(multiplier))
	}

	if board.Draw('w') || board.Draw('b') {
		return 0
	}

	return (material() + pawnStructure() + mobility()) * multiplier
}

func material() float64 {
	score := 0.0

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

func mobility() float64 {
	_, wn := board.GetAllValidMoves('w')
	_, bn := board.GetAllValidMoves('b')

	return float64(wn-bn) * mobilityWeight
}

func pawnStructure() float64 {
	score := 0.0

	for c := 0; c < 8; c++ {
		wCount := 0.0
		bCount := 0.0

		for r := 0; r < 8; r++ {
			piece := board.Board[r][c]
			if piece == "wP" {
				wCount++

				// blocked pawn (piece in front)
				if r > 0 && board.Board[r-1][c] != " " {
					score++
				}

				// isolated pawn (no pawns to help)
				if (c == 0 || board.Board[r][c-1] == " " || board.Board[r-1][c-1] == " ") &&
					(c == 7 || board.Board[r][c+1] == " " || board.Board[r-1][c+1] == " ") {
					score++
				}
			} else if piece == "bP" {
				bCount++

				// blocked pawn (piece in front)
				if r < 7 && board.Board[r+1][c] != " " {
					score--
				}

				// isolated pawn (no pawns to help)
				if (c == 0 || board.Board[r][c-1] == " " || board.Board[r+1][c-1] == " ") &&
					(c == 7 || board.Board[r][c+1] == " " || board.Board[r+1][c+1] == " ") {
					score--
				}
			}
		}

		// doubled pawn (more than one pawn in file)
		if wCount > 1 {
			score += wCount - 1
		}
		if bCount > 1 {
			score += bCount - 1
		}
	}

	return score * pawnStructureWeight
}
