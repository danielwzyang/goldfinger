package engine

import (
	"math"

	"danielyang.cc/chess/internal/board"
	"danielyang.cc/chess/internal/transposition"
)

var killerMoves = map[int]board.Move{}

func alphaBetaImpl(alpha float64, beta float64, depthLeft int, currentColor int) (board.Move, float64) {
	nextColor := currentColor ^ 1

	// reverse futility pruning
	eval := -Evaluate(nextColor)
	if depthLeft <= 3 {
		// margin based on depth
		margin := 150.0 * float64(depthLeft)
		if eval-margin >= beta {
			// prune this branch if the material balance + gain + margin doesn't improve alpha
			// fail soft
			return board.Move{}, eval
		}
	}

	// null move pruning
	if depthLeft >= 2 {
		// reduction factor
		const R = 2

		_, nullEval := alphaBetaImpl(-beta, -beta+1, depthLeft-1-R, nextColor)
		nullEval *= -1

		if nullEval >= beta {
			return board.Move{}, nullEval
		}
	}

	// stabilize with quiescence
	if depthLeft <= 0 {
		return board.Move{}, quiesce(alpha, beta, currentColor)
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
	var bestMove board.Move

	var moves []board.Move

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
					attacker := board.Board[move.From.Rank][move.From.File]
					victim := board.Board[move.To.Rank][move.To.File]
					score += int(pieceWeights[victim.Type]*13 - pieceWeights[attacker.Type])
				} else {
					// use psq for non capture
					pieceType := board.Board[move.From.Rank][move.From.File].Type
					index := move.To.Rank*8 + move.To.File
					if currentColor == board.BLACK {
						index = 63 - index
					}
					score += positionalPieceSquareTable[pieceType][index]
				}

				moveScores[i] = score
			}

			insertionSort(moves, moveScores)
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

	originalAlpha := alpha

	for i, move := range moves {
		isCapture := board.IsCapture(move)

		// make move
		board.MakeMove(move)

		// pawn promotion
		if board.Board[move.To.Rank][move.To.File].Type == board.PAWN && (move.To.Rank == 0 || move.To.Rank == 7) {
			// automatically promote to queen
			board.Board[move.To.Rank][move.To.File] = board.Piece{
				Type:  board.QUEEN,
				Color: color,
				Key:   color*6 + board.QUEEN + 1,
			}
		}

		// handle checkmate
		if board.Checkmate(currentColor) {
			return board.Move{}, math.Inf(1)
		}

		newDepth := depthLeft - 1

		if i >= 6 && !isCapture {
			newDepth = depthLeft / 3
		}

		_, score := alphaBetaImpl(-beta, -alpha, newDepth, nextColor)

		// one colors max is the other colors min
		score *= -1

		board.UndoMove()

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
			return board.Move{}, math.Inf(-1)
		} else {
			return board.Move{}, 0
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

	transposition.AddEntry(nodeType, bestMove, bestScore, depthLeft, moves, currentColor)

	return bestMove, bestScore
}

func insertionSort(moves []board.Move, moveScores []int) {
	for i := 1; i < len(moves); i++ {
		currentMove := moves[i]
		currentScore := moveScores[i]
		j := i - 1

		for j >= 0 && moveScores[j] < currentScore {
			moves[j+1] = moves[j]
			moveScores[j+1] = moveScores[j]
			j--
		}

		moves[j+1] = currentMove
		moveScores[j+1] = currentScore
	}
}
