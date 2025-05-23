package engine

import (
	"danielyang.cc/chess/internal/board"
)

const maxSeeDepth = 32
const deltaMargin = 1050

var SIMPLE_PIECE_VALUES = [12]int{100, 320, 330, 500, 900, 20000, 100, 320, 330, 500, 900, 20000}

func quiesce(alpha, beta int) int {
	if stopFlag {
		return 0
	}

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
		if stopFlag {
			return alpha
		}

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
		if seeCapture(board.GetSource(move), board.GetTarget(move), board.Side) < 0 {
			continue
		}

		if !board.MakeMove(move) {
			continue
		}

		score := -quiesce(-beta, -alpha)

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

func see(square int, side int, depth int) int {
	if depth < 0 {
		return 0
	}

	// get smallest attacker of square
	attacker, from := board.GetSmallestAttacker(square, side)
	if attacker == -1 {
		return 0
	}

	move := board.EncodeMove(from, square, attacker, 0, 1, 0, 0, 0)
	if !board.MakeMove(move) {
		return 0
	}

	// recursively call see on the square
	captured := board.LastCapturedValue()
	value := captured - see(square, side^1, depth-1)

	board.RestoreState()

	if value < 0 {
		return 0
	}

	return value
}

func seeCapture(from int, to int, side int) int {
	attacker := board.GetPieceOnSquare(from)

	move := board.EncodeMove(from, to, attacker, 0, 1, 0, 0, 0)
	if !board.MakeMove(move) {
		return 0
	}

	// call see on the target square to start recursion
	captured := board.LastCapturedValue()
	value := captured - see(to, side^1, maxSeeDepth)

	board.RestoreState()

	return value
}
