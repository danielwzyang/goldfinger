package board

import "math"

func MakeMove(move Move) {
	// add move to history
	BoardHistory = append(BoardHistory, BoardState{
		move,

		// piece a and b
		Board[move.From.Rank][move.From.File],
		Board[move.To.Rank][move.To.File],

		WCastleKS,
		WCastleQS,
		BCastleKS,
		BCastleQS,

		EnPassant,
	})
	HistoryLength++

	from := move.From
	to := move.To

	// update new position
	piece := Board[from.Rank][from.File]
	Board[to.Rank][to.File] = piece
	Board[from.Rank][from.File] = EMPTY_PIECE()

	// if moved king, update stored king position and invalidate castling
	if piece.Type == KING {
		if piece.Color == BLACK {
			BlackKing = Position{to.Rank, to.File}
			BCastleKS = false
			BCastleQS = false
		} else if piece.Color == WHITE {
			WhiteKing = Position{to.Rank, to.File}
			WCastleKS = false
			WCastleQS = false
		}
	}

	// if rook, invalidate castling for that side
	if piece.Type == ROOK {
		if piece.Color == BLACK {
			if from.File == 7 {
				BCastleKS = false
			}
			if from.File == 0 {
				BCastleQS = false
			}
		} else if piece.Color == WHITE {
			if from.File == 7 {
				WCastleKS = false
			}
			if from.File == 0 {
				WCastleQS = false
			}
		}
	}

	// if pawn handle en passant related things
	if piece.Type == PAWN {
		// pawn made double move
		if math.Abs(float64(to.Rank-from.Rank)) == 2 {
			// set en passant target to the space skipped over
			EnPassant = Position{(from.Rank + to.Rank) / 2, to.File}
			return
		} else if math.Abs(float64(to.Rank-from.Rank)) == 1 && math.Abs(float64(to.File-from.File)) == 1 &&
			to.Rank == EnPassant.Rank && to.File == EnPassant.File {
			// capturing en passant
			// moving diagonally so distance for rows and columns are both 1 and the space is the space for en passant

			var capturedPawnRow int
			if piece.Color == WHITE {
				// white pawn moves up so the captured pawn is one row below
				capturedPawnRow = to.Rank + 1
			} else {
				// black pawn moving down so the captured pawn is one row above
				capturedPawnRow = to.Rank - 1
			}

			// capture pawn
			Board[capturedPawnRow][to.File] = EMPTY_PIECE()

			// reset en passant target
			EnPassant = Position{-10, -10}
			return
		}
	}

	// any move that isn't double pawn invalidates en passant
	EnPassant = Position{-10, -10}
}

func UndoMove() {
	// remove last state
	lastState := BoardHistory[HistoryLength-1]
	move := lastState.LastMove
	BoardHistory = BoardHistory[:HistoryLength-1]
	HistoryLength--

	// reset king position
	if lastState.PieceA.Type == KING {
		if lastState.PieceA.Color == WHITE {
			WhiteKing = Position{move.From.Rank, move.From.File}
		} else {
			BlackKing = Position{move.From.Rank, move.From.File}
		}
	}

	// reset castling
	WCastleKS = lastState.WCastleKS
	WCastleQS = lastState.WCastleQS
	BCastleKS = lastState.BCastleKS
	BCastleQS = lastState.BCastleQS

	// en passant
	EnPassant = lastState.EnPassant

	// update board
	Board[move.From.Rank][move.From.File] = lastState.PieceA
	Board[move.To.Rank][move.To.File] = lastState.PieceB
}

func KingsideCastle(color int) {
	row := 7
	if color == BLACK {
		row = 0
	}

	MakeMove(Move{Position{row, 4}, Position{row, 6}}) // move king to g
	MakeMove(Move{Position{row, 7}, Position{row, 5}}) // move rook to f
}

func QueensideCastle(color int) {
	row := 7
	if color == BLACK {
		row = 0
	}

	MakeMove(Move{Position{row, 4}, Position{row, 2}}) // move king to c
	MakeMove(Move{Position{row, 0}, Position{row, 3}}) // move rook to d
}

func IsCapture(move Move) bool {
	return (Board[move.To.Rank][move.To.File].Type != EMPTY &&
		Board[move.From.Rank][move.From.File].Color != Board[move.To.Rank][move.To.File].Color)
}

func ContainsMove(moves []Move, move Move) bool {
	for _, item := range moves {
		if item.From.Rank == move.From.Rank && item.From.File == move.From.File &&
			item.To.Rank == move.To.Rank && item.To.File == move.To.File {
			return true
		}
	}

	return false
}

func ContainsPosition(positions []Position, position Position) bool {
	for _, item := range positions {
		if item.Rank == position.Rank && item.File == position.File {
			return true
		}
	}

	return false
}
