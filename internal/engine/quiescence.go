package engine

import (
	"danielyang.cc/chess/internal/board"
	"danielyang.cc/chess/internal/transposition"
)

var sortedMovesCache = map[uint64][][2][2]int{}

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

	// generate captures or get from cache
	var captures [][2][2]int

	cache, ok := sortedMovesCache[transposition.HashBoard(currentColor)]
	if ok {
		captures = cache
	} else {
		captures = board.GetCaptureMoves(currentColor)
	}

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
