package engine

import (
	"math/rand"

	"danielyang.cc/chess/internal/board"
)

var (
	killerHeuristic  = [2][64][2]int{}  // [side][depth][order] (side so no conflict when playing self)
	historyHeuristic = [2][12][64]int{} // [side][piece][square]
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

	// mvv lva for captures
	if board.GetCapture(move) > 0 {
		score += getMVVLVA(move)
	}

	// killer heuristic bonus
	if killerHeuristic[board.Side][depth][0] == move {
		score += 9000
	}
	if killerHeuristic[board.Side][depth][1] == move {
		score += 8000
	}

	// heuristic heuristic bonus
	score += historyHeuristic[board.Side][board.GetPiece(move)][board.GetTarget(move)]

	// small random factor for variation
	// only for non captures and non killers
	if board.GetCapture(move) == 0 &&
		killerHeuristic[board.Side][depth][0] != move &&
		killerHeuristic[board.Side][depth][1] != move {
		// random value between 50 and -50
		score += rand.Intn(101) - 50
	}

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
