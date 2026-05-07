package engine

import (
	"context"
	"math"

	"danielyang.cc/chess/internal/board"
)

func alphaBeta(ctx context.Context, alpha, beta, depth int, allowNull bool) (int, int) {
	select {
	case <-ctx.Done():
		return 0, 0
	default:
	}

	nodes++
	ply++
	defer func() { ply-- }()

	// root node
	root := ply == 1

	// draws
	if (!root && board.IsRepetition()) || board.Fifty >= 100 || board.InsufficientMaterial() {
		return 0, 0
	}

	// pv node
	pv := beta-alpha > 1

	// tt entry
	ttEntry, found := board.GetTTEntry(ply)
	if !root && !pv && found && ttEntry.Depth >= depth && board.Fifty < 90 {
		switch ttEntry.Type {
		case board.PVNode:
			return 0, ttEntry.Score
		case board.CutNode:
			if ttEntry.Score >= beta {
				return 0, ttEntry.Score
			}
		case board.AllNode:
			if ttEntry.Score <= alpha {
				return 0, ttEntry.Score
			}
		}
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

	staticEval := board.Evaluate()

	// reverse futility pruning
	if !root && !pv && !inCheck && depth <= 8 && staticEval-100*depth >= beta {
		return 0, staticEval
	}

	// null move pruning
	if allowNull && depth >= 3 && !root && !inCheck && hasMajorOrMinorPiece() {
		board.MakeNullMove()

		// reduction factor is default 2 but 3 when depth >= 6
		r := 2
		if depth >= 6 {
			r = 3
		}
		// cannot make null moves in a row
		_, nullEval := alphaBeta(ctx, -beta, -beta+1, depth-1-r, false)
		nullEval *= -1

		board.RestoreState()

		if ctx.Err() != nil {
			return 0, 0
		}

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
			scores[i] = PV_BONUS
			continue
		}

		scores[i] = scoreMove(moves.Moves[i], ply)
	}

	sortMoves(&moves, scores)
	legalMoves := 0

	for moveCount := 0; moveCount < moves.Count; moveCount++ {
		move := moves.Moves[moveCount]

		// late move pruning
		if depth <= 4 && !pv && !inCheck &&
			moveCount > 3+depth*depth &&
			board.GetCapture(move) == 0 &&
			board.GetPromotion(move) == 0 {
			continue
		}

		if !board.MakeMove(move) {
			continue
		}

		legalMoves++

		var score int

		if legalMoves == 1 {
			_, score = alphaBeta(ctx, -beta, -alpha, depth-1, true)
			score = -score
		} else {
			// late move reduction
			reduction := 0

			if depth >= 3 && legalMoves > 4 && !inCheck &&
			board.GetCapture(move) == 0 && board.GetPromotion(move) == 0 {
				reduction = int(1 + 0.5*math.Log1p(float64(depth)) + 0.7*math.Log1p(float64(moveCount)))
			}

			_, score = alphaBeta(ctx, -alpha-1, -alpha, depth-1-reduction, true)
			score = -score

			// principal variation search
			if score > alpha && score < beta {
				_, score = alphaBeta(ctx, -beta, -alpha, depth-1, true)
				score = -score
			}
		}

		board.RestoreState()

		if ctx.Err() != nil {
			return bestMove, bestScore
		}

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

			killerHeuristic[board.Side][ply][1] = killerHeuristic[board.Side][ply][0]
			killerHeuristic[board.Side][ply][0] = move

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

	board.AddTTEntry(bestMove, bestScore, depth, ply, nodeType)

	return bestMove, bestScore
}

func hasMajorOrMinorPiece() bool {
	side := board.Side
	for p := side*6 + 1; p <= side*6+4; p++ {
		if board.Bitboards[p] != 0 {
			return true
		}
	}
	return false
}

