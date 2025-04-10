package engine

import (
	"math"
	"sort"

	"danielyang.cc/chess/internal/board"
	"danielyang.cc/chess/internal/transposition"
)

var killerMoves = map[int][2][2]int{}

func alphaBetaImpl(alpha float64, beta float64, depthLeft int, currentColor byte) ([2][2]int, float64) {
	nextColor := byte('w')
	if currentColor == 'w' {
		nextColor = 'b'
	}

	// reverse futility pruning
	eval := -Evaluate(nextColor)
	if depthLeft <= 3 {
		// margin based on depth
		margin := 150.0 * float64(depthLeft)
		if eval-margin >= beta {
			// prune this branch if the material balance + gain + margin doesn't improve alpha
			// fail soft
			return [2][2]int{}, eval
		}
	}

	// null move pruning
	if depthLeft >= 2 {
		// reduction factor
		const R = 2

		_, nullEval := alphaBetaImpl(-beta, -beta+1, depthLeft-1-R, nextColor)
		nullEval *= -1

		if nullEval >= beta {
			return [2][2]int{}, nullEval
		}
	}

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

	bestScore := math.Inf(-1)
	var bestMove [2][2]int

	var moves [][2][2]int

	// if transposition table has entry then sorted moves are in there otherwise sort moves
	if ok {
		moves = entry.SortedMoves
	} else {
		var n int
		moves, n = board.GetAllValidMoves(currentColor)

		if len(moves) > 2 {
			// score moves based on killer heuristic and mvv lva
			moveScores := make([]int, n)

			for i, move := range moves {
				score := 0

				// killer heuristic
				if killerMoves[depthLeft] == move {
					score += 10000
					break
				}

				// if capture add mvv-lva score
				if board.IsCapture(move) {
					attacker := board.Board[move[0][0]][move[0][1]]
					victim := board.Board[move[1][0]][move[1][1]]
					score += int(pieceWeights[victim[1]]*13 - pieceWeights[attacker[1]])
				}

				moveScores[i] = score
			}

			sort.SliceStable(moves, func(i, j int) bool {
				return moveScores[i] > moveScores[j]
			})
		}

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
			// this is a new found killer so we can add it and use it in our move ordering
			if killerMoves[depthLeft] != move {
				killerMoves[depthLeft] = move
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

	transposition.AddEntry(nodeType, bestMove, bestScore, depthLeft, currentColor, moves)

	return bestMove, bestScore
}
