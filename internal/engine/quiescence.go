package engine

import (
	"context"

	"danielyang.cc/chess/internal/board"
)

const deltaMargin = 1050

var SEE_PIECE_VALUES = [13]int{100, 320, 330, 500, 900, 20000, 100, 320, 330, 500, 900, 20000, 0}

var qScores [256]int

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

	scores := qScores[:moves.Count]
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

			victimValue = SEE_PIECE_VALUES[victim]
		} else if victim >= 0 && victim < 12 {
			victimValue = SEE_PIECE_VALUES[victim]
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

		board.UndoMove(move)

		if score > alpha {
			alpha = score
			if score >= beta {
				return beta
			}
		}
	}

	return alpha
}


func seeCapture(move int) int {
	source := board.GetSource(move)
	target := board.GetTarget(move)
	piece := board.GetPiece(move)

	var victimValue int
	if board.GetEnPassant(move) > 0 {
		victimValue = SEE_PIECE_VALUES[board.WHITE_PAWN]
	} else {
		victim := board.GetPieceOnSquare(target)
		if victim < 0 {
			return 0
		}
		victimValue = SEE_PIECE_VALUES[victim]
	}

	// occupancy after the moving piece leaves its source square
	occupancy := board.Occupancies[board.BOTH] ^ (1 << source)

	return victimValue - seeFrom(target, piece, occupancy, board.Side^1)
}

// best gain from recapture
func seeFrom(square int, lastPiece int, occupancy uint64, side int) int {
	attackerPiece, attackerSquare := lvaAttacker(square, side, occupancy)
	if attackerPiece < 0 {
		// no recapture
		return 0
	}

	occupancy = occupancy ^ (1 << attackerSquare)

	gain := SEE_PIECE_VALUES[lastPiece] - seeFrom(square, attackerPiece, occupancy, side^1)

	if gain < 0 {
		// losing recapture
		return 0
	}

	return gain
}

// least valuable attacker
func lvaAttacker(square int, side int, occupancy uint64) (int, int) {
	pawn := side*6 + board.WHITE_PAWN
	var attack uint64
	if side == board.WHITE {
		attack = board.PAWN_ATTACKS[board.BLACK][square] & board.Bitboards[pawn] & occupancy
	} else {
		attack = board.PAWN_ATTACKS[board.WHITE][square] & board.Bitboards[pawn] & occupancy
	}
	if attack != 0 {
		return pawn, board.LS1B(attack)
	}

	knight := side*6 + board.WHITE_KNIGHT
	attack = board.KNIGHT_ATTACKS[square] & board.Bitboards[knight] & occupancy
	if attack != 0 {
		return knight, board.LS1B(attack)
	}

	diagonal := board.GetBishopAttacks(square, occupancy)
	orthogonal := board.GetRookAttacks(square, occupancy)

	bishop := side*6 + board.WHITE_BISHOP
	attack = diagonal & board.Bitboards[bishop] & occupancy
	if attack != 0 {
		return bishop, board.LS1B(attack)
	}

	rook := side*6 + board.WHITE_ROOK
	attack = orthogonal & board.Bitboards[rook] & occupancy
	if attack != 0 {
		return rook, board.LS1B(attack)
	}

	queen := side*6 + board.WHITE_QUEEN
	attack = (diagonal | orthogonal) & board.Bitboards[queen] & occupancy
	if attack != 0 {
		return queen, board.LS1B(attack)
	}

	king := side*6 + board.WHITE_KING
	attack = board.KING_ATTACKS[square] & board.Bitboards[king] & occupancy
	if attack != 0 {
		return king, board.LS1B(attack)
	}

	return -1, -1
}

