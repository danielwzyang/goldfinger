package engine

import (
	"danielyang.cc/chess/internal/board"
)

var (
	killerHeuristic  = [2][100][2]int{} // [side][depth][order] (side so no conflict when playing self)
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

const (
	PV_BONUS       = 20000
	CAPTURE_BONUS  = 10000
	KILLER_BONUS_1 = 9000
	KILLER_BONUS_2 = 8000
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
	if board.GetCapture(move) > 0 {
		return CAPTURE_BONUS + getMVVLVA(move)
	}

	if killerHeuristic[board.Side][depth][0] == move {
		return KILLER_BONUS_1
	}
	if killerHeuristic[board.Side][depth][1] == move {
		return KILLER_BONUS_2
	}

	return historyHeuristic[board.Side][board.GetPiece(move)][board.GetTarget(move)]
}

func getMVVLVA(move int) int {
	attacker := board.GetPiece(move)
	target := board.GetTarget(move)

	victim := board.GetPieceOnSquare(target)
	if victim == -1 {
		return 0
	}

	return MVV_LVA[attacker][victim]
}

func ResetHeuristics() {
	for i := range killerHeuristic {
		for j := range killerHeuristic[i] {
			for k := range killerHeuristic[i][j] {
				killerHeuristic[i][j][k] = 0
			}
		}
	}
	for i := range historyHeuristic {
		for j := range historyHeuristic[i] {
			for k := range historyHeuristic[i][j] {
				historyHeuristic[i][j][k] = 0
			}
		}
	}
}
