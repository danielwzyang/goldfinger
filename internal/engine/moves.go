package engine

import (
	"math"
	"math/rand"

	"danielyang.cc/chess/internal/board"
)

func makeRandomMove() string {
	options := 0

	kingsideCastle := board.ValidKingSideCastle(color)
	queensideCastle := board.ValidQueenSideCastle(color)

	if kingsideCastle {
		options++
	}

	if queensideCastle {
		options++
	}

	moves, n := board.GetAllValidMoves(color)

	options += n

	random := rand.Intn(options)

	if random == 0 && kingsideCastle {
		board.KingsideCastle(color)
		return "0-0"
	}

	if random == 1 && queensideCastle {
		board.QueensideCastle(color)
		return "0-0-0"
	}

	if kingsideCastle {
		random--
	}

	if queensideCastle {
		random--
	}

	move := moves[random]
	board.MakeMove(move[0][0], move[0][1], move[1][0], move[1][1])

	return numericToAlgebraic(move[1])
}

func alphaBeta() string {
	move, _ := alphaBetaImpl(math.MinInt, math.MaxInt, depth, color)
	board.MakeMove(move[0][0], move[0][1], move[1][0], move[1][1])
	// pawn automatically promote to queen
	if board.Board[move[1][0]][move[1][1]][1] == 'P' && (move[1][0] == 0 || move[1][0] == 7) {
		board.Board[move[1][0]][move[1][1]] = string(board.Board[move[1][0]][move[1][1]][0]) + "Q"
	}
	return numericToAlgebraic(move[1])
}
