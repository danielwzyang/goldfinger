package engine

import (
	"context"

	"danielyang.cc/chess/internal/board"
)

const maxSeeDepth = 32
const deltaMargin = 1050

var SIMPLE_PIECE_VALUES = [12]int{100, 320, 330, 500, 900, 20000, 100, 320, 330, 500, 900, 20000}

func quiesce(ctx context.Context, alpha, beta int) int {
	select {
	case <-ctx.Done():
		return 0
	default:
	}

	nodes++

	standpat := board.Evaluate()

	if standpat >= beta {
		return beta
	}

	if standpat > alpha {
		alpha = standpat
	}

	moves := board.MoveList{}
	board.GenerateAllCaptures(&moves)

	scores := make([]int, moves.Count)
	for i := 0; i < moves.Count; i++ {
		scores[i] = getMVVLVA(moves.Moves[i])
	}

	sortMoves(&moves, scores)

	for moveCount := 0; moveCount < moves.Count; moveCount++ {
		move := moves.Moves[moveCount]

		victim := board.GetPieceOnSquare(board.GetTarget(move))
		victimValue := 0
		if victim == -1 && board.GetEnPassant(move) > 0 {
			// en passant (target square is empty but still a capture)
			if board.Side == board.WHITE {
				victim = board.BLACK_PAWN
			} else {
				victim = board.WHITE_PAWN
			}

			victimValue = SIMPLE_PIECE_VALUES[victim]
		} else if victim >= 0 && victim < 12 {
			victimValue = SIMPLE_PIECE_VALUES[victim]
		}

		// delta pruning
		if standpat+victimValue+deltaMargin < alpha {
			continue
		}

		// static exchange evaluation
		if seeCapture(move) < 0 {
			continue
		}

		if !board.MakeMove(move) {
			continue
		}

		score := -quiesce(ctx, -beta, -alpha)

		board.RestoreState()

		if score > alpha {
			alpha = score
			if score >= beta {
				return beta
			}
		}
	}

	return alpha
}

func see(square int, depth int) int {
	if depth < 0 {
		return 0
	}

	// get smallest attacker of square
	attacker, from := board.GetSmallestAttacker(square)
	if attacker == -1 {
		return 0
	}

	move := board.EncodeMove(from, square, attacker, 0, 1, 0, 0, 0)
	if !board.MakeMove(move) {
		return 0
	}

	// recursively call see on the square
	captured := board.LastCapturedValue()
	value := captured - see(square, depth-1)

	board.RestoreState()

	if value < 0 {
		return 0
	}

	return value
}

func seeCapture(move int) int {
	if !board.MakeMove(move) {
		return 0
	}

	// call see on the target square to start recursion
	captured := board.LastCapturedValue()
	to := board.GetTarget(move)
	value := captured - see(to, maxSeeDepth)

	board.RestoreState()

	return value
}
