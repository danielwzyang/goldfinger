package board

func IsValidMove(move Move) bool {
	return ContainsMove(getValidMoves(move.From), move) && moveOutOfCheck(Move{move.From, move.To})
}

func getValidMoves(position Position) []Move {
	piece := Board[position.Rank][position.File]

	switch piece.Type {
	case PAWN:
		return getPawnMoves(position)
	case KNIGHT:
		return getKnightMoves(position)
	case BISHOP:
		return getBishopMoves(position)
	case ROOK:
		return getRookMoves(position)
	case QUEEN:
		return getQueenMoves(position)
	case KING:
		return getKingMoves(position)
	}

	return []Move{}
}

func moveOutOfCheck(move Move) bool {
	color := Board[move.From.Rank][move.From.File].Color

	MakeMove(move)

	outOfCheck := !InCheck(color)

	UndoMove()

	return outOfCheck
}

func InCheck(color int) bool {
	enemy := color ^ 1

	king := WhiteKing
	if color == BLACK {
		king = BlackKing
	}

	// instead of going through every opponent piece to see if it can attack the king,
	// we can check enemy pieces from the king's position to save time

	possibleKingAttacks := getKingMoves(king)
	for _, move := range possibleKingAttacks {
		to := Board[move.To.Rank][move.To.File]
		if to.Color == enemy && to.Type == KING {
			return true
		}
	}

	possiblePawnAttacks := getPawnMoves(king)
	for _, move := range possiblePawnAttacks {
		to := Board[move.To.Rank][move.To.File]
		if to.Color == enemy && to.Type == PAWN {
			return true
		}
	}

	possibleKnightAttacks := getKnightMoves(king)
	for _, move := range possibleKnightAttacks {
		to := Board[move.To.Rank][move.To.File]
		if to.Color == enemy && to.Type == KNIGHT {
			return true
		}
	}

	possibleBishopAttacks := getBishopMoves(king)
	for _, move := range possibleBishopAttacks {
		to := Board[move.To.Rank][move.To.File]
		if to.Color == enemy && (to.Type == BISHOP || to.Type == QUEEN) {
			return true
		}
	}

	possibleRookAttacks := getRookMoves(king)
	for _, move := range possibleRookAttacks {
		to := Board[move.To.Rank][move.To.File]
		if to.Color == enemy && (to.Type == ROOK || to.Type == QUEEN) {
			return true
		}
	}

	return false
}

func isValidPosition(position Position, color int) bool {
	r := position.Rank
	c := position.File

	// out of bounds
	if !(r >= 0 && r < 8 && c >= 0 && c < 8) {
		return false
	}

	// empty or capture
	return Board[r][c].Color != color
}

func ValidKingSideCastle(color int) bool {
	if InCheck(color) {
		return false
	}

	canCastle := WCastleKS
	row := 7
	if color == BLACK {
		row = 0
		canCastle = BCastleKS
	}

	// rook or king has moved or spaces aren't empty
	if !canCastle || Board[row][5].Type != EMPTY || Board[row][6].Type != EMPTY {
		return false
	}

	// checking if squares are in check

	// check right one
	if !moveOutOfCheck(Move{Position{row, 4}, Position{row, 5}}) {
		return false
	}

	// check right two
	if !moveOutOfCheck(Move{Position{row, 4}, Position{row, 6}}) {
		return false
	}

	return true
}

func ValidQueenSideCastle(color int) bool {
	if InCheck(color) {
		return false
	}

	canCastle := WCastleQS
	row := 7
	if color == BLACK {
		row = 0
		canCastle = BCastleQS
	}

	// rook or king has moved or spaces aren't empty
	if !canCastle || Board[row][1].Type != EMPTY || Board[row][2].Type != EMPTY || Board[row][3].Type != EMPTY {
		return false
	}

	// checking if squares are in check

	// check left one
	if !moveOutOfCheck(Move{Position{row, 4}, Position{row, 3}}) {
		return false
	}

	// check left two
	if !moveOutOfCheck(Move{Position{row, 4}, Position{row, 2}}) {
		return false
	}

	return true
}
