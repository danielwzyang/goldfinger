package engine

import (
	"math"

	"danielyang.cc/chess/internal/board"
)

func iterativeDeepening() int {
	move := 0
	score := 0

	alpha := -board.LIMIT_SCORE
	beta := board.LIMIT_SCORE

	for depth := 1; depth <= searchDepth; depth++ {
		move, score = alphaBeta(alpha, beta, searchDepth)

		if score <= alpha || score >= beta {
			alpha = -board.LIMIT_SCORE
			beta = board.LIMIT_SCORE
			continue
		}

		alpha = score - 50
		beta = score + 50
	}

	return move
}

func alphaBeta(alpha, beta, depth int) (int, int) {
	// draw
	if depth == searchDepth && board.IsRepetition() || board.Fifty >= 100 {
		return 0, 0
	}

	// pv node
	pv := beta-alpha > 1

	// tt entry
	ttEntry, found := board.GetTTEntry()
	if depth != searchDepth && !pv && found && ttEntry.Depth >= depth {
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

	// quiesce
	if depth == 0 {
		return 0, quiesce(alpha, beta)
	}

	originalAlpha := alpha
	bestScore := -board.LIMIT_SCORE
	bestMove := 0

	var king int
	if board.Side == board.WHITE {
		king = board.WHITE_KING
	} else {
		king = board.BLACK_KING
	}

	inCheck := board.IsSquareAttacked(board.LS1B(board.Bitboards[king]), board.Side^1)

	// increase depth in check
	if inCheck {
		depth++
	}

	// null move pruning
	if depth >= 3 && depth != searchDepth && !inCheck && board.HasNonPawnMaterial() {
		board.MakeNullMove()

		board.REPETITION_INDEX++
		board.REPETITION_TABLE[board.REPETITION_INDEX] = board.ZobristHash

		// reduction factor = 2
		_, nullEval := alphaBeta(-beta, -beta+1, depth-1-2)
		nullEval *= -1

		board.REPETITION_INDEX--

		board.RestoreState()

		if nullEval >= beta {
			return 0, beta
		}
	}

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

	for moveCount, move := range moves.Moves {
		board.REPETITION_INDEX++
		board.REPETITION_TABLE[board.REPETITION_INDEX] = board.ZobristHash

		if !board.MakeMove(move, board.ALL_MOVES) {
			board.REPETITION_INDEX--
			continue
		}
		legalMoves++

		var reduction int
		if depth < 3 || moveCount <= 3 || inCheck {
			reduction = 0
		} else if board.GetPromotion(move) > 0 || board.GetCapture(move) > 0 {
			reduction = int(0.7 + 0.3*math.Log1p(float64(depth)) + 0.3*math.Log1p(float64(moveCount)))
		} else {
			reduction = int(1 + 0.5*math.Log1p(float64(depth)) + 0.7*math.Log1p(float64(moveCount)))
		}

		if reduction >= depth {
			reduction = depth - 1
		}

		_, score := alphaBeta(-beta, -alpha, depth-1-reduction)
		score *= -1

		board.REPETITION_INDEX--

		board.RestoreState()

		if score > alpha {
			bestScore = score
			bestMove = move

			if board.GetCapture(move) == 0 {
				historyHeuristic[board.GetPiece(move)][board.GetTarget(move)] += depth
			}

			alpha = score

			if score >= beta {
				// store cutoff in tt
				board.AddTTEntry(bestMove, beta, depth, board.CutNode)

				// killer heuristic
				if board.GetCapture(move) == 0 {
					killerHeuristic[depth][1] = killerHeuristic[depth][0]
					killerHeuristic[depth][0] = move
				}

				// fail high
				return move, beta
			}

		}
	}

	// no legal moves found
	if legalMoves == 0 {
		// checkmate
		if inCheck {
			return 0, -board.MATE + depth
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
