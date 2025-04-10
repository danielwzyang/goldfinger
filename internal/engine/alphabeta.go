package engine

import (
	"math"
	"sort"

	"danielyang.cc/chess/internal/board"
	"danielyang.cc/chess/internal/transposition"
)

var (
	killerMoves      = map[int][][2][2]int{}
	sortedMovesCache = map[uint64][][2][2]int{}
)

func alphaBetaImpl(alpha float64, beta float64, depthLeft int, currentColor byte) ([2][2]int, float64) {
	// stabilize with quiescence
	if depthLeft <= 0 {
		return [2][2]int{}, quiesce(alpha, beta, currentColor)
	}

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

	nextColor := byte('w')
	if currentColor == 'w' {
		nextColor = 'b'
	}

	// null move pruning
	if depthLeft > 1 {
		_, nullMoveScore := alphaBetaImpl(-beta, -alpha, depthLeft-3, nextColor)
		nullMoveScore *= -1
		if nullMoveScore >= beta {
			// prune branch
			return [2][2]int{}, nullMoveScore
		}
	}

	bestScore := math.Inf(-1)
	var bestMove [2][2]int

	var moves [][2][2]int

	// see if moves are cached
	boardHash := transposition.HashBoard(currentColor)
	moves, found := sortedMovesCache[boardHash]

	if !found {
		moves, _ = board.GetAllValidMoves(currentColor)

		// if there's a pv move move it to the front
		if ok {
			pv := entry.BestMove
			for i, move := range moves {
				if move == pv {
					moves[0], moves[i] = moves[i], moves[0]
					break
				}
			}
		}

		if len(moves) > 1 {
			sort.SliceStable(moves[1:], func(i, j int) bool {
				moveI := moves[i+1]
				moveJ := moves[j+1]
				scoreI, scoreJ := 0, 0

				// killer heuristic bonus for i
				// if this move is stored in our killer moves then its a very strong move that cuts off possibilities for the enemy
				for _, km := range killerMoves[depthLeft] {
					if km == moveI {
						scoreI += 10000
						break
					}
				}

				// if capture add mvv-lva score
				if board.IsCapture(moveI) {
					attacker := board.Board[moveI[0][0]][moveI[0][1]]
					victim := board.Board[moveI[1][0]][moveI[1][1]]
					scoreI += int(pieceWeights[victim[1]]*10 - pieceWeights[attacker[1]])
				}

				// killer heuristic bonus for j
				for _, km := range killerMoves[depthLeft] {
					if km == moveJ {
						scoreJ += 10000
						break
					}
				}

				// if capture add mvv-lva score
				if board.IsCapture(moveJ) {
					attacker := board.Board[moveJ[0][0]][moveJ[0][1]]
					victim := board.Board[moveJ[1][0]][moveJ[1][1]]
					scoreJ += int(pieceWeights[victim[1]]*10 - pieceWeights[attacker[1]])
				}

				return scoreI > scoreJ
			})

			// cache sorted
			sortedMovesCache[boardHash] = moves
		}
	}

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

		// handle checkmate
		if board.Checkmate(currentColor) {
			return [2][2]int{}, math.Inf(1)
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
		}

		if bestScore > alpha {
			alpha = bestScore
		}

		// record killer move for current depth
		// killer move means that it causes a cutoff meaning
		// killer moves are strong moves because they reduces the number of possibilities
		if alpha >= beta {
			// this is a previously found killer so we don't have to add it
			foundKiller := false
			for _, km := range killerMoves[depthLeft] {
				if km == move {
					foundKiller = true
					break
				}
			}

			// this is a new found killer so we can add it and use it in our move ordering
			if !foundKiller {
				killerMoves[depthLeft] = append(killerMoves[depthLeft], move)
			}

			break
		}
	}

	// handle move not found

	if bestScore == math.Inf(-1) {
		// this is likely checkmate or stalemate
		if board.InCheck(currentColor) {
			return [2][2]int{}, math.Inf(-1)
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
