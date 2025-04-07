package engine

import (
	"math"
	"math/rand"

	"danielyang.cc/chess/internal/board"
)

func alphaBetaImpl(alpha float64, beta float64, depthLeft int, currentColor byte) ([2][2]int, float64) {
	if depthLeft == 0 {
		return [2][2]int{}, quiesce(alpha, beta, currentColor)
	}

	bestScore := math.Inf(-1)
	var bestMove [2][2]int

	moves, _ := board.GetAllValidMoves(currentColor)
	// shuffle moves
	for i := range moves {
		j := rand.Intn(i + 1)
		moves[i], moves[j] = moves[j], moves[i]
	}

	// save state
	boardState := board.Board
	kingPositions := [][2]int{board.WhiteKing, board.BlackKing}
	castleStates := [4]bool{board.WCastleKS, board.WCastleQS, board.BCastleKS, board.BCastleQS}
	enPassant := board.EnPassant

	nextColor := byte('w')
	if currentColor == 'w' {
		nextColor = 'b'
	}

	// null move pruning
	if depthLeft > 1 {
		boardState = board.Board
		_, nullMoveScore := alphaBetaImpl(-beta, -alpha, depthLeft-1, nextColor)
		nullMoveScore *= -1
		if nullMoveScore >= beta {
			// prune branch
			return [2][2]int{}, nullMoveScore
		}
	}

	for _, move := range moves {
		// make move while preserving things like castle state and en passant so they can be reset
		// naive approach is store copies of states

		// make move
		board.MakeMove(move[0][0], move[0][1], move[1][0], move[1][1])

		// pawn promotion
		if board.Board[move[1][0]][move[1][1]][1] == 'P' && (move[1][0] == 0 || move[1][0] == 7) {
			// automatically promote to queen
			board.Board[move[1][0]][move[1][1]] = string(currentColor) + "Q"
		}

		_, score := alphaBetaImpl(-beta, -alpha, depthLeft-1, nextColor)

		// one colors max is the other colors min
		score *= -1

		// reset states
		board.Board = boardState
		board.EnPassant = enPassant
		board.WhiteKing = kingPositions[0]
		board.BlackKing = kingPositions[1]
		board.WCastleKS = castleStates[0]
		board.WCastleQS = castleStates[1]
		board.BCastleKS = castleStates[2]
		board.BCastleQS = castleStates[3]

		// updating max
		if score > bestScore {
			bestScore = score
			if depthLeft == depth {
				bestMove = move
			}
		}

		if bestScore > alpha {
			alpha = bestScore
		}

		if alpha >= beta {
			break
		}
	}

	return bestMove, bestScore
}
