package engine

import (
	"fmt"
	"math"
	"math/rand"
	"os"

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
	board.MakeMove(move)

	return numericToAlgebraic(move.To)
}

func alphaBeta() string {
	move, _ := alphaBetaImpl(math.MinInt, math.MaxInt, depth, color)

	// all moves lead to a loss so move is essentially empty
	if board.Board[move.From.Rank][move.From.File].Type == board.EMPTY {
		fmt.Println("The engine resigns.")
		os.Exit(0)
	}

	board.MakeMove(move)

	// pawn promotion
	if board.Board[move.To.Rank][move.To.File].Type == board.PAWN && (move.To.Rank == 0 || move.To.Rank == 7) {
		// automatically promote to queen
		board.Board[move.To.Rank][move.To.File] = board.Piece{
			Type:  board.QUEEN,
			Color: color,
			Key:   board.GetKey(board.QUEEN, color),
		}
	}

	return numericToAlgebraic(move.To)
}
