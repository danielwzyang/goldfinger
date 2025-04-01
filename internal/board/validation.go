package board

func GetPossiblePieces(color byte, piece byte, rf int, cf int) [][2]int {
	possiblePieces := [][2]int{}
	for r, x := range Board {
		for c, y := range x {
			if y[0] == color && y[1] == piece && IsValidMove(r, c, rf, cf) {
				MakeMove(r, c, rf, cf, true)

				if !inCheck(color) {
					possiblePieces = append(possiblePieces, [2]int{r, c})
				}

				MakeMove(rf, cf, r, c, true)
			}
		}
	}

	return possiblePieces
}

func ContainsPosition(positions [][2]int, position [2]int) bool {
	for _, item := range positions {
		if item[0] == position[0] && item[1] == position[1] {
			return true
		}
	}

	return false
}

func IsValidMove(r1 int, c1 int, r2 int, c2 int) bool {
	return ContainsPosition(getValidMoves(r1, c1), [2]int{r2, c2})
}

func getValidMoves(r int, c int) [][2]int {
	piece := Board[r][c]

	switch piece[1] {
	case 'P':
		return getValidPawnMoves(piece[0], r, c)
	case 'N':
		return getValidKnightMoves(piece[0], r, c)
	case 'B':
		return getValidBishopMoves(piece[0], r, c)
	case 'R':
		return getValidRookMoves(piece[0], r, c)
	case 'Q':
		return getValidQueenMoves(piece[0], r, c)
	case 'K':
		return getValidKingMoves(piece[0], r, c)
	}

	return [][2]int{}
}

func inCheck(color byte) bool {
	enemy := "b"
	if color == 'b' {
		enemy = "w"
	}

	king := WhiteKing
	if color == 'b' {
		king = BlackKing
	}

	// instead of going through every opponent piece to see if it can attack the king,
	// we can check enemy pieces from the king's position to save time

	possiblePawnAttacks := getValidPawnMoves(color, king[0], king[1])
	for _, pos := range possiblePawnAttacks {
		if Board[pos[0]][pos[1]] == enemy+"P" {
			return true
		}
	}

	possibleKnightAttacks := getValidKnightMoves(color, king[0], king[1])
	for _, pos := range possibleKnightAttacks {
		if Board[pos[0]][pos[1]] == enemy+"N" {
			return true
		}
	}

	possibleBishopAttacks := getValidBishopMoves(color, king[0], king[1])
	for _, pos := range possibleBishopAttacks {
		if Board[pos[0]][pos[1]] == enemy+"B" || Board[pos[0]][pos[1]] == enemy+"Q" {
			return true
		}
	}

	possibleRookAttacks := getValidRookMoves(color, king[0], king[1])
	for _, pos := range possibleRookAttacks {
		if Board[pos[0]][pos[1]] == enemy+"R" || Board[pos[0]][pos[1]] == enemy+"Q" {
			return true
		}
	}

	return false
}

func isValidSquare(color byte, r int, c int) bool {
	// out of bounds
	if !(r >= 0 && r < 8 && c >= 0 && c < 8) {
		return false
	}

	// empty or capture
	if color == 'b' {
		return Board[r][c] == " " || Board[r][c][0] == 'w'
	} else {
		return Board[r][c] == " " || Board[r][c][0] == 'b'
	}
}

func ValidKingSideCastle(color byte) bool {
	if inCheck(color) {
		return false
	}

	canCastle := WCastleKS
	row := 7
	if color == 'b' {
		row = 0
		canCastle = BCastleKS
	}

	// rook or king has moved or spaces aren't empty
	if !canCastle || Board[row][5] != " " || Board[row][6] != " " {
		return false
	}

	// checking if squares are in check

	// move king right one
	MakeMove(row, 4, row, 5, true)
	if inCheck(color) {
		MakeMove(row, 5, row, 4, true)
		return false
	}

	// move king right one
	MakeMove(row, 5, row, 6, true)
	if inCheck(color) {
		MakeMove(row, 6, row, 4, true)
		return false
	}

	// move king back
	MakeMove(row, 6, row, 4, true)

	return true
}

func ValidQueenSideCastle(color byte) bool {
	if inCheck(color) {
		return false
	}

	canCastle := WCastleQS
	row := 7
	if color == 'b' {
		row = 0
		canCastle = BCastleQS
	}

	// rook or king has moved or spaces aren't empty
	if !canCastle || Board[row][1] != " " || Board[row][2] != " " || Board[row][3] != " " {
		return false
	}

	// checking if squares are in check

	// move king left one
	MakeMove(row, 4, row, 3, true)
	if inCheck(color) {
		MakeMove(row, 3, row, 4, true)
		return false
	}

	// move king left one
	MakeMove(row, 3, row, 2, true)
	if inCheck(color) {
		MakeMove(row, 2, row, 4, true)
		return false
	}

	// move king back
	MakeMove(row, 2, row, 4, true)

	return true
}
