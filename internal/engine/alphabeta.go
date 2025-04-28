package engine

import (
	"math"

	"danielyang.cc/chess/internal/board"
	"danielyang.cc/chess/internal/transposition"
)

var killerMoves = map[int]board.Move{}
var historyHeuristic = [6][64]int{}

func alphaBeta(alpha int, beta int, depthLeft int, currentColor int) (board.Move, int) {
	nextColor := currentColor ^ 1

	// stabilize with quiescence
	if depthLeft <= 0 {
		return board.Move{}, quiesce(alpha, beta, currentColor)
	}

	// check transposition table
	entry, pv := transposition.GetEntry(currentColor)

	// if entry exists and this current search will not search deeper than the entry
	if pv && entry.DepthLeft >= depthLeft {
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

	// null move pruning
	if !pv && searchDepth > depthLeft && beta-alpha > 1 && !board.InCheck(currentColor) {
		// reduction factor
		const R = 2

		_, nullEval := alphaBeta(-beta, -beta+1, depthLeft-1-R, nextColor)
		nullEval *= -1

		if nullEval >= beta {
			return board.Move{}, nullEval
		}
	}

	bestScore := -board.LIMIT_SCORE
	var bestMove board.Move

	var moves []board.Move

	// if transposition table has entry then sorted moves are in there otherwise sort moves
	if pv {
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
					score += pieceWeights[victim.Type]*13 - pieceWeights[attacker.Type]
				} else {
					// history heuristic for quiet moves
					piece := board.Board[move.From.Rank][move.From.File].Type
					index := move.To.Rank*8 + move.To.File
					score += historyHeuristic[piece][index]
				}

				moveScores[i] = score
			}

			insertionSort(moves, moveScores)
		}

		// if there's a pv move move it to the front
		if pv {
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

	// search
	for moveCount, move := range moves {
		capture := board.IsCapture(move)
		// make move
		board.MakeMove(move)

		// pawn promotion
		promotion := board.Board[move.To.Rank][move.To.File].Type == board.PAWN && (move.To.Rank == 0 || move.To.Rank == 7)
		if promotion {
			// automatically promote to queen
			board.Board[move.To.Rank][move.To.File] = board.Piece{
				Type:  board.QUEEN,
				Color: engineColor,
				Key:   engineColor*6 + board.QUEEN + 1,
			}
		}

		// late move reduction
		var reduction int
		if depthLeft < 3 || moveCount <= 3 || board.InCheck(currentColor) {
			reduction = 0
		} else if promotion || capture {
			reduction = int(0.20 + math.Log(float64(depthLeft))*math.Log(float64(moveCount))/3.35)
		} else {
			reduction = int(1.35 + math.Log(float64(depthLeft))*math.Log(float64(moveCount))/2.75)
		}

		_, score := alphaBeta(-beta, -alpha, depthLeft-reduction-1, nextColor)

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
		// use history heuristic if quiet killer
		if alpha >= beta {
			// record quiet killer in history
			if !capture {
				piece := board.Board[move.From.Rank][move.From.File].Type
				index := move.To.Rank*8 + move.To.File
				historyHeuristic[piece][index] += depthLeft * depthLeft
			}

			// this is a new found killer so we can add it and use it in our move ordering
			if killerMoves[depthLeft] != move {
				killerMoves[depthLeft] = move
			}

			break
		}
	}

	// handle move not found
	if bestScore == -board.LIMIT_SCORE {
		// this is likely checkmate or stalemate
		if board.InCheck(currentColor) {
			return board.Move{}, -board.MATE_SCORE + depthLeft
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
