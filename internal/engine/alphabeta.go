package engine

import (
	"context"
	"math"

	"danielyang.cc/chess/internal/board"
)

func alphaBeta(ctx context.Context, alpha, beta, depth int) (int, int) {
	select {
	case <-ctx.Done():
		return 0, 0
	default:
	}

	nodes++
	ply++
	defer func() { ply-- }()

	if ply != 0 && board.IsRepetition() || board.Fifty >= 100 {
		return 0, 0
	}

	// pv node
	pv := beta-alpha > 1

	root := ply == 0

	// tt entry
	ttEntry, found := board.GetTTEntry()
	if !root && !pv && found && ttEntry.Depth >= depth {
		switch ttEntry.Type {
		case board.PVNode:
			return ttEntry.Move, ttEntry.Score
		case board.CutNode:
			if ttEntry.Score >= beta {
				return ttEntry.Move, ttEntry.Score
			}
		case board.AllNode:
			if ttEntry.Score <= alpha {
				return ttEntry.Move, ttEntry.Score
			}
		}
	}

	// no tt entry so reduce depth by 1 to save time for next iteration
	if !found && depth >= 4 {
		depth--
	}

	inCheck := board.InCheck()

	// increase depth in check
	if inCheck {
		depth++
	}

	// quiesce
	if depth <= 0 {
		return 0, quiesce(ctx, alpha, beta)
	}

	// null move pruning
	if depth >= 3 && ply != 0 && !inCheck {
		ply++

		board.MakeNullMove()

		// reduction factor = 2
		_, nullEval := alphaBeta(ctx, -beta, -beta+1, depth-1-2)
		nullEval *= -1

		board.RestoreState()

		ply--

		if nullEval >= beta {
			return 0, beta
		}
	}

	originalAlpha := alpha
	bestScore := -board.LIMIT_SCORE
	bestMove := 0

	moves := board.MoveList{}
	board.GenerateAllMoves(&moves)
	scores := make([]int, moves.Count)

	for i := 0; i < moves.Count; i++ {
		if found && moves.Moves[i] == ttEntry.Move {
			// pv move
			scores[i] = 20000
			continue
		}

		scores[i] = scoreMove(moves.Moves[i], depth)
	}

	sortMoves(&moves, scores)
	legalMoves := 0

	for moveCount := 0; moveCount < moves.Count; moveCount++ {
		move := moves.Moves[moveCount]

		if !board.MakeMove(move) {
			continue
		}

		legalMoves++

		var score int

		if legalMoves == 1 {
			_, score = alphaBeta(ctx, -beta, -alpha, depth-1)
			score = -score
		} else {
			// late move reduction
			reduction := 0

			if depth < 3 || legalMoves <= 4 || inCheck {
				reduction = 0
			} else if board.GetPromotion(move) > 0 || board.GetCapture(move) > 0 {
				reduction = int(0.7 + 0.3*math.Log1p(float64(depth)) + 0.3*math.Log1p(float64(moveCount)))
			} else {
				reduction = int(1 + 0.5*math.Log1p(float64(depth)) + 0.7*math.Log1p(float64(moveCount)))
			}

			_, score = alphaBeta(ctx, -alpha-1, -alpha, depth-1-reduction)
			score = -score

			// principal variation search
			if score > alpha && score < beta {
				_, score = alphaBeta(ctx, -beta, -alpha, depth-1)
				score = -score
			}
		}

		board.RestoreState()

		if score > bestScore {
			bestScore = score
			bestMove = move
		}

		if bestScore > alpha {
			alpha = bestScore
		}

		if alpha >= beta {
			if board.GetCapture(move) == 0 {
				historyHeuristic[board.Side][board.GetPiece(move)][board.GetTarget(move)] += depth * depth
			}

			killerHeuristic[board.Side][depth][1] = killerHeuristic[board.Side][depth][0]
			killerHeuristic[board.Side][depth][0] = move

			break
		}
	}

	// no legal moves found
	if legalMoves == 0 {
		// checkmate
		if inCheck {
			return 0, -board.MATE + ply
		}
		// stalemate
		return 0, 0
	}

	// update tt
	nodeType := board.AllNode
	if bestScore <= originalAlpha {
		nodeType = board.AllNode
	} else if bestScore >= beta {
		nodeType = board.CutNode
	} else {
		nodeType = board.PVNode
	}

	board.AddTTEntry(bestMove, bestScore, depth, nodeType)

	return bestMove, bestScore
}
