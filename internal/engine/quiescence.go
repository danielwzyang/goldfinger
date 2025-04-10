package engine

import (
	"sort"

	"danielyang.cc/chess/internal/board"
)

func quiesce(alpha float64, beta float64, currentColor byte) float64 {
	standPat := Evaluate(currentColor)
	bestScore := standPat

	// fail hard
	if standPat >= beta {
		return standPat
	}

	// delta pruning; 200 centipawns
	delta := 2.0
	if standPat < alpha-delta {
		return alpha
	}

	if alpha < standPat {
		alpha = standPat
	}

	// generate captures
	captures := board.GetCaptureMoves(currentColor)

	sort.Slice(captures, func(i, j int) bool {
		// Example: prioritize higher value captures, for example: queen > rook > bishop > knight
		pieceValueI := pieceWeights[board.Board[captures[i][1][0]][captures[i][1][1]][1]]
		pieceValueJ := pieceWeights[board.Board[captures[j][1][0]][captures[j][1][1]][1]]
		return pieceValueI > pieceValueJ
	})

	// save state
	kingPositions := [][2]int{board.WhiteKing, board.BlackKing}
	castleStates := [4]bool{board.WCastleKS, board.WCastleQS, board.BCastleKS, board.BCastleQS}
	enPassant := board.EnPassant

	for _, capture := range captures {
		// save piece that's been captured
		movedPiece := board.Board[capture[0][0]][capture[0][1]]
		tempPiece := board.Board[capture[1][0]][capture[1][1]]

		board.MakeMove(capture[0][0], capture[0][1], capture[1][0], capture[1][1])

		nextColor := byte('w')
		if currentColor == 'w' {
			nextColor = 'b'
		}

		// pawn promotion
		if board.Board[capture[1][0]][capture[1][1]][1] == 'P' && (capture[1][0] == 0 || capture[1][0] == 7) {
			// automatically promote to queen
			board.Board[capture[1][0]][capture[1][1]] = string(currentColor) + "Q"
		}

		score := -quiesce(-beta, -alpha, nextColor)

		// reset states
		board.Board[capture[0][0]][capture[0][1]] = movedPiece
		board.Board[capture[1][0]][capture[1][1]] = tempPiece
		board.EnPassant = enPassant
		board.WhiteKing = kingPositions[0]
		board.BlackKing = kingPositions[1]
		board.WCastleKS = castleStates[0]
		board.WCastleQS = castleStates[1]
		board.BCastleKS = castleStates[2]
		board.BCastleQS = castleStates[3]

		if score >= beta {
			return score
		}
		if score > bestScore {
			bestScore = score
		}
		if score > alpha {
			alpha = score
		}
	}

	return bestScore
}
