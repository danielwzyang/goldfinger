package engine

import (
	"math"

	"danielyang.cc/chess/internal/board"
)

var (
	pieceWeights = map[int]float64{
		board.PAWN:   1.0,
		board.KNIGHT: 3.0,
		board.BISHOP: 3.0,
		board.ROOK:   5.0,
		board.QUEEN:  9.0,
	}
	pawnStructureWeight = 0.5
	mobilityWeight      = 0.1
)

func Evaluate(color int) float64 {
	if board.Draw(board.WHITE) || board.Draw(board.BLACK) {
		return 0
	}

	if board.Checkmate(board.WHITE) {
		if color == board.WHITE {
			return math.MinInt
		}
		return math.MaxInt
	}

	if board.Checkmate(board.BLACK) {
		if color == board.WHITE {
			return math.MaxInt
		}
		return math.MinInt
	}

	multiplier := 1.0
	if color == board.BLACK {
		multiplier = -1.0
	}

	return (material() + pawnStructure() + mobility()) * multiplier
}

func material() float64 {
	score := 0.0

	for _, row := range board.Board {
		for _, piece := range row {
			if piece.Type != board.KING {
				if piece.Color == board.WHITE {
					score += pieceWeights[piece.Type]
				} else {
					score -= pieceWeights[piece.Type]
				}
			}
		}
	}

	return score
}

func mobility() float64 {
	_, wn := board.GetAllValidMoves(board.WHITE)
	_, bn := board.GetAllValidMoves(board.BLACK)

	return float64(wn-bn) * mobilityWeight
}

func pawnStructure() float64 {
	score := 0.0

	for c := 0; c < 8; c++ {
		wCount := 0.0
		bCount := 0.0

		for r := 0; r < 8; r++ {
			piece := board.Board[r][c]
			if piece.Color == board.WHITE && piece.Type == board.PAWN {
				wCount++

				// blocked pawn (piece in front)
				if r > 0 && board.Board[r-1][c].Type != board.EMPTY {
					score++
				}

				// isolated pawn (no pawns to help)
				if (c == 0 || board.Board[r][c-1].Type != board.EMPTY || board.Board[r-1][c-1].Type != board.EMPTY) &&
					(c == 7 || board.Board[r][c+1].Type != board.EMPTY || board.Board[r-1][c+1].Type != board.EMPTY) {
					score++
				}
			} else if piece.Color == board.BLACK && piece.Type == board.PAWN {
				bCount++

				// blocked pawn (piece in front)
				if r < 7 && board.Board[r+1][c].Type != board.EMPTY {
					score--
				}

				// isolated pawn (no pawns to help)
				if (c == 0 || board.Board[r][c-1].Type != board.EMPTY || board.Board[r+1][c-1].Type != board.EMPTY) &&
					(c == 7 || board.Board[r][c+1].Type != board.EMPTY || board.Board[r+1][c+1].Type != board.EMPTY) {
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
