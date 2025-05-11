package engine

import (
	"danielyang.cc/chess/internal/board"
)

var (
	killerHeuristic  = [64][2]int{}  // [depth][order]
	historyHeuristic = [12][64]int{} // [piece][square]
	MVV_LVA          = [12][12]int{
		{105, 205, 305, 405, 505, 605, 105, 205, 305, 405, 505, 605},
		{104, 204, 304, 404, 504, 604, 104, 204, 304, 404, 504, 604},
		{103, 203, 303, 403, 503, 603, 103, 203, 303, 403, 503, 603},
		{102, 202, 302, 402, 502, 602, 102, 202, 302, 402, 502, 602},
		{101, 201, 301, 401, 501, 601, 101, 201, 301, 401, 501, 601},
		{100, 200, 300, 400, 500, 600, 100, 200, 300, 400, 500, 600},
		{105, 205, 305, 405, 505, 605, 105, 205, 305, 405, 505, 605},
		{104, 204, 304, 404, 504, 604, 104, 204, 304, 404, 504, 604},
		{103, 203, 303, 403, 503, 603, 103, 203, 303, 403, 503, 603},
		{102, 202, 302, 402, 502, 602, 102, 202, 302, 402, 502, 602},
		{101, 201, 301, 401, 501, 601, 101, 201, 301, 401, 501, 601},
		{100, 200, 300, 400, 500, 600, 100, 200, 300, 400, 500, 600},
	} // [attacker][victim]
)

// insertion sort
func sortMoves(moves *board.MoveList, scores []int) {
	for i := 1; i < moves.Count; i++ {
		keyScore := scores[i]
		keyMove := moves.Moves[i]
		j := i - 1

		for j >= 0 && scores[j] < keyScore {
			scores[j+1] = scores[j]
			moves.Moves[j+1] = moves.Moves[j]
			j--
		}
		scores[j+1] = keyScore
		moves.Moves[j+1] = keyMove
	}
}

func scoreMove(move int, depth int) int {
	score := 0

	if board.GetCapture(move) > 0 {
		score += getMVVLVA(move)
	}

	if killerHeuristic[depth][0] == move {
		score += 9000
	}
	if killerHeuristic[depth][1] == move {
		score += 8000
	}

	score += historyHeuristic[board.GetPiece(move)][board.GetTarget(move)]

	return score
}

func getMVVLVA(move int) int {
	attacker := board.GetPiece(move)
	target := board.GetTarget(move)

	var victim int

	// set piece range to see which piece is being captured
	var start, end int
	if board.Side == board.WHITE {
		start = board.BLACK_PAWN
		end = board.BLACK_KING
	} else {
		start = board.WHITE_PAWN
		end = board.WHITE_KING
	}

	for i := start; i <= end; i++ {
		if board.GetBit(board.Bitboards[i], target) != 0 {
			victim = i
			break
		}
	}

	return MVV_LVA[attacker][victim]
}
