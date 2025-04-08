package engine

import (
	"math"
	"sort"

	"danielyang.cc/chess/internal/board"
	"danielyang.cc/chess/internal/transposition"
)

func alphaBetaImpl(alpha float64, beta float64, depthLeft int, currentColor byte) ([2][2]int, float64) {
	// check transposition table
	entry, ok := transposition.GetEntry(currentColor)

	// if entry exists and the depth of the entry is at least at our depth or deeper
	if ok && entry.DepthLeft <= depthLeft {
		switch entry.Type {
		case transposition.PVNode:
			// exact score
			return entry.BestMove, entry.Score
		case transposition.AllNode:
			// upper bound
			if entry.Score <= alpha {
				return entry.BestMove, entry.Score
			}
		case transposition.CutNode:
			// lower bound
			if entry.Score >= beta {
				return entry.BestMove, entry.Score
			}
		}
	}

	// stabilize with quiescence
	if depthLeft == 0 {
		return [2][2]int{}, quiesce(alpha, beta, currentColor)
	}

	nextColor := byte('w')
	if currentColor == 'w' {
		nextColor = 'b'
	}

	// null move pruning
	if depthLeft > 1 {
		_, nullMoveScore := alphaBetaImpl(-beta, -alpha, depthLeft-1, nextColor)
		nullMoveScore *= -1
		if nullMoveScore >= beta {
			// prune branch
			return [2][2]int{}, nullMoveScore
		}
	}

	bestScore := math.Inf(-1)
	var bestMove [2][2]int
	foundMove := false

	moves, _ := board.GetAllValidMoves(currentColor)
	sort.Slice(moves, func(i int, j int) bool {
		return board.IsCapture(moves[i]) && !board.IsCapture(moves[j])
	})

	// save state
	kingPositions := [][2]int{board.WhiteKing, board.BlackKing}
	castleStates := [4]bool{board.WCastleKS, board.WCastleQS, board.BCastleKS, board.BCastleQS}
	enPassant := board.EnPassant
	originalAlpha := alpha

	for _, move := range moves {
		// make move while preserving things like castle state and en passant so they can be reset
		// naive approach is store copies of states

		// save piece
		movedPiece := board.Board[move[0][0]][move[0][1]]
		tempPiece := board.Board[move[1][0]][move[1][1]]

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
		board.Board[move[0][0]][move[0][1]] = movedPiece
		board.Board[move[1][0]][move[1][1]] = tempPiece
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
			bestMove = move
			foundMove = true
		}

		if bestScore > alpha {
			alpha = bestScore
		}

		if alpha >= beta {
			break
		}
	}

	// handle move not found

	if !foundMove {
		// this is likely checkmate or stalemate
		if board.InCheck(currentColor) {
			multiplier := 1
			if currentColor == 'b' {
				multiplier = -1
			}
			return [2][2]int{}, -math.Inf(multiplier)
		} else {
			return [2][2]int{}, 0
		}
	}

	var nodeType transposition.NodeType
	if bestScore <= originalAlpha {
		// upper bound
		nodeType = transposition.AllNode
	} else if bestScore >= beta {
		// lower bound
		nodeType = transposition.CutNode
	} else {
		// exact
		nodeType = transposition.PVNode
	}

	transposition.AddEntry(nodeType, bestMove, bestScore, depthLeft, currentColor)

	return bestMove, bestScore
}
